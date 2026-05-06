package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (repo *TicketRepository) GetAll(ctx context.Context) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	if err := repo.db.WithContext(ctx).Find(&tickets).Error; err != nil {
		return nil, utils.NewConflict("Failed to load tickets", err)
	}

	if len(tickets) == 0 {
		return nil, utils.NewNotFound(
			"Tickets not found",
			errors.New("TicketRepository.GetAll -> empty list"),
		)
	}

	return tickets, nil
}

func (repo *TicketRepository) GetByID(ctx context.Context, id uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}

	if err := repo.db.WithContext(ctx).First(ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Ticket not found", err)
		}

		return nil, utils.NewConflict("Failed to load tickets", err)
	}

	return ticket, nil
}

func (repo *TicketRepository) GetByUserID(ctx context.Context, id uint) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	err := repo.db.
		WithContext(ctx).
		Model(&models.Ticket{}).
		Joins("INNER JOIN orders on orders.id = tickets.order_id").
		Where("orders.user_id = ?", id).
		Find(&tickets).
		Error

	if err != nil {
		return nil, utils.NewConflict("Failed to load tickets", err)
	}

	if len(tickets) == 0 {
		return nil, utils.NewNotFound("Tickets not found", nil)
	}

	return tickets, nil
}

func (repo *TicketRepository) GetByProjectionID(ctx context.Context, id uint) ([]*models.Ticket, error) {
	var tickets []*models.Ticket

	err := repo.db.
		WithContext(ctx).
		Model(&models.Ticket{}).
		Joins("INNER JOIN projections on projection.id = tickets.projection_id").
		Find(&tickets).
		Error

	if err != nil {
		return nil, utils.NewConflict("Failed to load tickets", err)
	}

	if len(tickets) == 0 {
		return nil, utils.NewNotFound("Tickets not found", nil)
	}

	return tickets, nil
}

func (repo *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	exists, err := repo.existsSameProjection(ctx, 0, ticket)

	if exists {
		return nil, utils.NewConflict(
			"Ticket already exists",
			errors.New("TicketRepository.Create -> duplicate ticket"),
		)
	}

	if err = repo.db.WithContext(ctx).Create(ticket).Error; err != nil {
		return nil, utils.NewInvalidInput(
			"Failed to create ticket",
			err,
		)
	}

	return ticket, nil
}

func (repo *TicketRepository) Delete(ctx context.Context, id uint) error {
	exists, err := repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if exists != nil {
		err := repo.db.WithContext(ctx).Delete(exists, id).Error
		if err != nil {
			return utils.NewConflict(
				"Failed to delete ticket",
				err,
			)
		}
	}

	return nil
}

func (repo *TicketRepository) existsSameProjection(ctx context.Context, excludeID uint, ticket *models.Ticket) (bool, error) {
	existingTicket := &models.Ticket{}

	query := repo.db.WithContext(ctx).Where("projection_id = ? AND order_id = ?, AND seat_number = ?",
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
