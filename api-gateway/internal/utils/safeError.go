package utils

import (
	"fmt"
	"log/slog"
	"net/http"
)

type SafeError struct {
	Code     string
	UserMsg  string
	Status   int
	Internal error
	Metadata map[string]any
}

func (e *SafeError) Error() string {
	return e.UserMsg
}

func (e *SafeError) Unwrap() error {
	return e.Internal
}

func (e *SafeError) LogString() string {
	return fmt.Sprintf("Code: %s | Status: %d | Msg: %s | Cause: %v | Meta: %v",
		e.Code, e.Status, e.UserMsg, e.Internal, e.Metadata)
}

func (e *SafeError) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("code", e.Code),
		slog.Int("status", e.Status),
		slog.String("user_msg", e.UserMsg),
	}

	if e.Internal != nil {
		attrs = append(attrs, slog.String("internal", e.Internal.Error()))
	}

	if e.Metadata != nil {
		attrs = append(attrs, slog.Any("metadata", e.Metadata))
	}

	return slog.GroupValue(attrs...)
}

func NewSafeError(code string, userMsg string, status int, internal error, metadata map[string]any) *SafeError {
	return &SafeError{
		Code:     code,
		UserMsg:  userMsg,
		Status:   status,
		Internal: internal,
		Metadata: metadata,
	}
}

func NewInvalidInput(public string, internal error) *SafeError {
	return NewSafeError(
		"INVALID_INPUT",
		public,
		http.StatusBadRequest,
		internal,
		nil,
	)
}

func NewNotFound(public string, internal error) *SafeError {
	return NewSafeError(
		"NOT_FOUND",
		public,
		http.StatusNotFound,
		internal,
		nil,
	)
}

func NewConflict(public string, internal error) *SafeError {
	return NewSafeError(
		"CONFLICT",
		public,
		http.StatusConflict,
		internal,
		nil,
	)
}

func NewAuthFailed(public string, internal error) *SafeError {
	return NewSafeError(
		"AUTH_FAILED",
		public,
		http.StatusUnauthorized,
		internal,
		nil,
	)
}

func NewInternal(public string, internal error) *SafeError {
	return NewSafeError(
		"INTERNAL_ERROR",
		public,
		http.StatusInternalServerError,
		internal,
		nil,
	)
}
