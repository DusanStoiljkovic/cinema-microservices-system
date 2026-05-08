package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) GetAll(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order

	if err := repo.db.WithContext(ctx).Preload("Tickets").Find(&orders).Error; err != nil {
		return nil, utils.NewConflict("Failed to load orders", err)
	}

	if len(orders) == 0 {
		return nil, utils.NewNotFound(
			"Orders not found",
			errors.New("OrderRepository.GetAll -> empty list"),
		)
	}

	return orders, nil
}

func (repo *OrderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	order := &models.Order{}

	if err := repo.db.WithContext(ctx).First(order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Order not found", err)
		}

		return nil, utils.NewConflict("Failed to laod order", err)
	}

	return order, nil
}

func (repo *OrderRepository) GetByUserID(ctx context.Context, id uint) ([]*models.Order, error) {
	var orders []*models.Order

	if err := repo.db.WithContext(ctx).Where("user_id = ?", id).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Orders not found", err)
		}

		return nil, utils.NewConflict("Failed to laod order", err)
	}

	return orders, nil
}

func (repo *OrderRepository) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	if len(order.Tickets) == 0 {
		return nil, errors.New("order must have at least one ticket")
	}

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Provera da li ima dovoljno slobodnih mesta na toj projekciji
		var reservedTickets []*models.Ticket

		if err := tx.Where("projection_id = ?", order.Tickets[0].ProjectionID).Find(&reservedTickets).Error; err != nil {
			return err
		}

		projection := &models.Projection{}

		if err := tx.Select("id", "hall_id").First(projection, order.Tickets[0].ProjectionID).Error; err != nil {
			return err
		}

		hall := &models.Hall{}

		if err := tx.Select("id", "capacity").First(hall, projection.HallID).Error; err != nil {
			return err
		}

		if (int(hall.Capacity) - len(reservedTickets)) < len(order.Tickets) {
			return utils.NewConflict("There are no enough seats.", errors.New("OrderRepository.Create -> not enough seats for this projection"))
		}

		// Racuna TotalPrice
		var totalPrice uint = 0

		for i := range order.Tickets {
			ticket := &order.Tickets[i]

			projection := &models.Projection{}

			if err := tx.
				Select("id", "price").
				First(projection, ticket.ProjectionID).Error; err != nil {
				return err
			}

			totalPrice += uint(projection.Price)
		}

		order.TotalPrice = totalPrice

		// Kreiranje
		tickets := order.Tickets
		order.Tickets = nil

		// Pravimo prvo Order(kako bi u ticket mogli da ubacimo order_id)
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Dodeljujemo im order_id
		for i := range tickets {
			tickets[i].OrderID = order.ID
		}

		// Ovde radimo Batch Insert
		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		// Dodeljujemo tikete orderu
		order.Tickets = tickets

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (repo *OrderRepository) Pay(ctx context.Context, id uint) error {
	order, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	order.Status = models.OrderStatusPaid

	err = repo.db.WithContext(ctx).Save(order).Error
	if err != nil {
		return utils.NewConflict("Failed to pay", err)
	}

	return nil
}

func (repo *OrderRepository) Cancel(ctx context.Context, id uint) error {
	order, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	order.Status = models.OrderStatusCancelled

	err = repo.db.WithContext(ctx).Save(order).Error
	if err != nil {
		return utils.NewConflict("Failed to cancel", err)
	}

	return nil

}
func (repo *OrderRepository) Delete(ctx context.Context, id uint) error {

	order, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err = tx.Where("order_id = ?", order.ID).Delete(&order.Tickets).Error; err != nil {
			return err
		}

		if err = tx.Delete(&order, order.ID).Error; err != nil {
			return err
		}

		return nil
	})

	return nil
}

func (repo *OrderRepository) ExistsSameOrder(ctx context.Context, excludeID uint, order *models.Order) (bool, error) {
	var existingOrders []models.Order

	query := repo.db.
		WithContext(ctx).
		Preload("Tickets").
		Where(
			"user_id = ? AND status = ? AND total_price = ?",
			order.UserID,
			order.Status,
			order.TotalPrice,
		)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.Find(&existingOrders).Error

	if err == nil {
		counter := 0
		for _, ticket := range order.Tickets {

			_, err := utils.ExistsSameTicket(repo.db, ctx, 0, &ticket)
			if err == nil {
				counter++
			}

			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, err
			}
		}

		if len(order.Tickets) == counter {
			return true, nil
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, nil
}
