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

// TestTreasureCRUDIntegration tests the CRUD operations for treasures
func TestTreasureCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Treasure CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	// Create a test user
	log.Step("Creating a test user for authentication")
	testUser := CreateTestUserWithAuth(t, server)
	log.Success("Created test user with ID: %d, Username: %s", testUser.ID, testUser.Username)

	// First, we need to create a character for the treasure
	var characterID int64
	t.Run("Create Character for Treasure", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating character for treasure tests")

		characterData := models.CreateCharacterInput{
			UserID:       testUser.ID,
			Name:         "Aragorn",
			Class:        "Ranger",
			Level:        8,
			Strength:     16,
			Dexterity:    14,
			Constitution: 15,
			Wisdom:       12,
			Intelligence: 13,
			Charisma:     14,
			HitPoints:    75,
		}

		log.Info("Character Name: %s, Class: %s, Level: %d",
			characterData.Name, characterData.Class, characterData.Level)

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
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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
	})

	log.Separator()

	var treasureID int64

	t.Run("Create Treasure", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new treasure")

		treasureData := models.CreateTreasureInput{
			CharacterID:    &characterID,
			PlatinumCoins:  25,
			GoldCoins:      150,
			ElectrumCoins:  50,
			SilverCoins:    200,
			CopperCoins:    500,
			Gems:           "Ruby (50gp), Sapphire (100gp)",
			ArtObjects:     "Silver statue (25gp), Gold ring (75gp)",
			OtherValuables: "Ancient map to lost city",
			TotalValueGold: 425.0,
		}

		log.Info("Character ID: %d", *treasureData.CharacterID)
		log.Info("Coins: %d pp, %d gp, %d ep, %d sp, %d cp",
			treasureData.PlatinumCoins, treasureData.GoldCoins,
			treasureData.ElectrumCoins, treasureData.SilverCoins,
			treasureData.CopperCoins)
		log.Info("Total Value: %.1f gp", treasureData.TotalValueGold)

		payload, err := json.Marshal(treasureData)
		if !log.CheckNoError(err, "Marshal treasure data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/treasures"
		log.Info("Sending POST request to %s", endpoint)

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		var createdTreasure models.Treasure
		if err := json.NewDecoder(resp.Body).Decode(&createdTreasure); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		treasureID = createdTreasure.ID
		log.Success("Treasure created with ID: %d", treasureID)

		if createdTreasure.GoldCoins != treasureData.GoldCoins {
			log.Error("Gold coins mismatch. Expected: %d, Got: %d",
				treasureData.GoldCoins, createdTreasure.GoldCoins)
			t.Errorf("Expected gold coins %d, got %d",
				treasureData.GoldCoins, createdTreasure.GoldCoins)
		}

		if createdTreasure.TotalValueGold != treasureData.TotalValueGold {
			log.Error("Total value mismatch. Expected: %.1f, Got: %.1f",
				treasureData.TotalValueGold, createdTreasure.TotalValueGold)
			t.Errorf("Expected total value %.1f, got %.1f",
				treasureData.TotalValueGold, createdTreasure.TotalValueGold)
		}

		if createdTreasure.CharacterID == nil || *createdTreasure.CharacterID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %v",
				characterID, createdTreasure.CharacterID)
			t.Errorf("Expected character ID %d, got %v",
				characterID, createdTreasure.CharacterID)
		}

		log.Success("Treasure validation passed")
	})

	log.Separator()

	t.Run("Get Treasure", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created treasure")
		log.Info("Treasure ID: %d", treasureID)

		endpoint := fmt.Sprintf("%s/treasures/%d", server.URL, treasureID)
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		var treasure models.Treasure
		if err := json.NewDecoder(resp.Body).Decode(&treasure); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received treasure: ID=%d, Gold=%d, Total Value=%.1f",
			treasure.ID, treasure.GoldCoins, treasure.TotalValueGold)

		if treasure.ID != treasureID {
			log.Error("Treasure ID mismatch. Expected: %d, Got: %d", treasureID, treasure.ID)
			t.Errorf("Expected treasure ID %d, got %d", treasureID, treasure.ID)
		}

		if treasure.GoldCoins != 150 {
			log.Error("Gold coins mismatch. Expected: %d, Got: %d", 150, treasure.GoldCoins)
			t.Errorf("Expected gold coins 150, got %d", treasure.GoldCoins)
		}

		log.Success("Treasure data validation passed")
	})

	log.Separator()

	t.Run("Get Treasure By Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving treasure by character ID")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d/treasure", server.URL, characterID)
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		var treasure models.Treasure
		if err := json.NewDecoder(resp.Body).Decode(&treasure); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received treasure: ID=%d, Gold=%d, Total Value=%.1f",
			treasure.ID, treasure.GoldCoins, treasure.TotalValueGold)

		if treasure.ID != treasureID {
			log.Error("Treasure ID mismatch. Expected: %d, Got: %d", treasureID, treasure.ID)
			t.Errorf("Expected treasure ID %d, got %d", treasureID, treasure.ID)
		}

		log.Success("Treasure by character validation passed")
	})

	log.Separator()

	t.Run("List Treasures", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of treasures")

		endpoint := server.URL + "/treasures"
		log.Info("Sending GET request to %s", endpoint)

		req, err := http.NewRequest("GET", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		var treasures []*models.Treasure
		if err := json.NewDecoder(resp.Body).Decode(&treasures); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d treasures in total", len(treasures))

		found := false
		for _, tr := range treasures {
			if tr.ID == treasureID {
				found = true
				log.Info("Found our test treasure: ID=%d, Gold=%d", tr.ID, tr.GoldCoins)
				break
			}
		}

		if !found {
			log.Error("Treasure with ID %d not found in treasures list", treasureID)
			t.Errorf("Treasure with ID %d not found in treasures list", treasureID)
		} else {
			log.Success("Treasure found in the list")
		}
	})

	log.Separator()

	t.Run("Update Treasure", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating treasure")
		log.Info("Treasure ID: %d", treasureID)

		updateData := models.UpdateTreasureInput{
			PlatinumCoins:  50,
			GoldCoins:      300,
			ElectrumCoins:  100,
			SilverCoins:    400,
			CopperCoins:    1000,
			Gems:           "Ruby (50gp), Sapphire (100gp), Diamond (500gp)",
			ArtObjects:     "Silver statue (25gp), Gold ring (75gp), Ancient painting (250gp)",
			OtherValuables: "Ancient map to lost city, Magic key to ancient vault",
			TotalValueGold: 975.0,
		}

		log.Info("New coins: %d pp, %d gp, %d ep, %d sp, %d cp",
			updateData.PlatinumCoins, updateData.GoldCoins,
			updateData.ElectrumCoins, updateData.SilverCoins,
			updateData.CopperCoins)
		log.Info("New total value: %.1f gp", updateData.TotalValueGold)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/treasures/%d", server.URL, treasureID)
		log.Info("Sending PUT request to %s", endpoint)

		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(payload))
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		var updatedTreasure models.Treasure
		if err := json.NewDecoder(resp.Body).Decode(&updatedTreasure); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated treasure data: ID=%d, Gold=%d, Total Value=%.1f",
			updatedTreasure.ID, updatedTreasure.GoldCoins, updatedTreasure.TotalValueGold)

		if updatedTreasure.GoldCoins != updateData.GoldCoins {
			log.Error("Gold coins mismatch. Expected: %d, Got: %d",
				updateData.GoldCoins, updatedTreasure.GoldCoins)
			t.Errorf("Expected gold coins %d, got %d",
				updateData.GoldCoins, updatedTreasure.GoldCoins)
		}

		if updatedTreasure.PlatinumCoins != updateData.PlatinumCoins {
			log.Error("Platinum coins mismatch. Expected: %d, Got: %d",
				updateData.PlatinumCoins, updatedTreasure.PlatinumCoins)
			t.Errorf("Expected platinum coins %d, got %d",
				updateData.PlatinumCoins, updatedTreasure.PlatinumCoins)
		}

		if updatedTreasure.TotalValueGold != updateData.TotalValueGold {
			log.Error("Total value mismatch. Expected: %.1f, Got: %.1f",
				updateData.TotalValueGold, updatedTreasure.TotalValueGold)
			t.Errorf("Expected total value %.1f, got %.1f",
				updateData.TotalValueGold, updatedTreasure.TotalValueGold)
		}

		log.Success("Treasure update validation passed")
	})

	log.Separator()

	t.Run("Delete Treasure", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting treasure")
		log.Info("Treasure ID: %d", treasureID)

		endpoint := fmt.Sprintf("%s/treasures/%d", server.URL, treasureID)
		log.Info("Sending DELETE request to %s", endpoint)

		req, err := http.NewRequest("DELETE", endpoint, nil)
		if !log.CheckNoError(err, "Create request") {
			t.Fatal("Test failed")
		}
		req.Header.Set("Authorization", "Bearer "+testUser.Token)

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

		log.Info("Verifying treasure deletion by attempting to retrieve it")
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
		log.Success("Treasure confirmed deleted (received 404 Not Found)")
	})

	log.Section("TREASURE CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
