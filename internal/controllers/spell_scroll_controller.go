package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type SpellScrollController struct {
	spellScrollRepo repositories.SpellScrollRepository
	spellRepo       repositories.SpellRepository
	tmpl            *template.Template
}

func NewSpellScrollController(spellScrollRepo repositories.SpellScrollRepository, spellRepo repositories.SpellRepository, tmpl *template.Template) *SpellScrollController {
	return &SpellScrollController{
		spellScrollRepo: spellScrollRepo,
		spellRepo:       spellRepo,
		tmpl:            tmpl,
	}
}

func (c *SpellScrollController) GetSpellScroll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell scroll ID format"))
		return
	}

	spellScroll, err := c.spellScrollRepo.GetSpellScroll(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(spellScroll); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "spell_scroll.html", spellScroll); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellScrollController) ListSpellScrolls(w http.ResponseWriter, r *http.Request) {
	spellIDStr := r.URL.Query().Get("spell_id")

	var spellScrolls []*models.SpellScroll
	var err error

	if spellIDStr != "" {
		spellID, err := strconv.ParseInt(spellIDStr, 10, 64)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell ID format"))
			return
		}
		spellScrolls, err = c.spellScrollRepo.GetSpellScrollsBySpell(r.Context(), spellID)
	} else {
		spellScrolls, err = c.spellScrollRepo.ListSpellScrolls(r.Context())
	}

	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spellScrolls); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellScrollController) CreateSpellScroll(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSpellScrollInput
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

	// Verify that the referenced spell exists
	if c.spellRepo != nil {
		_, err := c.spellRepo.GetSpell(r.Context(), input.SpellID)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Spell with ID %d does not exist", input.SpellID)))
			return
		}
	}

	id, err := c.spellScrollRepo.CreateSpellScroll(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	spellScroll, err := c.spellScrollRepo.GetSpellScroll(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(spellScroll); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellScrollController) UpdateSpellScroll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell scroll ID format"))
		return
	}

	var input models.UpdateSpellScrollInput
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

	// Verify that the referenced spell exists
	if c.spellRepo != nil {
		_, err := c.spellRepo.GetSpell(r.Context(), input.SpellID)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Spell with ID %d does not exist", input.SpellID)))
			return
		}
	}

	if err := c.spellScrollRepo.UpdateSpellScroll(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedSpellScroll, err := c.spellScrollRepo.GetSpellScroll(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedSpellScroll); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *SpellScrollController) DeleteSpellScroll(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid spell scroll ID format"))
		return
	}

	if err := c.spellScrollRepo.DeleteSpellScroll(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
