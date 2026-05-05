package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"user-service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

var (
	clients = make(map[string]time.Time)
	mu      sync.Mutex
)

type ContextKey string

const (
	EmailKey  ContextKey = "email"
	UserIDKey ContextKey = "userID"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeSafeError(w http.ResponseWriter, r *http.Request, err error) {
	var safeErr *utils.SafeError

	if !errors.As(err, &safeErr) {
		safeErr = utils.NewInternal("Internal server error", err)
	}

	slog.Error(
		"middleware request failed",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("error", safeErr),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(safeErr.Status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    safeErr.Code,
		Message: safeErr.UserMsg,
	})
}

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

// JWT TOKEN LOGIC
func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("SECRET_KEY")
		if secret == "" {
			writeSafeError(w, r, utils.NewInternal(
				"Server configuration error",
				errors.New("JwtAuthMiddleware -> SECRET_KEY is empty"),
			))
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeSafeError(w, r, utils.NewAuthFailed(
				"Unauthorized",
				errors.New("JwtAuthMiddleware -> missing Authorization header"),
			))
			return
		}

		tokenString := strings.TrimSpace(authHeader)

		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			tokenString = strings.TrimSpace(tokenString[7:])
		}

		if tokenString == "" {
			writeSafeError(w, r, utils.NewAuthFailed(
				"Unauthorized",
				errors.New("JwtAuthMiddleware -> empty token"),
			))
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
			writeSafeError(w, r, utils.NewAuthFailed(
				"Invalid or expired token",
				err,
			))
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			writeSafeError(w, r, utils.NewAuthFailed(
				"Invalid token claims",
				errors.New("JwtAuthMiddleware -> email claim missing"),
			))
			return
		}

		userIDFloat, ok := claims[string(UserIDKey)].(float64)
		if !ok || userIDFloat <= 0 {
			writeSafeError(w, r, utils.NewAuthFailed(
				"Invalid token claims",
				errors.New("JwtAuthMiddleware -> userID claim missing"),
			))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, EmailKey, email)
		ctx = context.WithValue(ctx, UserIDKey, uint(userIDFloat))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateToken(ID uint, email string) (string, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return "", utils.NewInternal(
			"Server configuration error",
			errors.New("CreateToken -> SECRET_KEY is empty"),
		)
	}

	if ID == 0 {
		return "", utils.NewInvalidInput(
			"Invalid user id",
			errors.New("CreateToken -> user id is zero"),
		)
	}

	if strings.TrimSpace(email) == "" {
		return "", utils.NewInvalidInput(
			"Invalid email",
			errors.New("CreateToken -> email is empty"),
		)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			string(UserIDKey): ID,
			"email":           email,
			"exp":             time.Now().Add(time.Hour).Unix(),
		},
	)

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
			writeSafeError(w, r, utils.NewInternal(
				"Server configuration error",
				errors.New("AuthMiddleware -> SECRET_KEY is empty"),
			))
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		if apiKey != secret {
			writeSafeError(w, r, utils.NewAuthFailed(
				"Forbidden",
				errors.New("AuthMiddleware -> invalid api key"),
			))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// THROTTLE AUTH MIDDLEWARE
func ThrottleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		mu.Lock()
		lastRequest, found := clients[ip]
		if found && time.Since(lastRequest) < time.Second {
			mu.Unlock()

			writeSafeError(w, r, utils.NewConflict(
				"Too many requests",
				errors.New("ThrottleMiddleware -> rate limit exceeded"),
			))
			return
		}

		clients[ip] = time.Now()
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}
