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

type MagicItemController struct {
	magicItemRepo repositories.MagicItemRepository
	tmpl          *template.Template
}

func NewMagicItemController(magicItemRepo repositories.MagicItemRepository, tmpl *template.Template) *MagicItemController {
	return &MagicItemController{
		magicItemRepo: magicItemRepo,
		tmpl:          tmpl,
	}
}

func (c *MagicItemController) GetMagicItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid magic item ID format"))
		return
	}
	item, err := c.magicItemRepo.GetMagicItem(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(item); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "magic_item.html", item); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *MagicItemController) ListMagicItems(w http.ResponseWriter, r *http.Request) {
	itemType := r.URL.Query().Get("type")
	var items []*models.MagicItem
	var err error

	if itemType != "" {
		if !models.IsValidItemType(itemType) {
			apperrors.HandleError(w, apperrors.NewBadRequest("Invalid item type. Must be 'Rod', 'Wand', or 'Staff'"))
			return
		}
		items, err = c.magicItemRepo.ListMagicItemsByType(r.Context(), itemType)
	} else {
		items, err = c.magicItemRepo.ListMagicItems(r.Context())
	}

	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *MagicItemController) CreateMagicItem(w http.ResponseWriter, r *http.Request) {
	var input models.CreateMagicItemInput
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

	existingItem, err := c.magicItemRepo.GetMagicItemByName(r.Context(), input.Name)
	if err == nil && existingItem != nil {
		validationErrors := map[string]string{
			"name": "Magic item with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	id, err := c.magicItemRepo.CreateMagicItem(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	item, err := c.magicItemRepo.GetMagicItem(r.Context(), id)
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

func (c *MagicItemController) UpdateMagicItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid magic item ID format"))
		return
	}
	var input models.UpdateMagicItemInput
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

	existingItem, err := c.magicItemRepo.GetMagicItemByName(r.Context(), input.Name)
	if err == nil && existingItem != nil && existingItem.ID != id {
		validationErrors := map[string]string{
			"name": "Magic item with this name already exists",
		}
		apperrors.HandleValidationErrors(w, validationErrors)
		return
	}

	if err := c.magicItemRepo.UpdateMagicItem(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedItem, err := c.magicItemRepo.GetMagicItem(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedItem); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *MagicItemController) DeleteMagicItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid magic item ID format"))
		return
	}
	if err := c.magicItemRepo.DeleteMagicItem(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
