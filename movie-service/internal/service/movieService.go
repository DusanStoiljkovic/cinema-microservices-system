package service

import (
	"context"
	"errors"
	"strings"

	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/utils"
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
		return nil, utils.NewInvalidInput(
			"Invalid year range",
			errors.New("MovieService.GetMovies -> minYear is greater than maxYear"),
		)
	}

	if minRating < 0 || minRating > 10 {
		return nil, utils.NewInvalidInput(
			"Invalid minimum rating",
			errors.New("MovieService.GetMovies -> minRating must be between 0 and 10"),
		)
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
		return nil, utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.GetMovieByID -> id is zero"),
		)
	}

	return service.repo.GetMovieByID(ctx, id)
}

func (service *MovieService) GetRelationsByMovieID(ctx context.Context, id uint) ([]models.Genre, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.GetRelationsByMovieID -> id is zero"),
		)
	}

	return service.repo.GetRelationsByMovieID(ctx, id)
}

func (service *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	trimMovieSpaces(movie)

	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	return service.repo.Create(ctx, movie)
}

func (service *MovieService) CreateRelation(ctx context.Context, movieID, genreID uint) (*models.Movie, error) {
	if movieID == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.CreateRelation -> movie id is zero"),
		)
	}

	if genreID == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid genre id",
			errors.New("MovieService.CreateRelation -> genre id is zero"),
		)
	}

	existingMovie, err := service.repo.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}

	existingGenre, err := service.genreRepo.GetByFilter(ctx, &dto.GenreFilter{ID: &genreID})
	if err != nil {
		return nil, err
	}

	return service.repo.CreateRelation(ctx, existingMovie, existingGenre)
}

func (service *MovieService) UpdateMovie(ctx context.Context, id uint, movie *models.Movie) (*models.Movie, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.UpdateMovie -> id is zero"),
		)
	}

	trimMovieSpaces(movie)

	if err := validateMovie(movie); err != nil {
		return nil, err
	}

	return service.repo.Update(ctx, id, movie)
}

func (service *MovieService) DeleteMovie(ctx context.Context, id uint) error {
	if id == 0 {
		return utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.DeleteMovie -> id is zero"),
		)
	}

	return service.repo.Delete(ctx, id)
}

func (service *MovieService) DeleteRelation(ctx context.Context, movieID, genreID uint) error {
	if movieID == 0 {
		return utils.NewInvalidInput(
			"Invalid movie id",
			errors.New("MovieService.DeleteRelation -> movie id is zero"),
		)
	}

	if genreID == 0 {
		return utils.NewInvalidInput(
			"Invalid genre id",
			errors.New("MovieService.DeleteRelation -> genre id is zero"),
		)
	}

	if _, err := service.repo.GetMovieByID(ctx, movieID); err != nil {
		return err
	}

	if _, err := service.genreRepo.GetByFilter(ctx, &dto.GenreFilter{ID: &genreID}); err != nil {
		return err
	}

	return service.repo.DeleteRelation(ctx, movieID, genreID)
}

func validateMovie(movie *models.Movie) error {
	if movie == nil {
		return utils.NewInvalidInput(
			"Invalid movie data",
			errors.New("validateMovie -> movie is nil"),
		)
	}

	if strings.TrimSpace(movie.Title) == "" {
		return utils.NewInvalidInput(
			"Movie title is required",
			errors.New("validateMovie -> title is empty"),
		)
	}

	if strings.TrimSpace(movie.Description) == "" {
		return utils.NewInvalidInput(
			"Movie description is required",
			errors.New("validateMovie -> description is empty"),
		)
	}

	if strings.TrimSpace(movie.ImageURL) == "" {
		return utils.NewInvalidInput(
			"Movie image URL is required",
			errors.New("validateMovie -> image URL is empty"),
		)
	}

	if movie.Year < 1888 {
		return utils.NewInvalidInput(
			"Movie year is invalid",
			errors.New("validateMovie -> year is before 1888"),
		)
	}

	if movie.Duration <= 0 {
		return utils.NewInvalidInput(
			"Movie duration must be greater than zero",
			errors.New("validateMovie -> duration must be greater than zero"),
		)
	}

	if movie.Rating < 0 || movie.Rating > 10 {
		return utils.NewInvalidInput(
			"Movie rating must be between 0 and 10",
			errors.New("validateMovie -> rating must be between 0 and 10"),
		)
	}

	for _, genre := range movie.Genres {
		if genre.ID == 0 {
			return utils.NewInvalidInput(
				"Invalid genre id",
				errors.New("validateMovie -> genre id is zero"),
			)
		}
	}

	return nil
}

func trimMovieSpaces(movie *models.Movie) {
	if movie == nil {
		return
	}

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
		return "", utils.NewInvalidInput(
			"Invalid sort parameter",
			errors.New("mapSort -> unsupported sort value: "+sort),
		)
	}
}
