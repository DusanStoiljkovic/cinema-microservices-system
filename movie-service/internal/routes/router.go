package routes

import (
	"movie-service/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MovieHandler interface {
	HandleGetMovies(w http.ResponseWriter, r *http.Request) error
	HandleGetMovieByID(w http.ResponseWriter, r *http.Request) error
	HandleCreateMovie(w http.ResponseWriter, r *http.Request) error
	HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error
	HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(handler MovieHandler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/movies", func(r chi.Router) {
			r.Get("/", middleware.ErrorHandler(handler.HandleGetMovies))
			r.Get("/{id}", middleware.ErrorHandler(handler.HandleGetMovieByID))
			r.Post("/", middleware.ErrorHandler(handler.HandleCreateMovie))
			r.Put("/{id}", middleware.ErrorHandler(handler.HandleUpdateMovie))
			r.Delete("/{id}", middleware.ErrorHandler(handler.HandleDeleteMovie))
		})
	})
	return r
}
