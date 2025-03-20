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

// SpellController handles HTTP requests related to spells
type SpellController struct {
	spellRepo repositories.SpellRepository
	tmpl      *template.Template
}

// NewSpellController creates a new SpellController
func NewSpellController(spellRepo repositories.SpellRepository, tmpl *template.Template) *SpellController {
	return &SpellController{
		spellRepo: spellRepo,
		tmpl:      tmpl,
	}
}

// GetSpell handles requests to get a specific spell by ID
func (c *SpellController) GetSpell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	spell, err := c.spellRepo.GetSpell(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(spell); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "spell.html", spell); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// ListSpells handles requests to get all spells
func (c *SpellController) ListSpells(w http.ResponseWriter, r *http.Request) {
	spells, err := c.spellRepo.ListSpells(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spells); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// CreateSpell handles requests to create a new spell
func (c *SpellController) CreateSpell(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSpellInput
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

	id, err := c.spellRepo.CreateSpell(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	spell, err := c.spellRepo.GetSpell(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(spell); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// UpdateSpell handles requests to update an existing spell
func (c *SpellController) UpdateSpell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	var input models.UpdateSpellInput
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

	if err := c.spellRepo.UpdateSpell(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedSpell, err := c.spellRepo.GetSpell(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSpell); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// DeleteSpell handles requests to delete a spell
func (c *SpellController) DeleteSpell(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
		return
	}

	if err := c.spellRepo.DeleteSpell(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
