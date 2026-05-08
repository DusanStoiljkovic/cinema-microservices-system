package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		projectionID := order.Tickets[0].ProjectionID

		// Svi tiketi u jednom orderu moraju biti za istu projekciju
		for i := range order.Tickets {
			if order.Tickets[i].ProjectionID != projectionID {
				return utils.NewInvalidInput(
					"All tickets must be for the same projection",
					errors.New("OrderRepository.Create -> tickets have different projection_id"),
				)
			}
		}

		// Zakljucavamo konkretnu projekciju
		projection := &models.Projection{}

		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id", "hall_id", "price").
			First(projection, projectionID).Error; err != nil {
			return err
		}

		// Uzimamo kapacitet sale
		hall := &models.Hall{}

		if err := tx.
			Select("id", "capacity").
			First(hall, projection.HallID).Error; err != nil {
			return err
		}

		// Provera da li su poslata sedista validna i da nema duplikata u request
		requestedSeatsMap := make(map[uint]bool)
		requestedSeatNumbers := make([]uint, 0, len(order.Tickets))

		for i := range order.Tickets {
			seatNumber := order.Tickets[i].SeatNumber

			if seatNumber == 0 || seatNumber > hall.Capacity {
				return utils.NewInvalidInput(
					"Invalid seat number",
					errors.New("OrderRepository.Create -> invalid seat number"),
				)
			}

			if requestedSeatsMap[seatNumber] {
				return utils.NewInvalidInput(
					"Duplicate seat number in order",
					errors.New("OrderRepository.Create -> duplicate seat number in request"),
				)
			}

			requestedSeatsMap[seatNumber] = true
			requestedSeatNumbers = append(requestedSeatNumbers, seatNumber)
		}

		// Provera koliko vec ima rezervisanih tiketa za projekciju
		var reservedTicketsCount int64

		if err := tx.
			Model(&models.Ticket{}).
			Where("projection_id = ?", projectionID).
			Count(&reservedTicketsCount).Error; err != nil {
			return err
		}

		if int(hall.Capacity)-int(reservedTicketsCount) <= len(order.Tickets) {
			return utils.NewConflict(
				"There are not enough seats.",
				errors.New("OrderRepository.Create -> not enough seats for this projection"),
			)
		}

		// Provera da li su trazena sedišta vec zauzeta
		var takenSeats []uint

		if err := tx.
			Model(&models.Ticket{}).
			Where("projection_id = ? AND seat_number IN ?", projectionID, requestedSeatNumbers).
			Pluck("seat_number", &takenSeats).Error; err != nil {
			return err
		}

		if len(takenSeats) > 0 {
			return utils.NewConflict(
				fmt.Sprintf("Seat numbers are not available: %v", takenSeats),
				errors.New("OrderRepository.Create -> seat numbers are not available"),
			)
		}

		// Izracunavanje cene
		order.TotalPrice = uint(projection.Price) * uint(len(order.Tickets))

		// Kreiranje
		tickets := order.Tickets
		order.Tickets = nil

		// Prvo kreiramo order da bismo dobili order.ID
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Dodeljujemo order_id svakom tiketu
		for i := range tickets {
			tickets[i].OrderID = order.ID
		}

		// Batch insert tiketa
		if err := tx.Create(&tickets).Error; err != nil {
			return err
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
