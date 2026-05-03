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
<<<<<<< HEAD
func RegisterRouter(handler MovieHandler) http.Handler {
=======
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
type GenreHandler interface {
	HandleGetGenres(w http.ResponseWriter, r *http.Request) error
	HandleGetGenreByFilter(w http.ResponseWriter, r *http.Request) error
	HandleCreateGenre(w http.ResponseWriter, r *http.Request) error
	HandleUpdateGenre(w http.ResponseWriter, r *http.Request) error
	HandleDeleteGenre(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(movieHandler MovieHandler, genreHandler GenreHandler) http.Handler {
<<<<<<< HEAD
>>>>>>> feature/movieService
=======
>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/movies", func(r chi.Router) {
<<<<<<< HEAD
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

>>>>>>> da5f31b (feat(movie-service): implement genre management with repository, service, and handler layers; enhance movie handler and routes)
	return r
}
