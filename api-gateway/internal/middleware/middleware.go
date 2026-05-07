package middleware

import (
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

// JWT TOKEN LOGIC
func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			utils.NewInternal(
				"Server configuration error",
				errors.New("JwtAuthMiddleware -> SECRET_KEY is empty"),
			)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.NewAuthFailed(
				"Unauthorized",
				errors.New("JwtAuthMiddleware -> missing Authorization header"),
			)
			return
		}

		tokenString := strings.TrimSpace(authHeader)

		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = strings.TrimSpace(tokenString[7:])
		}

		if tokenString == "" {
			utils.NewAuthFailed(
				"Unauthorized",
				errors.New("JwtAuthMiddleware -> empty token"),
			)
			return
		}

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		})

		if err != nil || token == nil || !token.Valid {
			utils.NewAuthFailed(
				"Invalid or expired token",
				err,
			)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			utils.NewAuthFailed(
				"Invalid token claims",
				errors.New("JwtAuthMiddleware -> email claim missing"),
			)
			return
		}

		userID, ok := claims[string(UserIDKey)].(uint)
		if !ok || userID <= 0 {
			utils.NewAuthFailed(
				"Invalid token claims",
				errors.New("JwtAuthMiddleware -> userID claim missing"),
			)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, EmailKey, email)
		ctx = context.WithValue(ctx, UserIDKey, uint(userID))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
