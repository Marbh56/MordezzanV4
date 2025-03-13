package middleware

import (
	"fmt"
	"mordezzanV4/internal/logger"
	"net/http"
	"runtime/debug"

	apperrors "mordezzanV4/internal/errors"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()

				// Log the panic with structured fields
				logger.With(
					"error", fmt.Sprintf("%v", err),
					"stack", string(stack),
					"url", r.URL.String(),
					"method", r.Method,
				).Error("PANIC in HTTP handler")

				// Return a generic error to the client
				serverErr := apperrors.NewInternalError(fmt.Errorf("%v", err))
				apperrors.HandleError(w, serverErr)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
