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

func TestEquipmentCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Equipment CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var equipmentID int64

	t.Run("Create Equipment", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new equipment")

		equipmentData := models.CreateEquipmentInput{
			Name:        "Rope (50 ft)",
			Description: "Strong hemp rope that can support up to 400 pounds",
			Cost:        1.0,
			Weight:      10,
		}

		log.Info("Equipment Name: %s", equipmentData.Name)
		log.Info("Description: %s, Cost: %.2f gp, Weight: %d lb",
			equipmentData.Description, equipmentData.Cost, equipmentData.Weight)

		payload, err := json.Marshal(equipmentData)
		if !log.CheckNoError(err, "Marshal equipment data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/equipment"
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

		var createdEquipment models.Equipment
		if err := json.NewDecoder(resp.Body).Decode(&createdEquipment); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		equipmentID = createdEquipment.ID
		log.Success("Equipment created with ID: %d", equipmentID)

		if createdEquipment.Name != equipmentData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", equipmentData.Name, createdEquipment.Name)
			t.Errorf("Expected equipment name %s, got %s", equipmentData.Name, createdEquipment.Name)
		}

		if createdEquipment.Description != equipmentData.Description {
			log.Error("Description mismatch. Expected: %s, Got: %s", equipmentData.Description, createdEquipment.Description)
			t.Errorf("Expected description %s, got %s", equipmentData.Description, createdEquipment.Description)
		}

		if createdEquipment.Cost != equipmentData.Cost {
			log.Error("Cost mismatch. Expected: %.2f, Got: %.2f", equipmentData.Cost, createdEquipment.Cost)
			t.Errorf("Expected cost %.2f, got %.2f", equipmentData.Cost, createdEquipment.Cost)
		}

		if createdEquipment.Weight != equipmentData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", equipmentData.Weight, createdEquipment.Weight)
			t.Errorf("Expected weight %d, got %d", equipmentData.Weight, createdEquipment.Weight)
		}

		log.Success("Equipment validation passed")
	})

	log.Separator()

	t.Run("Get Equipment", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created equipment")
		log.Info("Equipment ID: %d", equipmentID)

		endpoint := fmt.Sprintf("%s/equipment/%d", server.URL, equipmentID)
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

		var equipment models.Equipment
		if err := json.NewDecoder(resp.Body).Decode(&equipment); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received equipment: ID=%d, Name=%s, Cost=%.2f, Weight=%d",
			equipment.ID, equipment.Name, equipment.Cost, equipment.Weight)

		if equipment.ID != equipmentID {
			log.Error("Equipment ID mismatch. Expected: %d, Got: %d", equipmentID, equipment.ID)
			t.Errorf("Expected equipment ID %d, got %d", equipmentID, equipment.ID)
		}

		if equipment.Name != "Rope (50 ft)" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Rope (50 ft)", equipment.Name)
			t.Errorf("Expected equipment name 'Rope (50 ft)', got '%s'", equipment.Name)
		}

		log.Success("Equipment data validation passed")
	})

	log.Separator()

	t.Run("List Equipment", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of equipment")

		endpoint := server.URL + "/equipment"
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

		var equipmentItems []*models.Equipment
		if err := json.NewDecoder(resp.Body).Decode(&equipmentItems); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d equipment items in total", len(equipmentItems))

		found := false
		for _, eq := range equipmentItems {
			if eq.ID == equipmentID {
				found = true
				log.Info("Found our test equipment: ID=%d, Name=%s", eq.ID, eq.Name)
				break
			}
		}

		if !found {
			log.Error("Equipment with ID %d not found in equipment list", equipmentID)
			t.Errorf("Equipment with ID %d not found in equipment list", equipmentID)
		} else {
			log.Success("Equipment found in the list")
		}
	})

	log.Separator()

	t.Run("Update Equipment", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating equipment")
		log.Info("Equipment ID: %d", equipmentID)

		updateData := models.UpdateEquipmentInput{
			Name:        "Silk Rope (50 ft)",
			Description: "Lightweight silk rope that can support up to 350 pounds",
			Cost:        10.0,
			Weight:      5,
		}

		log.Info("New equipment name: %s", updateData.Name)
		log.Info("New description: %s, Cost: %.2f, Weight: %d",
			updateData.Description, updateData.Cost, updateData.Weight)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/equipment/%d", server.URL, equipmentID)
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

		var updatedEquipment models.Equipment
		if err := json.NewDecoder(resp.Body).Decode(&updatedEquipment); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated equipment data: ID=%d, Name=%s, Cost=%.2f, Weight=%d",
			updatedEquipment.ID, updatedEquipment.Name, updatedEquipment.Cost, updatedEquipment.Weight)

		if updatedEquipment.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedEquipment.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedEquipment.Name)
		}

		if updatedEquipment.Description != updateData.Description {
			log.Error("Description mismatch. Expected: %s, Got: %s",
				updateData.Description, updatedEquipment.Description)
			t.Errorf("Expected description %s, got %s",
				updateData.Description, updatedEquipment.Description)
		}

		if updatedEquipment.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.2f, Got: %.2f",
				updateData.Cost, updatedEquipment.Cost)
			t.Errorf("Expected cost %.2f, got %.2f",
				updateData.Cost, updatedEquipment.Cost)
		}

		log.Success("Equipment update validation passed")
	})

	log.Separator()

	t.Run("Delete Equipment", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting equipment")
		log.Info("Equipment ID: %d", equipmentID)

		endpoint := fmt.Sprintf("%s/equipment/%d", server.URL, equipmentID)
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

		log.Info("Verifying equipment deletion by attempting to retrieve it")
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
		log.Success("Equipment confirmed deleted (received 404 Not Found)")
	})

	t.Run("Equipment With Invalid Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing validation with invalid equipment data")

		invalidEquipment := models.CreateEquipmentInput{
			Name:        "", // Invalid: empty name
			Description: "", // Invalid: empty description
			Cost:        -5, // Invalid: negative cost
			Weight:      0,  // Invalid: zero weight
		}

		log.Info("Invalid data: empty name and description, negative cost, zero weight")

		payload, err := json.Marshal(invalidEquipment)
		if !log.CheckNoError(err, "Marshal invalid equipment data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/equipment"
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

	log.Section("EQUIPMENT CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
