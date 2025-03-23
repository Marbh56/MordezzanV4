package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mordezzanV4/internal/contextkeys"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"mordezzanV4/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type CharacterController struct {
	repo          repositories.CharacterRepository
	userRepo      repositories.UserRepository
	characterRepo repositories.CharacterRepository
	classService  *services.ClassService
	Templates     *template.Template
}

type UpdateHPInput struct {
	CurrentHitPoints   int `json:"current_hit_points"`
	MaxHitPoints       int `json:"max_hit_points"`
	TemporaryHitPoints int `json:"temporary_hit_points"`
}

func NewCharacterController(repo repositories.CharacterRepository, userRepo repositories.UserRepository, classService *services.ClassService, tmpl *template.Template) *CharacterController {
	return &CharacterController{
		repo:         repo,
		userRepo:     userRepo,
		classService: classService,
		Templates:    tmpl,
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

	// Enrich with class data using the unified class service
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Return the character data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (c *CharacterController) RenderCharacterDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	// Enrich with class data using the unified class service
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Get experience for next level if not at max level
	nextLevelExp, err := c.classService.GetExperienceForNextLevel(r.Context(), character.Class, character.Level)
	if err == nil && nextLevelExp > character.ExperiencePoints {
		data := struct {
			*models.Character
			NextLevelExperience int
			ExperienceNeeded    int
		}{
			Character:           character,
			NextLevelExperience: nextLevelExp,
			ExperienceNeeded:    nextLevelExp - character.ExperiencePoints,
		}
		err = c.Templates.ExecuteTemplate(w, "character_detail.html", data)
		if err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	err = c.Templates.ExecuteTemplate(w, "character_detail.html", character)
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

	// For class-specific level adjustments based on experience points
	if input.ExperiencePoints != existingCharacter.ExperiencePoints {
		// Get all class level data
		classLevelData, err := c.classService.GetAllClassLevelData(r.Context(), input.Class)
		if err == nil {
			// Find the appropriate level for the XP
			for i := len(classLevelData) - 1; i >= 0; i-- {
				if input.ExperiencePoints >= classLevelData[i].ExperiencePoints {
					input.Level = classLevelData[i].Level
					break
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

func (c *CharacterController) handleError(w http.ResponseWriter, err error, statusCode int) {
	logger.Error("Error in character controller: %v", err)

	// Set content type and status code
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	// Render error template or default error message
	errorMsg := err.Error()
	if statusCode == http.StatusInternalServerError {
		errorMsg = "An internal server error occurred"
	}

	// Try to render an error template if it exists
	if c.Templates != nil {
		errorData := map[string]interface{}{
			"Error":  errorMsg,
			"Status": statusCode,
			"Title":  "Error",
		}

		templateErr := c.Templates.ExecuteTemplate(w, "error", errorData)
		if templateErr == nil {
			return
		}
	}

	// Fallback to plain text response
	http.Error(w, errorMsg, statusCode)
}

func (c *CharacterController) RenderDashboard(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	userID, ok := r.Context().Value(contextkeys.UserIDKey).(int64)

	// If not authenticated, render the home page instead of dashboard
	if !ok {
		data := map[string]interface{}{
			"IsAuthenticated": false,
		}

		err := c.Templates.ExecuteTemplate(w, "home", data)
		if err != nil {
			c.handleError(w, err, http.StatusInternalServerError)
			return
		}
		return
	}

	// Fetch the user
	user, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		c.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Fetch the user's characters
	characters, err := c.repo.GetCharactersByUser(r.Context(), userID)
	if err != nil {
		c.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := map[string]interface{}{
		"IsAuthenticated": true,
		"User":            user,
		"Characters":      characters,
		"Title":           "Dashboard",
	}

	// Render the dashboard template
	err = c.Templates.ExecuteTemplate(w, "dashboard", data)
	if err != nil {
		c.handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (c *CharacterController) RenderCreateForm(w http.ResponseWriter, r *http.Request) {
	err := c.Templates.ExecuteTemplate(w, "character_create.html", nil)
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
	err = c.Templates.ExecuteTemplate(w, "character_edit.html", nil)
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

	var input UpdateHPInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate input - Max HP must be positive
	if input.MaxHitPoints <= 0 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Max hit points must be positive"))
		return
	}

	// Current HP cannot be less than -10 (death)
	if input.CurrentHitPoints < -10 {
		input.CurrentHitPoints = -10
	}

	// Temp HP cannot be negative
	if input.TemporaryHitPoints < 0 {
		input.TemporaryHitPoints = 0
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

	// Create update input with just the HP fields changed
	updateInput := models.UpdateCharacterInput{
		Name:               existingChar.Name,
		Class:              existingChar.Class,
		Level:              existingChar.Level,
		ExperiencePoints:   existingChar.ExperiencePoints,
		Strength:           existingChar.Strength,
		Dexterity:          existingChar.Dexterity,
		Constitution:       existingChar.Constitution,
		Wisdom:             existingChar.Wisdom,
		Intelligence:       existingChar.Intelligence,
		Charisma:           existingChar.Charisma,
		MaxHitPoints:       input.MaxHitPoints,
		CurrentHitPoints:   input.CurrentHitPoints,
		TemporaryHitPoints: input.TemporaryHitPoints,
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

func (c *CharacterController) ModifyCharacterHP(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	var input struct {
		Delta int  `json:"delta"` // Positive for healing, negative for damage
		Temp  bool `json:"temp"`  // If true, adds temporary HP instead of healing
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Get existing character
	existingChar, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	newCurrentHP := existingChar.CurrentHitPoints
	newTempHP := existingChar.TemporaryHitPoints

	if input.Delta < 0 {
		// Taking damage
		damageAmount := -input.Delta

		// Apply to temp HP first
		if newTempHP > 0 {
			if damageAmount <= newTempHP {
				newTempHP -= damageAmount
				damageAmount = 0
			} else {
				damageAmount -= newTempHP
				newTempHP = 0
			}
		}

		// Apply remaining damage to current HP
		if damageAmount > 0 {
			if newCurrentHP-damageAmount < -10 {
				newCurrentHP = -10
			} else {
				newCurrentHP = newCurrentHP - damageAmount
			}
		}
	} else if input.Delta > 0 {
		// Healing or temp HP
		if input.Temp {
			// Adding temporary hit points
			newTempHP += input.Delta
		} else {
			// Regular healing
			if newCurrentHP+input.Delta > existingChar.MaxHitPoints {
				newCurrentHP = existingChar.MaxHitPoints
			} else {
				newCurrentHP = newCurrentHP + input.Delta
			}
		}
	}

	// Create update input with the new HP value
	updateInput := models.UpdateCharacterInput{
		Name:               existingChar.Name,
		Class:              existingChar.Class,
		Level:              existingChar.Level,
		ExperiencePoints:   existingChar.ExperiencePoints,
		Strength:           existingChar.Strength,
		Dexterity:          existingChar.Dexterity,
		Constitution:       existingChar.Constitution,
		Wisdom:             existingChar.Wisdom,
		Intelligence:       existingChar.Intelligence,
		Charisma:           existingChar.Charisma,
		MaxHitPoints:       existingChar.MaxHitPoints,
		CurrentHitPoints:   newCurrentHP,
		TemporaryHitPoints: newTempHP,
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

	// Get all class level data
	newLevel := existingChar.Level
	classLevelData, err := c.classService.GetAllClassLevelData(r.Context(), existingChar.Class)
	if err == nil {
		// Find the appropriate level for the XP
		for i := len(classLevelData) - 1; i >= 0; i-- {
			if input.ExperiencePoints >= classLevelData[i].ExperiencePoints {
				newLevel = classLevelData[i].Level
				break
			}
		}
	}

	// Create update input with updated XP and potentially level
	updateInput := models.UpdateCharacterInput{
		Name:               existingChar.Name,
		Class:              existingChar.Class,
		Level:              newLevel,
		ExperiencePoints:   input.ExperiencePoints,
		Strength:           existingChar.Strength,
		Dexterity:          existingChar.Dexterity,
		Constitution:       existingChar.Constitution,
		Wisdom:             existingChar.Wisdom,
		Intelligence:       existingChar.Intelligence,
		Charisma:           existingChar.Charisma,
		MaxHitPoints:       existingChar.MaxHitPoints,
		CurrentHitPoints:   existingChar.CurrentHitPoints,
		TemporaryHitPoints: existingChar.TemporaryHitPoints,
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

	// Enrich with class data
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (c *CharacterController) GetCharacterClassData(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	// Get the character to determine the class
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", id))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	// Get all level data for this class
	levelData, err := c.classService.GetAllClassLevelData(r.Context(), character.Class)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Enrich character with class-specific data
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Get class abilities
	abilities, err := c.classService.GetClassAbilitiesByLevel(r.Context(), character.Class, character.Level)
	if err != nil {
		abilities = []*models.ClassAbility{} // Set to empty array if error
	}

	// Create a response with both the full level progression and current abilities
	classData := map[string]interface{}{
		"class_type": character.Class,
		"level_data": levelData,
		"current_level_data": map[string]interface{}{
			"level":            character.Level,
			"hit_dice":         character.HitDice,
			"saving_throw":     character.SavingThrow,
			"fighting_ability": character.FightingAbility,
			"casting_ability":  character.CastingAbility,
			"spell_slots":      character.SpellSlots,
			"abilities":        abilities,
		},
	}

	// Return the class data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classData)
}
