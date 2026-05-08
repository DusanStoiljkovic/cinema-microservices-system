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
	"user-service/internal/mapper"
	"user-service/internal/models"
	"user-service/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUserByFilter(ctx context.Context, filter *dto.UserFilter) (*models.User, error)
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, user *models.User) (map[string]string, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uint) error
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

func (handler *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := handler.service.GetAllUsers(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.UsersToResponse(users))
}

func (h *UserHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(auth.UserIDKey).(uint)
	if !ok {
		return utils.NewAuthFailed(
			"Unauthorized",
			errors.New(fmt.Sprintf("UserHandler.GetMe -> user id missing from context -> userID: %d", userID)),
		)
	}

	user, err := h.service.GetUserByFilter(r.Context(), &dto.UserFilter{ID: userID})
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.UserToResponse(user))
}

func (h *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	userID, err := parseIDParam(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid user id", err)
	}

	user, err := h.service.GetUserByFilter(r.Context(), &dto.UserFilter{ID: userID})
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.UserToResponse(user))
}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) error {
	req := &dto.RegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	user, err := h.service.RegisterUser(r.Context(), mapper.UserFromRegisterRequest(req))
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, mapper.UserToResponse(user))
}

func (h *UserHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) error {
	req := &dto.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	token, err := h.service.LoginUser(r.Context(), mapper.UserFromLoginRequest(req))
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, token)
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

func (handler *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	req := &models.User{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid body input", err)
	}

	updatedUser, err := handler.service.UpdateUser(r.Context(), req)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, mapper.UserToResponse(updatedUser))
	return nil
}

func (handler *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	userID, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	if err := handler.service.DeleteUser(r.Context(), userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
