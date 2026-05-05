package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProjectionRepository struct {
	db *gorm.DB
}

func NewProjectionRepository(db *gorm.DB) *ProjectionRepository {
	return &ProjectionRepository{db: db}
}

func (repo *ProjectionRepository) GetAll(ctx context.Context) ([]*models.Projection, error) {
	var projections []*models.Projection

	if err := repo.db.WithContext(ctx).Preload("Hall").Find(&projections).Error; err != nil {
		return nil, utils.ErrConflict
	}

	if len(projections) == 0 {
		return nil, utils.ErrNotFound
	}

	return projections, nil
}

func (repo *ProjectionRepository) GetByID(ctx context.Context, id uint) (*models.Projection, error) {
	projection := &models.Projection{}

	if err := repo.db.WithContext(ctx).First(projection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}
		return nil, utils.ErrConflict
	}

	return projection, nil
}

func (repo *ProjectionRepository) GetByMovieID(ctx context.Context, id uint) ([]*models.Projection, error) {
	projections := []*models.Projection{}

	if err := repo.db.WithContext(ctx).Where("movie_id = ?", id).Find(&projections).Error; err != nil {
		return nil, utils.ErrConflict
	}

	if len(projections) == 0 {
		return nil, utils.ErrNotFound
	}

	return projections, nil
}

func (repo *ProjectionRepository) Create(ctx context.Context, projection *models.Projection) (*models.Projection, error) {
	existingProjection := &models.Projection{}

	err := repo.db.WithContext(ctx).
		Where("movie_id = ? AND hall_id = ? AND start_time = ? AND end_time = ?",
			projection.MovieID,
			projection.HallID,
			projection.StartTime,
			projection.EndTime,
		).First(existingProjection).Error

	if err == nil {
		return nil, utils.ErrConflict
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrConflict
	}

	if err := repo.db.WithContext(ctx).Create(projection).Error; err != nil {
		return nil, utils.ErrInvalidInput
	}

	return projection, nil
}

func (repo *ProjectionRepository) Update(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error) {
	existingProjection := &models.Projection{}

	err := repo.db.WithContext(ctx).First(existingProjection, id).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrConflict
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrNotFound
	}

	existingProjection.MovieID = projection.MovieID
	existingProjection.HallID = projection.HallID
	existingProjection.StartTime = projection.StartTime
	existingProjection.EndTime = projection.EndTime
	existingProjection.Price = projection.Price

	if err := repo.db.WithContext(ctx).Save(existingProjection).Error; err != nil {
		return nil, utils.ErrInvalidInput
	}

	return existingProjection, nil
}

func (repo *ProjectionRepository) Delete(ctx context.Context, id uint) error {
	existingProjection := &models.Projection{}

	err := repo.db.WithContext(ctx).First(existingProjection, id).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrConflict
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrNotFound
	}

	return repo.db.WithContext(ctx).Delete(existingProjection, id).Error
}
