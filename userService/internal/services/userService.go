package services

import (
	"context"
	"errors"
	"fmt"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrEmailInUse = errors.New("email already in use")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByFilter(ctx context.Context, req *models.User) (*models.User, error) {
	user, err := s.repo.GetUserByFilter(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Register(ctx context.Context, req *models.User) (*models.User, error) {
	_, err := s.repo.GetUserByFilter(ctx, &models.User{Email: req.Email})
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword
	req.Role = "user"

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
