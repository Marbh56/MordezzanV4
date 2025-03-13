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

// TestMagicItemCRUDIntegration tests the CRUD operations for magic items
func TestMagicItemCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Magic Item CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	var magicItemID int64

	t.Run("Create Magic Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new magic item")

		charges := 10
		magicItemData := models.CreateMagicItemInput{
			Name:        "Staff of Power",
			ItemType:    "Staff",
			Description: "A powerful staff that can cast various spells",
			Charges:     &charges,
			Cost:        5000,
			Weight:      4,
		}

		log.Info("Magic Item Name: %s", magicItemData.Name)
		log.Info("Item Type: %s, Charges: %d",
			magicItemData.ItemType, *magicItemData.Charges)

		payload, err := json.Marshal(magicItemData)
		if !log.CheckNoError(err, "Marshal magic item data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/magic-items"
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

		var createdMagicItem models.MagicItem
		if err := json.NewDecoder(resp.Body).Decode(&createdMagicItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		magicItemID = createdMagicItem.ID
		log.Success("Magic item created with ID: %d", magicItemID)

		if createdMagicItem.Name != magicItemData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", magicItemData.Name, createdMagicItem.Name)
			t.Errorf("Expected magic item name %s, got %s", magicItemData.Name, createdMagicItem.Name)
		}

		if createdMagicItem.ItemType != magicItemData.ItemType {
			log.Error("Item type mismatch. Expected: %s, Got: %s", magicItemData.ItemType, createdMagicItem.ItemType)
			t.Errorf("Expected item type %s, got %s", magicItemData.ItemType, createdMagicItem.ItemType)
		}

		if *createdMagicItem.Charges != *magicItemData.Charges {
			log.Error("Charges mismatch. Expected: %d, Got: %d", *magicItemData.Charges, *createdMagicItem.Charges)
			t.Errorf("Expected charges %d, got %d", *magicItemData.Charges, *createdMagicItem.Charges)
		}

		log.Success("Magic item validation passed")
	})

	log.Separator()

	t.Run("Get Magic Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created magic item")
		log.Info("Magic Item ID: %d", magicItemID)

		endpoint := fmt.Sprintf("%s/magic-items/%d", server.URL, magicItemID)
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

		var magicItem models.MagicItem
		if err := json.NewDecoder(resp.Body).Decode(&magicItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received magic item: ID=%d, Name=%s, Type=%s",
			magicItem.ID, magicItem.Name, magicItem.ItemType)

		if magicItem.ID != magicItemID {
			log.Error("Magic item ID mismatch. Expected: %d, Got: %d", magicItemID, magicItem.ID)
			t.Errorf("Expected magic item ID %d, got %d", magicItemID, magicItem.ID)
		}

		if magicItem.Name != "Staff of Power" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Staff of Power", magicItem.Name)
			t.Errorf("Expected magic item name 'Staff of Power', got '%s'", magicItem.Name)
		}

		log.Success("Magic item data validation passed")
	})

	log.Separator()

	t.Run("List Magic Items", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of magic items")

		endpoint := server.URL + "/magic-items"
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

		var magicItems []*models.MagicItem
		if err := json.NewDecoder(resp.Body).Decode(&magicItems); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d magic items in total", len(magicItems))

		found := false
		for _, a := range magicItems {
			if a.ID == magicItemID {
				found = true
				log.Info("Found our test magic item: ID=%d, Name=%s", a.ID, a.Name)
				break
			}
		}

		if !found {
			log.Error("Magic item with ID %d not found in magic items list", magicItemID)
			t.Errorf("Magic item with ID %d not found in magic items list", magicItemID)
		} else {
			log.Success("Magic item found in the list")
		}
	})

	log.Separator()

	t.Run("List Magic Items By Type", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of magic items by type")

		endpoint := server.URL + "/magic-items?type=Staff"
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

		var staffItems []*models.MagicItem
		if err := json.NewDecoder(resp.Body).Decode(&staffItems); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d staff items in total", len(staffItems))

		found := false
		for _, a := range staffItems {
			if a.ID == magicItemID {
				found = true
				log.Info("Found our test staff item: ID=%d, Name=%s", a.ID, a.Name)

				if a.ItemType != "Staff" {
					log.Error("Item type mismatch. Expected: %s, Got: %s", "Staff", a.ItemType)
					t.Errorf("Expected item type %s, got %s", "Staff", a.ItemType)
				}
				break
			}
		}

		if !found {
			log.Error("Magic item with ID %d not found in staff items list", magicItemID)
			t.Errorf("Magic item with ID %d not found in staff items list", magicItemID)
		} else {
			log.Success("Staff item found in the filtered list")
		}
	})

	log.Separator()

	t.Run("Update Magic Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating magic item")
		log.Info("Magic Item ID: %d", magicItemID)

		charges := 15
		updateData := models.UpdateMagicItemInput{
			Name:        "Enhanced Staff of Power",
			ItemType:    "Staff",
			Description: "An enhanced staff with even more power",
			Charges:     &charges,
			Cost:        7500,
			Weight:      5,
		}

		log.Info("New magic item name: %s", updateData.Name)
		log.Info("New charges: %d", *updateData.Charges)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/magic-items/%d", server.URL, magicItemID)
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

		var updatedMagicItem models.MagicItem
		if err := json.NewDecoder(resp.Body).Decode(&updatedMagicItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated magic item data: ID=%d, Name=%s, Type=%s, Charges=%d",
			updatedMagicItem.ID, updatedMagicItem.Name, updatedMagicItem.ItemType, *updatedMagicItem.Charges)

		if updatedMagicItem.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s",
				updateData.Name, updatedMagicItem.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedMagicItem.Name)
		}

		if *updatedMagicItem.Charges != *updateData.Charges {
			log.Error("Charges mismatch. Expected: %d, Got: %d",
				*updateData.Charges, *updatedMagicItem.Charges)
			t.Errorf("Expected charges %d, got %d",
				*updateData.Charges, *updatedMagicItem.Charges)
		}

		log.Success("Magic item update validation passed")
	})

	log.Separator()

	t.Run("Delete Magic Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting magic item")
		log.Info("Magic Item ID: %d", magicItemID)

		endpoint := fmt.Sprintf("%s/magic-items/%d", server.URL, magicItemID)
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

		log.Info("Verifying magic item deletion by attempting to retrieve it")
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
		log.Success("Magic item confirmed deleted (received 404 Not Found)")
	})

	log.Section("MAGIC ITEM CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
