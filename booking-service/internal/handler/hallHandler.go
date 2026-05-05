package handler

import (
	"booking-service/internal/dto"
	"booking-service/internal/mapper"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"encoding/json"
	"log"
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

	utils.WriteJSON(w, http.StatusOK, halls)
	return nil
}

func (handler *HallHandler) HandleGetHallByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	hall, err := handler.service.GetHallByID(r.Context(), id)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, hall)
	return nil
}

func (handler *HallHandler) HandleCreateHall(w http.ResponseWriter, r *http.Request) error {
	req := &dto.HallRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return utils.ErrInvalidInput
	}

	hall := mapper.HallFromRequest(req)
	createdHall, err := handler.service.CreateHall(r.Context(), hall)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusCreated, createdHall)
	return nil
}

func (handler *HallHandler) HandleUpdateHall(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	req := &dto.HallRequest{}

	err = json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return utils.ErrInvalidInput
	}

	hall := mapper.HallFromRequest(req)

	log.Println("Hall iz Handlera: ", hall)

	updatedHall, err := handler.service.UpdateHall(r.Context(), id, hall)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, updatedHall)
	return nil
}

func (handler *HallHandler) HandleDeleteHall(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	err = handler.service.DeleteHall(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func parseParamID(r *http.Request, param string) (uint, error) {
	id := chi.URLParam(r, param)

	idParam, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return uint(idParam), nil
}
