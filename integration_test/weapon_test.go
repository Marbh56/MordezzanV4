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

func TestWeaponCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Weapon CRUD Integration Test")
	log.Step("Setting up test environment")
	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()
	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()
	log.Success("Test server started at %s", server.URL)

	var weaponID int64

	t.Run("Create Weapon", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new weapon")
		weaponData := models.CreateWeaponInput{
			Name:            "Longsword +1",
			Category:        "Melee",
			WeaponClass:     2,
			Cost:            150,
			Weight:          3,
			Damage:          "1d8",
			DamageTwoHanded: "1d10",
			Properties:      "Versatile, Magic",
		}
		log.Info("Weapon Name: %s", weaponData.Name)
		log.Info("Weapon Category: %s, Class: %d, Damage: %s",
			weaponData.Category, weaponData.WeaponClass, weaponData.Damage)

		payload, err := json.Marshal(weaponData)
		if !log.CheckNoError(err, "Marshal weapon data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/weapons"
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

		var createdWeapon models.Weapon
		if err := json.NewDecoder(resp.Body).Decode(&createdWeapon); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}
		weaponID = createdWeapon.ID
		log.Success("Weapon created with ID: %d", weaponID)

		if createdWeapon.Name != weaponData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", weaponData.Name, createdWeapon.Name)
			t.Errorf("Expected weapon name %s, got %s", weaponData.Name, createdWeapon.Name)
		}
		if createdWeapon.WeaponClass != weaponData.WeaponClass {
			log.Error("Weapon class mismatch. Expected: %d, Got: %d", weaponData.WeaponClass, createdWeapon.WeaponClass)
			t.Errorf("Expected weapon class %d, got %d", weaponData.WeaponClass, createdWeapon.WeaponClass)
		}
		if createdWeapon.Properties != weaponData.Properties {
			log.Error("Properties mismatch. Expected: %s, Got: %s", weaponData.Properties, createdWeapon.Properties)
			t.Errorf("Expected properties %s, got %s", weaponData.Properties, createdWeapon.Properties)
		}
		log.Success("Weapon validation passed")
	})

	log.Separator()

	t.Run("Get Weapon", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created weapon")
		log.Info("Weapon ID: %d", weaponID)

		endpoint := fmt.Sprintf("%s/weapons/%d", server.URL, weaponID)
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

		var weapon models.Weapon
		if err := json.NewDecoder(resp.Body).Decode(&weapon); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}
		log.Info("Received weapon: ID=%d, Name=%s, Category=%s",
			weapon.ID, weapon.Name, weapon.Category)

		if weapon.ID != weaponID {
			log.Error("Weapon ID mismatch. Expected: %d, Got: %d", weaponID, weapon.ID)
			t.Errorf("Expected weapon ID %d, got %d", weaponID, weapon.ID)
		}
		if weapon.Name != "Longsword +1" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Longsword +1", weapon.Name)
			t.Errorf("Expected weapon name 'Longsword +1', got '%s'", weapon.Name)
		}
		log.Success("Weapon data validation passed")
	})

	log.Separator()

	t.Run("List Weapons", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of weapons")

		endpoint := server.URL + "/weapons"
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

		var weapons []*models.Weapon
		if err := json.NewDecoder(resp.Body).Decode(&weapons); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}
		log.Info("Received %d weapons in total", len(weapons))

		found := false
		for _, w := range weapons {
			if w.ID == weaponID {
				found = true
				log.Info("Found our test weapon: ID=%d, Name=%s", w.ID, w.Name)
				break
			}
		}
		if !found {
			log.Error("Weapon with ID %d not found in weapons list", weaponID)
			t.Errorf("Weapon with ID %d not found in weapons list", weaponID)
		} else {
			log.Success("Weapon found in the list")
		}
	})

	log.Separator()

	t.Run("Update Weapon", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating weapon")
		log.Info("Weapon ID: %d", weaponID)

		updateData := models.UpdateWeaponInput{
			Name:            "Longsword +2",
			Category:        "Melee",
			WeaponClass:     2,
			Cost:            300,
			Weight:          3,
			Damage:          "1d8+2",
			DamageTwoHanded: "1d10+2",
			Properties:      "Versatile, Magic, Enchanted",
		}
		log.Info("New weapon name: %s", updateData.Name)
		log.Info("New damage: %s, Properties: %s", updateData.Damage, updateData.Properties)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/weapons/%d", server.URL, weaponID)
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

		var updatedWeapon models.Weapon
		if err := json.NewDecoder(resp.Body).Decode(&updatedWeapon); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}
		log.Info("Updated weapon data: ID=%d, Name=%s, Damage=%s",
			updatedWeapon.ID, updatedWeapon.Name, updatedWeapon.Damage)

		if updatedWeapon.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedWeapon.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedWeapon.Name)
		}
		if updatedWeapon.Damage != updateData.Damage {
			log.Error("Damage mismatch. Expected: %s, Got: %s",
				updateData.Damage, updatedWeapon.Damage)
			t.Errorf("Expected damage %s, got %s",
				updateData.Damage, updatedWeapon.Damage)
		}
		if updatedWeapon.Properties != updateData.Properties {
			log.Error("Properties mismatch. Expected: %s, Got: %s",
				updateData.Properties, updatedWeapon.Properties)
			t.Errorf("Expected properties %s, got %s",
				updateData.Properties, updatedWeapon.Properties)
		}
		log.Success("Weapon update validation passed")
	})

	log.Separator()

	t.Run("Delete Weapon", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting weapon")
		log.Info("Weapon ID: %d", weaponID)

		endpoint := fmt.Sprintf("%s/weapons/%d", server.URL, weaponID)
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

		log.Info("Verifying weapon deletion by attempting to retrieve it")
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
		log.Success("Weapon confirmed deleted (received 404 Not Found)")
	})

	t.Run("Weapon With Invalid Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing validation with invalid weapon data")

		invalidWeapon := models.CreateWeaponInput{
			Name:        "", // Invalid: empty name
			Category:    "Melee",
			WeaponClass: -1,  // Invalid: negative weapon class
			Cost:        -50, // Invalid: negative cost
			Weight:      0,   // Invalid: zero weight
			Damage:      "",  // Invalid: empty damage
		}
		log.Info("Invalid data: empty name, negative weapon class, negative cost, zero weight, empty damage")

		payload, err := json.Marshal(invalidWeapon)
		if !log.CheckNoError(err, "Marshal invalid weapon data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/weapons"
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

	log.Section("WEAPON CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}

// Helper function to create a test weapon
func createTestWeapon(t *testing.T, server *httptest.Server) int64 {
	weaponData := models.CreateWeaponInput{
		Name:        "Test Sword",
		Category:    "Melee",
		WeaponClass: 1,
		Cost:        15,
		Weight:      2,
		Damage:      "1d6",
	}

	payload, err := json.Marshal(weaponData)
	if err != nil {
		t.Fatalf("Failed to marshal weapon data: %v", err)
	}

	req, err := http.NewRequest("POST", server.URL+"/weapons", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var createdWeapon models.Weapon
	if err := json.NewDecoder(resp.Body).Decode(&createdWeapon); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return createdWeapon.ID
}
