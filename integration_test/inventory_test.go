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

// TestInventoryCRUDIntegration tests the CRUD operations for inventory
func TestInventoryCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Inventory CRUD Integration Test")
	log.Step("Setting up test environment")

	// Setup test environment with app, server, authenticated user and cleanup function
	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()
	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	// Create a test user with authentication
	testUser := CreateTestUserWithAuth(t, server)
	log.AuthInfo(testUser)

	// Create a character for the test user to own the inventory
	log.Step("Creating a character for the test user")
	characterData := models.CreateCharacterInput{
		UserID:       testUser.ID,
		Name:         "Aragorn",
		Class:        "Ranger",
		Level:        7,
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Wisdom:       12,
		Intelligence: 10,
		Charisma:     13,
		HitPoints:    56,
	}

	payload, err := json.Marshal(characterData)
	if !log.CheckNoError(err, "Marshal character data") {
		t.Fatal("Test failed")
	}

	req := AuthenticatedRequest(t, "POST", server.URL+"/characters", bytes.NewBuffer(payload), testUser)
	resp, err := http.DefaultClient.Do(req)
	if !log.CheckNoError(err, "Send character creation request") {
		t.Fatal("Test failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var character models.Character
	if err := json.NewDecoder(resp.Body).Decode(&character); err != nil {
		log.Error("Failed to decode character response: %v", err)
		t.Fatalf("Failed to decode character response: %v", err)
	}

	characterID := character.ID
	log.Success("Character created with ID: %d", characterID)

	// Create test items
	log.Step("Creating test equipment item")
	equipment := CreateTestEquipment(t, server, testUser)
	log.Success("Created test equipment with ID: %d", equipment.ID)

	log.Step("Creating test weapon item")
	weapon := CreateTestWeapon(t, server, testUser)
	log.Success("Created test weapon with ID: %d", weapon.ID)

	var inventoryID int64

	t.Run("Create Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new inventory for character")

		inventoryData := models.CreateInventoryInput{
			CharacterID: characterID,
			MaxWeight:   100.0,
		}

		payload, err := json.Marshal(inventoryData)
		if !log.CheckNoError(err, "Marshal inventory data") {
			t.Fatal("Test failed")
		}

		req := AuthenticatedRequest(t, "POST", server.URL+"/inventories", bytes.NewBuffer(payload), testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send inventory creation request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var createdInventory models.Inventory
		if err := json.NewDecoder(resp.Body).Decode(&createdInventory); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		inventoryID = createdInventory.ID
		log.Success("Inventory created with ID: %d", inventoryID)

		if createdInventory.CharacterID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %d", characterID, createdInventory.CharacterID)
			t.Errorf("Expected character ID %d, got %d", characterID, createdInventory.CharacterID)
		}

		if createdInventory.MaxWeight != inventoryData.MaxWeight {
			log.Error("Max weight mismatch. Expected: %.2f, Got: %.2f", inventoryData.MaxWeight, createdInventory.MaxWeight)
			t.Errorf("Expected max weight %.2f, got %.2f", inventoryData.MaxWeight, createdInventory.MaxWeight)
		}

		log.Success("Inventory validation passed")
	})

	if inventoryID <= 0 {
		log.Error("Cannot continue tests without valid inventory ID")
		t.Fatal("Cannot continue tests without valid inventory ID")
	}

	log.Separator()

	var equipmentItemID int64
	t.Run("Add Equipment Item to Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Adding equipment item to inventory")

		addItemData := models.AddItemInput{
			ItemType:   "equipment",
			ItemID:     equipment.ID,
			Quantity:   1,
			IsEquipped: false,
			Notes:      "Test equipment",
		}

		payload, err := json.Marshal(addItemData)
		if !log.CheckNoError(err, "Marshal add item data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/inventories/%d/items", server.URL, inventoryID)
		req := AuthenticatedRequest(t, "POST", endpoint, bytes.NewBuffer(payload), testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send add item request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var addedItem models.InventoryItem
		if err := json.NewDecoder(resp.Body).Decode(&addedItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		equipmentItemID = addedItem.ID
		log.Success("Item added to inventory with ID: %d", equipmentItemID)

		if addedItem.ItemType != "equipment" {
			log.Error("Item type mismatch. Expected: %s, Got: %s", "equipment", addedItem.ItemType)
			t.Errorf("Expected item type %s, got %s", "equipment", addedItem.ItemType)
		}

		if addedItem.ItemID != equipment.ID {
			log.Error("Item ID mismatch. Expected: %d, Got: %d", equipment.ID, addedItem.ItemID)
			t.Errorf("Expected item ID %d, got %d", equipment.ID, addedItem.ItemID)
		}

		log.Success("Item validation passed")
	})

	log.Separator()

	var weaponItemID int64
	t.Run("Add Weapon Item to Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Adding weapon item to inventory")

		addItemData := models.AddItemInput{
			ItemType:   "weapon",
			ItemID:     weapon.ID,
			Quantity:   1,
			IsEquipped: true,
			Notes:      "Test weapon",
		}

		payload, err := json.Marshal(addItemData)
		if !log.CheckNoError(err, "Marshal add item data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/inventories/%d/items", server.URL, inventoryID)
		req := AuthenticatedRequest(t, "POST", endpoint, bytes.NewBuffer(payload), testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send add item request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Error("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var addedItem models.InventoryItem
		if err := json.NewDecoder(resp.Body).Decode(&addedItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		weaponItemID = addedItem.ID
		log.Success("Weapon added to inventory with ID: %d", weaponItemID)

		if addedItem.ItemType != "weapon" {
			log.Error("Item type mismatch. Expected: %s, Got: %s", "weapon", addedItem.ItemType)
			t.Errorf("Expected item type %s, got %s", "weapon", addedItem.ItemType)
		}

		if addedItem.ItemID != weapon.ID {
			log.Error("Item ID mismatch. Expected: %d, Got: %d", weapon.ID, addedItem.ItemID)
			t.Errorf("Expected item ID %d, got %d", weapon.ID, addedItem.ItemID)
		}

		if !addedItem.IsEquipped {
			log.Error("Item should be equipped, but it's not")
			t.Error("Item should be equipped, but it's not")
		}

		log.Success("Weapon item validation passed")
	})

	log.Separator()

	t.Run("Get Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving inventory")
		log.Info("Inventory ID: %d", inventoryID)

		endpoint := fmt.Sprintf("%s/inventories/%d", server.URL, inventoryID)
		req := AuthenticatedRequest(t, "GET", endpoint, nil, testUser)
		req.Header.Set("Accept", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send get inventory request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var inventory models.Inventory
		if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received inventory data: ID=%d, CharacterID=%d, Items=%d",
			inventory.ID, inventory.CharacterID, len(inventory.Items))

		if inventory.ID != inventoryID {
			log.Error("Inventory ID mismatch. Expected: %d, Got: %d", inventoryID, inventory.ID)
			t.Errorf("Expected inventory ID %d, got %d", inventoryID, inventory.ID)
		}

		if len(inventory.Items) != 2 {
			log.Error("Expected 2 items in inventory, got %d", len(inventory.Items))
			t.Errorf("Expected 2 items in inventory, got %d", len(inventory.Items))
		}

		log.Success("Inventory data validation passed")
	})

	log.Separator()

	t.Run("Get Inventory Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving inventory item")
		log.Info("Item ID: %d", equipmentItemID)

		endpoint := fmt.Sprintf("%s/inventories/%d/items/%d", server.URL, inventoryID, equipmentItemID)
		req := AuthenticatedRequest(t, "GET", endpoint, nil, testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send get item request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var item struct {
			ID          int64       `json:"id"`
			InventoryID int64       `json:"inventory_id"`
			ItemType    string      `json:"item_type"`
			ItemID      int64       `json:"item_id"`
			ItemDetails interface{} `json:"item_details"`
			Quantity    int         `json:"quantity"`
			IsEquipped  bool        `json:"is_equipped"`
			Notes       string      `json:"notes"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received item data: ID=%d, Type=%s, ItemID=%d",
			item.ID, item.ItemType, item.ItemID)

		if item.ID != equipmentItemID {
			log.Error("Item ID mismatch. Expected: %d, Got: %d", equipmentItemID, item.ID)
			t.Errorf("Expected item ID %d, got %d", equipmentItemID, item.ID)
		}

		log.Success("Item data validation passed")
	})

	log.Separator()

	t.Run("Get Inventory By Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving inventory by character ID")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d/inventory", server.URL, characterID)
		req := AuthenticatedRequest(t, "GET", endpoint, nil, testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send get inventory by character request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var response struct {
			Inventory models.Inventory         `json:"inventory"`
			Items     []map[string]interface{} `json:"items"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received inventory data: ID=%d, CharacterID=%d",
			response.Inventory.ID, response.Inventory.CharacterID)
		log.Info("Received %d enriched items", len(response.Items))

		if response.Inventory.ID != inventoryID {
			log.Error("Inventory ID mismatch. Expected: %d, Got: %d", inventoryID, response.Inventory.ID)
			t.Errorf("Expected inventory ID %d, got %d", inventoryID, response.Inventory.ID)
		}

		if len(response.Items) != 2 {
			log.Error("Expected 2 enriched items, got %d", len(response.Items))
			t.Errorf("Expected 2 enriched items, got %d", len(response.Items))
		}

		log.Success("Inventory by character validation passed")
	})

	log.Separator()

	t.Run("Update Inventory Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating inventory item")
		log.Info("Item ID: %d", equipmentItemID)

		newQuantity := 5
		isEquipped := true
		newNotes := "Updated equipment notes"

		updateData := models.UpdateItemInput{
			Quantity:   &newQuantity,
			IsEquipped: &isEquipped,
			Notes:      &newNotes,
		}

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/inventories/%d/items/%d", server.URL, inventoryID, equipmentItemID)
		req := AuthenticatedRequest(t, "PUT", endpoint, bytes.NewBuffer(payload), testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send update item request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var updatedItem models.InventoryItem
		if err := json.NewDecoder(resp.Body).Decode(&updatedItem); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received updated item data: Quantity=%d, IsEquipped=%v",
			updatedItem.Quantity, updatedItem.IsEquipped)

		if updatedItem.Quantity != newQuantity {
			log.Error("Quantity mismatch. Expected: %d, Got: %d", newQuantity, updatedItem.Quantity)
			t.Errorf("Expected quantity %d, got %d", newQuantity, updatedItem.Quantity)
		}

		if !updatedItem.IsEquipped {
			log.Error("IsEquipped mismatch. Expected: %v, Got: %v", true, updatedItem.IsEquipped)
			t.Errorf("Expected isEquipped %v, got %v", true, updatedItem.IsEquipped)
		}

		if updatedItem.Notes != newNotes {
			log.Error("Notes mismatch. Expected: %s, Got: %s", newNotes, updatedItem.Notes)
			t.Errorf("Expected notes %s, got %s", newNotes, updatedItem.Notes)
		}

		log.Success("Item update validation passed")
	})

	log.Separator()

	t.Run("Update Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating inventory")
		log.Info("Inventory ID: %d", inventoryID)

		newMaxWeight := 150.0
		newCurrentWeight := 45.5

		updateData := models.UpdateInventoryInput{
			MaxWeight:     &newMaxWeight,
			CurrentWeight: &newCurrentWeight,
		}

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/inventories/%d", server.URL, inventoryID)
		req := AuthenticatedRequest(t, "PUT", endpoint, bytes.NewBuffer(payload), testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send update inventory request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var updatedInventory models.Inventory
		if err := json.NewDecoder(resp.Body).Decode(&updatedInventory); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received updated inventory data: MaxWeight=%.2f, CurrentWeight=%.2f",
			updatedInventory.MaxWeight, updatedInventory.CurrentWeight)

		if updatedInventory.MaxWeight != newMaxWeight {
			log.Error("MaxWeight mismatch. Expected: %.2f, Got: %.2f", newMaxWeight, updatedInventory.MaxWeight)
			t.Errorf("Expected maxWeight %.2f, got %.2f", newMaxWeight, updatedInventory.MaxWeight)
		}

		if updatedInventory.CurrentWeight != newCurrentWeight {
			log.Error("CurrentWeight mismatch. Expected: %.2f, Got: %.2f", newCurrentWeight, updatedInventory.CurrentWeight)
			t.Errorf("Expected currentWeight %.2f, got %.2f", newCurrentWeight, updatedInventory.CurrentWeight)
		}

		log.Success("Inventory update validation passed")
	})

	log.Separator()

	t.Run("Remove Inventory Item", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Removing inventory item")
		log.Info("Item ID: %d", weaponItemID)

		endpoint := fmt.Sprintf("%s/inventories/%d/items/%d", server.URL, inventoryID, weaponItemID)
		req := AuthenticatedRequest(t, "DELETE", endpoint, nil, testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send remove item request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			log.Error("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		// Verify item was removed by trying to fetch it
		verifyEndpoint := fmt.Sprintf("%s/inventories/%d/items/%d", server.URL, inventoryID, weaponItemID)
		verifyReq := AuthenticatedRequest(t, "GET", verifyEndpoint, nil, testUser)
		verifyResp, err := http.DefaultClient.Do(verifyReq)
		if !log.CheckNoError(err, "Send verification request") {
			t.Fatal("Test failed")
		}
		defer verifyResp.Body.Close()

		if verifyResp.StatusCode != http.StatusNotFound {
			log.Error("Expected status %d after deletion, got %d", http.StatusNotFound, verifyResp.StatusCode)
			t.Fatalf("Expected status %d after deletion, got %d", http.StatusNotFound, verifyResp.StatusCode)
		}
		log.Success("Item confirmed removed (received 404 Not Found)")
	})

	log.Separator()

	t.Run("Delete Inventory", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting inventory")
		log.Info("Inventory ID: %d", inventoryID)

		endpoint := fmt.Sprintf("%s/inventories/%d", server.URL, inventoryID)
		req := AuthenticatedRequest(t, "DELETE", endpoint, nil, testUser)
		resp, err := http.DefaultClient.Do(req)
		if !log.CheckNoError(err, "Send delete inventory request") {
			t.Fatal("Test failed")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			log.Error("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
			t.Fatalf("Expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		// Verify inventory was deleted by trying to fetch it
		verifyEndpoint := fmt.Sprintf("%s/inventories/%d", server.URL, inventoryID)
		verifyReq := AuthenticatedRequest(t, "GET", verifyEndpoint, nil, testUser)
		verifyReq.Header.Set("Accept", "application/json")
		verifyResp, err := http.DefaultClient.Do(verifyReq)
		if !log.CheckNoError(err, "Send verification request") {
			t.Fatal("Test failed")
		}
		defer verifyResp.Body.Close()

		if verifyResp.StatusCode != http.StatusNotFound {
			log.Error("Expected status %d after deletion, got %d", http.StatusNotFound, verifyResp.StatusCode)
			t.Fatalf("Expected status %d after deletion, got %d", http.StatusNotFound, verifyResp.StatusCode)
		}
		log.Success("Inventory confirmed deleted (received 404 Not Found)")
	})

	log.Section("INVENTORY CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
