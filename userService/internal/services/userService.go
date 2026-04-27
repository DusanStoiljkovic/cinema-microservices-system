package services

import (
	"context"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/secure"
	"user-service/internal/utils"
)

type UserRepository interface {
	GetUserByFilter(ctx context.Context, filter *dto.UserFilter) (*models.User, error)
	Create(ctx context.Context, user *models.User) (*models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByFilter(ctx context.Context, req *models.User) (*models.User, error) {
	user, err := s.repo.GetUserByFilter(ctx, &dto.UserFilter{Name: &req.Name, Email: &req.Email})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existedUser, err := s.repo.GetUserByFilter(ctx, &dto.UserFilter{Email: &user.Email})
	if err == nil && existedUser != nil {
		return nil, secure.NewAuthFailed("User already exist", err, nil)
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	user.Role = "user"

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	existedUser, err := s.repo.GetUserByFilter(ctx, &dto.UserFilter{Email: &user.Email})
	if err != nil {
		return nil, secure.NewAuthFailed("Invalid credentials", err, nil)
	}

	err = utils.VerifyPassword(existedUser.Password, user.Password)
	if err != nil {
		return nil, secure.NewAuthFailed("Invalid credentials", err, nil)
	}

	return existedUser, nil

}
