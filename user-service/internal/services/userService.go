package services

import (
	"context"
	"database/sql"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"
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
	user, err := s.GetUserByFilter(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Register(ctx context.Context, req *models.User) (*models.User, error) {
	_, err := s.repo.GetUserByFilter(ctx, &models.User{Email: req.Email})
	if err != nil {
		return nil, ErrEmailInUse
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashedPassword

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
