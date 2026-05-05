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

	if err := repo.db.
		WithContext(ctx).
		Preload("Hall").
		Find(&projections).Error; err != nil {
		return nil, utils.NewConflict("Failed to load projections", err)
	}

	if len(projections) == 0 {
		return nil, utils.NewNotFound(
			"Projections not found",
			errors.New("ProjectionRepository.GetAll -> empty list"),
		)
	}

	return projections, nil
}

func (repo *ProjectionRepository) GetByID(ctx context.Context, id uint) (*models.Projection, error) {
	projection := &models.Projection{}

	if err := repo.db.
		WithContext(ctx).
		Preload("Hall").
		First(projection, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound(
				"Projection not found",
				err,
			)
		}

		return nil, utils.NewConflict(
			"Failed to load projection",
			err,
		)
	}

	return projection, nil
}

func (repo *ProjectionRepository) GetByMovieID(ctx context.Context, movieID uint) ([]*models.Projection, error) {
	var projections []*models.Projection

	if err := repo.db.
		WithContext(ctx).
		Preload("Hall").
		Where("movie_id = ?", movieID).
		Find(&projections).Error; err != nil {
		return nil, utils.NewConflict(
			"Failed to load projections by movie",
			err,
		)
	}

	if len(projections) == 0 {
		return nil, utils.NewNotFound(
			"Projections not found for this movie",
			errors.New("ProjectionRepository.GetByMovieID -> empty list"),
		)
	}

	return projections, nil
}

func (repo *ProjectionRepository) Create(ctx context.Context, projection *models.Projection) (*models.Projection, error) {
	exists, err := repo.existsSameProjection(ctx, 0, projection)
	if err != nil {
		return nil, utils.NewConflict(
			"Failed to check existing projection",
			err,
		)
	}

	if exists {
		return nil, utils.NewConflict(
			"Projection already exists",
			errors.New("ProjectionRepository.Create -> duplicate projection"),
		)
	}

	if err := repo.db.WithContext(ctx).Create(projection).Error; err != nil {
		return nil, utils.NewInvalidInput(
			"Failed to create projection",
			err,
		)
	}

	if err := repo.db.
		WithContext(ctx).
		Preload("Hall").
		First(projection, projection.ID).Error; err != nil {
		return nil, utils.NewConflict(
			"Projection created, but failed to load hall data",
			err,
		)
	}

	return projection, nil
}

func (repo *ProjectionRepository) Update(ctx context.Context, id uint, projection *models.Projection) (*models.Projection, error) {
	existingProjection := &models.Projection{}

	if err := repo.db.WithContext(ctx).First(existingProjection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound(
				"Projection not found",
				err,
			)
		}

		return nil, utils.NewConflict(
			"Failed to load projection",
			err,
		)
	}

	exists, err := repo.existsSameProjection(ctx, id, projection)
	if err != nil {
		return nil, utils.NewConflict(
			"Failed to check existing projection",
			err,
		)
	}

	if exists {
		return nil, utils.NewConflict(
			"Projection already exists",
			errors.New("ProjectionRepository.Update -> duplicate projection"),
		)
	}

	existingProjection.MovieID = projection.MovieID
	existingProjection.HallID = projection.HallID
	existingProjection.StartTime = projection.StartTime
	existingProjection.EndTime = projection.EndTime
	existingProjection.Price = projection.Price

	if err := repo.db.WithContext(ctx).Save(existingProjection).Error; err != nil {
		return nil, utils.NewInvalidInput(
			"Failed to update projection",
			err,
		)
	}

	if err := repo.db.
		WithContext(ctx).
		Preload("Hall").
		First(existingProjection, id).Error; err != nil {
		return nil, utils.NewConflict(
			"Projection updated, but failed to load hall data",
			err,
		)
	}

	return existingProjection, nil
}

func (repo *ProjectionRepository) Delete(ctx context.Context, id uint) error {
	existingProjection := &models.Projection{}

	if err := repo.db.WithContext(ctx).First(existingProjection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewNotFound(
				"Projection not found",
				err,
			)
		}

		return utils.NewConflict(
			"Failed to load projection",
			err,
		)
	}

	if err := repo.db.WithContext(ctx).Delete(existingProjection).Error; err != nil {
		return utils.NewConflict(
			"Failed to delete projection",
			err,
		)
	}

	return nil
}

func (repo *ProjectionRepository) existsSameProjection(ctx context.Context, excludeID uint, projection *models.Projection) (bool, error) {
	existingProjection := &models.Projection{}

	query := repo.db.WithContext(ctx).
		Where(
			"movie_id = ? AND hall_id = ? AND start_time = ? AND end_time = ?",
			projection.MovieID,
			projection.HallID,
			projection.StartTime,
			projection.EndTime,
		)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.First(existingProjection).Error

	if err == nil {
		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}
