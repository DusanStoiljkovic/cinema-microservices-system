package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"user-service/internal/dto"
	"user-service/internal/middleware"
	"user-service/internal/models"
	"user-service/internal/secure"
	"user-service/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	GetUserByFilter(ctx context.Context, req *models.User) (*models.User, error)
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		return secure.NewAuthFailed("Invalid user ID", err, nil)
	}

	user, err := h.service.GetUserByFilter(r.Context(), &models.User{ID: uint(userID)})
	if err != nil {
		return secure.NewAuthFailed("User does not exist", err, nil)
	}

	json.NewEncoder(w).Encode(&dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})

	return nil
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	var request dto.RegisterRequest
	var validate = validator.New()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return secure.NewAuthFailed("Invalid request body", err, nil)
	}

	err = validate.Struct(request)
	if err != nil {
		return secure.NewAuthFailed("Validation failed", err, nil)
	}

	user, err := h.service.Register(r.Context(), &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, user)
	return nil
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	request := &dto.LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return secure.NewAuthFailed(
			"Invalid request body",
			err,
			nil,
		)
	}

	_, err = h.service.Login(r.Context(), &models.User{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return err
	}

	token, err := middleware.CreateToken(request.Email)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"JWT": token,
	})

	return nil

}
