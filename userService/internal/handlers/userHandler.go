package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/services"
	"user-service/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	userID, _ := strconv.Atoi(idParam)

	user, err := h.service.GetUserByFilter(r.Context(), &models.User{ID: uint(userID)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&dto.UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})

}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	var validate = validator.New()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = validate.Struct(request)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.GetUserByFilter(r.Context(), &models.User{
		Email: request.Email,
	})
	if err == nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	request := &dto.LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	existedUser, err := h.service.GetUserByFilter(r.Context(), &models.User{Email: request.Email})
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	err = utils.VerifyPassword(existedUser.Password, request.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "JWT TOKEN 412nemkodsanaook3n1"}`))
}
