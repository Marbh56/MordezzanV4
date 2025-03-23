package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/repositories"

	"github.com/go-chi/chi"
)

// BaseController provides generic CRUD operations
type BaseController[
	T any,
	CreateInput any,
	UpdateInput any,
	Validator interface {
		Validate() error
	},
] struct {
	Repository   repositories.Repository[T, CreateInput, UpdateInput]
	TemplateName string
	Tmpl         *template.Template
}

// Get handles fetching a single entity
func (c *BaseController[T, CreateInput, UpdateInput, Validator]) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ID format"))
		return
	}

	item, err := c.Repository.Get(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// List handles listing all entities
func (c *BaseController[T, CreateInput, UpdateInput, Validator]) List(w http.ResponseWriter, r *http.Request) {
	items, err := c.Repository.List(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// Create handles creating a new entity
func (c *BaseController[T, CreateInput, UpdateInput, Validator]) Create(w http.ResponseWriter, r *http.Request) {
	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate if input implements Validator
	if validator, ok := any(&input).(Validator); ok {
		if err := validator.Validate(); err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	id, err := c.Repository.Create(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	item, err := c.Repository.Get(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(item); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// Update handles updating an entity
func (c *BaseController[T, CreateInput, UpdateInput, Validator]) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ID format"))
		return
	}

	var input UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate if input implements Validator
	if validator, ok := any(&input).(Validator); ok {
		if err := validator.Validate(); err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	if err := c.Repository.Update(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	item, err := c.Repository.Get(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

// Delete handles deleting an entity
func (c *BaseController[T, CreateInput, UpdateInput, Validator]) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid ID format"))
		return
	}

	if err := c.Repository.Delete(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
