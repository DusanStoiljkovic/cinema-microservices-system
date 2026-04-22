package handlers

import (
	"net/http"
	"strconv"
	"user-service/internal/models"
	"user-service/internal/services"

	"github.com/go-chi/chi/v5"
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

	h.service.GetUserByFilter(r.Context(), &models.User{ID: uint(userID)})
}
