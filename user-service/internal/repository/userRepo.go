package repository

import (
	"context"
	"errors"
	"user-service/internal/dto"
	"user-service/internal/models"

	"gorm.io/gorm"
)

var (
	ErrInputFilter = errors.New("invalid user filter")
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

	if filter == nil {
		return nil, ErrInputFilter
	}

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
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
