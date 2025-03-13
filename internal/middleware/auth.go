package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") || r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		// Add request-scoped deadline
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// For authenticated routes, you could add user info to context
		// ctx = context.WithValue(ctx, userContextKey, userInfo)

		// Use the enhanced context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
