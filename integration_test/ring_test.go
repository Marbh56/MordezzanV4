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

// TestRingCRUDIntegration tests the CRUD operations for rings
func TestRingCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Ring CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var ringID int64

	t.Run("Create Ring", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new ring")

		ringData := models.CreateRingInput{
			Name:        "Ring of Protection",
			Description: "A magical ring that provides enhanced protection to the wearer",
			Cost:        250.0,
			Weight:      1,
		}

		log.Info("Ring Name: %s", ringData.Name)
		log.Info("Description: %s", ringData.Description)
		log.Info("Cost: %.1f gp", ringData.Cost)
		log.Info("Weight: %d lb", ringData.Weight)

		payload, err := json.Marshal(ringData)
		if !log.CheckNoError(err, "Marshal ring data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/rings"
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

		var createdRing models.Ring
		if err := json.NewDecoder(resp.Body).Decode(&createdRing); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		ringID = createdRing.ID
		log.Success("Ring created with ID: %d", ringID)

		if createdRing.Name != ringData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", ringData.Name, createdRing.Name)
			t.Errorf("Expected ring name %s, got %s", ringData.Name, createdRing.Name)
		}

		if createdRing.Description != ringData.Description {
			log.Error("Description mismatch. Expected: %s, Got: %s", ringData.Description, createdRing.Description)
			t.Errorf("Expected description %s, got %s", ringData.Description, createdRing.Description)
		}

		if createdRing.Cost != ringData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", ringData.Cost, createdRing.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", ringData.Cost, createdRing.Cost)
		}

		if createdRing.Weight != ringData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", ringData.Weight, createdRing.Weight)
			t.Errorf("Expected weight %d, got %d", ringData.Weight, createdRing.Weight)
		}

		log.Success("Ring validation passed")
	})

	log.Separator()

	t.Run("Get Ring", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created ring")
		log.Info("Ring ID: %d", ringID)

		endpoint := fmt.Sprintf("%s/rings/%d", server.URL, ringID)
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

		var ring models.Ring
		if err := json.NewDecoder(resp.Body).Decode(&ring); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received ring: ID=%d, Name=%s, Cost=%.1f gp", ring.ID, ring.Name, ring.Cost)

		if ring.ID != ringID {
			log.Error("Ring ID mismatch. Expected: %d, Got: %d", ringID, ring.ID)
			t.Errorf("Expected ring ID %d, got %d", ringID, ring.ID)
		}

		if ring.Name != "Ring of Protection" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Ring of Protection", ring.Name)
			t.Errorf("Expected ring name 'Ring of Protection', got '%s'", ring.Name)
		}

		if ring.Cost != 250.0 {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", 250.0, ring.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", 250.0, ring.Cost)
		}

		log.Success("Ring data validation passed")
	})

	log.Separator()

	t.Run("List Rings", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of rings")

		endpoint := server.URL + "/rings"
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

		var rings []*models.Ring
		if err := json.NewDecoder(resp.Body).Decode(&rings); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d rings in total", len(rings))

		found := false
		for _, r := range rings {
			if r.ID == ringID {
				found = true
				log.Info("Found our test ring: ID=%d, Name=%s, Cost=%.1f gp", r.ID, r.Name, r.Cost)
				break
			}
		}

		if !found {
			log.Error("Ring with ID %d not found in rings list", ringID)
			t.Errorf("Ring with ID %d not found in rings list", ringID)
		} else {
			log.Success("Ring found in the list")
		}
	})

	log.Separator()

	t.Run("Update Ring", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating ring")
		log.Info("Ring ID: %d", ringID)

		updateData := models.UpdateRingInput{
			Name:        "Ring of Greater Protection",
			Description: "An enhanced magical ring that provides superior protection to the wearer",
			Cost:        500.0,
			Weight:      1,
		}

		log.Info("New ring name: %s", updateData.Name)
		log.Info("New description: %s", updateData.Description)
		log.Info("New cost: %.1f gp", updateData.Cost)
		log.Info("Weight: %d lb", updateData.Weight)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/rings/%d", server.URL, ringID)
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

		var updatedRing models.Ring
		if err := json.NewDecoder(resp.Body).Decode(&updatedRing); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated ring data: ID=%d, Name=%s, Cost=%.1f gp", updatedRing.ID, updatedRing.Name, updatedRing.Cost)

		if updatedRing.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", updateData.Name, updatedRing.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedRing.Name)
		}

		if updatedRing.Description != updateData.Description {
			log.Error("Description mismatch. Expected: %s, Got: %s", updateData.Description, updatedRing.Description)
			t.Errorf("Expected description %s, got %s", updateData.Description, updatedRing.Description)
		}

		if updatedRing.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", updateData.Cost, updatedRing.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", updateData.Cost, updatedRing.Cost)
		}

		log.Success("Ring update validation passed")
	})

	log.Separator()

	t.Run("Delete Ring", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting ring")
		log.Info("Ring ID: %d", ringID)

		endpoint := fmt.Sprintf("%s/rings/%d", server.URL, ringID)
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

		log.Info("Verifying ring deletion by attempting to retrieve it")
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
			log.Error("Expected status %d after deletion, got %d", http.StatusNotFound, getResp.StatusCode)
			t.Fatalf("Expected status %d after deletion, got %d", http.StatusNotFound, getResp.StatusCode)
		}
		log.Success("Ring confirmed deleted (received 404 Not Found)")
	})

	log.Section("RING CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
