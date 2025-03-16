package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CharacterController struct {
	characterRepo repositories.CharacterRepository
	userRepo      repositories.UserRepository
	tmpl          *template.Template
}

func NewCharacterController(characterRepo repositories.CharacterRepository, userRepo repositories.UserRepository, tmpl *template.Template) *CharacterController {
	return &CharacterController{
		characterRepo: characterRepo,
		userRepo:      userRepo,
		tmpl:          tmpl,
	}
}

func (c *CharacterController) GetCharacter(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	// Repository call
	character, err := c.characterRepo.GetCharacter(r.Context(), id)

	// Add back error handling
	if err != nil {
		fmt.Printf("DEBUG: Error from repository: %v\n", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	// Return the character data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (c *CharacterController) GetCharactersByUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid user ID: %s", idParam)))
		return
	}

	// Verify the user exists
	user, err := c.userRepo.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("user", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	characters, err := c.characterRepo.GetCharactersByUser(r.Context(), user.ID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

func (c *CharacterController) ListCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := c.characterRepo.ListCharacters(r.Context())
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(characters)
}

func (c *CharacterController) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCharacterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate the input
	if err := input.Validate(); err != nil {
		var validationErr *models.ValidationError
		if errors.As(err, &validationErr) {
			fields := map[string]string{
				validationErr.Field: validationErr.Message,
			}
			apperrors.HandleValidationErrors(w, fields)
			return
		}
		apperrors.HandleError(w, apperrors.NewBadRequest(err.Error()))
		return
	}

	// Check if the user exists
	_, err := c.userRepo.GetUser(r.Context(), input.UserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("user", input.UserID))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Create the character
	id, err := c.characterRepo.CreateCharacter(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Get the created character
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(character)
}

func (c *CharacterController) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	var input models.UpdateCharacterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Check if the character exists
	_, err = c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Update the character
	err = c.characterRepo.UpdateCharacter(r.Context(), id, &input)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Get the updated character
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (c *CharacterController) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	// Check if the character exists
	_, err = c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Delete the character
	err = c.characterRepo.DeleteCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CharacterController) RenderDashboard(w http.ResponseWriter, r *http.Request) {
	err := c.tmpl.ExecuteTemplate(w, "dashboard.html", nil)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}

func (c *CharacterController) RenderCreateForm(w http.ResponseWriter, r *http.Request) {
	err := c.tmpl.ExecuteTemplate(w, "character_create.html", nil)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}

func (c *CharacterController) RenderCharacterDetail(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
	character.CalculateDerivedStats()

	err = c.tmpl.ExecuteTemplate(w, "character_detail.html", character)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}
