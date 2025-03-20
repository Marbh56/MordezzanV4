package controllers

import (
	"encoding/json"
	"errors"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type SpellbookController struct {
	spellbookRepo repositories.SpellbookRepository
	spellRepo     repositories.SpellRepository
	tmpl          *template.Template
}

func NewSpellbookController(spellbookRepo repositories.SpellbookRepository, spellRepo repositories.SpellRepository, tmpl *template.Template) *SpellbookController {
	return &SpellbookController{
		spellbookRepo: spellbookRepo,
		spellRepo:     spellRepo,
		tmpl:          tmpl,
	}
}

func (c *SpellbookController) GetSpellbook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	spellbook, err := c.spellbookRepo.GetSpellbook(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(spellbook); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "spellbook.html", spellbook); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) ListSpellbooks(w http.ResponseWriter, r *http.Request) {
	spellbooks, err := c.spellbookRepo.ListSpellbooks(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellbooks); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) CreateSpellbook(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSpellbookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

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

	id, err := c.spellbookRepo.CreateSpellbook(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	spellbook, err := c.spellbookRepo.GetSpellbook(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(spellbook); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) UpdateSpellbook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	var input models.UpdateSpellbookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

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

	if err := c.spellbookRepo.UpdateSpellbook(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedSpellbook, err := c.spellbookRepo.GetSpellbook(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSpellbook); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) DeleteSpellbook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	if err := c.spellbookRepo.DeleteSpellbook(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *SpellbookController) AddSpellToSpellbook(w http.ResponseWriter, r *http.Request) {
	spellbookID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	var input struct {
		SpellID        int64  `json:"spell_id"`
		CharacterClass string `json:"character_class"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body format"))
		return
	}

	if input.SpellID <= 0 {
		apperrors.HandleValidationErrors(w, map[string]string{
			"spell_id": "Spell ID must be positive",
		})
		return
	}

	// Validate the character class
	validClasses := []string{"Magician", "Cryo-mancer", "Illusionist", "Necromancer",
		"Pyromancer", "Witch", "Cleric", "Druid"}

	isValidClass := false
	for _, class := range validClasses {
		if input.CharacterClass == class {
			isValidClass = true
			break
		}
	}

	if !isValidClass {
		apperrors.HandleValidationErrors(w, map[string]string{
			"character_class": "Invalid character class. Must be one of: " + strings.Join(validClasses, ", "),
		})
		return
	}

	// Verify that the spell exists
	_, err = c.spellRepo.GetSpell(r.Context(), input.SpellID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Add the spell to the spellbook
	if err := c.spellbookRepo.AddSpellToSpellbook(r.Context(), spellbookID, input.SpellID, input.CharacterClass); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get the updated spellbook
	updatedSpellbook, err := c.spellbookRepo.GetSpellbook(r.Context(), spellbookID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSpellbook); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) RemoveSpellFromSpellbook(w http.ResponseWriter, r *http.Request) {
	spellbookID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	spellID, err := strconv.ParseInt(chi.URLParam(r, "spellId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	// Remove the spell from the spellbook
	if err := c.spellbookRepo.RemoveSpellFromSpellbook(r.Context(), spellbookID, spellID); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get the updated spellbook
	updatedSpellbook, err := c.spellbookRepo.GetSpellbook(r.Context(), spellbookID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSpellbook); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellbookController) GetSpellsInSpellbook(w http.ResponseWriter, r *http.Request) {
	spellbookID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spellbook ID format"))
		return
	}

	// Check if spellbook exists
	_, err = c.spellbookRepo.GetSpellbook(r.Context(), spellbookID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get the spell IDs
	spellIDs, err := c.spellbookRepo.GetSpellsInSpellbook(r.Context(), spellbookID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get the full spell objects
	spells := make([]*models.Spell, 0, len(spellIDs))
	for _, spellID := range spellIDs {
		spell, err := c.spellRepo.GetSpell(r.Context(), spellID)
		if err != nil {
			// Skip spells that can't be found
			continue
		}
		spells = append(spells, spell)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spells); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}
