package middleware

import (
	"context"
	"mordezzanV4/internal/auth"
	"mordezzanV4/internal/contextkeys"
	"mordezzanV4/internal/logger"
	"net/http"
	"strings"
	"time"
)

type AuthConfig struct {
	JWTSecret string
	Issuer    string
}

// AuthMiddleware protects routes by validating JWT tokens
func JWTAuthMiddleware(config AuthConfig) func(next http.Handler) http.Handler {
	jwtConfig := auth.JWTConfig{
		Secret:     config.JWTSecret,
		Issuer:     config.Issuer,
		Expiration: 24 * time.Hour,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicRoute(r.URL.Path) ||
				(r.Method == "POST" && r.URL.Path == "/users") ||
				(r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/users/")) {
				next.ServeHTTP(w, r)
				return
			}

			// Rest of your authentication code
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := auth.ValidateToken(tokenString, jwtConfig)
			if err != nil {
				logger.With("error", err).Warn("Invalid token")
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, contextkeys.UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/static/",
		"/auth/login",
		"/auth/register",
		"/auth/login-page",
		"/auth/register-page",
		"/health",
		"/users",
		"/",
	}

	// Special case for user GET routes if they should be public
	if strings.HasPrefix(path, "/users/") && strings.Count(path, "/") == 2 {
		return true // Allow GET /users/{id}
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}

	return false
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(contextkeys.UserRoleKey).(string)
		if !ok || role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
