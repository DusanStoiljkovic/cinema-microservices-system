package service

import (
	"context"
	"movie-service/internal/models"
)

type MovieRepository interface {
	GetMovies(ctx context.Context,
		limit, offset int,
		sort string,
		genre string,
		minYear, maxYear int,
		minRating float64) ([]*models.Movie, error)
	GetMovieByID(ctx context.Context, id uint) (*models.Movie, error)
	Create(ctx context.Context, movie *models.Movie) (*models.Movie, error)
	Update(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error)
	Delete(ctx context.Context, id uint) error
}

type MovieService struct {
	repo MovieRepository
}

func NewMovieService(repo MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

func (service *MovieService) GetMovies(ctx context.Context,
	limit, offset int,
	sort string,
	genre string,
	minYear, maxYear int,
	minRating float64) ([]*models.Movie, error) {
	return service.repo.GetMovies(ctx, limit, offset, sort, genre, minYear, maxYear, minRating)
}

func (service *MovieService) GetMovieByID(ctx context.Context, id uint) (*models.Movie, error) {
	return nil, nil
}

func (service *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	return nil, nil
}

func (service *MovieService) UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error) {
	return nil, nil
}

func (service *MovieService) DeleteMovie(ctx context.Context, id uint) error {
	return nil
}
