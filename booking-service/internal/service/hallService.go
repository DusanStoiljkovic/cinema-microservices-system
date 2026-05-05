package service

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"
	"strings"
)

type HallRepository interface {
	GetAll(ctx context.Context) ([]*models.Hall, error)
	GetByID(ctx context.Context, id uint) (*models.Hall, error)
	Create(ctx context.Context, hall *models.Hall) (*models.Hall, error)
	Update(ctx context.Context, id uint, hall *models.Hall) (*models.Hall, error)
	Delete(ctx context.Context, id uint) error
}

type HallService struct {
	repo HallRepository
}

func NewHallService(repo HallRepository) *HallService {
	return &HallService{repo: repo}
}

func (service *HallService) GetAllHalls(ctx context.Context) ([]*models.Hall, error) {
	return service.repo.GetAll(ctx)
}

func (service *HallService) GetHallByID(ctx context.Context, id uint) (*models.Hall, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid hall id",
			errors.New("HallService.GetHallByID -> id is zero"),
		)
	}

	return service.repo.GetByID(ctx, id)
}

func (service *HallService) CreateHall(ctx context.Context, hall *models.Hall) (*models.Hall, error) {
	trimHallSpaces(hall)

	if err := validateHall(hall); err != nil {
		return nil, err
	}

	return service.repo.Create(ctx, hall)
}

func (service *HallService) UpdateHall(ctx context.Context, id uint, hall *models.Hall) (*models.Hall, error) {
	if id == 0 {
		return nil, utils.NewInvalidInput(
			"Invalid hall id",
			errors.New("HallService.UpdateHall -> id is zero"),
		)
	}

	trimHallSpaces(hall)

	if err := validateHall(hall); err != nil {
		return nil, err
	}

	return service.repo.Update(ctx, id, hall)
}

func (service *HallService) DeleteHall(ctx context.Context, id uint) error {
	if id == 0 {
		return utils.NewInvalidInput(
			"Invalid hall id",
			errors.New("HallService.DeleteHall -> id is zero"),
		)
	}

	return service.repo.Delete(ctx, id)
}

func trimHallSpaces(hall *models.Hall) {
	if hall == nil {
		return
	}

	hall.Name = strings.TrimSpace(hall.Name)
	hall.Location = strings.TrimSpace(hall.Location)
}

func validateHall(hall *models.Hall) error {
	if hall == nil {
		return utils.NewInvalidInput(
			"Invalid hall data",
			errors.New("validateHall -> hall is nil"),
		)
	}

	if hall.Name == "" {
		return utils.NewInvalidInput(
			"Hall name is required",
			errors.New("validateHall -> name is empty"),
		)
	}

	if hall.Location == "" {
		return utils.NewInvalidInput(
			"Hall location is required",
			errors.New("validateHall -> location is empty"),
		)
	}

	if hall.Capacity <= 0 {
		return utils.NewInvalidInput(
			"Hall capacity must be greater than zero",
			errors.New("validateHall -> capacity must be greater than zero"),
		)
	}

	return nil
}
