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

type ShieldController struct {
	shieldRepo repositories.ShieldRepository
	tmpl       *template.Template
}

func NewShieldController(shieldRepo repositories.ShieldRepository, tmpl *template.Template) *ShieldController {
	return &ShieldController{
		shieldRepo: shieldRepo,
		tmpl:       tmpl,
	}
}

func (c *ShieldController) GetShield(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid shield ID format"))
		return
	}

	shield, err := c.shieldRepo.GetShield(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(shield); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "shield.html", shield); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ShieldController) ListShields(w http.ResponseWriter, r *http.Request) {
	shields, err := c.shieldRepo.ListShields(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shields); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ShieldController) CreateShield(w http.ResponseWriter, r *http.Request) {
	var input models.CreateShieldInput
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

	existingShield, err := c.shieldRepo.GetShieldByName(r.Context(), input.Name)
	if err == nil && existingShield != nil {
		validationErrors := map[string]string{
			"name": "Shield with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	id, err := c.shieldRepo.CreateShield(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	shield, err := c.shieldRepo.GetShield(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(shield); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ShieldController) UpdateShield(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid shield ID format"))
		return
	}

	var input models.UpdateShieldInput
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

	existingShield, err := c.shieldRepo.GetShieldByName(r.Context(), input.Name)
	if err == nil && existingShield != nil && existingShield.ID != id {
		validationErrors := map[string]string{
			"name": "Shield with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	if err := c.shieldRepo.UpdateShield(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedShield, err := c.shieldRepo.GetShield(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedShield); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ShieldController) DeleteShield(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid shield ID format"))
		return
	}

	if err := c.shieldRepo.DeleteShield(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
