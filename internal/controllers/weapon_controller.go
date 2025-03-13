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

type WeaponController struct {
	weaponRepo repositories.WeaponRepository
	tmpl       *template.Template
}

func NewWeaponController(weaponRepo repositories.WeaponRepository, tmpl *template.Template) *WeaponController {
	return &WeaponController{
		weaponRepo: weaponRepo,
		tmpl:       tmpl,
	}
}

func (c *WeaponController) GetWeapon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid weapon ID format"))
		return
	}
	weapon, err := c.weaponRepo.GetWeapon(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(weapon); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	// If we had a weapon.html template, we would use it here
	// For now, we'll just return JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weapon); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *WeaponController) ListWeapons(w http.ResponseWriter, r *http.Request) {
	weapons, err := c.weaponRepo.ListWeapons(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weapons); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *WeaponController) CreateWeapon(w http.ResponseWriter, r *http.Request) {
	var input models.CreateWeaponInput
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

	existingWeapon, err := c.weaponRepo.GetWeaponByName(r.Context(), input.Name)
	if err == nil && existingWeapon != nil {
		validationErrors := map[string]string{
			"name": "Weapon with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	id, err := c.weaponRepo.CreateWeapon(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	weapon, err := c.weaponRepo.GetWeapon(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(weapon); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *WeaponController) UpdateWeapon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid weapon ID format"))
		return
	}

	var input models.UpdateWeaponInput
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

	existingWeapon, err := c.weaponRepo.GetWeaponByName(r.Context(), input.Name)
	if err == nil && existingWeapon != nil && existingWeapon.ID != id {
		validationErrors := map[string]string{
			"name": "Weapon with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	if err := c.weaponRepo.UpdateWeapon(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedWeapon, err := c.weaponRepo.GetWeapon(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedWeapon); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *WeaponController) DeleteWeapon(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid weapon ID format"))
		return
	}

	if err := c.weaponRepo.DeleteWeapon(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
