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

type ProjectionService interface {
	GetAllProjections(ctx context.Context) ([]*models.Projection, error)
	GetProjectionByID(ctx context.Context, id uint) (*models.Projection, error)
	GetProjectionsByMovieID(ctx context.Context, id uint) ([]*models.Projection, error)
	CreateProjection(ctx context.Context, projection *models.Projection) (*models.Projection, error)
	UpdateProjection(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error)
	DeleteProjection(ctx context.Context, id uint) error
}

type ProjectionHandler struct {
	service ProjectionService
}

func NewProjectionHandler(service ProjectionService) *ProjectionHandler {
	return &ProjectionHandler{service: service}
}

func (handler *ProjectionHandler) HandleGetAllProjections(w http.ResponseWriter, r *http.Request) error {
	projections, err := handler.service.GetAllProjections(r.Context())
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, projections)
	return nil
}

func (handler *ProjectionHandler) HandleGetProjectionByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	projection, err := handler.service.GetProjectionByID(r.Context(), id)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, projection)
	return nil
}

func (handler *ProjectionHandler) HandleGetProjectionsByMovieID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	projections, err := handler.service.GetProjectionsByMovieID(r.Context(), id)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, projections)
	return nil
}

func (handler *ProjectionHandler) HandleCreateProjection(w http.ResponseWriter, r *http.Request) error {
	req := &dto.ProjectionRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.ErrInvalidInput
	}

	projection := mapper.ProjectionFromRequest(req)

	createdProjection, err := handler.service.CreateProjection(r.Context(), projection)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusCreated, createdProjection)
	return nil
}

func (handler *ProjectionHandler) HandleUpdateProjection(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	req := &dto.ProjectionRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.ErrInvalidInput
	}

	projection := mapper.ProjectionFromRequest(req)

	updatedProjection, err := handler.service.UpdateProjection(r.Context(), id, projection)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusCreated, updatedProjection)
	return nil
}

func (handler *ProjectionHandler) HandleDeleteProjection(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.ErrInvalidInput
	}

	return handler.service.DeleteProjection(r.Context(), id)
}
