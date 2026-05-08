package routes

import (
	"net/http"
	"user-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	HandleGetAllUsers(w http.ResponseWriter, r *http.Request) error
	HandleGetUserByID(w http.ResponseWriter, r *http.Request) error
	HandleGetMe(w http.ResponseWriter, r *http.Request) error
	HandleRegisterUser(w http.ResponseWriter, r *http.Request) error
	HandleLoginUser(w http.ResponseWriter, r *http.Request) error
	HandleUpdateUser(w http.ResponseWriter, r *http.Request) error
	HandleDeleteUser(w http.ResponseWriter, r *http.Request) error
}

func RegisterRouter(userHandler UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", middleware.ErrorHandler(userHandler.HandleRegisterUser))
			r.Post("/login", middleware.ErrorHandler(userHandler.HandleLoginUser))

			r.Group(func(r chi.Router) {
				r.Use(middleware.JwtAuthMiddleware)
				r.Get("/me", middleware.ErrorHandler(userHandler.HandleGetMe))
				r.Get("/{id}", middleware.ErrorHandler(userHandler.HandleGetUserByID))
				r.Put("/", middleware.ErrorHandler(userHandler.HandleUpdateUser))

				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireAdmin)
					r.Get("/", middleware.ErrorHandler(userHandler.HandleGetAllUsers))
					r.Delete("/{id}", middleware.ErrorHandler(userHandler.HandleDeleteUser))
				})
			})

		})
	})

	return r
}
