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

type PotionController struct {
	potionRepo repositories.PotionRepository
	tmpl       *template.Template
}

func NewPotionController(potionRepo repositories.PotionRepository, tmpl *template.Template) *PotionController {
	return &PotionController{
		potionRepo: potionRepo,
		tmpl:       tmpl,
	}
}

func (c *PotionController) GetPotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid potion ID format"))
		return
	}
	potion, err := c.potionRepo.GetPotion(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(potion); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "potion.html", potion); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *PotionController) ListPotions(w http.ResponseWriter, r *http.Request) {
	potions, err := c.potionRepo.ListPotions(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(potions); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *PotionController) CreatePotion(w http.ResponseWriter, r *http.Request) {
	var input models.CreatePotionInput
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
	existingPotion, err := c.potionRepo.GetPotionByName(r.Context(), input.Name)
	if err == nil && existingPotion != nil {
		validationErrors := map[string]string{
			"name": "Potion with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	id, err := c.potionRepo.CreatePotion(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	potion, err := c.potionRepo.GetPotion(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(potion); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *PotionController) UpdatePotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid potion ID format"))
		return
	}
	var input models.UpdatePotionInput
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
	existingPotion, err := c.potionRepo.GetPotionByName(r.Context(), input.Name)
	if err == nil && existingPotion != nil && existingPotion.ID != id {
		validationErrors := map[string]string{
			"name": "Potion with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	if err := c.potionRepo.UpdatePotion(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedPotion, err := c.potionRepo.GetPotion(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedPotion); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *PotionController) DeletePotion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid potion ID format"))
		return
	}
	if err := c.potionRepo.DeletePotion(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
