package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"user-service/internal/auth"
	"user-service/internal/models"
	"user-service/internal/utils"

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

func CreateToken(user *models.User) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "", utils.NewInternal(
			"Server configuration error",
			errors.New("auth.CreateToken -> SECRET_KEY is empty"),
		)
	}

	if user == nil {
		return "", utils.NewInvalidInput(
			"Invalid user",
			errors.New("auth.CreateToken -> user is nil"),
		)
	}

	if user.ID == 0 {
		return "", utils.NewInvalidInput(
			"Invalid user id",
			errors.New("auth.CreateToken() -> user id is zero"),
		)
	}

	role := strings.TrimSpace(user.Role)
	if role == "" {
		role = "user"
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"sub": strconv.FormatUint(uint64(user.ID), 10),
		"iss": "cinema-user-service",
		"aud": "cinema-api",
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(1 * time.Hour).Unix(),

		// Custom claims
		"userID": user.ID,
		"role":   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", utils.NewInternal(
			"Failed to create token",
			err,
		)
	}

	return tokenString, nil
}

// API KEY MIDDLEWARE
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			authErr := utils.NewInternal("Server configuration error", errors.New("SECRET_KEY is empty"))
			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != secret {
			authErr := utils.NewInternal(
				"Server configuration error",
				errors.New("AuthMiddleware -> invalid api key"))

			utils.WriteJSON(w, authErr.Status, map[string]string{
				"error": authErr.UserMsg,
				"code":  authErr.Code,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
