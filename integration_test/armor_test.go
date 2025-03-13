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

// TestArmorCRUDIntegration tests the CRUD operations for armors
func TestArmorCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Armor CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var armorID int64

	t.Run("Create Armor", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new armor")

		armorData := models.CreateArmorInput{
			Name:            "Test Plate Mail",
			ArmorType:       "Plate Mail",
			AC:              3,
			Cost:            350,
			DamageReduction: 2,
			Weight:          40,
			WeightClass:     "Heavy",
			MovementRate:    20,
		}

		log.Info("Armor Name: %s", armorData.Name)
		log.Info("Armor Type: %s, AC: %d, DR: %d",
			armorData.ArmorType, armorData.AC, armorData.DamageReduction)

		payload, err := json.Marshal(armorData)
		if !log.CheckNoError(err, "Marshal armor data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/armors"
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

		var createdArmor models.Armor
		if err := json.NewDecoder(resp.Body).Decode(&createdArmor); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		armorID = createdArmor.ID
		log.Success("Armor created with ID: %d", armorID)

		if createdArmor.Name != armorData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", armorData.Name, createdArmor.Name)
			t.Errorf("Expected armor name %s, got %s", armorData.Name, createdArmor.Name)
		}

		if createdArmor.AC != armorData.AC {
			log.Error("AC mismatch. Expected: %d, Got: %d", armorData.AC, createdArmor.AC)
			t.Errorf("Expected AC %d, got %d", armorData.AC, createdArmor.AC)
		}

		if createdArmor.WeightClass != armorData.WeightClass {
			log.Error("Weight class mismatch. Expected: %s, Got: %s",
				armorData.WeightClass, createdArmor.WeightClass)
			t.Errorf("Expected weight class %s, got %s",
				armorData.WeightClass, createdArmor.WeightClass)
		}

		log.Success("Armor validation passed")
	})

	log.Separator()

	t.Run("Get Armor", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created armor")
		log.Info("Armor ID: %d", armorID)

		endpoint := fmt.Sprintf("%s/armors/%d", server.URL, armorID)
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

		var armor models.Armor
		if err := json.NewDecoder(resp.Body).Decode(&armor); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received armor: ID=%d, Name=%s, AC=%d",
			armor.ID, armor.Name, armor.AC)

		if armor.ID != armorID {
			log.Error("Armor ID mismatch. Expected: %d, Got: %d", armorID, armor.ID)
			t.Errorf("Expected armor ID %d, got %d", armorID, armor.ID)
		}

		if armor.Name != "Test Plate Mail" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Test Plate Mail", armor.Name)
			t.Errorf("Expected armor name 'Test Plate Mail', got '%s'", armor.Name)
		}

		log.Success("Armor data validation passed")
	})

	log.Separator()

	t.Run("List Armors", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of armors")

		endpoint := server.URL + "/armors"
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

		var armors []*models.Armor
		if err := json.NewDecoder(resp.Body).Decode(&armors); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d armors in total", len(armors))

		found := false
		for _, a := range armors {
			if a.ID == armorID {
				found = true
				log.Info("Found our test armor: ID=%d, Name=%s", a.ID, a.Name)
				break
			}
		}

		if !found {
			log.Error("Armor with ID %d not found in armors list", armorID)
			t.Errorf("Armor with ID %d not found in armors list", armorID)
		} else {
			log.Success("Armor found in the list")
		}
	})

	log.Separator()

	t.Run("Update Armor", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating armor")
		log.Info("Armor ID: %d", armorID)

		updateData := models.UpdateArmorInput{
			Name:            "Enhanced Plate Mail",
			ArmorType:       "Plate Mail",
			AC:              2,
			Cost:            500,
			DamageReduction: 3,
			Weight:          45,
			WeightClass:     "Heavy",
			MovementRate:    15,
		}

		log.Info("New armor name: %s", updateData.Name)
		log.Info("New AC: %d, DR: %d", updateData.AC, updateData.DamageReduction)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/armors/%d", server.URL, armorID)
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

		var updatedArmor models.Armor
		if err := json.NewDecoder(resp.Body).Decode(&updatedArmor); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated armor data: ID=%d, Name=%s, AC=%d, DR=%d",
			updatedArmor.ID, updatedArmor.Name, updatedArmor.AC, updatedArmor.DamageReduction)

		if updatedArmor.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedArmor.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedArmor.Name)
		}

		if updatedArmor.AC != updateData.AC {
			log.Error("AC mismatch. Expected: %d, Got: %d",
				updateData.AC, updatedArmor.AC)
			t.Errorf("Expected AC %d, got %d",
				updateData.AC, updatedArmor.AC)
		}

		if updatedArmor.DamageReduction != updateData.DamageReduction {
			log.Error("DR mismatch. Expected: %d, Got: %d",
				updateData.DamageReduction, updatedArmor.DamageReduction)
			t.Errorf("Expected DR %d, got %d",
				updateData.DamageReduction, updatedArmor.DamageReduction)
		}

		log.Success("Armor update validation passed")
	})

	log.Separator()

	t.Run("Delete Armor", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting armor")
		log.Info("Armor ID: %d", armorID)

		endpoint := fmt.Sprintf("%s/armors/%d", server.URL, armorID)
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

		log.Info("Verifying armor deletion by attempting to retrieve it")
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
		log.Success("Armor confirmed deleted (received 404 Not Found)")
	})

	log.Section("ARMOR CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
