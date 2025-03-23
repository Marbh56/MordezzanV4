package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"mordezzanV4/internal/logger"
	"net/http"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrBadRequest    = errors.New("invalid request")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal server error")
	ErrConflict      = errors.New("resource conflict")
	ErrValidation    = errors.New("validation error")
	ErrDatabaseError = errors.New("database error")
)

type AppError struct {
	Err     error
	Message string
	Code    int
	Field   string
}

func (e *AppError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFound(resource string, id interface{}) *AppError {
	return &AppError{
		Err:     ErrNotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
		Code:    http.StatusNotFound,
	}
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		Err:     ErrBadRequest,
		Message: msg,
		Code:    http.StatusBadRequest,
	}
}

func NewValidationError(field, msg string) *AppError {
	return &AppError{
		Err:     ErrValidation,
		Message: msg,
		Code:    http.StatusBadRequest,
		Field:   field,
	}
}

func NewDatabaseError(err error) *AppError {
	return &AppError{
		Err:     ErrDatabaseError,
		Message: "Database operation failed",
		Code:    http.StatusInternalServerError,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Err:     ErrInternal,
		Message: "An internal error occurred",
		Code:    http.StatusInternalServerError,
	}
}

func IsNotFound(err error) bool {
	var appErr *AppError
	return (errors.As(err, &appErr) && errors.Is(appErr.Err, ErrNotFound))
}

func IsValidation(err error) bool {
	var appErr *AppError
	return (errors.As(err, &appErr) && errors.Is(appErr.Err, ErrValidation))
}

func IsDatabase(err error) bool {
	var appErr *AppError
	return (errors.As(err, &appErr) && errors.Is(appErr.Err, ErrDatabaseError))
}

type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func HandleError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = &AppError{
			Err:     ErrInternal,
			Message: "An unexpected error occurred",
			Code:    http.StatusInternalServerError,
		}
	}

	resp := ErrorResponse{
		Status:  appErr.Code,
		Message: appErr.Message,
	}

	if appErr.Field != "" {
		resp.Fields = map[string]string{
			appErr.Field: appErr.Message,
		}
	}

	// Use structured logging with Zap
	if appErr.Code >= 500 {
		logger.With(
			"status_code", appErr.Code,
			"error_type", appErr.Err.Error(),
			"message", appErr.Message,
		).Error("Server error occurred")
	} else {
		logger.With(
			"status_code", appErr.Code,
			"error_type", appErr.Err.Error(),
			"message", appErr.Message,
		).Info("Client error occurred")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(appErr.Message))
	}
}

func HandleValidationErrors(w http.ResponseWriter, validationErrors map[string]string) {
	resp := ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Fields:  validationErrors,
	}

	// Log validation errors with structured fields
	logger.With(
		"status_code", http.StatusBadRequest,
		"error_type", "validation_error",
		"validation_errors", validationErrors,
	).Info("Validation error occurred")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: msg,
		Code:    http.StatusUnauthorized,
	}
}
