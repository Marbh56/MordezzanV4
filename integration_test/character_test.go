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

// TestCharacterCRUDIntegration tests the CRUD operations for characters
func TestCharacterCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Character CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	// Create a test user with authentication
	log.Step("Creating test user with authentication")
	testUser := CreateTestUserWithAuth(t, server)
	log.Success("Test user created with ID: %d", testUser.ID)
	log.Success("Authentication token generated for user")

	var characterID int64

	t.Run("Create_Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new character")

		characterData := models.CreateCharacterInput{
			UserID:       testUser.ID,
			Name:         "Aragorn",
			Class:        "Ranger",
			Level:        5, // Added level field
			Strength:     16,
			Dexterity:    14,
			Constitution: 15,
			Wisdom:       13,
			Intelligence: 12,
			Charisma:     15,
			HitPoints:    20,
		}

		log.Info("Character Name: %s", characterData.Name)
		log.Info("Stats: Level=%d, STR=%d, DEX=%d, CON=%d, WIS=%d, INT=%d, CHA=%d, HP=%d",
			characterData.Level, characterData.Strength, characterData.Dexterity, characterData.Constitution,
			characterData.Wisdom, characterData.Intelligence, characterData.Charisma,
			characterData.HitPoints)

		payload, err := json.Marshal(characterData)
		if !log.CheckNoError(err, "Marshal character data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/characters"
		log.Info("Sending authenticated POST request to %s", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "POST", endpoint, bytes.NewBuffer(payload), testUser)

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

		var createdCharacter models.Character
		if err := json.NewDecoder(resp.Body).Decode(&createdCharacter); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		characterID = createdCharacter.ID
		log.Success("Character created with ID: %d", characterID)

		if createdCharacter.Name != characterData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", characterData.Name, createdCharacter.Name)
			t.Errorf("Expected character name %s, got %s", characterData.Name, createdCharacter.Name)
		}

		if createdCharacter.Class != characterData.Class {
			log.Error("Class mismatch. Expected: %s, Got: %s", characterData.Class, createdCharacter.Class)
			t.Errorf("Expected character class %s, got %s", characterData.Class, createdCharacter.Class)
		}

		if createdCharacter.Level != characterData.Level {
			log.Error("Level mismatch. Expected: %d, Got: %d", characterData.Level, createdCharacter.Level)
			t.Errorf("Expected level %d, got %d", characterData.Level, createdCharacter.Level)
		}

		if createdCharacter.Strength != characterData.Strength {
			log.Error("Strength mismatch. Expected: %d, Got: %d", characterData.Strength, createdCharacter.Strength)
			t.Errorf("Expected strength %d, got %d", characterData.Strength, createdCharacter.Strength)
		}

		log.Success("Character validation passed")
	})

	log.Separator()

	t.Run("Get_Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created character")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
		log.Info("Sending authenticated GET request to %s", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "GET", endpoint, nil, testUser)
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

		var character models.Character
		if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received character: ID=%d, Name=%s, Class=%s, Level=%d",
			character.ID, character.Name, character.Class, character.Level)

		if character.ID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %d", characterID, character.ID)
			t.Errorf("Expected character ID %d, got %d", characterID, character.ID)
		}

		if character.Name != "Aragorn" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Aragorn", character.Name)
			t.Errorf("Expected character name 'Aragorn', got '%s'", character.Name)
		}

		if character.Class != "Ranger" {
			log.Error("Class mismatch. Expected: %s, Got: %s", "Ranger", character.Class)
			t.Errorf("Expected character class 'Ranger', got '%s'", character.Class)
		}

		if character.Level != 5 {
			log.Error("Level mismatch. Expected: %d, Got: %d", 5, character.Level)
			t.Errorf("Expected character level 5, got %d", character.Level)
		}

		log.Success("Character data validation passed")
	})

	log.Separator()

	t.Run("Get_Characters_By_User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving characters for user")
		log.Info("User ID: %d", testUser.ID)

		endpoint := fmt.Sprintf("%s/users/%d/characters", server.URL, testUser.ID)
		log.Info("Sending authenticated GET request to %s", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "GET", endpoint, nil, testUser)

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var characters []*models.Character
		if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d characters for user", len(characters))

		if len(characters) < 1 {
			log.Error("Expected at least 1 character, got %d", len(characters))
			t.Errorf("Expected at least 1 character, got %d", len(characters))
		} else {
			found := false
			for _, c := range characters {
				if c.ID == characterID {
					found = true
					log.Info("Found our test character: ID=%d, Name=%s, Class=%s, Level=%d",
						c.ID, c.Name, c.Class, c.Level)
					break
				}
			}

			if !found {
				log.Error("Character with ID %d not found in user's characters list", characterID)
				t.Errorf("Character with ID %d not found in user's characters list", characterID)
			} else {
				log.Success("Character found in the list")
			}
		}
	})

	log.Separator()

	t.Run("Update_Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating character")
		log.Info("Character ID: %d", characterID)

		updateData := models.UpdateCharacterInput{
			Name:         "Strider",
			Class:        "Ranger",
			Level:        6, // Increase the level
			Strength:     17,
			Dexterity:    15,
			Constitution: 16,
			Wisdom:       14,
			Intelligence: 13,
			Charisma:     16,
			HitPoints:    25,
		}

		log.Info("New character name: %s", updateData.Name)
		log.Info("New stats: Level=%d, STR=%d, DEX=%d, CON=%d, WIS=%d, INT=%d, CHA=%d, HP=%d",
			updateData.Level, updateData.Strength, updateData.Dexterity, updateData.Constitution,
			updateData.Wisdom, updateData.Intelligence, updateData.Charisma,
			updateData.HitPoints)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
		log.Info("Sending authenticated PUT request to %s", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "PUT", endpoint, bytes.NewBuffer(payload), testUser)

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

		var updatedCharacter models.Character
		if err := json.NewDecoder(resp.Body).Decode(&updatedCharacter); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated character data: ID=%d, Name=%s, Class=%s, Level=%d, Strength=%d",
			updatedCharacter.ID, updatedCharacter.Name, updatedCharacter.Class,
			updatedCharacter.Level, updatedCharacter.Strength)

		if updatedCharacter.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedCharacter.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedCharacter.Name)
		}

		if updatedCharacter.Class != updateData.Class {
			log.Error("Class mismatch. Expected: %s, Got: %s",
				updateData.Class, updatedCharacter.Class)
			t.Errorf("Expected class %s, got %s", updateData.Class, updatedCharacter.Class)
		}

		if updatedCharacter.Level != updateData.Level {
			log.Error("Level mismatch. Expected: %d, Got: %d",
				updateData.Level, updatedCharacter.Level)
			t.Errorf("Expected level %d, got %d", updateData.Level, updatedCharacter.Level)
		}

		if updatedCharacter.Strength != updateData.Strength {
			log.Error("Strength mismatch. Expected: %d, Got: %d",
				updateData.Strength, updatedCharacter.Strength)
			t.Errorf("Expected strength %d, got %d", updateData.Strength, updatedCharacter.Strength)
		}

		log.Success("Character update validation passed")
	})

	log.Separator()

	t.Run("Delete_Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting character")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
		log.Info("Sending authenticated DELETE request to %s", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "DELETE", endpoint, nil, testUser)

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

		log.Info("Verifying character deletion by attempting to retrieve it")
		getReq, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create verification request") {
			t.Fatal("Test failed")
		}
		getReq.Header.Set("Accept", "application/json")
		getReq.Header.Set("Authorization", "Bearer "+testUser.Token)

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
		log.Success("Character confirmed deleted (received 404 Not Found)")
	})

	t.Run("Character_With_Invalid_Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing validation with invalid character data")

		log.Info("Invalid data: empty name, level=0, strength=25")

		invalidData := models.CreateCharacterInput{
			UserID:       testUser.ID,
			Name:         "", // Invalid: empty name
			Class:        "Warrior",
			Level:        0,  // Invalid: level=0
			Strength:     25, // Invalid: strength > 18
			Dexterity:    14,
			Constitution: 15,
			Wisdom:       13,
			Intelligence: 12,
			Charisma:     15,
			HitPoints:    20,
		}

		payload, err := json.Marshal(invalidData)
		if !log.CheckNoError(err, "Marshal invalid data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/characters"
		log.Info("Sending authenticated POST request to %s with invalid data", endpoint)

		// Use AuthenticatedRequest helper
		req := AuthenticatedRequest(t, "POST", endpoint, bytes.NewBuffer(payload), testUser)

		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			log.Error("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
			t.Fatalf("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
		}
		log.Success("Correctly received BadRequest response: %d", resp.StatusCode)

		var errorResp struct {
			Status  int               `json:"status"`
			Message string            `json:"message"`
			Fields  map[string]string `json:"fields"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Error("Failed to decode error response: %v", err)
			t.Fatalf("Failed to decode error response: %v", err)
		}

		log.Info("Error response: %s", errorResp.Message)

		if errorResp.Message != "Validation failed" || len(errorResp.Fields) == 0 {
			log.Error("Expected validation error with field details")
			t.Errorf("Expected validation error message with field details")
		} else {
			log.Success("Validation fields returned properly")
			for field, msg := range errorResp.Fields {
				log.Info("Validation error: %s - %s", field, msg)
			}
		}
	})

	log.Section("CHARACTER CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
