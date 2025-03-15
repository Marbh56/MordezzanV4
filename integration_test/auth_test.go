package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthMiddlewareIntegration tests the authentication middleware functionality
func TestAuthMiddlewareIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Authentication Middleware Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var userID int64

	t.Run("Create Test User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating test user for authentication tests")

		userData := models.CreateUserInput{
			Username: "authuser",
			Email:    "authuser@example.com",
			Password: "securepassword123",
		}

		log.Info("Username: %s", userData.Username)
		log.Info("Email: %s", userData.Email)

		payload, err := json.Marshal(userData)
		if !log.CheckNoError(err, "Marshal user data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/users"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var createdUser models.User
		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		userID = createdUser.ID
		log.Success("User created with ID: %d", userID)

		if createdUser.Username != userData.Username {
			log.Error("Username mismatch. Expected: %s, Got: %s", userData.Username, createdUser.Username)
			t.Errorf("Expected username %s, got %s", userData.Username, createdUser.Username)
		}

		if createdUser.Email != userData.Email {
			log.Error("Email mismatch. Expected: %s, Got: %s", userData.Email, createdUser.Email)
			t.Errorf("Expected email %s, got %s", userData.Email, createdUser.Email)
		}

		log.Success("User validation passed")
	})

	log.Separator()

	t.Run("Access Static Resources", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing access to static resources")

		endpoint := server.URL + "/static/test.css"
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// Static resources should not return 401 Unauthorized
		// They might return 404 Not Found since we're not creating actual static files
		if resp.StatusCode == http.StatusUnauthorized {
			log.Error("Static resource access failed with Unauthorized status")
			t.Fatalf("Static resource should be accessible without authentication")
		}

		log.Success("Static resources correctly accessible without authentication")
	})

	log.Separator()

	t.Run("Access Public GET Endpoint", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing access to public GET endpoint")
		log.Info("User ID: %d", userID)

		endpoint := fmt.Sprintf("%s/users/%d", server.URL, userID)
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var user models.User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Retrieved user: ID=%d, Username=%s", user.ID, user.Username)

		if user.ID != userID {
			log.Error("User ID mismatch. Expected: %d, Got: %d", userID, user.ID)
			t.Errorf("Expected user ID %d, got %d", userID, user.ID)
		}

		log.Success("Public GET endpoint access validated")
	})

	log.Separator()

	t.Run("Create Resource (POST)", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing POST request to create resource")

		characterData := models.CreateCharacterInput{
			UserID:       userID,
			Name:         "AuthTestCharacter",
			Class:        "Warrior",
			Level:        5,
			Strength:     15,
			Dexterity:    14,
			Constitution: 16,
			Wisdom:       10,
			Intelligence: 12,
			Charisma:     13,
			HitPoints:    45,
		}

		log.Info("Character Name: %s", characterData.Name)
		log.Info("Character Class: %s, Level: %d", characterData.Class, characterData.Level)

		payload, err := json.Marshal(characterData)
		if !log.CheckNoError(err, "Marshal character data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/characters"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// With current middleware implementation, POST should be allowed
		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var createdCharacter models.Character
		if err := json.NewDecoder(resp.Body).Decode(&createdCharacter); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Success("Character created with ID: %d", createdCharacter.ID)
		log.Success("POST request passed through auth middleware as expected")
	})

	log.Separator()

	t.Run("Update Resource (PUT)", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing PUT request to update resource")
		log.Info("User ID: %d", userID)

		updateData := models.UpdateUserInput{
			Username: "updatedauthuser",
			Email:    "updated@example.com",
		}

		log.Info("New Username: %s", updateData.Username)
		log.Info("New Email: %s", updateData.Email)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/users/%d", server.URL, userID)
		log.Info("Sending PUT request to %s", endpoint)

		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// With current middleware implementation, PUT should be allowed
		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var updatedUser models.User
		if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated user data: ID=%d, Username=%s, Email=%s",
			updatedUser.ID, updatedUser.Username, updatedUser.Email)

		if updatedUser.Username != updateData.Username {
			log.Error("Username mismatch. Expected: %s, Got: %s",
				updateData.Username, updatedUser.Username)
			t.Errorf("Expected username %s, got %s", updateData.Username, updatedUser.Username)
		}

		log.Success("PUT request passed through auth middleware as expected")
	})

	log.Separator()

	t.Run("Delete Resource (DELETE)", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing DELETE request")
		log.Info("User ID: %d", userID)

		endpoint := fmt.Sprintf("%s/users/%d", server.URL, userID)
		log.Info("Sending DELETE request to %s", endpoint)

		req, err := http.NewRequest("DELETE", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// With current middleware implementation, DELETE should be allowed
		if resp.StatusCode != http.StatusNoContent {
			log.Error("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		log.Info("Verifying user deletion by attempting to retrieve it")
		getReq, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create verification request") {
			t.Fatal("Test failed")
		}
		getReq.Header.Set("Accept", "application/json")

		getResp, err := http.DefaultClient.Do(getReq)
		if !log.CheckNoError(err, "Send verification request") {
			t.Fatal("Test failed")
		}
		defer getResp.Body.Close()

		if getResp.StatusCode != http.StatusNotFound {
			log.Error("Expected status %d after deletion, got %d",
				http.StatusNotFound, getResp.StatusCode)
			t.Fatalf("Expected status %d after deletion, got %d",
				http.StatusNotFound, getResp.StatusCode)
		}
		log.Success("User confirmed deleted (received 404 Not Found)")
		log.Success("DELETE request passed through auth middleware as expected")
	})

	log.Section("AUTH MIDDLEWARE INTEGRATION TEST COMPLETED SUCCESSFULLY")
}

// TestUserRegistrationIntegration tests the user registration process
func TestUserRegistrationIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("User Registration Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestApp(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	t.Run("Register Valid User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Registering new user with valid data")

		userData := models.CreateUserInput{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "securepassword123",
		}

		log.Info("Username: %s", userData.Username)
		log.Info("Email: %s", userData.Email)
		log.Info("Password length: %d characters", len(userData.Password))

		payload, err := json.Marshal(userData)
		if !log.CheckNoError(err, "Marshal user data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/users"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var createdUser models.User
		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Success("User created with ID: %d", createdUser.ID)

		if createdUser.Username != userData.Username {
			log.Error("Username mismatch. Expected: %s, Got: %s", userData.Username, createdUser.Username)
			t.Errorf("Expected username %s, got %s", userData.Username, createdUser.Username)
		}

		if createdUser.Email != userData.Email {
			log.Error("Email mismatch. Expected: %s, Got: %s", userData.Email, createdUser.Email)
			t.Errorf("Expected email %s, got %s", userData.Email, createdUser.Email)
		}

		log.Success("User registration validation passed")
	})

	log.Separator()

	t.Run("Register User with Invalid Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing registration with invalid data")

		invalidData := models.CreateUserInput{
			Username: "a", // Too short
			Email:    "notanemail",
			Password: "short",
		}

		log.Info("Invalid Username: %s (too short)", invalidData.Username)
		log.Info("Invalid Email: %s (not a valid email)", invalidData.Email)
		log.Info("Invalid Password length: %d characters (too short)", len(invalidData.Password))

		payload, err := json.Marshal(invalidData)
		if !log.CheckNoError(err, "Marshal invalid data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/users"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// Should get a validation error (400 Bad Request)
		if resp.StatusCode != http.StatusBadRequest {
			log.Error("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
			t.Fatalf("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
		}
		log.Success("Correctly received status %d for invalid data", resp.StatusCode)

		// Check error response structure
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Error("Failed to decode error response: %v", err)
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if fields, ok := errorResp["fields"].(map[string]interface{}); ok {
			log.Info("Validation errors received:")
			for field, message := range fields {
				log.Info("  %s: %s", field, message)
			}
		}

		log.Success("Invalid registration properly rejected")
	})

	log.Separator()

	t.Run("Register Duplicate User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing registration with duplicate username/email")

		// First, create a user
		userData := models.CreateUserInput{
			Username: "duplicateuser",
			Email:    "duplicate@example.com",
			Password: "securepassword123",
		}

		log.Info("First creating a user to test duplication")
		log.Info("Username: %s", userData.Username)
		log.Info("Email: %s", userData.Email)

		payload, err := json.Marshal(userData)
		if !log.CheckNoError(err, "Marshal user data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/users"
		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		resp.Body.Close()

		// Now try to create a duplicate
		log.Info("Now trying to register with the same username")

		req, err = http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create duplicate request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send duplicate request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// Should get a validation error (400 Bad Request) for duplicate username
		if resp.StatusCode != http.StatusBadRequest {
			log.Error("Expected status %d for duplicate username, got %d", http.StatusBadRequest, resp.StatusCode)
			t.Fatalf("Expected status %d for duplicate username, got %d", http.StatusBadRequest, resp.StatusCode)
		}
		log.Success("Correctly received status %d for duplicate username", resp.StatusCode)

		// Check error response
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Error("Failed to decode error response: %v", err)
			t.Fatalf("Failed to decode error response: %v", err)
		}

		if fields, ok := errorResp["fields"].(map[string]interface{}); ok {
			if msg, exists := fields["username"]; exists {
				log.Info("Received expected validation error: %s", msg)
			} else {
				log.Warning("Expected username validation error not found in response")
			}
		}

		log.Success("Duplicate user registration properly rejected")
	})

	log.Section("USER REGISTRATION INTEGRATION TEST COMPLETED SUCCESSFULLY")
}

func TestJWTMiddleware(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("JWT Authentication Middleware Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	// Create a test user with authentication
	testUser := CreateTestUserWithAuth(t, server)
	log.Success("Test user created with ID: %d and JWT token", testUser.ID)

	// Test accessing a protected endpoint without authentication
	t.Run("Access_Protected_Endpoint_Without_Auth", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Accessing protected endpoint without authentication")

		// Choose an endpoint that should require authentication
		endpoint := fmt.Sprintf("%s/users/%d", server.URL, testUser.ID)

		// Create request without authentication
		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(`{"username":"updated"}`)))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// Should get unauthorized
		if resp.StatusCode != http.StatusUnauthorized {
			log.Error("Expected status %d (Unauthorized), got %d", http.StatusUnauthorized, resp.StatusCode)
			t.Fatalf("Expected status %d (Unauthorized), got %d", http.StatusUnauthorized, resp.StatusCode)
		}
		log.Success("Correctly received status %d (Unauthorized)", resp.StatusCode)
	})

	// Test accessing the same endpoint with authentication
	t.Run("Access_Protected_Endpoint_With_Auth", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Accessing protected endpoint with authentication")

		// Choose the same endpoint
		endpoint := fmt.Sprintf("%s/users/%d", server.URL, testUser.ID)

		// Create authenticated request
		updateData := map[string]string{"username": "updated_name"}
		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		req := AuthenticatedRequest(t, "PUT", endpoint, bytes.NewBuffer(payload), testUser)

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		// Should be successful
		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d (OK), got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d (OK), got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Correctly received status %d (OK) with authentication", resp.StatusCode)
	})
	log.Section("JWT MIDDLEWARE INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
