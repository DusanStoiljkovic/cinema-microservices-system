package repository

import (
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"

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
		return nil, utils.ErrNotFound
	}

	return halls, nil
}

func (repo *HallRepository) GetByID(ctx context.Context, id uint) (*models.Hall, error) {
	hall := &models.Hall{}

	if err := repo.db.WithContext(ctx).First(hall, id).Error; err != nil {
		return nil, utils.ErrNotFound
	}

	return hall, nil
}

func (repo *HallRepository) Create(ctx context.Context, hall *models.Hall) (*models.Hall, error) {
	err := repo.db.WithContext(ctx).Create(hall).Error
	if err != nil {
		return nil, utils.ErrRecordAlreadyExist
	}

	return hall, nil
}

func (repo *HallRepository) Update(ctx context.Context, id uint, hall *models.Hall) (*models.Hall, error) {
	existHall := &models.Hall{}

	if err := repo.db.WithContext(ctx).First(existHall, id).Error; err != nil {
		return nil, utils.ErrNotFound
	}

	existHall.Name = hall.Name
	existHall.Location = hall.Location
	existHall.Capacity = hall.Capacity

	if err := repo.db.WithContext(ctx).Save(existHall).Error; err != nil {
		return nil, utils.ErrConflict
	}

	return existHall, nil
}

func (repo *HallRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&models.Hall{}, id).Error
}
