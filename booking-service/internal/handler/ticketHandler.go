package handler

import (
	"booking-service/internal/dto"
	"booking-service/internal/mapper"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"encoding/json"
	"net/http"
)

type TicketService interface {
	GetAllTickets(ctx context.Context) ([]*models.Ticket, error)
	GetTicketByID(ctx context.Context, id uint) (*models.Ticket, error)
	GetTicketByUserID(ctx context.Context, id uint) ([]*models.Ticket, error)
	GetTicketByProjectionID(ctx context.Context, id uint) ([]*models.Ticket, error)
	CreateTicket(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error)
	DeleteTicket(ctx context.Context, id uint) error
}

type TicketHandler struct {
	service TicketService
}

func NewTicketHandler(service TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (handler *TicketHandler) HandleGetAllTickets(w http.ResponseWriter, r *http.Request) error {

	tickets, err := handler.service.GetAllTickets(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.TicketsToResponse(tickets))
}

func (handler *TicketHandler) HandleGetTicketByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid ticket id", err)
	}

	ticket, err := handler.service.GetTicketByID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.TicketToResponse(ticket))
}

func (handler *TicketHandler) HandleGetTicketsByUserID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid user id", err)
	}

	ticket, err := handler.service.GetTicketByUserID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.TicketsToResponse(ticket))
}

func (handler *TicketHandler) HandleGetTicketsByProjectionID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid projection id", err)
	}

	tickets, err := handler.service.GetTicketByProjectionID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.TicketsToResponse(tickets))
}

func (handler *TicketHandler) HandleCreateTicket(w http.ResponseWriter, r *http.Request) error {
	req := &dto.TicketRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	ticket := mapper.TicketFromRequest(req)

	createdTicket, err := handler.service.CreateTicket(r.Context(), ticket)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, mapper.TicketToResponse(createdTicket))
}

func (handler *TicketHandler) HandleDeleteTicket(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid projection id", err)
	}

	err = handler.service.DeleteTicket(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
