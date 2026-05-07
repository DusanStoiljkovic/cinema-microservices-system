package utils

import (
	"booking-service/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

func ExistsSameTicket(tx *gorm.DB, ctx context.Context, excludeID uint, ticket *models.Ticket) (bool, error) {
	existingTicket := &models.Ticket{}

	query := tx.WithContext(ctx).Where("projection_id = ? AND order_id = ? AND seat_number = ?",
		ticket.ProjectionID, ticket.OrderID, ticket.SeatNumber)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.First(existingTicket).Error

	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}
