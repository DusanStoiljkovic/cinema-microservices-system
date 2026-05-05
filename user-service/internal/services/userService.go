package services

import (
	"context"
	"errors"
	"strings"

	"user-service/internal/dto"
	"user-service/internal/models"
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
	if req == nil {
		return nil, utils.NewInvalidInput(
			"Invalid user filter",
			errors.New("UserService.GetUserByFilter -> request is nil"),
		)
	}

	filter := &dto.UserFilter{}

	if req.ID != 0 {
		filter.ID = &req.ID
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name != "" {
		filter.Name = &req.Name
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email != "" {
		filter.Email = &req.Email
	}

	if filter.ID == nil && filter.Name == nil && filter.Email == nil {
		return nil, utils.NewInvalidInput(
			"At least one user filter is required",
			errors.New("UserService.GetUserByFilter -> empty filter"),
		)
	}

	return s.repo.GetUserByFilter(ctx, filter)
}

func (s *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	trimUserSpaces(user)

	if err := validateRegisterUser(user); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetUserByFilter(ctx, &dto.UserFilter{
		Email: &user.Email,
	})

	if err != nil && !isSafeNotFound(err) {
		return nil, err
	}

	if existingUser != nil {
		return nil, utils.NewConflict(
			"User already exists",
			errors.New("UserService.Register -> duplicate email"),
		)
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return nil, utils.NewInternal(
			"Failed to process password",
			err,
		)
	}

	user.Password = hashedPassword
	user.Role = "user"

	return s.repo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, user *models.User) (*models.User, error) {
	trimUserSpaces(user)

	if err := validateLoginUser(user); err != nil {
		return nil, err
	}

	existingUser, err := s.repo.GetUserByFilter(ctx, &dto.UserFilter{
		Email: &user.Email,
	})
	if err != nil {
		return nil, utils.NewAuthFailed(
			"Invalid credentials",
			err,
		)
	}

	if err := utils.VerifyPassword(existingUser.Password, user.Password); err != nil {
		return nil, utils.NewAuthFailed(
			"Invalid credentials",
			err,
		)
	}

	return existingUser, nil
}

func validateRegisterUser(user *models.User) error {
	if user == nil {
		return utils.NewInvalidInput(
			"Invalid user data",
			errors.New("validateRegisterUser -> user is nil"),
		)
	}

	if user.Name == "" {
		return utils.NewInvalidInput(
			"Name is required",
			errors.New("validateRegisterUser -> name is empty"),
		)
	}

	if user.Email == "" {
		return utils.NewInvalidInput(
			"Email is required",
			errors.New("validateRegisterUser -> email is empty"),
		)
	}

	if user.Password == "" {
		return utils.NewInvalidInput(
			"Password is required",
			errors.New("validateRegisterUser -> password is empty"),
		)
	}

	return nil
}

func validateLoginUser(user *models.User) error {
	if user == nil {
		return utils.NewInvalidInput(
			"Invalid login data",
			errors.New("validateLoginUser -> user is nil"),
		)
	}

	if user.Email == "" {
		return utils.NewInvalidInput(
			"Email is required",
			errors.New("validateLoginUser -> email is empty"),
		)
	}

	if user.Password == "" {
		return utils.NewInvalidInput(
			"Password is required",
			errors.New("validateLoginUser -> password is empty"),
		)
	}

	return nil
}

func trimUserSpaces(user *models.User) {
	if user == nil {
		return
	}

	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
}

func isSafeNotFound(err error) bool {
	var safeErr *utils.SafeError

	if !errors.As(err, &safeErr) {
		return false
	}

	return safeErr.Code == "NOT_FOUND"
}
