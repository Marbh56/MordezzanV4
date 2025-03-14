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

type RingController struct {
	ringRepo repositories.RingRepository
	tmpl     *template.Template
}

func NewRingController(ringRepo repositories.RingRepository, tmpl *template.Template) *RingController {
	return &RingController{
		ringRepo: ringRepo,
		tmpl:     tmpl,
	}
}

func (c *RingController) GetRing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ring ID format"))
		return
	}
	ring, err := c.ringRepo.GetRing(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ring); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "ring.html", ring); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *RingController) ListRings(w http.ResponseWriter, r *http.Request) {
	rings, err := c.ringRepo.ListRings(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rings); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *RingController) CreateRing(w http.ResponseWriter, r *http.Request) {
	var input models.CreateRingInput
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
	existingRing, err := c.ringRepo.GetRingByName(r.Context(), input.Name)
	if err == nil && existingRing != nil {
		validationErrors := map[string]string{
			"name": "Ring with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	id, err := c.ringRepo.CreateRing(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	ring, err := c.ringRepo.GetRing(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ring); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *RingController) UpdateRing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ring ID format"))
		return
	}
	var input models.UpdateRingInput
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
	existingRing, err := c.ringRepo.GetRingByName(r.Context(), input.Name)
	if err == nil && existingRing != nil && existingRing.ID != id {
		validationErrors := map[string]string{
			"name": "Ring with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}
	if err := c.ringRepo.UpdateRing(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedRing, err := c.ringRepo.GetRing(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedRing); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *RingController) DeleteRing(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ring ID format"))
		return
	}
	if err := c.ringRepo.DeleteRing(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
