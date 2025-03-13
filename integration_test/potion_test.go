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

// TestPotionCRUDIntegration tests the CRUD operations for potions
func TestPotionCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Potion CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var potionID int64

	t.Run("Create Potion", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new potion")

		potionData := models.CreatePotionInput{
			Name:        "Healing Potion",
			Description: "Restores 20 hit points when consumed.",
			Uses:        3,
			Weight:      1,
		}

		log.Info("Potion Name: %s", potionData.Name)
		log.Info("Description: %s, Uses: %d, Weight: %d",
			potionData.Description, potionData.Uses, potionData.Weight)

		payload, err := json.Marshal(potionData)
		if !log.CheckNoError(err, "Marshal potion data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/potions"
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

		var createdPotion models.Potion
		if err := json.NewDecoder(resp.Body).Decode(&createdPotion); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		potionID = createdPotion.ID
		log.Success("Potion created with ID: %d", potionID)

		if createdPotion.Name != potionData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", potionData.Name, createdPotion.Name)
			t.Errorf("Expected potion name %s, got %s", potionData.Name, createdPotion.Name)
		}

		if createdPotion.Uses != potionData.Uses {
			log.Error("Uses mismatch. Expected: %d, Got: %d", potionData.Uses, createdPotion.Uses)
			t.Errorf("Expected uses %d, got %d", potionData.Uses, createdPotion.Uses)
		}

		if createdPotion.Weight != potionData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", potionData.Weight, createdPotion.Weight)
			t.Errorf("Expected weight %d, got %d", potionData.Weight, createdPotion.Weight)
		}

		log.Success("Potion validation passed")
	})

	log.Separator()

	t.Run("Get Potion", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created potion")
		log.Info("Potion ID: %d", potionID)

		endpoint := fmt.Sprintf("%s/potions/%d", server.URL, potionID)
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

		var potion models.Potion
		if err := json.NewDecoder(resp.Body).Decode(&potion); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received potion: ID=%d, Name=%s, Uses=%d",
			potion.ID, potion.Name, potion.Uses)

		if potion.ID != potionID {
			log.Error("Potion ID mismatch. Expected: %d, Got: %d", potionID, potion.ID)
			t.Errorf("Expected potion ID %d, got %d", potionID, potion.ID)
		}

		if potion.Name != "Healing Potion" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Healing Potion", potion.Name)
			t.Errorf("Expected potion name 'Healing Potion', got '%s'", potion.Name)
		}

		log.Success("Potion data validation passed")
	})

	log.Separator()

	t.Run("List Potions", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of potions")

		endpoint := server.URL + "/potions"
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

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var potions []*models.Potion
		if err := json.NewDecoder(resp.Body).Decode(&potions); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d potions in total", len(potions))

		found := false
		for _, p := range potions {
			if p.ID == potionID {
				found = true
				log.Info("Found our test potion: ID=%d, Name=%s", p.ID, p.Name)
				break
			}
		}

		if !found {
			log.Error("Potion with ID %d not found in potions list", potionID)
			t.Errorf("Potion with ID %d not found in potions list", potionID)
		} else {
			log.Success("Potion found in the list")
		}
	})

	log.Separator()

	t.Run("Update Potion", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating potion")
		log.Info("Potion ID: %d", potionID)

		updateData := models.UpdatePotionInput{
			Name:        "Greater Healing Potion",
			Description: "Restores 50 hit points when consumed.",
			Uses:        2,
			Weight:      2,
		}

		log.Info("New potion name: %s", updateData.Name)
		log.Info("New description: %s, Uses: %d", updateData.Description, updateData.Uses)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/potions/%d", server.URL, potionID)
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

		var updatedPotion models.Potion
		if err := json.NewDecoder(resp.Body).Decode(&updatedPotion); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated potion data: ID=%d, Name=%s, Uses=%d, Weight=%d",
			updatedPotion.ID, updatedPotion.Name, updatedPotion.Uses, updatedPotion.Weight)

		if updatedPotion.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedPotion.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedPotion.Name)
		}

		if updatedPotion.Uses != updateData.Uses {
			log.Error("Uses mismatch. Expected: %d, Got: %d",
				updateData.Uses, updatedPotion.Uses)
			t.Errorf("Expected uses %d, got %d",
				updateData.Uses, updatedPotion.Uses)
		}

		if updatedPotion.Weight != updateData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d",
				updateData.Weight, updatedPotion.Weight)
			t.Errorf("Expected weight %d, got %d",
				updateData.Weight, updatedPotion.Weight)
		}

		log.Success("Potion update validation passed")
	})

	log.Separator()

	t.Run("Delete Potion", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting potion")
		log.Info("Potion ID: %d", potionID)

		endpoint := fmt.Sprintf("%s/potions/%d", server.URL, potionID)
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

		log.Info("Verifying potion deletion by attempting to retrieve it")
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
		log.Success("Potion confirmed deleted (received 404 Not Found)")
	})

	log.Section("POTION CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
