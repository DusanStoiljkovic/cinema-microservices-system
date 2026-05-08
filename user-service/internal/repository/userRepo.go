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

func (repo *UserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	users := &[]*models.User{}

	if err := repo.db.WithContext(ctx).Find(users).Error; err != nil {
		return nil, utils.NewInternal("Failed to load users", err)
	}

	if len(*users) == 0 {
		return nil, utils.NewNotFound(
			"Users not found",
			errors.New("UserRepository.GetAll -> empty list"),
		)
	}

	return *users, nil
}

func (r *UserRepository) GetByFilter(ctx context.Context, filter *dto.UserFilter) (*models.User, error) {
	if filter == nil {
		return nil, utils.NewInvalidInput(
			"Invalid user filter",
			errors.New("UserRepository.GetUserByFilter -> filter is nil"),
		)
	}

	if filter.ID == 0 &&
		(filter.Email == "" || filter.Email == "") &&
		(filter.Name == "" || filter.Name == "") {
		return nil, utils.NewInvalidInput(
			"At least one user filter is required",
			errors.New("UserRepository.GetUserByFilter -> empty filter"),
		)
	}

	user := &models.User{}

	query := r.db.WithContext(ctx).Model(&models.User{})

	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Email != "" && filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	if filter.Name != "" && filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if err := query.First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewNotFound("User not found", err)
		}

		return nil, utils.NewConflict("Failed to load user", err)
	}

	return user, nil
}

func (r *UserRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, utils.NewConflict("Failed to create user", err)
	}

	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser, err := repo.GetByFilter(ctx, &dto.UserFilter{ID: user.ID})
	if err != nil {
		return nil, err
	}

	updatedUser.Name = user.Name
	updatedUser.Email = user.Email

	if err := repo.db.WithContext(ctx).Save(updatedUser).Error; err != nil {
		return nil, utils.NewInternal("Failed to update user", err)
	}

	return updatedUser, nil
}

func (repo *UserRepository) Delete(ctx context.Context, id uint) error {
	deleteUser, err := repo.GetByFilter(ctx, &dto.UserFilter{ID: id})
	if err != nil {
		return err
	}

	if err := repo.db.WithContext(ctx).Delete(deleteUser, id).Error; err != nil {
		return utils.NewInternal("Failed to delete user", err)
	}

	return nil
}
