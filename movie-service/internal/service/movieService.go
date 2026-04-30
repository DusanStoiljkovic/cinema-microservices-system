package service

import (
	"context"
	"errors"
	"movie-service/internal/models"
	"movie-service/internal/repository"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
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
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	if minYear != 0 && maxYear != 0 && minYear > maxYear {
		return nil, ErrInvalidInput
	}

	if minRating < 0 || minRating > 10 {
		return nil, ErrInvalidInput
	}

	safeSort, err := mapSort(sort)
	if err != nil {
		return nil, err
	}

	genre = strings.TrimSpace(genre)

	return service.repo.GetMovies(ctx, limit, offset, safeSort, genre, minYear, maxYear, minRating)
}

func (service *MovieService) GetMovieByID(ctx context.Context, id uint) (*models.Movie, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	movie, err := service.repo.GetMovieByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return movie, err
}

func (service *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	normalizeMovie(movie)

	createdMovie, err := service.repo.Create(ctx, movie)
	if errors.Is(err, repository.ErrGenreNotFound) {
		return nil, ErrInvalidInput
	}

	return createdMovie, nil
}

func (service *MovieService) UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	normalizeMovie(movie)

	updatedMovie, err := service.repo.Update(ctx, id, movie)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	if errors.Is(err, repository.ErrGenreNotFound) {
		return nil, ErrNotFound
	}

	return updatedMovie, nil
}

func (service *MovieService) DeleteMovie(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}

	err := service.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrRecordNotFound) {
		return ErrNotFound
	}

	return err

}

func validateMovie(movie *models.Movie) error {
	if movie == nil {
		return ErrInvalidInput
	}

	if strings.TrimSpace(movie.Title) == "" {
		return ErrInvalidInput
	}

	if strings.TrimSpace(movie.Description) == "" {
		return ErrInvalidInput
	}

	if strings.TrimSpace(movie.ImageURL) == "" {
		return ErrInvalidInput
	}

	if movie.Year < 1888 {
		return ErrInvalidInput
	}

	if movie.Duration <= 0 {
		return ErrInvalidInput
	}

	if movie.Rating < 0 || movie.Rating > 10 {
		return ErrInvalidInput
	}

	for _, genre := range movie.Genres {
		if genre.ID == 0 {
			return ErrInvalidInput
		}
	}

	return nil
}

func normalizeMovie(movie *models.Movie) {
	movie.Title = strings.TrimSpace(movie.Title)
	movie.Description = strings.TrimSpace(movie.Description)
	movie.ImageURL = strings.TrimSpace(movie.ImageURL)

	movie.Genres = uniqueGenres(movie.Genres)
}

func uniqueGenres(genres []models.Genre) []models.Genre {
	if genres == nil {
		return nil
	}

	seen := make(map[uint]bool)
	unique := make([]models.Genre, 0, len(genres))

	for _, genre := range genres {
		if !seen[genre.ID] {
			seen[genre.ID] = true
			unique = append(unique, genre)
		}
	}

	return unique
}

func mapSort(sort string) (string, error) {
	switch sort {
	case "", "newest":
		return "created_at DESC", nil
	case "oldest":
		return "created_at ASC", nil
	case "rating_desc":
		return "rating DESC", nil
	case "rating_asc":
		return "rating ASC", nil
	case "year_desc":
		return "year DESC", nil
	case "year_asc":
		return "year ASC", nil
	case "title_asc":
		return "title ASC", nil
	case "title_dedsc":
		return "title DESC", nil
	default:
		return "", ErrInvalidInput
	}
}
