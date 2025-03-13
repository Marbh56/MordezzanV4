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

type ArmorController struct {
	armorRepo repositories.ArmorRepository
	tmpl      *template.Template
}

func NewArmorController(armorRepo repositories.ArmorRepository, tmpl *template.Template) *ArmorController {
	return &ArmorController{
		armorRepo: armorRepo,
		tmpl:      tmpl,
	}
}

func (c *ArmorController) GetArmor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid armor ID format"))
		return
	}
	armor, err := c.armorRepo.GetArmor(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(armor); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "armor.html", armor); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ArmorController) ListArmors(w http.ResponseWriter, r *http.Request) {
	armors, err := c.armorRepo.ListArmors(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(armors); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ArmorController) CreateArmor(w http.ResponseWriter, r *http.Request) {
	var input models.CreateArmorInput
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

	// Check if armor with the same name already exists
	existingArmor, err := c.armorRepo.GetArmorByName(r.Context(), input.Name)
	if err == nil && existingArmor != nil {
		validationErrors := map[string]string{
			"name": "Armor with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	id, err := c.armorRepo.CreateArmor(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	armor, err := c.armorRepo.GetArmor(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(armor); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ArmorController) UpdateArmor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid armor ID format"))
		return
	}
	var input models.UpdateArmorInput
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

	// Check if we're trying to update to a name that's already taken
	existingArmor, err := c.armorRepo.GetArmorByName(r.Context(), input.Name)
	if err == nil && existingArmor != nil && existingArmor.ID != id {
		validationErrors := map[string]string{
			"name": "Armor with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	if err := c.armorRepo.UpdateArmor(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedArmor, err := c.armorRepo.GetArmor(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedArmor); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *ArmorController) DeleteArmor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid armor ID format"))
		return
	}
	if err := c.armorRepo.DeleteArmor(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
