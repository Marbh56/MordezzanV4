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

// TestSpellCRUDIntegration tests the CRUD operations for spells
func TestSpellCRUDIntegration(t *testing.T) {
	log := NewTestLogger(t)
	log.Section("Spell CRUD Integration Test")
	log.Step("Setting up test environment")

	testApp, cleanup := setupTestAppWithFullSchema(t)
	defer cleanup()

	handler := testApp.SetupRoutes()
	server := httptest.NewServer(handler)
	defer server.Close()

	log.Success("Test server started at %s", server.URL)

	log.Step("Creating test user")
	userID := createTestUser(t, server)
	log.Success("Test user created with ID: %d", userID)

	log.Step("Creating test character")
	characterID := createTestCharacter(t, server, userID)
	log.Success("Test character created with ID: %d", characterID)

	var spellID int64

	t.Run("Create Spell", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Creating new spell")

		spellData := models.CreateSpellInput{
			CharacterID:  characterID,
			Name:         "Magic Missile",
			MagLevel:     1,
			CryLevel:     0,
			IllLevel:     0,
			NecLevel:     0,
			PyrLevel:     0,
			WchLevel:     0,
			ClrLevel:     0,
			DrdLevel:     0,
			Range:        "120 feet",
			Duration:     "Instantaneous",
			AreaOfEffect: "",
			Components:   "V, S",
			Description:  "You create three glowing darts of magical force. Each dart hits a creature of your choice that you can see within range. A dart deals 1d4+1 force damage to its target.",
		}

		log.Info("Spell Name: %s", spellData.Name)
		log.Info("Spell School Levels: MAG=%d, CRY=%d, ILL=%d, NEC=%d, PYR=%d, WCH=%d, CLR=%d, DRD=%d",
			spellData.MagLevel, spellData.CryLevel, spellData.IllLevel, spellData.NecLevel,
			spellData.PyrLevel, spellData.WchLevel, spellData.ClrLevel, spellData.DrdLevel)

		payload, err := json.Marshal(spellData)
		if !log.CheckNoError(err, "Marshal spell data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/spells"
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

		var createdSpell models.Spell
		if err := json.NewDecoder(resp.Body).Decode(&createdSpell); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		spellID = createdSpell.ID
		log.Success("Spell created with ID: %d", spellID)

		if createdSpell.Name != spellData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", spellData.Name, createdSpell.Name)
			t.Errorf("Expected spell name %s, got %s", spellData.Name, createdSpell.Name)
		}

		if createdSpell.MagLevel != spellData.MagLevel {
			log.Error("Magician level mismatch. Expected: %d, Got: %d", spellData.MagLevel, createdSpell.MagLevel)
			t.Errorf("Expected magician level %d, got %d", spellData.MagLevel, createdSpell.MagLevel)
		}

		if createdSpell.CharacterID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %d", characterID, createdSpell.CharacterID)
			t.Errorf("Expected character ID %d, got %d", characterID, createdSpell.CharacterID)
		}

		log.Success("Spell validation passed")
	})

	log.Separator()

	t.Run("Get Spell", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving created spell")
		log.Info("Spell ID: %d", spellID)

		endpoint := fmt.Sprintf("%s/spells/%d", server.URL, spellID)
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

		var spell models.Spell
		if err := json.NewDecoder(resp.Body).Decode(&spell); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received spell: ID=%d, Name=%s, Character ID=%d",
			spell.ID, spell.Name, spell.CharacterID)

		if spell.ID != spellID {
			log.Error("Spell ID mismatch. Expected: %d, Got: %d", spellID, spell.ID)
			t.Errorf("Expected spell ID %d, got %d", spellID, spell.ID)
		}

		if spell.Name != "Magic Missile" {
			log.Error("Name mismatch. Expected: %s, Got: %s", "Magic Missile", spell.Name)
			t.Errorf("Expected spell name 'Magic Missile', got '%s'", spell.Name)
		}

		if spell.CharacterID != characterID {
			log.Error("Character ID mismatch. Expected: %d, Got: %d", characterID, spell.CharacterID)
			t.Errorf("Expected character ID %d, got %d", characterID, spell.CharacterID)
		}

		log.Success("Spell data validation passed")
	})

	log.Separator()

	t.Run("Get Spells By Character", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Retrieving spells for character")
		log.Info("Character ID: %d", characterID)

		endpoint := fmt.Sprintf("%s/characters/%d/spells", server.URL, characterID)
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
			log.Error("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}
		log.Success("Received correct status code: %d", resp.StatusCode)

		var spells []*models.Spell
		if err := json.NewDecoder(resp.Body).Decode(&spells); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Received %d spells for character", len(spells))

		if len(spells) < 1 {
			log.Error("Expected at least 1 spell, got %d", len(spells))
			t.Fatalf("Expected at least 1 spell, got %d", len(spells))
		}

		found := false
		for _, s := range spells {
			if s.ID == spellID {
				found = true
				log.Info("Found our test spell: ID=%d, Name=%s", s.ID, s.Name)
				if s.Name != "Magic Missile" {
					log.Error("Name mismatch. Expected: %s, Got: %s", "Magic Missile", s.Name)
					t.Errorf("Expected spell name 'Magic Missile', got '%s'", s.Name)
				}
				break
			}
		}

		if !found {
			log.Error("Spell with ID %d not found in spells list", spellID)
			t.Errorf("Spell with ID %d not found in spells list", spellID)
		} else {
			log.Success("Spell found in the list")
		}
	})

	log.Separator()

	t.Run("Update Spell", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Updating spell")
		log.Info("Spell ID: %d", spellID)

		updateData := models.UpdateSpellInput{
			Name:         "Greater Magic Missile",
			MagLevel:     3,
			CryLevel:     0,
			IllLevel:     0,
			NecLevel:     0,
			PyrLevel:     0,
			WchLevel:     0,
			ClrLevel:     0,
			DrdLevel:     0,
			Range:        "150 feet",
			Duration:     "Instantaneous",
			AreaOfEffect: "",
			Components:   "V, S, M (a small piece of quartz)",
			Description:  "You create five glowing darts of magical force. Each dart hits a creature of your choice that you can see within range. A dart deals 1d6+1 force damage to its target.",
		}

		log.Info("New spell name: %s", updateData.Name)
		log.Info("New spell level: %d", updateData.MagLevel)

		payload, err := json.Marshal(updateData)
		if !log.CheckNoError(err, "Marshal update data") {
			t.Fatal("Test failed")
		}

		endpoint := fmt.Sprintf("%s/spells/%d", server.URL, spellID)
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

		var updatedSpell models.Spell
		if err := json.NewDecoder(resp.Body).Decode(&updatedSpell); err != nil {
			log.Error("Failed to decode response: %v", err)
			t.Fatalf("Failed to decode response: %v", err)
		}

		log.Info("Updated spell data: ID=%d, Name=%s", updatedSpell.ID, updatedSpell.Name)

		if updatedSpell.Name != updateData.Name {
			log.Error("Name mismatch. Expected: %s, Got: %s", updateData.Name, updatedSpell.Name)
			t.Errorf("Expected name %s, got %s", updateData.Name, updatedSpell.Name)
		}

		if updatedSpell.MagLevel != updateData.MagLevel {
			log.Error("Level mismatch. Expected: %d, Got: %d", updateData.MagLevel, updatedSpell.MagLevel)
			t.Errorf("Expected level %d, got %d", updateData.MagLevel, updatedSpell.MagLevel)
		}

		if updatedSpell.Components != updateData.Components {
			log.Error("Components mismatch. Expected: %s, Got: %s", updateData.Components, updatedSpell.Components)
			t.Errorf("Expected components %s, got %s", updateData.Components, updatedSpell.Components)
		}

		log.Success("Spell update validation passed")
	})

	log.Separator()

	t.Run("Delete Spell", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Deleting spell")
		log.Info("Spell ID: %d", spellID)

		endpoint := fmt.Sprintf("%s/spells/%d", server.URL, spellID)
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

		log.Info("Verifying spell deletion by attempting to retrieve it")
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
		log.Success("Spell confirmed deleted (received 404 Not Found)")
	})

	t.Run("Spell With Invalid Data", func(t *testing.T) {
		log := NewTestLogger(t)
		log.Step("Testing validation with invalid spell data")

		invalidSpell := models.CreateSpellInput{
			CharacterID: characterID,
			Name:        "", // Invalid: empty name
			MagLevel:    -1, // Invalid: negative level
			CryLevel:    0,
			IllLevel:    0,
			NecLevel:    0,
			PyrLevel:    0,
			WchLevel:    0,
			ClrLevel:    0,
			DrdLevel:    0,
			Range:       "", // Invalid: empty range
			Duration:    "", // Invalid: empty duration
			Description: "Test spell description",
		}

		log.Info("Invalid data: empty name, negative level, empty range and duration")

		payload, err := json.Marshal(invalidSpell)
		if !log.CheckNoError(err, "Marshal invalid spell data") {
			t.Fatal("Test failed")
		}

		endpoint := server.URL + "/spells"
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
			log.Error("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
			t.Fatalf("Expected status %d for invalid data, got %d", http.StatusBadRequest, resp.StatusCode)
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
			log.Error("Expected error status %d, got %d", http.StatusBadRequest, errorResp.Status)
			t.Errorf("Expected error status %d, got %d", http.StatusBadRequest, errorResp.Status)
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

	log.Section("SPELL CRUD INTEGRATION TEST COMPLETED SUCCESSFULLY")
}
