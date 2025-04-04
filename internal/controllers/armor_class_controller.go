package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/services"

	"github.com/go-chi/chi"
)

// ACController handles armor class related requests
type ACController struct {
	acService *services.ACService
}

// NewACController creates a new AC controller
func NewACController(acService *services.ACService) *ACController {
	return &ACController{
		acService: acService,
	}
}

// GetCharacterAC returns a character's armor class
func (c *ACController) GetCharacterAC(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Calculate AC using the service
	acDetails, err := c.acService.CalculateCharacterAC(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acDetails)
}
