package service

import (
	"context"
	"errors"
	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/utils"
	"strings"

	"gorm.io/gorm"
)

type MovieRepository interface {
	GetMovies(
		ctx context.Context,
		limit, offset int,
		sort string,
		genre string,
		minYear, maxYear int,
		minRating float64,
	) ([]*models.Movie, error)

	GetMovieByID(ctx context.Context, id uint) (*models.Movie, error)
	GetRelationsByMovieID(ctx context.Context, id uint) ([]models.Genre, error)
	Create(ctx context.Context, movie *models.Movie) (*models.Movie, error)
	CreateRelation(ctx context.Context, movie *models.Movie, genreID *models.Genre) (*models.Movie, error)
	Update(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error)
	Delete(ctx context.Context, id uint) error
	DeleteRelation(ctx context.Context, movieID, genreID uint) error
}

type MovieService struct {
	repo      MovieRepository
	genreRepo GenreRepository
}

func NewMovieService(repo MovieRepository, genreRepo GenreRepository) *MovieService {
	return &MovieService{repo: repo, genreRepo: genreRepo}
}

func (service *MovieService) GetMovies(
	ctx context.Context,
	limit, offset int,
	sort string,
	genre string,
	minYear, maxYear int,
	minRating float64,
) ([]*models.Movie, error) {
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
		return nil, utils.ErrInvalidInput
	}

	if minRating < 0 || minRating > 10 {
		return nil, utils.ErrInvalidInput
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
		return nil, utils.ErrInvalidInput
	}

	movie, err := service.repo.GetMovieByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, utils.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}

		return nil, err
	}

	return movie, nil
}

func (service *MovieService) GetRelationsByMovieID(ctx context.Context, id uint) ([]models.Genre, error) {
	if id == 0 {
		return nil, utils.ErrGenreNotFound
	}

	relations, err := service.repo.GetRelationsByMovieID(ctx, id)
	if err != nil {
		return nil, err
	}

	return relations, nil
}

func (service *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	trimSpace(movie)

	createdMovie, err := service.repo.Create(ctx, movie)
	if err != nil {
		if errors.Is(err, utils.ErrGenreNotFound) {
			return nil, utils.ErrInvalidInput
		}

		return nil, err
	}

	return createdMovie, nil
}

func (service *MovieService) CreateRelation(ctx context.Context, movieID, genreID uint) (*models.Movie, error) {
	if movieID == 0 || genreID == 0 {
		return nil, utils.ErrInvalidInput
	}

	existMovie, err := service.repo.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}

	existGenre, err := service.genreRepo.GetByFilter(ctx, &dto.GenreFilter{ID: &genreID})
	if err != nil {
		return nil, err
	}

	if existMovie == nil || existGenre == nil {
		return nil, utils.ErrRecordNotFound
	}

	movie, err := service.repo.CreateRelation(ctx, existMovie, existGenre)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (service *MovieService) UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error) {
	if id == 0 {
		return nil, utils.ErrInvalidInput
	}

	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	trimSpace(movie)

	updatedMovie, err := service.repo.Update(ctx, id, movie)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, utils.ErrRecordNotFound) {
			return nil, utils.ErrNotFound
		}

		if errors.Is(err, utils.ErrGenreNotFound) {
			return nil, utils.ErrInvalidInput
		}

		return nil, err
	}

	return updatedMovie, nil
}

func (service *MovieService) DeleteMovie(ctx context.Context, id uint) error {
	if id == 0 {
		return utils.ErrInvalidInput
	}

	err := service.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, utils.ErrRecordNotFound) {
			return utils.ErrNotFound
		}

		return err
	}

	return nil
}

func (service *MovieService) DeleteRelation(ctx context.Context, movieID, genreID uint) error {
	if movieID == 0 || genreID == 0 {
		return utils.ErrInvalidInput
	}

	existMovie, err := service.repo.GetMovieByID(ctx, movieID)
	if err != nil {
		return err
	}

	existGenre, err := service.genreRepo.GetByFilter(ctx, &dto.GenreFilter{ID: &genreID})
	if err != nil {
		return err
	}

	if existMovie == nil || existGenre == nil {
		return utils.ErrRecordNotFound
	}

	err = service.repo.DeleteRelation(ctx, movieID, genreID)
	if err != nil {
		return err
	}

	return nil
}

func validateMovie(movie *models.Movie) error {
	if movie == nil {
		return utils.ErrInvalidInput
	}

	if strings.TrimSpace(movie.Title) == "" {
		return utils.ErrInvalidInput
	}

	if strings.TrimSpace(movie.Description) == "" {
		return utils.ErrInvalidInput
	}

	if strings.TrimSpace(movie.ImageURL) == "" {
		return utils.ErrInvalidInput
	}

	if movie.Year < 1888 {
		return utils.ErrInvalidInput
	}

	if movie.Duration <= 0 {
		return utils.ErrInvalidInput
	}

	if movie.Rating < 0 || movie.Rating > 10 {
		return utils.ErrInvalidInput
	}

	for _, genre := range movie.Genres {
		if genre.ID == 0 {
			return utils.ErrInvalidInput
		}
	}

	return nil
}

func trimSpace(movie *models.Movie) {
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
	case "title_desc":
		return "title DESC", nil
	default:
		return "", utils.ErrInvalidInput
	}
}
