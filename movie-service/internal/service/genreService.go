package service

import (
	"context"
	"errors"
	"movie-service/internal/dto"
	"movie-service/internal/models"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrRecordAlreadyExist = errors.New("record already exist")
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
	genres, err := service.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return genres, nil
}

func (service *GenreService) GetGenreByFilter(ctx context.Context, filter *dto.GenreFilter) (*models.Genre, error) {
	genre, err := service.repo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return genre, nil
}

func (service *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error) {
	if err := validateGenre(genre); err != nil {
		return nil, err
	}

	name := strings.TrimSpace(genre.Name)

	existGenre, err := service.repo.GetByFilter(ctx, &dto.GenreFilter{Name: &name})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existGenre != nil {
		return nil, ErrRecordAlreadyExist
	}

	createdGenre, err := service.repo.Create(ctx, genre)
	if err != nil {
		return nil, err
	}

	return createdGenre, nil
}

func (service *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre, id uint) (*models.Genre, error) {
	if err := validateGenre(genre); err != nil {
		return nil, err
	}

	updatedGenre, err := service.repo.Update(ctx, genre, id)
	if err != nil {
		return nil, err
	}

	return updatedGenre, nil
}

func (service *GenreService) DeleteGenre(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrInvalidInput
	}

	err := service.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func validateGenre(genre *models.Genre) error {
	if genre == nil {
		return ErrInvalidInput
	}

	if strings.TrimSpace(genre.Name) == "" {
		return ErrInvalidInput
	}

	return nil
}
