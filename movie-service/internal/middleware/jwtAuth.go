package middleware

import (
	"context"
	"errors"
	"log/slog"
	"movie-service/internal/auth"
	"movie-service/internal/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LOGGING MIDDLEWARE
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info(
			"request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
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
			utils.NewInternal("Server configuration error", errors.New("SECRET_KEY is empty"))
			return
		}

		tokenString, err := extractBearerToken(r)
		if err != nil {
			utils.NewAuthFailed("Unauthorized", err)
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
			utils.NewAuthFailed("Invalid or expired token", err)
			return
		}

		if claims.UserID == 0 {
			utils.NewInternal("Invalid token claims", errors.New("Invalid token claims"))
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

	if strings.HasPrefix(strings.ToLower(authHeader), "bearer") {
		authHeader = strings.TrimSpace(authHeader[7:])
	}

	if authHeader == "" {
		return "", errors.New("empty token")
	}

	return authHeader, nil
}
