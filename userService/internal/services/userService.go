package services

import (
	"context"
	"errors"
	"log"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/utils"

	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
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

func (s *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := s.repo.GetUserByFilter(ctx, &models.User{Email: user.Email})
	if err == nil {
		return nil, errors.New(ErrUserAlreadyExists.Error())
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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
	existedUser, err := s.repo.GetUserByFilter(ctx, &models.User{Email: user.Email})
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	err = utils.VerifyPassword(existedUser.Password, user.Password)
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	log.Print("User registered: %s", user.Email)
	return user, nil

}
