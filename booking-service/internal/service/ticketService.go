package service

import (
	"booking-service/internal/models"
	"context"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]*models.Ticket, error)
	GetByID(ctx context.Context, id uint) (*models.Ticket, error)
	GetByUserID(ctx context.Context, id uint) ([]*models.Ticket, error)
	GetByProjectionID(ctx context.Context, id uint) ([]*models.Ticket, error)
	Create(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	Delete(ctx context.Context, id uint) error
}

type TicketService struct {
	repo TicketRepository
}

func NewTicketService(repo TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (service *TicketService) GetAllTickets(ctx context.Context) ([]*models.Ticket, error) {
	return service.repo.GetAll(ctx)
}

func (service *TicketService) GetTicketByID(ctx context.Context, id uint) (*models.Ticket, error) {
	return service.repo.GetByID(ctx, id)
}

func (service *TicketService) GetTicketByUserID(ctx context.Context, id uint) ([]*models.Ticket, error) {
	return service.repo.GetByUserID(ctx, id)
}

func (service *TicketService) GetTicketByProjectionID(ctx context.Context, id uint) ([]*models.Ticket, error) {
	return service.repo.GetByProjectionID(ctx, id)
}

func (service *TicketService) CreateTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	return service.repo.Create(ctx, ticket)
}

func (service *TicketService) DeleteTicket(ctx context.Context, id uint) error {
	return service.repo.Delete(ctx, id)
}
