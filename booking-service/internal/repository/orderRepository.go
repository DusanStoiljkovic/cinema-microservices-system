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
	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tickets := order.Tickets

		if err := tx.Omit("Tickets").Create(order).Error; err != nil {
			return utils.NewInvalidInput("Failed to create order", err)
		}

		for i := range tickets {
			tickets[i].OrderID = order.ID

			exists, err := utils.ExistsSameTicket(tx, ctx, 0, &tickets[i])
			if err != nil {
				return utils.NewConflict("Failed to check existing ticket", err)
			}

			if exists {
				return utils.NewConflict(
					"Seat already reserved",
					errors.New("OrderRepository.Create -> duplicate ticket"),
				)
			}

			if err := tx.Create(&tickets[i]).Error; err != nil {
				return utils.NewInvalidInput("Failed to create ticket", err)
			}
		}

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

	err = repo.db.WithContext(ctx).Delete(order, id).Error
	if err != nil {
		return utils.NewConflict("Failed to delete", err)
	}

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
