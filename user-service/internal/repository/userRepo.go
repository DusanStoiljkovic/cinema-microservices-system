package repository

import (
	"context"
	"user-service/internal/dto"
	"user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByFilter(ctx context.Context, filter *dto.UserFilter) (*models.User, error) {
	user := &models.User{}

	query := r.db.WithContext(ctx).Model(&models.User{})

	if filter.ID != nil {
		query = query.Where("id = ?", *filter.ID)
	}

	if filter.Email != nil && *filter.Email != "" {
		query = query.Where("email = ?", *filter.Email)
	}

	if filter.Name != nil && *filter.Name != "" {
		query = query.Where("name = ?", *filter.Name)
	}

	err := query.First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
