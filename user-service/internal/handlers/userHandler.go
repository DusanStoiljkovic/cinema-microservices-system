package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"user-service/internal/auth"
	"user-service/internal/dto"
	"user-service/internal/middleware"
	"user-service/internal/models"
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
	service  UserService
	validate *validator.Validate
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
	if !ok {
		return utils.NewAuthFailed(
			"Unauthorized",
			errors.New(fmt.Sprintf("UserHandler.GetMe -> user id missing from context -> userID: %d", userID)),
		)
	}

	user, err := h.service.GetUserByFilter(r.Context(), &models.User{ID: userID})
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	userID, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid user id", err)
	}

	user, err := h.service.GetUserByFilter(r.Context(), &models.User{ID: userID})
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) error {
	req := &dto.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.NewInvalidInput("Validation failed", err)
	}

	user, err := h.service.Register(r.Context(), &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	req := &dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.NewInvalidInput("Validation failed", err)
	}

	user, err := h.service.Login(r.Context(), &models.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	token, err := middleware.CreateToken(user)
	if err != nil {
		return utils.NewAuthFailed(
			"Failed to create authentication token",
			err,
		)
	}

	return utils.WriteJSON(w, http.StatusOK, map[string]string{
		"jwt": token,
	})
}

func parseIDParam(r *http.Request, param string) (uint, error) {
	value := chi.URLParam(r, param)

	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id must be greater than zero")
	}

	return uint(id), nil
}
