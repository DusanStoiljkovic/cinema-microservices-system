package service

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"
)

type ProjectionRepository interface {
	GetAll(ctx context.Context) ([]*models.Projection, error)
	GetByID(ctx context.Context, id uint) (*models.Projection, error)
	GetByMovieID(ctx context.Context, id uint) ([]*models.Projection, error)
	Create(ctx context.Context, projection *models.Projection) (*models.Projection, error)
	Update(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error)
	Delete(ctx context.Context, id uint) error
}

type ProjectionService struct {
	repo ProjectionRepository
}

func NewProjectionService(repo ProjectionRepository) *ProjectionService {
	return &ProjectionService{repo: repo}
}

func (service *ProjectionService) GetAllProjections(ctx context.Context) ([]*models.Projection, error) {
	return service.repo.GetAll(ctx)
}

func (service *ProjectionService) GetProjectionByID(ctx context.Context, id uint) (*models.Projection, error) {
	return service.repo.GetByID(ctx, id)
}

func (service *ProjectionService) GetProjectionsByMovieID(ctx context.Context, movieID uint) ([]*models.Projection, error) {
	return service.repo.GetByMovieID(ctx, movieID)
}

func (service *ProjectionService) CreateProjection(ctx context.Context, projection *models.Projection) (*models.Projection, error) {
	if err := validateProjection(projection); err != nil {
		return nil, err
	}

	return service.repo.Create(ctx, projection)
}

func (service *ProjectionService) UpdateProjection(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid projection id",
			errors.New("ProjectionService.UpdateProjection -> id is zero"),
		)
	}

	if err := validateProjection(projection); err != nil {
		return nil, err
	}

	return service.repo.Update(ctx, id, projection)
}

func (service *ProjectionService) DeleteProjection(ctx context.Context, id uint) error {
	if id == 0 {
		return utils.NewInvalidInput(
			"Invalid projection id",
			errors.New("ProjectionService.DeleteProjection -> id is zero"),
		)
	}

	return service.repo.Delete(ctx, id)
}

func validateProjection(projection *models.Projection) error {
	if projection == nil {
		return utils.NewInvalidInput(
			"Invalid projection data",
			errors.New("validateProjection -> projection is nil"),
		)
	}

	if projection.MovieID == 0 {
		return utils.NewInvalidInput(
			"Movie id is required",
			errors.New("validateProjection -> movie id is zero"),
		)
	}

	if projection.HallID == 0 {
		return utils.NewInvalidInput(
			"Hall id is required",
			errors.New("validateProjection -> hall id is zero"),
		)
	}

	if projection.StartTime.IsZero() {
		return utils.NewInvalidInput(
			"Start time is required",
			errors.New("validateProjection -> start time is zero"),
		)
	}

	if projection.EndTime.IsZero() {
		return utils.NewInvalidInput(
			"End time is required",
			errors.New("validateProjection -> end time is zero"),
		)
	}

	if !projection.EndTime.After(projection.StartTime) {
		return utils.NewInvalidInput(
			"End time must be after start time",
			errors.New("validateProjection -> end time is not after start time"),
		)
	}

	if projection.Price <= 0 {
		return utils.NewInvalidInput(
			"Price must be greater than zero",
			errors.New("validateProjection -> price must be greater than zero"),
		)
	}

	return nil
}
