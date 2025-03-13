package middleware

import (
	"context"
	"mordezzanV4/internal/contextkeys"
	"mordezzanV4/internal/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), contextkeys.RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)

		start := time.Now()
		crw := &customResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process the request
		next.ServeHTTP(crw, r.WithContext(ctx))

		duration := time.Since(start)

		// Use structured logging with zap
		logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"status", crw.statusCode,
			"duration", duration,
			"request_id", requestID,
		).Infof("HTTP %s %s", r.Method, r.URL.Path)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}
