package middleware

import (
	"api-gateway/internal/auth"
	"api-gateway/internal/utils"
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

const (
	EmailKey  ContextKey = "email"
	UserIDKey ContextKey = "userID"
)

// API KEY MIDDLEWARE
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			utils.NewInternal(
				"Server configuration error",
				errors.New("AuthMiddleware -> SECRET_KEY is empty"))
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != secret {
			utils.NewAuthFailed(
				"Forbidden",
				errors.New("AuthMiddleware -> invalid api key"),
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type JwtClaims struct {
	UserID uint   `json:"userID"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			authErr := utils.NewInternal(
				"Server configuration error",
				errors.New("SECRET_KEY is empty"))

			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return
		}

		tokenString, err := extractBearerToken(r)
		if err != nil {
			authErr := utils.NewAuthFailed("Unauthorized", err)
			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})

			return
		}

		claims := &JwtClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		})
		if err != nil || token == nil || !token.Valid {
			authErr := utils.NewAuthFailed("Invalid or expired token", err)
			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return

		}

		if claims.UserID == 0 {
			authErr := utils.NewAuthFailed(
				"Invalid token claims",
				errors.New("Invalid token claims"))

			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserIDKey, claims.UserID)

		if claims.Role != "" {
			ctx = context.WithValue(ctx, auth.RoleKey, claims.Role)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))

	if authHeader == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 {
		return "", errors.New("Invalid Authorization header format")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid Authorization header scheme")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}
