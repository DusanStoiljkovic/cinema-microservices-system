package repository

import (
	"context"
	"log"
	"movie-service/internal/dto"
	"movie-service/internal/models"

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

	err := repo.db.WithContext(ctx).Find(&genres).Error
	if err != nil {
		return nil, err
	}

	log.Println("Genres from repo: ", genres)

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
		return nil, err
	}

	return genre, nil
}

func (repo *GenreRepository) Create(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := repo.db.WithContext(ctx).Create(genre).Error; err != nil {
		return nil, err
	}

	return genre, nil
}

func (repo *GenreRepository) Update(ctx context.Context, updatedGenre *models.Genre, id uint) (*models.Genre, error) {
	genre := &models.Genre{}

	if err := repo.db.WithContext(ctx).First(genre, id).Error; err != nil {
		return nil, err
	}

	genre.Name = updatedGenre.Name

	if err := repo.db.WithContext(ctx).Save(genre).Error; err != nil {
		return nil, err
	}

	return genre, nil
}

func (repo *GenreRepository) Delete(ctx context.Context, id uint) error {
	if err := repo.db.WithContext(ctx).Delete(&models.Genre{}, id).Error; err != nil {
		return err
	}

	return nil
}
