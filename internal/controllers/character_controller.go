package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"mordezzanV4/internal/services"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
)

type CharacterController struct {
	userRepo       repositories.UserRepository
	characterRepo  repositories.CharacterRepository
	classService   *services.ClassService
	Templates      *template.Template
	sessionManager *scs.SessionManager
}

type UpdateHPInput struct {
	CurrentHitPoints   int `json:"current_hit_points"`
	MaxHitPoints       int `json:"max_hit_points"`
	TemporaryHitPoints int `json:"temporary_hit_points"`
}

func NewCharacterController(repo repositories.CharacterRepository, userRepo repositories.UserRepository, classService *services.ClassService, tmpl *template.Template, sessionManager *scs.SessionManager) *CharacterController {
	return &CharacterController{
		characterRepo:  repo,
		userRepo:       userRepo,
		classService:   classService,
		Templates:      tmpl,
		sessionManager: sessionManager,
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
	logger.Info("Starting RenderCharacterDetail function")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	idParam := chi.URLParam(r, "id")
	logger.Debug("Character ID from URL: %s", idParam)

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		logger.Error("Failed to parse character ID: %v", err)
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}
	logger.Debug("Parsed character ID: %d", id)

	logger.Debug("Fetching character from repository")
	character, err := c.characterRepo.GetCharacter(r.Context(), id)
	if err != nil {
		logger.Error("Failed to get character from repository: %v", err)
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
	logger.Debug("Successfully fetched character: %s (ID: %d)", character.Name, character.ID)

	// Log character details to identify potential nil values
	logger.Debug("Character details - Class: %s, Level: %d", character.Class, character.Level)

	logger.Debug("Enriching character with class data")
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		logger.Error("Failed to enrich character with class data: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
	logger.Debug("Successfully enriched character with class data")

	logger.Debug("Getting experience for next level")
	nextLevelExp, err := c.classService.GetExperienceForNextLevel(r.Context(), character.Class, character.Level)
	if err != nil {
		logger.Error("Failed to get experience for next level: %v", err)
		// Continue without next level experience info
		logger.Debug("Executing template without next level experience info")
		err = c.Templates.ExecuteTemplate(w, "character_detail", character)
		if err != nil {
			logger.Error("Failed to execute template: %v", err)
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}
	logger.Debug("Experience for next level: %d", nextLevelExp)

	if nextLevelExp > character.ExperiencePoints {
		logger.Debug("Character needs %d more XP to level up", nextLevelExp-character.ExperiencePoints)
		data := struct {
			*models.Character
			NextLevelExperience int
			ExperienceNeeded    int
		}{
			Character:           character,
			NextLevelExperience: nextLevelExp,
			ExperienceNeeded:    nextLevelExp - character.ExperiencePoints,
		}
		logger.Debug("Executing template with next level experience info")
		err = c.Templates.ExecuteTemplate(w, "character_detail", data)
		if err != nil {
			logger.Error("Failed to execute template with next level experience: %v", err)
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	logger.Debug("Character is at max level or has enough XP to level up")
	logger.Debug("Executing template without next level experience info")
	err = c.Templates.ExecuteTemplate(w, "character_detail", character)
	if err != nil {
		logger.Error("Failed to execute template without next level experience: %v", err)
		apperrors.HandleError(w, apperrors.NewInternalError(err))
	}
	logger.Info("Completed RenderCharacterDetail function")
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
	userID := c.sessionManager.GetInt64(r.Context(), "userID")

	// Fetch the user
	user, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		c.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// Fetch the user's characters
	characters, err := c.characterRepo.GetCharactersByUser(r.Context(), userID)
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
	// Get current user from session
	userID := c.sessionManager.GetInt64(r.Context(), "userID")
	if userID == 0 {
		http.Redirect(w, r, "/auth/login-page", http.StatusSeeOther)
		return
	}

	user, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Pass user data to the template
	data := map[string]interface{}{
		"User": user,
	}

	err = c.Templates.ExecuteTemplate(w, "character_create", data)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}
}

func (c *CharacterController) RenderEditForm(w http.ResponseWriter, r *http.Request) {
	// Get current user from session
	userID := c.sessionManager.GetInt64(r.Context(), "userID")
	if userID == 0 {
		http.Redirect(w, r, "/auth/login-page", http.StatusSeeOther)
		return
	}

	user, err := c.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

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

	// Check if the character belongs to the current user
	if character.UserID != userID {
		http.Error(w, "Unauthorized access to this character", http.StatusForbidden)
		return
	}

	data := map[string]interface{}{
		"Character": character,
		"IsEdit":    true,
		"User":      user,
	}

	err = c.Templates.ExecuteTemplate(w, "character_create", data)
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

	// Enrich character with class-specific data (including abilities)
	if err := c.classService.EnrichCharacterWithClassData(r.Context(), character); err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Extract abilities from the character object
	var abilities []*models.ClassAbility
	if character.Abilities != nil {
		if abilitiesMap, ok := character.Abilities.(map[string]interface{}); ok {
			if classAbilities, ok := abilitiesMap["class_abilities"]; ok {
				if cas, ok := classAbilities.([]*models.ClassAbility); ok {
					abilities = cas
				}
			}
		}
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
