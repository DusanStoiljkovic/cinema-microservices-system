package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"

	"gorm.io/gorm"
)

type HallRepository struct {
	db *gorm.DB
}

func NewHallRepository(db *gorm.DB) *HallRepository {
	return &HallRepository{db: db}
}

func (repo *HallRepository) GetAll(ctx context.Context) ([]*models.Hall, error) {
	var halls []*models.Hall

	if err := repo.db.WithContext(ctx).Find(&halls).Error; err != nil {
		return nil, utils.NewConflict("Failed to load halls", err)
	}

	if len(halls) == 0 {
		return nil, utils.NewNotFound(
			"Halls not found",
			errors.New("HallRepository.GetAll -> empty list"),
		)
	}

	return halls, nil
}

func (repo *HallRepository) GetByID(ctx context.Context, id uint) (*models.Hall, error) {
	hall := &models.Hall{}

	if err := repo.db.WithContext(ctx).First(hall, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Hall not found", err)
		}

		return nil, utils.NewConflict("Failed to load hall", err)
	}

	return hall, nil
}

func (repo *HallRepository) Create(ctx context.Context, hall *models.Hall) (*models.Hall, error) {
	exists, err := repo.existsSameHall(ctx, 0, hall)
	if err != nil {
		return nil, utils.NewConflict(
			"Failed to check existing hall",
			err,
		)
	}

	if exists {
		return nil, utils.NewConflict(
			"Hall already exists",
			errors.New("HallRepository.Create -> duplicate hall"),
		)
	}

	if err := repo.db.WithContext(ctx).Create(hall).Error; err != nil {
		return nil, utils.NewInvalidInput("Failed to create hall", err)
	}

	return hall, nil
}

func (repo *HallRepository) Update(ctx context.Context, id uint, hall *models.Hall) (*models.Hall, error) {
	existingHall := &models.Hall{}

	if err := repo.db.WithContext(ctx).First(existingHall, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("Hall not found", err)
		}

		return nil, utils.NewConflict("Failed to load hall", err)
	}

	exists, err := repo.existsSameHall(ctx, id, hall)
	if err != nil {
		return nil, utils.NewConflict(
			"Failed to check existing hall",
			err,
		)
	}

	if exists {
		return nil, utils.NewConflict(
			"Hall already exists",
			errors.New("HallRepository.Update -> duplicate hall"),
		)
	}

	existingHall.Name = hall.Name
	existingHall.Location = hall.Location
	existingHall.Capacity = hall.Capacity

	if err := repo.db.WithContext(ctx).Save(existingHall).Error; err != nil {
		return nil, utils.NewConflict("Failed to update hall", err)
	}

	return existingHall, nil
}

func (repo *HallRepository) Delete(ctx context.Context, id uint) error {
	existingHall := &models.Hall{}

	if err := repo.db.WithContext(ctx).First(existingHall, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewNotFound("Hall not found", err)
		}

		return utils.NewConflict("Failed to load hall", err)
	}

	if err := repo.db.WithContext(ctx).Delete(existingHall).Error; err != nil {
		return utils.NewConflict("Failed to delete hall", err)
	}

	return nil
}

func (repo *HallRepository) existsSameHall(ctx context.Context, excludeID uint, hall *models.Hall) (bool, error) {
	existingHall := &models.Hall{}

	query := repo.db.WithContext(ctx).
		Where("name = ?",
			hall.Name,
		)

	if excludeID != 0 {
		query = query.Where("id <> ?", excludeID)
	}

	err := query.First(existingHall).Error

	if err == nil {
		return true, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}
