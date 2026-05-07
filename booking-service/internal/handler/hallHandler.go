package handler

import (
	"booking-service/internal/dto"
	"booking-service/internal/mapper"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HallService interface {
	GetAllHalls(ctx context.Context) ([]*models.Hall, error)
	GetHallByID(ctx context.Context, id uint) (*models.Hall, error)
	CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	UpdateHall(ctx context.Context, id uint, hall *models.Hall) (*models.Hall, error)
	DeleteHall(ctx context.Context, id uint) error
}

type HallHandler struct {
	service HallService
}

func NewHallHandler(service HallService) *HallHandler {
	return &HallHandler{service: service}
}

func (handler *HallHandler) HandleGetAllHalls(w http.ResponseWriter, r *http.Request) error {
	halls, err := handler.service.GetAllHalls(r.Context())
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, mapper.HallsToResponse(halls))
	return nil
}

func (handler *HallHandler) HandleGetHallByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid hall id", err)
	}

	hall, err := handler.service.GetHallByID(r.Context(), id)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, mapper.HallToResponse(hall))
	return nil
}

func (handler *HallHandler) HandleCreateHall(w http.ResponseWriter, r *http.Request) error {
	req := &dto.HallRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	hall := mapper.HallFromRequest(req)

	createdHall, err := handler.service.CreateHall(r.Context(), hall)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusCreated, mapper.HallToResponse(createdHall))
	return nil
}

func (handler *HallHandler) HandleUpdateHall(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid hall id", err)
	}

	req := &dto.HallRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	hall := mapper.HallFromRequest(req)

	updatedHall, err := handler.service.UpdateHall(r.Context(), id, hall)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, mapper.HallToResponse(updatedHall))
	return nil
}

func (handler *HallHandler) HandleDeleteHall(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid hall id", err)
	}

	if err := handler.service.DeleteHall(r.Context(), id); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parseParamID(r *http.Request, param string) (uint, error) {
	id := chi.URLParam(r, param)

	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}

	if parsedID == 0 {
		return 0, strconv.ErrSyntax
	}

	return uint(parsedID), nil
}
