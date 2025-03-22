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
	characterRepo   repositories.CharacterRepository
	userRepo        repositories.UserRepository
	fighterService  *services.FighterService
	magicianService *services.MagicianService
	clericService   *services.ClericService
	tmpl            *template.Template
}

func NewCharacterController(
	characterRepo repositories.CharacterRepository,
	userRepo repositories.UserRepository,
	fighterService *services.FighterService,
	magicianService *services.MagicianService,
	clericService *services.ClericService,
	tmpl *template.Template,
) *CharacterController {
	return &CharacterController{
		characterRepo:   characterRepo,
		userRepo:        userRepo,
		fighterService:  fighterService,
		magicianService: magicianService,
		clericService:   clericService,
		tmpl:            tmpl,
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

	character.CalculateDerivedStats()

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

	// First calculate basic derived stats
	character.CalculateDerivedStats()

	// Then enrich with class-specific data (including save bonuses)
	switch character.Class {
	case "Fighter":
		if err := c.fighterService.EnrichCharacterWithFighterData(r.Context(), character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		nextLevelExp, err := c.fighterService.GetExperienceForNextLevel(r.Context(), character.Level)
		if err == nil && character.Level < 12 {
			data := struct {
				*models.Character
				NextLevelExperience int
				ExperienceNeeded    int
			}{
				Character:           character,
				NextLevelExperience: nextLevelExp,
				ExperienceNeeded:    nextLevelExp - character.ExperiencePoints,
			}
			err = c.tmpl.ExecuteTemplate(w, "character_detail.html", data)
			if err != nil {
				apperrors.HandleError(w, apperrors.NewInternalError(err))
			}
			return
		}
	case "Magician":
		if err := c.magicianService.EnrichCharacterWithMagicianData(r.Context(), character); err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}
	}

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
		CurrentHitPoints   int `json:"current_hit_points"`
		MaxHitPoints       int `json:"max_hit_points,omitempty"`       // Optional, only if changing max HP
		TemporaryHitPoints int `json:"temporary_hit_points,omitempty"` // Optional, only if changing temp HP
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// Validate the new hit points values
	if input.CurrentHitPoints < -10 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Current hit points cannot be less than -10"))
		return
	}

	if input.MaxHitPoints != 0 && input.MaxHitPoints < 1 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Max hit points must be positive"))
		return
	}

	if input.TemporaryHitPoints < 0 {
		apperrors.HandleError(w, apperrors.NewBadRequest("Temporary hit points cannot be negative"))
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

	// Create update input with HP fields changed
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
		CurrentHitPoints: input.CurrentHitPoints,
		// Only update max_hit_points if it was provided
		MaxHitPoints: existingChar.MaxHitPoints,
		// Only update temporary_hit_points if it was provided, otherwise keep existing
		TemporaryHitPoints: existingChar.TemporaryHitPoints,
	}

	// Update max hit points if provided
	if input.MaxHitPoints != 0 {
		updateInput.MaxHitPoints = input.MaxHitPoints
	}

	// Update temporary hit points if provided
	if r.FormValue("temporary_hit_points") != "" || r.ContentLength > 0 {
		updateInput.TemporaryHitPoints = input.TemporaryHitPoints
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
		DamageAmount    int  `json:"damage_amount,omitempty"`
		HealAmount      int  `json:"heal_amount,omitempty"`
		TemporaryHP     int  `json:"temporary_hp,omitempty"`
		IgnoreTemporary bool `json:"ignore_temporary,omitempty"`
		MaxHPChange     int  `json:"max_hp_change,omitempty"`
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

	// Create a copy to work with
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
		CurrentHitPoints:   existingChar.CurrentHitPoints,
		TemporaryHitPoints: existingChar.TemporaryHitPoints,
	}

	// Apply changes according to the input

	// Handle damage
	if input.DamageAmount > 0 {
		// First check if we can absorb damage with temporary HP
		if existingChar.TemporaryHitPoints > 0 && !input.IgnoreTemporary {
			if existingChar.TemporaryHitPoints >= input.DamageAmount {
				// Temp HP can absorb all the damage
				updateInput.TemporaryHitPoints = existingChar.TemporaryHitPoints - input.DamageAmount
			} else {
				// Temp HP absorbs part of the damage
				remainingDamage := input.DamageAmount - existingChar.TemporaryHitPoints
				updateInput.TemporaryHitPoints = 0
				updateInput.CurrentHitPoints = existingChar.CurrentHitPoints - remainingDamage
			}
		} else {
			// No temp HP or ignoring it, apply damage directly to current HP
			updateInput.CurrentHitPoints = existingChar.CurrentHitPoints - input.DamageAmount
		}
	}

	// Handle healing
	if input.HealAmount > 0 {
		// Healing applies to current HP, but cannot exceed max HP
		newHP := existingChar.CurrentHitPoints + input.HealAmount
		if newHP > existingChar.MaxHitPoints {
			newHP = existingChar.MaxHitPoints
		}
		updateInput.CurrentHitPoints = newHP
	}

	// Handle temporary HP (overwrites existing temp HP rather than stacking)
	if input.TemporaryHP > 0 {
		updateInput.TemporaryHitPoints = input.TemporaryHP
	}

	// Handle max HP changes
	if input.MaxHPChange != 0 {
		newMaxHP := existingChar.MaxHitPoints + input.MaxHPChange
		if newMaxHP < 1 {
			newMaxHP = 1 // Ensure max HP doesn't go below 1
		}
		updateInput.MaxHitPoints = newMaxHP

		// If current HP is higher than new max HP, reduce it to match
		if updateInput.CurrentHitPoints > newMaxHP {
			updateInput.CurrentHitPoints = newMaxHP
		}
	}

	// Ensure current HP doesn't go below -10
	if updateInput.CurrentHitPoints < -10 {
		updateInput.CurrentHitPoints = -10
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

	// Determine which service to use based on the character's class
	var classData interface{}
	var enrichErr error

	switch character.Class {
	case "Fighter":
		// Get fighter-specific data
		fighterData, err := c.fighterService.GetAllFighterLevelData(r.Context())
		if err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		// Enrich character with fighter-specific data for its current level
		enrichErr = c.fighterService.EnrichCharacterWithFighterData(r.Context(), character)

		// Create a response with both the full level progression and current abilities
		classData = map[string]interface{}{
			"class_type": "Fighter",
			"level_data": fighterData,
			"current_level_data": map[string]interface{}{
				"level":            character.Level,
				"hit_dice":         character.HitDice,
				"saving_throw":     character.SavingThrow,
				"fighting_ability": character.FightingAbility,
				"abilities":        character.Abilities,
			},
		}
	case "Magician", "Wizard": // Handle both terms in case you use them interchangeably
		// Get magician-specific data
		magicianData, err := c.magicianService.GetAllMagicianLevelData(r.Context())
		if err != nil {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
			return
		}

		// Enrich character with magician-specific data for its current level
		enrichErr = c.magicianService.EnrichCharacterWithMagicianData(r.Context(), character)

		// Create a response with both the full level progression and current abilities
		classData = map[string]interface{}{
			"class_type": "Magician",
			"level_data": magicianData,
			"current_level_data": map[string]interface{}{
				"level":           character.Level,
				"hit_dice":        character.HitDice,
				"saving_throw":    character.SavingThrow,
				"casting_ability": character.CastingAbility,
				"spell_slots":     character.SpellSlots,
				"abilities":       character.Abilities,
			},
		}
	default:
		// For other classes or if no specific handler exists
		classData = map[string]interface{}{
			"class_type": character.Class,
			"message":    "No specific class data available for this character class",
		}
	}

	if enrichErr != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(enrichErr))
		return
	}

	// Return the class data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classData)
}
