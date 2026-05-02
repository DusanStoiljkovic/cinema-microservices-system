package handler

import (
	"context"
	"encoding/json"
	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/utils"
	"net/http"
)

type GenreService interface {
	GetGenres(ctx context.Context) ([]*models.Genre, error)
	GetGenreByFilter(ctx context.Context, req *dto.GenreFilter) (*models.Genre, error)
	CreateGenre(ctx context.Context, genre *models.Genre) (*models.Genre, error)
	UpdateGenre(ctx context.Context, genre *models.Genre, id uint) (*models.Genre, error)
	DeleteGenre(ctx context.Context, id uint) error
}

type GenreHandler struct {
	service GenreService
}

func NewGenreHandler(service GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

func (handler *GenreHandler) HandleGetGenres(w http.ResponseWriter, r *http.Request) error {
	genres, err := handler.service.GetGenres(r.Context())
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, genres)
	return nil
}

func (handler *GenreHandler) HandleGetGenreByFilter(w http.ResponseWriter, r *http.Request) error {
	var req *dto.GenreFilter

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	genre, err := handler.service.GetGenreByFilter(r.Context(), req)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, genre)
	return nil
}

func (handler *GenreHandler) HandleCreateGenre(w http.ResponseWriter, r *http.Request) error {
	var req *dto.GenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	createdGenre, err := handler.service.CreateGenre(r.Context(), &models.Genre{Name: req.Name})
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusCreated, createdGenre)
	return nil
}

func (handler *GenreHandler) HandleUpdateGenre(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	var req dto.GenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	createdGenre, err := handler.service.UpdateGenre(r.Context(), &models.Genre{Name: req.Name}, id)
	if err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusOK, createdGenre)
	return nil
}

func (handler *GenreHandler) HandleDeleteGenre(w http.ResponseWriter, r *http.Request) error {
	id, err := parseIDParam(r, "id")
	if err != nil {
		return err
	}

	if err := handler.service.DeleteGenre(r.Context(), id); err != nil {
		return err
	}

	utils.WriteJSON(w, http.StatusNoContent, "")
	return nil
}
