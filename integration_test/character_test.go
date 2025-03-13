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

	log.Step("Creating test user")
	userID := createTestUser(t, server)
	log.Success("Test user created with ID: %d", userID)

	var characterID int64

	t.Run("Create Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new character")

		characterData := models.CreateCharacterInput{
			UserID:       userID,
			Name:         "Aragorn",
			Strength:     16,
			Dexterity:    14,
			Constitution: 15,
			Wisdom:       13,
			Intelligence: 12,
			Charisma:     15,
			HitPoints:    20,
		}

		log.Info("Character Name: %s", characterData.Name)
		log.Info("Stats: STR=%d, DEX=%d, CON=%d, WIS=%d, INT=%d, CHA=%d, HP=%d",
			characterData.Strength, characterData.Dexterity, characterData.Constitution,
			characterData.Wisdom, characterData.Intelligence, characterData.Charisma,
			characterData.HitPoints)

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
			log.Error("Name mismatch. Expected: %s, Got: %s",
				characterData.Name, createdCharacter.Name)
			t.Errorf("Expected character name %s, got %s",
				characterData.Name, createdCharacter.Name)
		}

		if createdCharacter.Strength != characterData.Strength {
			log.Error("Strength mismatch. Expected: %d, Got: %d",
				characterData.Strength, createdCharacter.Strength)
			t.Errorf("Expected strength %d, got %d",
				characterData.Strength, createdCharacter.Strength)
		}

		if createdCharacter.UserID != userID {
			log.Error("User ID mismatch. Expected: %d, Got: %d",
				userID, createdCharacter.UserID)
			t.Errorf("Expected user ID %d, got %d", userID, createdCharacter.UserID)
		}

		log.Success("Character validation passed")
	})

	log.Separator()

	t.Run("Get Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created character")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
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

		var character models.Character
		if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received character: ID=%d, Name=%s, User ID=%d",
			character.ID, character.Name, character.UserID)

		if character.ID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %d",
				characterID, character.ID)
			t.Errorf("Expected character ID %d, got %d", characterID, character.ID)
		}

		if character.Name != "Aragorn" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Aragorn", character.Name)
			t.Errorf("Expected character name 'Aragorn', got '%s'", character.Name)
		}

		if character.UserID != userID {
			log.Error("User ID mismatch. Expected: %d, Got: %d",
				userID, character.UserID)
			t.Errorf("Expected user ID %d, got %d", userID, character.UserID)
		}

		log.Success("Character data validation passed")
	})

	log.Separator()

	t.Run("Get Characters By User", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving characters for user")
		log.Info("User ID: %d", userID)

		endpoint := fmt.Sprintf("%s/users/%d/characters", server.URL, userID)
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
			t.Fatalf("Expected at least 1 character, got %d", len(characters))
		}

		found := false
		for _, c := range characters {
			if c.ID == characterID {
				found = true
				log.Info("Found our test character: ID=%d, Name=%s", c.ID, c.Name)
				if c.Name != "Aragorn" {
					log.Error("Name mismatch. Expected: %s, Got: %s",
						"Aragorn", c.Name)
					t.Errorf("Expected character name 'Aragorn', got '%s'", c.Name)
				}
				break
			}
		}

		if !found {
			log.Error("Character with ID %d not found in characters list", characterID)
			t.Errorf("Character with ID %d not found in characters list", characterID)
		} else {
			log.Success("Character found in the list")
		}
	})

	log.Separator()

	t.Run("Update Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating character")
		log.Info("Character ID: %d", characterID)

		updateData := models.UpdateCharacterInput{
			Name:         "Strider",
			Strength:     17,
			Dexterity:    15,
			Constitution: 16,
			Wisdom:       14,
			Intelligence: 13,
			Charisma:     16,
			HitPoints:    25,
		}

		log.Info("New character name: %s", updateData.Name)
		log.Info("New stats: STR=%d, DEX=%d, CON=%d, WIS=%d, INT=%d, CHA=%d, HP=%d",
			updateData.Strength, updateData.Dexterity, updateData.Constitution,
			updateData.Wisdom, updateData.Intelligence, updateData.Charisma,
			updateData.HitPoints)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
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

		log.Info("Updated character data: ID=%d, Name=%s",
			updatedCharacter.ID, updatedCharacter.Name)
		log.Info("Updated stats: STR=%d, DEX=%d, CON=%d, WIS=%d, INT=%d, CHA=%d, HP=%d",
			updatedCharacter.Strength, updatedCharacter.Dexterity, updatedCharacter.Constitution,
			updatedCharacter.Wisdom, updatedCharacter.Intelligence, updatedCharacter.Charisma,
			updatedCharacter.HitPoints)

		if updatedCharacter.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedCharacter.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedCharacter.Name)
		}

		if updatedCharacter.Strength != updateData.Strength {
			log.Error("Strength mismatch. Expected: %d, Got: %d",
				updateData.Strength, updatedCharacter.Strength)
			t.Errorf("Expected strength %d, got %d",
				updateData.Strength, updatedCharacter.Strength)
		}

		if updatedCharacter.HitPoints != updateData.HitPoints {
			log.Error("Hit points mismatch. Expected: %d, Got: %d",
				updateData.HitPoints, updatedCharacter.HitPoints)
			t.Errorf("Expected hit points %d, got %d",
				updateData.HitPoints, updatedCharacter.HitPoints)
		}

		log.Success("Character update validation passed")
	})

	log.Separator()

	t.Run("Delete Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting character")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d", server.URL, characterID)
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

	t.Run("Character With Invalid Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing validation with invalid character data")

		invalidCharacter := models.CreateCharacterInput{
			UserID:       userID,
			Name:         "", // Invalid: empty name
			Strength:     25, // Invalid: too high (must be <= 18)
			Dexterity:    14,
			Constitution: 15,
			Wisdom:       13,
			Intelligence: 12,
			Charisma:     15,
			HitPoints:    20,
		}

		log.Info("Invalid data: empty name, strength=25")

		payload, err := json.Marshal(invalidCharacter)
		if !log.CheckNoError(err, "Marshal invalid character data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/characters"
		log.Info("Sending POST request to %s with invalid data", endpoint)

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

		if resp.StatusCode != http.StatusBadRequest {
			log.Error("Expected status %d for invalid data, got %d",
				http.StatusBadRequest, resp.StatusCode)
			t.Fatalf("Expected status %d for invalid data, got %d",
				http.StatusBadRequest, resp.StatusCode)
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

		if errorResp.Status != http.StatusBadRequest {
			log.Error("Expected error status %d, got %d",
				http.StatusBadRequest, errorResp.Status)
			t.Errorf("Expected error status %d, got %d",
				http.StatusBadRequest, errorResp.Status)
		}

		if errorResp.Fields == nil || len(errorResp.Fields) == 0 {
			log.Error("Expected validation error fields, got none")
			t.Error("Expected validation error fields, got none")
		} else {
			log.Success("Validation fields returned properly")
			for field, message := range errorResp.Fields {
				log.Info("Validation error: %s - %s", field, message)
			}
		}
	})

	log.Section("CHARACTER CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
