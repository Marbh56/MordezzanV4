// controllers/spell_casting_controller.go

package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/services"
)

// SpellCastingController handles HTTP requests related to character spells
type SpellCastingController struct {
	spellService *services.SpellService
}

// NewSpellCastingController creates a new spell casting controller
func NewSpellCastingController(spellService *services.SpellService) *SpellCastingController {
	return &SpellCastingController{
		spellService: spellService,
	}
}

// GetCharacterSpellsInfo retrieves all spell-related information for a character
func (c *SpellCastingController) GetCharacterSpellsInfo(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellsInfo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// AddKnownSpell adds a spell to a character's known spells
func (c *SpellCastingController) AddKnownSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	var input models.AddKnownSpellInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	// Set character ID from URL param
	input.CharacterID = characterID

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	// Learn the spell using the level-up functionality to respect limits
	err = c.spellService.LearnSpellOnLevelUp(r.Context(), characterID, input.SpellID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return updated spell info
	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(spellsInfo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// RemoveKnownSpell removes a spell from a character's known spells
func (c *SpellCastingController) RemoveKnownSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	spellID, err := strconv.ParseInt(chi.URLParam(r, "spellId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	err = c.spellService.RemoveKnownSpell(r.Context(), characterID, spellID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PrepareSpell prepares a spell for a character
func (c *SpellCastingController) PrepareSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	var input models.PrepareSpellInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	// Set character ID from URL param
	input.CharacterID = characterID

	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			validationErrors := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
		apperrors.HandleError(w, err)
		return
	}

	id, err := c.spellService.PrepareSpell(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return updated spell info
	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          id,
		"spells_info": spellsInfo,
	}); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// UnprepareSpell removes a prepared spell
func (c *SpellCastingController) UnprepareSpell(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	spellID, err := strconv.ParseInt(chi.URLParam(r, "spellId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	err = c.spellService.UnprepareSpell(r.Context(), characterID, spellID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return updated spell info
	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellsInfo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// ClearPreparedSpells removes all prepared spells for a character
func (c *SpellCastingController) ClearPreparedSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	err = c.spellService.ClearPreparedSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PrepareAllSpells prepares all spells a character can prepare
func (c *SpellCastingController) PrepareAllSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	err = c.spellService.PrepareAllSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return updated spell info
	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellsInfo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// GetSpellsLearnableOnLevelUp gets spells a character can learn when leveling up
func (c *SpellCastingController) GetSpellsLearnableOnLevelUp(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	newLevelStr := r.URL.Query().Get("new_level")
	newLevel, err := strconv.Atoi(newLevelStr)
	if err != nil || newLevel <= 0 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid or missing new_level parameter"))
		return
	}

	spells, err := c.spellService.GetSpellsLearnableOnLevelUp(r.Context(), characterID, newLevel)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spells); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// AddInitialSpellsForNewCharacter adds the starting spells for a new character
func (c *SpellCastingController) AddInitialSpellsForNewCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	err = c.spellService.AddInitialSpellsForNewCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return updated spell info
	spellsInfo, err := c.spellService.GetCharacterSpellsInfo(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellsInfo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}
