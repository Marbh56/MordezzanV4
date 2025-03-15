package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mordezzanV4/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestUserCRUDIntegration tests the CRUD operations for users
func TestUserCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("User CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestApp(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var userID int64
	var authToken string

	t.Run("Create User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new user")

		username := fmt.Sprintf("integrationuser_%s", time.Now().Format("150405"))
		email := fmt.Sprintf("integration_%s@example.com", time.Now().Format("150405"))
		userData := models.CreateUserInput{
			Username: username,
			Email:    email,
			Password: "securepassword123",
		}

		log.Info("Username: %s", userData.Username)
		log.Info("Email: %s", userData.Email)

		payload, err := json.Marshal(userData)
		if !log.CheckNoError(err, "Marshal user data") {
			t.Fatal("Test failed")
		}

		req, err := http.NewRequest("POST", server.URL+"/users", bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")

		log.Info("Sending POST request to %s/users", server.URL)
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

		// Generate JWT token for newly created user
		authToken = generateTestToken(t, userID)
		log.Success("Generated authentication token for user")
	})

	if userID <= 0 {
		log.Error("Cannot continue tests without valid user ID")
		t.Fatal("Cannot continue tests without valid user ID")
	}

	log.Separator()

	t.Run("Get User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created user")
		log.Info("User ID: %d", userID)

		endpoint := fmt.Sprintf("%s/users/%d", server.URL, userID)
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Accept", "application/json")
		// Add the authorization header
		req.Header.Set("Authorization", "Bearer "+authToken)

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

		log.Info("Received user data: ID=%d, Username=%s, Email=%s", user.ID, user.Username, user.Email)

		if user.ID != userID {
			log.Error("User ID mismatch. Expected: %d, Got: %d", userID, user.ID)
			t.Errorf("Expected user ID %d, got %d", userID, user.ID)
		}

		log.Success("User data validation passed")
	})

	log.Separator()

	t.Run("List Users", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of users")

		endpoint := server.URL + "/users"
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		// Add authorization for the list endpoint
		req.Header.Set("Authorization", "Bearer "+authToken)

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

		var users []*models.User
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d users in total", len(users))

		found := false
		for _, u := range users {
			if u.ID == userID {
				found = true
				log.Info("Found our test user: ID=%d, Username=%s", u.ID, u.Username)
				break
			}
		}

		if !found {
			log.Error("User with ID %d not found in user list", userID)
			t.Errorf("User with ID %d not found in user list", userID)
		} else {
			log.Success("User found in the list")
		}
	})

	log.Separator()

	t.Run("Update User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating user")
		log.Info("User ID: %d", userID)

		updateData := struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			Username: "updateduser",
			Email:    "updated@example.com",
		}
		log.Info("New username: %s", updateData.Username)
		log.Info("New email: %s", updateData.Email)

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
		// Add authorization for the update endpoint
		req.Header.Set("Authorization", "Bearer "+authToken)

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
			t.Errorf("Expected username %s, got %s",
				updateData.Username, updatedUser.Username)
		}

		if updatedUser.Email != updateData.Email {
			log.Error("Email mismatch. Expected: %s, Got: %s",
				updateData.Email, updatedUser.Email)
			t.Errorf("Expected email %s, got %s",
				updateData.Email, updatedUser.Email)
		}

		log.Success("Update validation passed")
	})

	log.Separator()

	t.Run("Delete User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting user")
		log.Info("User ID: %d", userID)

		endpoint := fmt.Sprintf("%s/users/%d", server.URL, userID)
		log.Info("Sending DELETE request to %s", endpoint)

		req, err := http.NewRequest("DELETE", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		// Add authorization for the delete endpoint
		req.Header.Set("Authorization", "Bearer "+authToken)

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

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
		// No need for authorization here as we expect a 404 anyway

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
	})

	log.Section("USER CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
