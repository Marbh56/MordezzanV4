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
	"strings"

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
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	if err := c.tmpl.ExecuteTemplate(w, "character.html", character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *CharacterController) GetCharactersByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid user ID format"))
		return
	}
	characters, err := c.characterRepo.GetCharactersByUser(r.Context(), userID)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(characters); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *CharacterController) ListCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := c.characterRepo.ListCharacters(r.Context())
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(characters); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *CharacterController) CreateCharacter(w http.ResponseWriter, r *http.Request) {
	var input models.CreateCharacterInput
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
	if c.userRepo != nil {
		_, err := c.userRepo.GetUser(r.Context(), input.UserID)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("User with ID %d does not exist", input.UserID)))
			return
		}
	}
	id, err := c.characterRepo.CreateCharacter(r.Context(), &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *CharacterController) UpdateCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}
	var input models.UpdateCharacterInput
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
	if err := c.characterRepo.UpdateCharacter(r.Context(), id, &input); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	updatedCharacter, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedCharacter); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
}

func (c *CharacterController) DeleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid character ID format"))
		return
	}
	if err := c.characterRepo.DeleteCharacter(r.Context(), id); err != nil {
		apperrors.HandleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
