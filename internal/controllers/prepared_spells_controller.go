package controllers

import (
	"encoding/json"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// PreparedSpellController handles HTTP requests related to spell preparation
type PreparedSpellController struct {
	spellPreparationService *services.SpellPreparationService
	tmpl                    *template.Template
}

// NewPreparedSpellController creates a new prepared spell controller
func NewPreparedSpellController(
	spellPreparationService *services.SpellPreparationService,
	tmpl *template.Template,
) *PreparedSpellController {
	return &PreparedSpellController{
		spellPreparationService: spellPreparationService,
		tmpl:                    tmpl,
	}
}

// GetPreparedSpells handles requests to get all prepared spells for a character
func (c *PreparedSpellController) GetPreparedSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	preparedSpells, err := c.spellPreparationService.GetPreparedSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(preparedSpells); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// GetAvailableSpellSlots handles requests to get available spell slots for a character
func (c *PreparedSpellController) GetAvailableSpellSlots(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	slots, err := c.spellPreparationService.GetAvailableSpellSlots(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(slots); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// PrepareSpell handles requests to prepare a spell
func (c *PreparedSpellController) PrepareSpell(w http.ResponseWriter, r *http.Request) {
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

	// Set character ID from URL parameter
	input.CharacterID = characterID

	if err := c.spellPreparationService.PrepareSpell(r.Context(), characterID, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get updated prepared spells and available slots for response
	preparedSpells, err := c.spellPreparationService.GetPreparedSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	slots, err := c.spellPreparationService.GetAvailableSpellSlots(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Return both the prepared spells and available slots
	response := struct {
		PreparedSpells []*models.PreparedSpell `json:"prepared_spells"`
		AvailableSlots *models.SpellSlots      `json:"available_slots"`
	}{
		PreparedSpells: preparedSpells,
		AvailableSlots: slots,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// UnprepareSpell handles requests to unprepare a spell
func (c *PreparedSpellController) UnprepareSpell(w http.ResponseWriter, r *http.Request) {
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

	if err := c.spellPreparationService.UnprepareSpell(r.Context(), characterID, spellID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get updated prepared spells and available slots for response
	preparedSpells, err := c.spellPreparationService.GetPreparedSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	slots, err := c.spellPreparationService.GetAvailableSpellSlots(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Return both the prepared spells and available slots
	response := struct {
		PreparedSpells []*models.PreparedSpell `json:"prepared_spells"`
		AvailableSlots *models.SpellSlots      `json:"available_slots"`
	}{
		PreparedSpells: preparedSpells,
		AvailableSlots: slots,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// ClearPreparedSpells handles requests to remove all prepared spells for a character
func (c *PreparedSpellController) ClearPreparedSpells(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	if err := c.spellPreparationService.ClearPreparedSpells(r.Context(), characterID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Return empty slots data
	slots, err := c.spellPreparationService.GetAvailableSpellSlots(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		PreparedSpells []interface{}      `json:"prepared_spells"`
		AvailableSlots *models.SpellSlots `json:"available_slots"`
	}{
		PreparedSpells: []interface{}{},
		AvailableSlots: slots,
	}); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// RenderSpellPrepPage renders the spell preparation page for a character
func (c *PreparedSpellController) RenderSpellPrepPage(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Get the prepared spells
	preparedSpells, err := c.spellPreparationService.GetPreparedSpells(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get available spell slots
	slots, err := c.spellPreparationService.GetAvailableSpellSlots(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Render the spell preparation template
	data := struct {
		CharacterID    int64                   `json:"character_id"`
		PreparedSpells []*models.PreparedSpell `json:"prepared_spells"`
		AvailableSlots *models.SpellSlots      `json:"available_slots"`
	}{
		CharacterID:    characterID,
		PreparedSpells: preparedSpells,
		AvailableSlots: slots,
	}

	err = c.tmpl.ExecuteTemplate(w, "spell_preparation.html", data)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}
