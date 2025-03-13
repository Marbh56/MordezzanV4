package controllers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"mordezzanV4/internal/controllers"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// MockUserRepository implements the UserRepository interface for testing
type MockUserRepository struct {
	users     map[int64]*models.User
	nextID    int64
	callCount map[string]int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:     make(map[int64]*models.User),
		nextID:    1,
		callCount: make(map[string]int),
	}
}

func (m *MockUserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	m.callCount["GetUser"]++
	user, exists := m.users[id]
	if !exists {
		return nil, sql.ErrNoRows
	}
	return user, nil
}

func (m *MockUserRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	m.callCount["ListUsers"]++
	userList := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		userList = append(userList, user)
	}
	return userList, nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, username, email, passwordHash string) (int64, error) {
	m.callCount["CreateUser"]++
	id := m.nextID
	m.nextID++
	now := time.Now()
	m.users[id] = &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return id, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id int64, username, email string) error {
	m.callCount["UpdateUser"]++
	user, exists := m.users[id]
	if !exists {
		return sql.ErrNoRows
	}
	user.Username = username
	user.Email = email
	user.UpdatedAt = time.Now()
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int64) error {
	m.callCount["DeleteUser"]++
	if _, exists := m.users[id]; !exists {
		return sql.ErrNoRows
	}
	delete(m.users, id)
	return nil
}

func setupTemplates(t *testing.T) *template.Template {
	// Create a temporary template file
	tempDir, err := os.MkdirTemp("", "template_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create a simple user template
	templateContent := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>User Profile</title>
	</head>
	<body>
		<div class="user-profile">
			<h1>User Profile</h1>
			<div>Username: {{.Username}}</div>
			<div>Email: {{.Email}}</div>
		</div>
	</body>
	</html>
	`
	templatePath := filepath.Join(tempDir, "user.html")
	if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
		t.Fatalf("Failed to write template file: %v", err)
	}

	// Parse the template
	tmpl, err := template.ParseGlob(filepath.Join(tempDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}

	// Return the template and a cleanup function
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tmpl
}

func TestUserController(t *testing.T) {
	mockRepo := NewMockUserRepository()
	tmpl := setupTemplates(t)
	controller := controllers.NewUserController(mockRepo, tmpl)

	// Seed the mock repository with some test data
	ctx := context.Background()
	testUserID1, _ := mockRepo.CreateUser(ctx, "testuser1", "test1@example.com", "hashedpw1")
	testUserID2, _ := mockRepo.CreateUser(ctx, "testuser2", "test2@example.com", "hashedpw2")

	// Reset call counts after setup
	mockRepo.callCount = make(map[string]int)

	t.Run("GetUser - JSON response", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/"+strconv.FormatInt(testUserID1, 10), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		controller.GetUser(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Verify content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type %s, got %s", "application/json", contentType)
		}

		// Parse the response
		var user models.User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// Verify user data
		if user.ID != testUserID1 {
			t.Errorf("Expected user ID %d, got %d", testUserID1, user.ID)
		}
		if user.Username != "testuser1" {
			t.Errorf("Expected username 'testuser1', got '%s'", user.Username)
		}
		if user.Email != "test1@example.com" {
			t.Errorf("Expected email 'test1@example.com', got '%s'", user.Email)
		}

		// Verify repository was called
		if mockRepo.callCount["GetUser"] != 1 {
			t.Errorf("Expected GetUser to be called once, got %d", mockRepo.callCount["GetUser"])
		}
	})

	t.Run("GetUser - HTML response", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/"+strconv.FormatInt(testUserID1, 10), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.GetUser(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Verify HTML content
		body := w.Body.String()
		if !strings.Contains(body, "testuser1") || !strings.Contains(body, "test1@example.com") {
			t.Errorf("HTML response doesn't contain expected user data")
		}
	})

	t.Run("GetUser - Invalid ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/invalid", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.GetUser(w, req)

		// Check status code
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("GetUser - Nonexistent user", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/999", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.GetUser(w, req)

		// Check status code
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.ListUsers(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Verify content type
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type %s, got %s", "application/json", contentType)
		}

		// Parse the response
		var users []*models.User
		if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// Verify we got at least 2 users
		if len(users) < 2 {
			t.Errorf("Expected at least 2 users, got %d", len(users))
		}

		// Verify repository was called
		if mockRepo.callCount["ListUsers"] != 1 {
			t.Errorf("Expected ListUsers to be called once, got %d", mockRepo.callCount["ListUsers"])
		}
	})

	t.Run("CreateUser", func(t *testing.T) {
		input := models.CreateUserInput{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
		}

		body, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.CreateUser(w, req)

		// Check status code
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		// Parse the response
		var user models.User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// Verify user data
		if user.Username != "newuser" {
			t.Errorf("Expected username 'newuser', got '%s'", user.Username)
		}
		if user.Email != "new@example.com" {
			t.Errorf("Expected email 'new@example.com', got '%s'", user.Email)
		}

		// Verify repository was called
		if mockRepo.callCount["CreateUser"] != 1 {
			t.Errorf("Expected CreateUser to be called once, got %d", mockRepo.callCount["CreateUser"])
		}
	})

	t.Run("CreateUser - Invalid input", func(t *testing.T) {
		// Missing required fields
		input := struct {
			Username string `json:"username"`
		}{
			Username: "incomplete",
		}

		body, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.CreateUser(w, req)

		// Check status code
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		input := struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		body, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("PUT", "/users/"+strconv.FormatInt(testUserID2, 10), bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.UpdateUser(w, req)

		// Check status code
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Parse the response
		var user models.User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		// Verify user data
		if user.ID != testUserID2 {
			t.Errorf("Expected user ID %d, got %d", testUserID2, user.ID)
		}
		if user.Username != "updateduser" {
			t.Errorf("Expected username 'updateduser', got '%s'", user.Username)
		}
		if user.Email != "updated@example.com" {
			t.Errorf("Expected email 'updated@example.com', got '%s'", user.Email)
		}

		// Verify repository was called
		if mockRepo.callCount["UpdateUser"] != 1 {
			t.Errorf("Expected UpdateUser to be called once, got %d", mockRepo.callCount["UpdateUser"])
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/users/"+strconv.FormatInt(testUserID2, 10), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()

		controller.DeleteUser(w, req)

		// Check status code
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}

		// Verify repository was called
		if mockRepo.callCount["DeleteUser"] != 1 {
			t.Errorf("Expected DeleteUser to be called once, got %d", mockRepo.callCount["DeleteUser"])
		}

		// Verify user was deleted by trying to get it
		req, _ = http.NewRequest("GET", "/users/"+strconv.FormatInt(testUserID2, 10), nil)
		w = httptest.NewRecorder()
		controller.GetUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d after deletion, got %d", http.StatusNotFound, w.Code)
		}
	})
}
