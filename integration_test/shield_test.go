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

// TestShieldCRUDIntegration tests the CRUD operations for shields
func TestShieldCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Shield CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var shieldID int64

	t.Run("Create Shield", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new shield")

		shieldData := models.CreateShieldInput{
			Name:            "Wooden Shield",
			Cost:            25.0,
			Weight:          7,
			DefenseModifier: 1,
		}

		log.Info("Shield Name: %s", shieldData.Name)
		log.Info("Cost: %.1f gp, Weight: %d lb, Defense Modifier: +%d",
			shieldData.Cost, shieldData.Weight, shieldData.DefenseModifier)

		payload, err := json.Marshal(shieldData)
		if !log.CheckNoError(err, "Marshal shield data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/shields"
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

		var createdShield models.Shield
		if err := json.NewDecoder(resp.Body).Decode(&createdShield); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		shieldID = createdShield.ID
		log.Success("Shield created with ID: %d", shieldID)

		if createdShield.Name != shieldData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", shieldData.Name, createdShield.Name)
			t.Errorf("Expected shield name %s, got %s", shieldData.Name, createdShield.Name)
		}

		if createdShield.Cost != shieldData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", shieldData.Cost, createdShield.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", shieldData.Cost, createdShield.Cost)
		}

		if createdShield.Weight != shieldData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", shieldData.Weight, createdShield.Weight)
			t.Errorf("Expected weight %d, got %d", shieldData.Weight, createdShield.Weight)
		}

		if createdShield.DefenseModifier != shieldData.DefenseModifier {
			log.Error("Defense modifier mismatch. Expected: %d, Got: %d",
				shieldData.DefenseModifier, createdShield.DefenseModifier)
			t.Errorf("Expected defense modifier %d, got %d",
				shieldData.DefenseModifier, createdShield.DefenseModifier)
		}

		log.Success("Shield validation passed")
	})

	log.Separator()

	t.Run("Get Shield", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created shield")
		log.Info("Shield ID: %d", shieldID)

		endpoint := fmt.Sprintf("%s/shields/%d", server.URL, shieldID)
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

		var shield models.Shield
		if err := json.NewDecoder(resp.Body).Decode(&shield); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received shield: ID=%d, Name=%s, Defense Modifier=+%d",
			shield.ID, shield.Name, shield.DefenseModifier)

		if shield.ID != shieldID {
			log.Error("Shield ID mismatch. Expected: %d, Got: %d", shieldID, shield.ID)
			t.Errorf("Expected shield ID %d, got %d", shieldID, shield.ID)
		}

		if shield.Name != "Wooden Shield" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Wooden Shield", shield.Name)
			t.Errorf("Expected shield name 'Wooden Shield', got '%s'", shield.Name)
		}

		log.Success("Shield data validation passed")
	})

	log.Separator()

	t.Run("List Shields", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of shields")

		endpoint := server.URL + "/shields"
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

		var shields []*models.Shield
		if err := json.NewDecoder(resp.Body).Decode(&shields); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d shields in total", len(shields))

		found := false
		for _, s := range shields {
			if s.ID == shieldID {
				found = true
				log.Info("Found our test shield: ID=%d, Name=%s", s.ID, s.Name)
				break
			}
		}

		if !found {
			log.Error("Shield with ID %d not found in shields list", shieldID)
			t.Errorf("Shield with ID %d not found in shields list", shieldID)
		} else {
			log.Success("Shield found in the list")
		}
	})

	log.Separator()

	t.Run("Update Shield", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating shield")
		log.Info("Shield ID: %d", shieldID)

		updateData := models.UpdateShieldInput{
			Name:            "Enhanced Wooden Shield",
			Cost:            50.0,
			Weight:          8,
			DefenseModifier: 2,
		}

		log.Info("New shield name: %s", updateData.Name)
		log.Info("New Cost: %.1f gp, Weight: %d lb, Defense Modifier: +%d",
			updateData.Cost, updateData.Weight, updateData.DefenseModifier)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/shields/%d", server.URL, shieldID)
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

		var updatedShield models.Shield
		if err := json.NewDecoder(resp.Body).Decode(&updatedShield); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated shield data: ID=%d, Name=%s, Cost=%.1f, Defense=+%d",
			updatedShield.ID, updatedShield.Name, updatedShield.Cost, updatedShield.DefenseModifier)

		if updatedShield.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedShield.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedShield.Name)
		}

		if updatedShield.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f",
				updateData.Cost, updatedShield.Cost)
			t.Errorf("Expected cost %.1f, got %.1f",
				updateData.Cost, updatedShield.Cost)
		}

		if updatedShield.DefenseModifier != updateData.DefenseModifier {
			log.Error("Defense modifier mismatch. Expected: %d, Got: %d",
				updateData.DefenseModifier, updatedShield.DefenseModifier)
			t.Errorf("Expected defense modifier %d, got %d",
				updateData.DefenseModifier, updatedShield.DefenseModifier)
		}

		log.Success("Shield update validation passed")
	})

	log.Separator()

	t.Run("Delete Shield", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting shield")
		log.Info("Shield ID: %d", shieldID)

		endpoint := fmt.Sprintf("%s/shields/%d", server.URL, shieldID)
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

		log.Info("Verifying shield deletion by attempting to retrieve it")
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
		log.Success("Shield confirmed deleted (received 404 Not Found)")
	})

	log.Section("SHIELD CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
