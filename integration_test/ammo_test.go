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

// TestAmmoCRUDIntegration tests the CRUD operations for ammo
func TestAmmoCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Ammo CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var ammoID int64

	t.Run("Create Ammo", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new ammo")

		ammoData := models.CreateAmmoInput{
			Name:   "Arrows (20)",
			Cost:   1.0,
			Weight: 1,
		}

		log.Info("Ammo Name: %s", ammoData.Name)
		log.Info("Cost: %.1f gp, Weight: %d lb", ammoData.Cost, ammoData.Weight)

		payload, err := json.Marshal(ammoData)
		if !log.CheckNoError(err, "Marshal ammo data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/ammo"
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

		var createdAmmo models.Ammo
		if err := json.NewDecoder(resp.Body).Decode(&createdAmmo); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		ammoID = createdAmmo.ID
		log.Success("Ammo created with ID: %d", ammoID)

		if createdAmmo.Name != ammoData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", ammoData.Name, createdAmmo.Name)
			t.Errorf("Expected ammo name %s, got %s", ammoData.Name, createdAmmo.Name)
		}

		if createdAmmo.Cost != ammoData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", ammoData.Cost, createdAmmo.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", ammoData.Cost, createdAmmo.Cost)
		}

		if createdAmmo.Weight != ammoData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", ammoData.Weight, createdAmmo.Weight)
			t.Errorf("Expected weight %d, got %d", ammoData.Weight, createdAmmo.Weight)
		}

		log.Success("Ammo validation passed")
	})

	log.Separator()

	t.Run("Get Ammo", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created ammo")
		log.Info("Ammo ID: %d", ammoID)

		endpoint := fmt.Sprintf("%s/ammo/%d", server.URL, ammoID)
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

		var ammo models.Ammo
		if err := json.NewDecoder(resp.Body).Decode(&ammo); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received ammo: ID=%d, Name=%s",
			ammo.ID, ammo.Name)

		if ammo.ID != ammoID {
			log.Error("Ammo ID mismatch. Expected: %d, Got: %d", ammoID, ammo.ID)
			t.Errorf("Expected ammo ID %d, got %d", ammoID, ammo.ID)
		}

		if ammo.Name != "Arrows (20)" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Arrows (20)", ammo.Name)
			t.Errorf("Expected ammo name 'Arrows (20)', got '%s'", ammo.Name)
		}

		log.Success("Ammo data validation passed")
	})

	log.Separator()

	t.Run("List Ammo", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of ammo")

		endpoint := server.URL + "/ammo"
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

		var ammoList []*models.Ammo
		if err := json.NewDecoder(resp.Body).Decode(&ammoList); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d ammo in total", len(ammoList))

		found := false
		for _, a := range ammoList {
			if a.ID == ammoID {
				found = true
				log.Info("Found our test ammo: ID=%d, Name=%s", a.ID, a.Name)
				break
			}
		}

		if !found {
			log.Error("Ammo with ID %d not found in ammo list", ammoID)
			t.Errorf("Ammo with ID %d not found in ammo list", ammoID)
		} else {
			log.Success("Ammo found in the list")
		}
	})

	log.Separator()

	t.Run("Update Ammo", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating ammo")
		log.Info("Ammo ID: %d", ammoID)

		updateData := models.UpdateAmmoInput{
			Name:   "Bolts (20)",
			Cost:   2.0,
			Weight: 2,
		}

		log.Info("New ammo name: %s", updateData.Name)
		log.Info("New Cost: %.1f gp, Weight: %d lb",
			updateData.Cost, updateData.Weight)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/ammo/%d", server.URL, ammoID)
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

		var updatedAmmo models.Ammo
		if err := json.NewDecoder(resp.Body).Decode(&updatedAmmo); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated ammo data: ID=%d, Name=%s, Cost=%.1f, Weight=%d",
			updatedAmmo.ID, updatedAmmo.Name, updatedAmmo.Cost, updatedAmmo.Weight)

		if updatedAmmo.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedAmmo.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedAmmo.Name)
		}

		if updatedAmmo.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f",
				updateData.Cost, updatedAmmo.Cost)
			t.Errorf("Expected cost %.1f, got %.1f",
				updateData.Cost, updatedAmmo.Cost)
		}

		if updatedAmmo.Weight != updateData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d",
				updateData.Weight, updatedAmmo.Weight)
			t.Errorf("Expected weight %d, got %d",
				updateData.Weight, updatedAmmo.Weight)
		}

		log.Success("Ammo update validation passed")
	})

	log.Separator()

	t.Run("Delete Ammo", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting ammo")
		log.Info("Ammo ID: %d", ammoID)

		endpoint := fmt.Sprintf("%s/ammo/%d", server.URL, ammoID)
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

		log.Info("Verifying ammo deletion by attempting to retrieve it")
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
		log.Success("Ammo confirmed deleted (received 404 Not Found)")
	})

	log.Section("AMMO CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
