package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"mordezzanV4/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CharacterController struct {
	characterRepo  repositories.CharacterRepository
	userRepo       repositories.UserRepository
	fighterService *services.FighterService
	tmpl           *template.Template
}

func NewCharacterController(
	characterRepo repositories.CharacterRepository,
	userRepo repositories.UserRepository,
	fighterService *services.FighterService,
	tmpl *template.Template,
) *CharacterController {
	return &CharacterController{
		characterRepo:  characterRepo,
		userRepo:       userRepo,
		fighterService: fighterService,
		tmpl:           tmpl,
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

	// Enrich with fighter class data if applicable
	if character.Class == "Fighter" {
		if err := c.fighterService.EnrichCharacterWithFighterData(r.Context(), character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}
	}

	// Return the character data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
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

	// Enrich with fighter class data if applicable
	if character.Class == "Fighter" {
		if err := c.fighterService.EnrichCharacterWithFighterData(r.Context(), character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		// Get next level experience for the character
		nextLevelExp, err := c.fighterService.GetExperienceForNextLevel(r.Context(), character.Level)
		if err == nil && character.Level < 12 {
			// Create a data wrapper to pass to the template
			data := struct {
				*models.Character
				NextLevelExperience int
				ExperienceNeeded    int
			}{
				Character:           character,
				NextLevelExperience: nextLevelExp,
				ExperienceNeeded:    nextLevelExp - character.ExperiencePoints,
			}

			character.CalculateDerivedStats()
			err = c.tmpl.ExecuteTemplate(w, "character_detail.html", data)
			if err != nil {
				apperrors.HandleError(w, apperrors.NewInternalError(err))
			}
			return
		}
	}

	// If not fighter or error getting next level exp, just render the character
	err = c.tmpl.ExecuteTemplate(w, "character_detail.html", character)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
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
	existingCharacter, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// For Fighter class, update level based on experience points
	if existingCharacter.Class == "Fighter" && input.Class == "Fighter" {
		// Check if experience points have changed
		if input.ExperiencePoints != existingCharacter.ExperiencePoints {
			// Get all fighter level data
			fighterLevels, err := c.fighterService.GetAllFighterLevelData(r.Context())
			if err == nil {
				// Find the appropriate level for the XP
				for i := len(fighterLevels) - 1; i >= 0; i-- {
					if input.ExperiencePoints >= fighterLevels[i].ExperiencePoints {
						input.Level = fighterLevels[i].Level
						break
					}
				}
			}
		}
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

func (c *CharacterController) RenderEditForm(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	// Verify the character exists
	_, err = c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Render the edit form template
	err = c.tmpl.ExecuteTemplate(w, "character_edit.html", nil)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}

func (c *CharacterController) UpdateCharacterHP(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	var input struct {
		HitPoints int `json:"hit_points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Get existing character to preserve other fields
	existingChar, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Create update input with just the HP changed
	updateInput := models.UpdateCharacterInput{
		Name:             existingChar.Name,
		Class:            existingChar.Class,
		Level:            existingChar.Level,
		ExperiencePoints: existingChar.ExperiencePoints,
		Strength:         existingChar.Strength,
		Dexterity:        existingChar.Dexterity,
		Constitution:     existingChar.Constitution,
		Wisdom:           existingChar.Wisdom,
		Intelligence:     existingChar.Intelligence,
		Charisma:         existingChar.Charisma,
		HitPoints:        input.HitPoints,
	}

	// Update the character
	err = c.characterRepo.UpdateCharacter(r.Context(), id, &updateInput)
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

func (c *CharacterController) UpdateCharacterXP(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	var input struct {
		ExperiencePoints int `json:"experience_points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Ensure XP is not negative
	if input.ExperiencePoints < 0 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Experience points cannot be negative"))
		return
	}

	// Get existing character to preserve other fields
	existingChar, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// For Fighter class, also update level based on new XP
	newLevel := existingChar.Level
	if existingChar.Class == "Fighter" {
		// Get fighter levels to determine new level based on XP
		fighterLevels, err := c.fighterService.GetAllFighterLevelData(r.Context())
		if err == nil {
			// Find the appropriate level for the XP
			for i := len(fighterLevels) - 1; i >= 0; i-- {
				if input.ExperiencePoints >= fighterLevels[i].ExperiencePoints {
					newLevel = fighterLevels[i].Level
					break
				}
			}
		}
	}

	// Create update input with updated XP and potentially level
	updateInput := models.UpdateCharacterInput{
		Name:             existingChar.Name,
		Class:            existingChar.Class,
		Level:            newLevel,
		ExperiencePoints: input.ExperiencePoints,
		Strength:         existingChar.Strength,
		Dexterity:        existingChar.Dexterity,
		Constitution:     existingChar.Constitution,
		Wisdom:           existingChar.Wisdom,
		Intelligence:     existingChar.Intelligence,
		Charisma:         existingChar.Charisma,
		HitPoints:        existingChar.HitPoints,
	}

	// Update the character
	err = c.characterRepo.UpdateCharacter(r.Context(), id, &updateInput)
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

	// Apply fighter class data if applicable
	if character.Class == "Fighter" {
		if err := c.fighterService.EnrichCharacterWithFighterData(r.Context(), character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}
