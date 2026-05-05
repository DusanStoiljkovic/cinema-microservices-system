package repository

import (
	"context"
	"errors"

	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/utils"

	"gorm.io/gorm"
)

type GenreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

func (repo *GenreRepository) GetAll(ctx context.Context) ([]*models.Genre, error) {
	var genres []*models.Genre

	if err := repo.db.WithContext(ctx).Find(&genres).Error; err != nil {
		return nil, utils.NewConflict("Failed to load genres", err)
	}

	if len(genres) == 0 {
		return nil, utils.NewNotFound(
			"Genres not found",
			errors.New("GenreRepository.GetAll -> empty list"),
		)
	}

	return genres, nil
}

func (repo *GenreRepository) GetByFilter(ctx context.Context, filter *dto.GenreFilter) (*models.Genre, error) {
	genre := &models.Genre{}

	query := repo.db.WithContext(ctx).Model(&models.Genre{})

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.Name != nil {
		query = query.Where("LOWER(name) = LOWER(?)", *filter.Name)
	}

	if err := query.First(genre).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Genre not found", err)
		}

		return nil, utils.NewConflict("Failed to load genre", err)
	}

	return genre, nil
}

func (repo *GenreRepository) Create(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := repo.db.WithContext(ctx).Create(genre).Error; err != nil {
		return nil, utils.NewInvalidInput("Failed to create genre", err)
	}

	return genre, nil
}

func (repo *GenreRepository) Update(ctx context.Context, updatedGenre *models.Genre, id uint) (*models.Genre, error) {
	genre := &models.Genre{}

	if err := repo.db.WithContext(ctx).First(genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Genre not found", err)
		}

		return nil, utils.NewConflict("Failed to load genre", err)
	}

	genre.Name = updatedGenre.Name

	if err := repo.db.WithContext(ctx).Save(genre).Error; err != nil {
		return nil, utils.NewConflict("Failed to update genre", err)
	}

	return genre, nil
}

func (repo *GenreRepository) Delete(ctx context.Context, id uint) error {
	genre := &models.Genre{}

	if err := repo.db.WithContext(ctx).First(genre, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewNotFound("Genre not found", err)
		}

		return utils.NewConflict("Failed to load genre", err)
	}

	if err := repo.db.WithContext(ctx).Delete(genre).Error; err != nil {
		return utils.NewConflict("Failed to delete genre", err)
	}

	return nil
}
