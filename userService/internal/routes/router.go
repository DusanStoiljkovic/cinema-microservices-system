package routes

import (
	"net/http"
	"user-service/internal/handlers"
	"user-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRouter(UserHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", UserHandler.RegisterUser)
			r.Post("/login", UserHandler.LoginUser)

			r.Group(func(r chi.Router) {
				r.Use(middleware.JwtAuthMiddleware)
				r.Get("/{id}", UserHandler.GetUserByID)
			})
		})
	})

	return r
}
