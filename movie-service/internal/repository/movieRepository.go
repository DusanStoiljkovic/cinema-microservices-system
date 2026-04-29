package repository

import (
	"context"
	"movie-service/internal/models"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (repo *MovieRepository) GetMovies(
	ctx context.Context,
	limit, offset int,
	sort string,
	genre string,
	minYear, maxYear int,
	minRating float64,
) ([]*models.Movie, error) {
	var movies []*models.Movie

	query := repo.db.WithContext(ctx).Model(&models.Movie{}).Preload("Genres")

	if genre != "" {
		query = query.Joins("JOIN movie_genres mg ON mg.movie_id = movies.id").
			Joins("JOIN genres g ON g.id = mg.genre_id").
			Where("g.name = ?", genre)
	}

	if minYear != 0 && maxYear != 0 {
		query = query.Where("year BETWEEN ? AND ?", minYear, maxYear)
	}

	if minRating != 0 {
		query = query.Where("rating >= ?", minRating)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	if offset >= 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (repo *MovieRepository) GetMovieByID(ctx context.Context, id uint) (*models.Movie, error) {
	return nil, nil
}

func (repo *MovieRepository) Create(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	return nil, nil
}

func (repo *MovieRepository) Update(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error) {
	return nil, nil
}

func (repo *MovieRepository) Delete(ctx context.Context, id uint) error {
	return nil
}
