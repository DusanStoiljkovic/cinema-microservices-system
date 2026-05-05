package repository

import (
	"context"
	"errors"

	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/utils"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByFilter(ctx context.Context, filter *dto.UserFilter) (*models.User, error) {
	if filter == nil {
		return nil, utils.NewInvalidInput(
			"Invalid user filter",
			errors.New("UserRepository.GetUserByFilter -> filter is nil"),
		)
	}

	if filter.ID == nil &&
		(filter.Email == nil || *filter.Email == "") &&
		(filter.Name == nil || *filter.Name == "") {
		return nil, utils.NewInvalidInput(
			"At least one user filter is required",
			errors.New("UserRepository.GetUserByFilter -> empty filter"),
		)
	}

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

	if err := query.First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("User not found", err)
		}

		return nil, utils.NewConflict("Failed to load user", err)
	}

	return user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, utils.NewConflict("Failed to create user", err)
	}

	return user, nil
}
