package routes

import (
	"movie-service/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type MovieHandler interface {
	HandleGetMovies(w http.ResponseWriter, r *http.Request) error
	HandleGetMovieByID(w http.ResponseWriter, r *http.Request) error
	HandleGetRelationsByMovieID(w http.ResponseWriter, r *http.Request) error
	HandleCreateMovie(w http.ResponseWriter, r *http.Request) error
	HandleCreateRelation(w http.ResponseWriter, r *http.Request) error
	HandleUpdateMovie(w http.ResponseWriter, r *http.Request) error
	HandleDeleteMovie(w http.ResponseWriter, r *http.Request) error
}

type GenreHandler interface {
	HandleGetGenres(w http.ResponseWriter, r *http.Request) error
	HandleGetGenreByFilter(w http.ResponseWriter, r *http.Request) error
	HandleCreateGenre(w http.ResponseWriter, r *http.Request) error
	HandleUpdateGenre(w http.ResponseWriter, r *http.Request) error
	HandleDeleteGenre(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(movieHandler MovieHandler, genreHandler GenreHandler) http.Handler {
	r := chi.NewRouter()

	r.Route("/movies", func(r chi.Router) {
		r.Get("/", middleware.ErrorHandler(movieHandler.HandleGetMovies))
		r.Get("/{id}/genres", middleware.ErrorHandler(movieHandler.HandleGetRelationsByMovieID))
		r.Get("/{id}", middleware.ErrorHandler(movieHandler.HandleGetMovieByID))
		r.Post("/", middleware.ErrorHandler(movieHandler.HandleCreateMovie))
		r.Post("/{movieId}/genres/{genreId}", middleware.ErrorHandler(movieHandler.HandleCreateRelation))
		r.Put("/{id}", middleware.ErrorHandler(movieHandler.HandleUpdateMovie))
		r.Delete("/{id}", middleware.ErrorHandler(movieHandler.HandleDeleteMovie))

		r.Group(func(r chi.Router) {
			r.Route("/genres", func(r chi.Router) {
				r.Get("/", middleware.ErrorHandler(genreHandler.HandleGetGenres))
				r.Get("/{id}", middleware.ErrorHandler(genreHandler.HandleGetGenreByFilter))
				r.Post("/", middleware.ErrorHandler(genreHandler.HandleCreateGenre))
				r.Put("/{id}", middleware.ErrorHandler(genreHandler.HandleUpdateGenre))
				r.Delete("/{id}", middleware.ErrorHandler(genreHandler.HandleDeleteGenre))
			})
		})
	})

	return r
}
