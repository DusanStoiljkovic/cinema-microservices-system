package service

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
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

func (service *ProjectionService) GetProjectionsByMovieID(ctx context.Context, id uint) ([]*models.Projection, error) {
	return service.repo.GetByMovieID(ctx, id)
}

func (service *ProjectionService) CreateProjection(ctx context.Context, projection *models.Projection) (*models.Projection, error) {
	if err := validateProjection(projection); err != nil {
		return nil, err
	}

	return service.repo.Create(ctx, projection)
}

func (service *ProjectionService) UpdateProjection(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error) {
	if err := validateProjection(projection); err != nil {
		return nil, err
	}

	return service.repo.Update(ctx, id, projection)
}

func (service *ProjectionService) DeleteProjection(ctx context.Context, id uint) error {
	return service.repo.Delete(ctx, id)
}

func validateProjection(projection *models.Projection) error {
	if projection == nil {
		return utils.ErrInvalidInput
	}

	if projection.MovieID == 0 {
		return utils.ErrInvalidInput
	}

	if projection.HallID == 0 {
		return utils.ErrInvalidInput
	}

	if projection.StartTime.IsZero() {
		return utils.ErrInvalidInput
	}

	if projection.EndTime.IsZero() {
		return utils.ErrInvalidInput
	}

	if projection.Price <= 0 {
		return utils.ErrInvalidInput
	}

	return nil
}
