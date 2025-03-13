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

type EquipmentController struct {
	equipmentRepo repositories.EquipmentRepository
	tmpl          *template.Template
}

func NewEquipmentController(equipmentRepo repositories.EquipmentRepository, tmpl *template.Template) *EquipmentController {
	return &EquipmentController{
		equipmentRepo: equipmentRepo,
		tmpl:          tmpl,
	}
}

func (c *EquipmentController) GetEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid equipment ID format"))
		return
	}
	equipment, err := c.equipmentRepo.GetEquipment(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(equipment); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "equipment.html", equipment); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *EquipmentController) ListEquipment(w http.ResponseWriter, r *http.Request) {
	equipment, err := c.equipmentRepo.ListEquipment(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(equipment); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *EquipmentController) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	var input models.CreateEquipmentInput
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
	existingEquipment, err := c.equipmentRepo.GetEquipmentByName(r.Context(), input.Name)
	if err == nil && existingEquipment != nil {
		validationErrors := map[string]string{
			"name": "Equipment with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	id, err := c.equipmentRepo.CreateEquipment(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	equipment, err := c.equipmentRepo.GetEquipment(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(equipment); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *EquipmentController) UpdateEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid equipment ID format"))
		return
	}
	var input models.UpdateEquipmentInput
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
	existingEquipment, err := c.equipmentRepo.GetEquipmentByName(r.Context(), input.Name)
	if err == nil && existingEquipment != nil && existingEquipment.ID != id {
		validationErrors := map[string]string{
			"name": "Equipment with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	if err := c.equipmentRepo.UpdateEquipment(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedEquipment, err := c.equipmentRepo.GetEquipment(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedEquipment); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *EquipmentController) DeleteEquipment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid equipment ID format"))
		return
	}
	if err := c.equipmentRepo.DeleteEquipment(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
