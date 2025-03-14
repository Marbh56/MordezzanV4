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

type AmmoController struct {
	ammoRepo repositories.AmmoRepository
	tmpl     *template.Template
}

func NewAmmoController(ammoRepo repositories.AmmoRepository, tmpl *template.Template) *AmmoController {
	return &AmmoController{
		ammoRepo: ammoRepo,
		tmpl:     tmpl,
	}
}

func (c *AmmoController) GetAmmo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ammo ID format"))
		return
	}
	ammo, err := c.ammoRepo.GetAmmo(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ammo); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "ammo.html", ammo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AmmoController) ListAmmo(w http.ResponseWriter, r *http.Request) {
	ammoList, err := c.ammoRepo.ListAmmo(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ammoList); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AmmoController) CreateAmmo(w http.ResponseWriter, r *http.Request) {
	var input models.CreateAmmoInput
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
	existingAmmo, err := c.ammoRepo.GetAmmoByName(r.Context(), input.Name)
	if err == nil && existingAmmo != nil {
		validationErrors := map[string]string{
			"name": "Ammo with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	id, err := c.ammoRepo.CreateAmmo(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	ammo, err := c.ammoRepo.GetAmmo(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ammo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AmmoController) UpdateAmmo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ammo ID format"))
		return
	}
	var input models.UpdateAmmoInput
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
	existingAmmo, err := c.ammoRepo.GetAmmoByName(r.Context(), input.Name)
	if err == nil && existingAmmo != nil && existingAmmo.ID != id {
		validationErrors := map[string]string{
			"name": "Ammo with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	if err := c.ammoRepo.UpdateAmmo(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedAmmo, err := c.ammoRepo.GetAmmo(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedAmmo); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *AmmoController) DeleteAmmo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ammo ID format"))
		return
	}
	if err := c.ammoRepo.DeleteAmmo(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
