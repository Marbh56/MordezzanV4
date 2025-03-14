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

// TestSpellScrollCRUDIntegration tests the CRUD operations for spell scrolls
func TestSpellScrollCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Spell Scroll CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	// First, we need to create a user
	log.Step("Creating a test user")
	userID := createTestUser(t, server)
	log.Success("Created test user with ID %d", userID)

	// Next, create a character
	log.Step("Creating a test character")
	characterID := createTestCharacter(t, server, userID)
	log.Success("Created test character with ID %d", characterID)

	// Create a spell
	log.Step("Creating a test spell")
	spellData := models.CreateSpellInput{
		CharacterID:  characterID,
		Name:         "Fireball",
		MagLevel:     3,
		Range:        "150 feet",
		Duration:     "Instantaneous",
		AreaOfEffect: "20-foot radius",
		Components:   "V, S, M (a tiny ball of bat guano and sulfur)",
		Description:  "A bright streak flashes from your pointing finger to a point you choose within range then blossoms with a low roar into an explosion of flame.",
	}

	log.Info("Spell Name: %s, Level: %d", spellData.Name, spellData.MagLevel)

	spellPayload, err := json.Marshal(spellData)
	if !log.CheckNoError(err, "Marshal spell data") {
		t.Fatal("Test failed")
	}

	spellEndpoint := server.URL + "/spells"
	log.Info("Sending POST request to %s", spellEndpoint)

	spellReq, err := http.NewRequest("POST", spellEndpoint, bytes.NewBuffer(spellPayload))
	if !log.CheckNoError(err, "Create spell request") {
		t.Fatal("Test failed")
	}
	spellReq.Header.Set("Content-Type", "application/json")

	spellResp, err := http.DefaultClient.Do(spellReq)
	if !log.CheckNoError(err, "Send spell request") {
		t.Fatal("Test failed")
	}
	defer spellResp.Body.Close()

	if spellResp.StatusCode != http.StatusCreated {
		log.Error("Expected status %d, got %d for spell creation", http.StatusCreated, spellResp.StatusCode)
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, spellResp.StatusCode)
	}
	log.Success("Created spell successfully")

	var createdSpell models.Spell
	if err := json.NewDecoder(spellResp.Body).Decode(&createdSpell); err != nil {
		log.Error("Failed to decode response: %v", err)
		t.Fatalf("Failed to decode response: %v", err)
	}

	spellID := createdSpell.ID
	log.Success("Spell created with ID: %d", spellID)

	var spellScrollID int64

	t.Run("Create Spell Scroll", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new spell scroll")

		scrollData := models.CreateSpellScrollInput{
			SpellID:      spellID,
			CastingLevel: 3,
			Cost:         150.0,
			Weight:       1,
			Description:  "A handwritten scroll containing the Fireball spell, ready to be cast.",
		}

		log.Info("Spell Scroll for: %s", createdSpell.Name)
		log.Info("Casting Level: %d, Cost: %.1f gp, Weight: %d lb",
			scrollData.CastingLevel, scrollData.Cost, scrollData.Weight)

		payload, err := json.Marshal(scrollData)
		if !log.CheckNoError(err, "Marshal scroll data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/spell-scrolls"
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

		var createdScroll models.SpellScroll
		if err := json.NewDecoder(resp.Body).Decode(&createdScroll); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		spellScrollID = createdScroll.ID
		log.Success("Spell scroll created with ID: %d", spellScrollID)

		if createdScroll.SpellID != scrollData.SpellID {
			log.Error("Spell ID mismatch. Expected: %d, Got: %d", scrollData.SpellID, createdScroll.SpellID)
			t.Errorf("Expected spell ID %d, got %d", scrollData.SpellID, createdScroll.SpellID)
		}

		if createdScroll.CastingLevel != scrollData.CastingLevel {
			log.Error("Casting level mismatch. Expected: %d, Got: %d", scrollData.CastingLevel, createdScroll.CastingLevel)
			t.Errorf("Expected casting level %d, got %d", scrollData.CastingLevel, createdScroll.CastingLevel)
		}

		if createdScroll.Cost != scrollData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f", scrollData.Cost, createdScroll.Cost)
			t.Errorf("Expected cost %.1f, got %.1f", scrollData.Cost, createdScroll.Cost)
		}

		if createdScroll.Weight != scrollData.Weight {
			log.Error("Weight mismatch. Expected: %d, Got: %d", scrollData.Weight, createdScroll.Weight)
			t.Errorf("Expected weight %d, got %d", scrollData.Weight, createdScroll.Weight)
		}

		log.Success("Spell scroll validation passed")
	})

	log.Separator()

	t.Run("Get Spell Scroll", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created spell scroll")
		log.Info("Spell Scroll ID: %d", spellScrollID)

		endpoint := fmt.Sprintf("%s/spell-scrolls/%d", server.URL, spellScrollID)
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

		var spellScroll models.SpellScroll
		if err := json.NewDecoder(resp.Body).Decode(&spellScroll); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received spell scroll: ID=%d, Spell=%s, Level=%d",
			spellScroll.ID, spellScroll.SpellName, spellScroll.CastingLevel)

		if spellScroll.ID != spellScrollID {
			log.Error("Spell scroll ID mismatch. Expected: %d, Got: %d", spellScrollID, spellScroll.ID)
			t.Errorf("Expected spell scroll ID %d, got %d", spellScrollID, spellScroll.ID)
		}

		if spellScroll.SpellName != "Fireball" {
			log.Error("Spell name mismatch. Expected: %s, Got: %s", "Fireball", spellScroll.SpellName)
			t.Errorf("Expected spell name 'Fireball', got '%s'", spellScroll.SpellName)
		}

		log.Success("Spell scroll data validation passed")
	})

	log.Separator()

	t.Run("List Spell Scrolls", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving list of spell scrolls")

		endpoint := server.URL + "/spell-scrolls"
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

		var scrolls []*models.SpellScroll
		if err := json.NewDecoder(resp.Body).Decode(&scrolls); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d spell scrolls in total", len(scrolls))

		found := false
		for _, s := range scrolls {
			if s.ID == spellScrollID {
				found = true
				log.Info("Found our test spell scroll: ID=%d, Spell=%s", s.ID, s.SpellName)
				break
			}
		}

		if !found {
			log.Error("Spell scroll with ID %d not found in spell scrolls list", spellScrollID)
			t.Errorf("Spell scroll with ID %d not found in spell scrolls list", spellScrollID)
		} else {
			log.Success("Spell scroll found in the list")
		}
	})

	log.Separator()

	t.Run("Update Spell Scroll", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating spell scroll")
		log.Info("Spell Scroll ID: %d", spellScrollID)

		updateData := models.UpdateSpellScrollInput{
			SpellID:      spellID,
			CastingLevel: 5,
			Cost:         300.0,
			Weight:       1,
			Description:  "An enhanced Fireball scroll, prepared by a master mage to cast at a higher level.",
		}

		log.Info("New casting level: %d", updateData.CastingLevel)
		log.Info("New cost: %.1f gp, Weight: %d lb", updateData.Cost, updateData.Weight)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/spell-scrolls/%d", server.URL, spellScrollID)
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

		var updatedScroll models.SpellScroll
		if err := json.NewDecoder(resp.Body).Decode(&updatedScroll); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated spell scroll data: ID=%d, Level=%d, Cost=%.1f",
			updatedScroll.ID, updatedScroll.CastingLevel, updatedScroll.Cost)

		if updatedScroll.CastingLevel != updateData.CastingLevel {
			log.Error("Casting level mismatch. Expected: %d, Got: %d",
				updateData.CastingLevel, updatedScroll.CastingLevel)
			t.Errorf("Expected casting level %d, got %d",
				updateData.CastingLevel, updatedScroll.CastingLevel)
		}

		if updatedScroll.Cost != updateData.Cost {
			log.Error("Cost mismatch. Expected: %.1f, Got: %.1f",
				updateData.Cost, updatedScroll.Cost)
			t.Errorf("Expected cost %.1f, got %.1f",
				updateData.Cost, updatedScroll.Cost)
		}

		log.Success("Spell scroll update validation passed")
	})

	log.Separator()

	t.Run("Delete Spell Scroll", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting spell scroll")
		log.Info("Spell Scroll ID: %d", spellScrollID)

		endpoint := fmt.Sprintf("%s/spell-scrolls/%d", server.URL, spellScrollID)
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

		log.Info("Verifying spell scroll deletion by attempting to retrieve it")
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
		log.Success("Spell scroll confirmed deleted (received 404 Not Found)")
	})

	log.Section("SPELL SCROLL CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
