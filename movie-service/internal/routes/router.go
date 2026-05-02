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

<<<<<<< HEAD
func RegisterRouter(handler MovieHandler) http.Handler {
=======
type GenreHandler interface {
	HandleGetGenres(w http.ResponseWriter, r *http.Request) error
	HandleGetGenreByFilter(w http.ResponseWriter, r *http.Request) error
	HandleCreateGenre(w http.ResponseWriter, r *http.Request) error
	HandleUpdateGenre(w http.ResponseWriter, r *http.Request) error
	HandleDeleteGenre(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(movieHandler MovieHandler, genreHandler GenreHandler) http.Handler {
>>>>>>> feature/movieService
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/movies", func(r chi.Router) {
<<<<<<< HEAD
			r.Get("/", middleware.ErrorHandler(handler.HandleGetMovies))
			r.Get("/{id}", middleware.ErrorHandler(handler.HandleGetMovieByID))
			r.Post("/", middleware.ErrorHandler(handler.HandleCreateMovie))
			r.Put("/{id}", middleware.ErrorHandler(handler.HandleUpdateMovie))
			r.Delete("/{id}", middleware.ErrorHandler(handler.HandleDeleteMovie))
		})
	})
=======
			r.Get("/", middleware.ErrorHandler(movieHandler.HandleGetMovies))
			r.Get("/{id}", middleware.ErrorHandler(movieHandler.HandleGetMovieByID))
			r.Post("/", middleware.ErrorHandler(movieHandler.HandleCreateMovie))
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
	})

>>>>>>> feature/movieService
	return r
}
