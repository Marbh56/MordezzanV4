package middleware_test

import (
	"mordezzanV4/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a simple handler for testing
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test handler"))
	})

	// Wrap it with our middleware
	wrapped := middleware.AuthMiddleware(testHandler)

	t.Run("Static files bypass auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/static/css/main.css", nil)
		recorder := httptest.NewRecorder()

		wrapped.ServeHTTP(recorder, req)

		// Should pass through without auth
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d for static file, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("GET requests bypass auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/users/1", nil)
		recorder := httptest.NewRecorder()

		wrapped.ServeHTTP(recorder, req)

		// Should pass through without auth
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d for GET request, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("POST requests go through auth", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/users", nil)
		recorder := httptest.NewRecorder()

		wrapped.ServeHTTP(recorder, req)

		// In our current implementation, POST requests still go through
		// but the middleware just logs it and allows it
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d for POST request, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("PUT requests go through auth", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/users/1", nil)
		recorder := httptest.NewRecorder()

		wrapped.ServeHTTP(recorder, req)

		// In our current implementation, PUT requests still go through
		// but the middleware just logs it and allows it
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d for PUT request, got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("DELETE requests go through auth", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/users/1", nil)
		recorder := httptest.NewRecorder()

		wrapped.ServeHTTP(recorder, req)

		// In our current implementation, DELETE requests still go through
		// but the middleware just logs it and allows it
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code %d for DELETE request, got %d", http.StatusOK, recorder.Code)
		}
	})
}
