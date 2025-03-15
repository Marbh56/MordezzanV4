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

type TreasureController struct {
	treasureRepo  repositories.TreasureRepository
	characterRepo repositories.CharacterRepository
	tmpl          *template.Template
}

func NewTreasureController(treasureRepo repositories.TreasureRepository, characterRepo repositories.CharacterRepository, tmpl *template.Template) *TreasureController {
	return &TreasureController{
		treasureRepo:  treasureRepo,
		characterRepo: characterRepo,
		tmpl:          tmpl,
	}
}

func (c *TreasureController) GetTreasure(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid treasure ID format"))
		return
	}

	treasure, err := c.treasureRepo.GetTreasure(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(treasure); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	if err := c.tmpl.ExecuteTemplate(w, "treasure.html", treasure); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *TreasureController) GetTreasureByCharacter(w http.ResponseWriter, r *http.Request) {
	characterID, err := strconv.ParseInt(chi.URLParam(r, "characterId"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}

	// Check if character exists
	if c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), characterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
	}

	treasure, err := c.treasureRepo.GetTreasureByCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(treasure); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *TreasureController) ListTreasures(w http.ResponseWriter, r *http.Request) {
	treasures, err := c.treasureRepo.ListTreasures(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(treasures); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *TreasureController) CreateTreasure(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTreasureInput
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

	// Check if character exists if character_id provided
	if input.CharacterID != nil && c.characterRepo != nil {
		_, err := c.characterRepo.GetCharacter(r.Context(), *input.CharacterID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}

		// Check if treasure already exists for this character
		existingTreasure, err := c.treasureRepo.GetTreasureByCharacter(r.Context(), *input.CharacterID)
		if err == nil && existingTreasure != nil {
			validationErrors := map[string]string{
				"character_id": "Treasure already exists for this character",
			}
			apperrors.HandleValidationErrors(w, validationErrors)
			return
		}
	}

	id, err := c.treasureRepo.CreateTreasure(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	treasure, err := c.treasureRepo.GetTreasure(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(treasure); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *TreasureController) UpdateTreasure(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid treasure ID format"))
		return
	}

	var input models.UpdateTreasureInput
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

	if err := c.treasureRepo.UpdateTreasure(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	updatedTreasure, err := c.treasureRepo.GetTreasure(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedTreasure); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *TreasureController) DeleteTreasure(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid treasure ID format"))
		return
	}

	if err := c.treasureRepo.DeleteTreasure(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
