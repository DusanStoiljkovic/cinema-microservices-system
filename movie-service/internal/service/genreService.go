package service

import (
	"context"
	"errors"
	"strings"

	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/utils"
)

type GenreRepository interface {
	GetAll(ctx context.Context) ([]*models.Genre, error)
	GetByFilter(ctx context.Context, req *dto.GenreFilter) (*models.Genre, error)
	Create(ctx context.Context, genre *models.Genre) (*models.Genre, error)
	Update(ctx context.Context, genre *models.Genre, id uint) (*models.Genre, error)
	Delete(ctx context.Context, id uint) error
}

type GenreService struct {
	repo GenreRepository
}

func NewGenreService(repo GenreRepository) *GenreService {
	return &GenreService{repo: repo}
}

func (service *GenreService) GetGenres(ctx context.Context) ([]*models.Genre, error) {
	return service.repo.GetAll(ctx)
}

func (service *GenreService) GetGenreByFilter(ctx context.Context, filter *dto.GenreFilter) (*models.Genre, error) {
	if filter == nil {
		return nil, utils.NewInvalidInput(
			"Invalid genre filter",
			errors.New("GenreService.GetGenreByFilter -> filter is nil"),
		)
	}

	return service.repo.GetByFilter(ctx, filter)
}

func (service *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	trimGenreSpaces(genre)

	if err := validateGenre(genre); err != nil {
		return nil, err
	}

	existingGenre, err := service.repo.GetByFilter(ctx, &dto.GenreFilter{
		Name: &genre.Name,
	})

	if err != nil && !isSafeNotFound(err) {
		return nil, err
	}

	if existingGenre != nil {
		return nil, utils.NewConflict(
			"Genre already exists",
			errors.New("GenreService.CreateGenre -> duplicate genre name"),
		)
	}

	return service.repo.Create(ctx, genre)
}

func (service *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre, id uint) (*models.Genre, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid genre id",
			errors.New("GenreService.UpdateGenre -> id is zero"),
		)
	}

	trimGenreSpaces(genre)

	if err := validateGenre(genre); err != nil {
		return nil, err
	}

	existingGenre, err := service.repo.GetByFilter(ctx, &dto.GenreFilter{
		Name: &genre.Name,
	})

	if err != nil && !isSafeNotFound(err) {
		return nil, err
	}

	if existingGenre != nil && existingGenre.ID != id {
		return nil, utils.NewConflict(
			"Genre already exists",
			errors.New("GenreService.UpdateGenre -> duplicate genre name"),
		)
	}

	return service.repo.Update(ctx, genre, id)
}

func (service *GenreService) DeleteGenre(ctx context.Context, id uint) error {
	if id == 0 {
		return utils.NewInvalidInput(
			"Invalid genre id",
			errors.New("GenreService.DeleteGenre -> id is zero"),
		)
	}

	return service.repo.Delete(ctx, id)
}

func validateGenre(genre *models.Genre) error {
	if genre == nil {
		return utils.NewInvalidInput(
			"Invalid genre data",
			errors.New("validateGenre -> genre is nil"),
		)
	}

	if strings.TrimSpace(genre.Name) == "" {
		return utils.NewInvalidInput(
			"Genre name is required",
			errors.New("validateGenre -> name is empty"),
		)
	}

	return nil
}

func trimGenreSpaces(genre *models.Genre) {
	if genre == nil {
		return
	}

	genre.Name = strings.TrimSpace(genre.Name)
}

func isSafeNotFound(err error) bool {
	var safeErr *utils.SafeError

	if !errors.As(err, &safeErr) {
		return false
	}

	return safeErr.Code == "NOT_FOUND"
}
