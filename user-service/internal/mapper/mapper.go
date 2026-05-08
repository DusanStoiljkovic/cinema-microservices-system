package mapper

import (
	"user-service/internal/dto"
	"user-service/internal/models"
)

func UserToResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func UsersToResponse(users []*models.User) []*dto.UserResponse {
	var usersResponse []*dto.UserResponse

	for _, val := range users {
		usersResponse = append(usersResponse, UserToResponse(val))
	}

	return usersResponse
}

func UserFromRegisterRequest(req *dto.RegisterRequest) *models.User {
	return &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func UserFromLoginRequest(req *dto.LoginRequest) *models.User {
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
}
