package middleware

import (
	"errors"
	"log"
	"net/http"
	"user-service/internal/auth"
	"user-service/internal/utils"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(auth.RoleKey).(string)
		if !ok || role == "" {
			authErr := utils.NewAuthFailed("Unauthorized", errors.New("role missing from context"))

			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return
		}

		if role != "admin" {
			log.Print("ROLE: ", role)
			authErr := utils.NewAuthFailed("Forbidden", errors.New("user is not admin"))
			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": "Forbidden",
				"code":  "FORBIDDEN",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
