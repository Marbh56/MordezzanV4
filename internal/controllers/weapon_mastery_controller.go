package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type WeaponMasteryController struct {
	weaponMasteryRepo repositories.WeaponMasteryRepository
	characterRepo     repositories.CharacterRepository
	weaponRepo        repositories.WeaponRepository
}

func NewWeaponMasteryController(
	weaponMasteryRepo repositories.WeaponMasteryRepository,
	characterRepo repositories.CharacterRepository,
	weaponRepo repositories.WeaponRepository,
) *WeaponMasteryController {
	return &WeaponMasteryController{
		weaponMasteryRepo: weaponMasteryRepo,
		characterRepo:     characterRepo,
		weaponRepo:        weaponRepo,
	}
}

func (c *WeaponMasteryController) GetWeaponMasteriesByCharacter(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	masteries, err := c.weaponMasteryRepo.GetWeaponMasteriesByCharacter(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masteries)
}

func (c *WeaponMasteryController) AddWeaponMastery(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	characterID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	var input struct {
		WeaponID       int64  `json:"weapon_id,omitempty"`
		WeaponBaseName string `json:"weapon_base_name,omitempty"`
		MasteryLevel   string `json:"mastery_level"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
		return
	}

	// If weapon_id is provided, get the base name
	if input.WeaponID > 0 && input.WeaponBaseName == "" {
		baseName, err := c.weaponMasteryRepo.GetWeaponBaseNameFromID(r.Context(), input.WeaponID)
		if err != nil {
			apperrors.HandleError(w, err)
			return
		}
		input.WeaponBaseName = baseName
	}

	if input.WeaponBaseName == "" {
		apperrors.HandleError(w, apperrors.NewBadRequest("Weapon base name is required"))
		return
	}

	// Check if the character can have more masteries
	character, err := c.characterRepo.GetCharacter(r.Context(), characterID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", characterID))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	// Check if the class can have weapon mastery
	if !models.CanHaveWeaponMastery(character.Class) {
		apperrors.HandleError(w, apperrors.NewBadRequest("Character class does not have weapon mastery ability"))
		return
	}

	// Calculate available slots
	availableSlots := models.GetAvailableMasterySlots(character.Class, character.Level)

	// Get current mastery count
	masteredCount, err := c.weaponMasteryRepo.CountWeaponMasteries(r.Context(), characterID, "mastered")
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	grandMasteryCount, err := c.weaponMasteryRepo.CountWeaponMasteries(r.Context(), characterID, "grand_mastery")
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	totalMasteries := masteredCount + grandMasteryCount

	if totalMasteries >= availableSlots {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Character can only have %d weapon masteries at level %d", availableSlots, character.Level)))
		return
	}

	// Check grand mastery requirements
	if input.MasteryLevel == "grand_mastery" {
		if character.Level < 4 {
			apperrors.HandleError(w, apperrors.NewBadRequest("Grand mastery is only available at level 4 or higher"))
			return
		}

		if grandMasteryCount >= 1 {
			apperrors.HandleError(w, apperrors.NewBadRequest("Character can only have one grand mastery weapon"))
			return
		}
	}

	// Create the weapon mastery
	masteryInput := &models.AddWeaponMasteryInput{
		CharacterID:    characterID,
		WeaponBaseName: input.WeaponBaseName,
		MasteryLevel:   input.MasteryLevel,
	}

	id, err := c.weaponMasteryRepo.AddWeaponMastery(r.Context(), masteryInput)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	mastery, err := c.weaponMasteryRepo.GetWeaponMasteryByID(r.Context(), id)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mastery)
}

func (c *WeaponMasteryController) UpdateWeaponMastery(w http.ResponseWriter, r *http.Request) {
	characterIDParam := chi.URLParam(r, "id")
	characterID, err := strconv.ParseInt(characterIDParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", characterIDParam)))
		return
	}

	weaponBaseNameParam := chi.URLParam(r, "weaponBaseName")
	if weaponBaseNameParam == "" {
		apperrors.HandleError(w, apperrors.NewBadRequest("Weapon base name is required"))
		return
	}

	var input models.UpdateWeaponMasteryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest("Invalid request body"))
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

	// If upgrading to grand mastery, check requirements
	if input.MasteryLevel == "grand_mastery" {
		character, err := c.characterRepo.GetCharacter(r.Context(), characterID)
		if err != nil {
			if errors.Is(err, apperrors.ErrNotFound) {
				apperrors.HandleError(w, apperrors.NewNotFound("character", characterID))
			} else {
				apperrors.HandleError(w, apperrors.NewInternalError(err))
			}
			return
		}

		if character.Level < 4 {
			apperrors.HandleError(w, apperrors.NewBadRequest("Grand mastery is only available at level 4 or higher"))
			return
		}
	}

	err = c.weaponMasteryRepo.UpdateWeaponMastery(r.Context(), characterID, weaponBaseNameParam, &input)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	// Get updated masteries for response
	masteries, err := c.weaponMasteryRepo.GetWeaponMasteriesByCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masteries)
}

func (c *WeaponMasteryController) DeleteWeaponMastery(w http.ResponseWriter, r *http.Request) {
	characterIDParam := chi.URLParam(r, "id")
	characterID, err := strconv.ParseInt(characterIDParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", characterIDParam)))
		return
	}

	weaponBaseNameParam := chi.URLParam(r, "weaponBaseName")
	if weaponBaseNameParam == "" {
		apperrors.HandleError(w, apperrors.NewBadRequest("Weapon base name is required"))
		return
	}

	err = c.weaponMasteryRepo.DeleteWeaponMastery(r.Context(), characterID, weaponBaseNameParam)
	if err != nil {
		apperrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *WeaponMasteryController) GetAvailableWeaponsForMastery(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	characterID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewBadRequest(fmt.Sprintf("Invalid character ID: %s", idParam)))
		return
	}

	// Check character exists and get class/level
	character, err := c.characterRepo.GetCharacter(r.Context(), characterID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			apperrors.HandleError(w, apperrors.NewNotFound("character", characterID))
		} else {
			apperrors.HandleError(w, apperrors.NewInternalError(err))
		}
		return
	}

	// Check if class has weapon mastery
	hasWeaponMastery := models.CanHaveWeaponMastery(character.Class)
	if !hasWeaponMastery {
		apperrors.HandleError(w, apperrors.NewBadRequest("Character class does not have weapon mastery ability"))
		return
	}

	// Get all weapons
	weapons, err := c.weaponRepo.ListWeapons(r.Context())
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Get character's current masteries
	masteries, err := c.weaponMasteryRepo.GetWeaponMasteriesByCharacter(r.Context(), characterID)
	if err != nil {
		apperrors.HandleError(w, apperrors.NewInternalError(err))
		return
	}

	// Create a map of base weapon names that are already mastered
	masteredBaseNames := make(map[string]bool)
	for _, mastery := range masteries {
		masteredBaseNames[mastery.WeaponBaseName] = true
	}

	// Group weapons by their base name
	weaponsByBase := make(map[string][]*models.Weapon)
	baseWeapons := make([]*models.Weapon, 0)

	for _, weapon := range weapons {
		// Extract base name
		baseName := extractBaseWeaponName(weapon.Name)

		// Add to the group
		if _, exists := weaponsByBase[baseName]; !exists {
			weaponsByBase[baseName] = make([]*models.Weapon, 0)

			// Create a base weapon entry for the dropdown
			baseWeapon := *weapon // Copy the struct
			baseWeapon.Name = baseName
			baseWeapons = append(baseWeapons, &baseWeapon)
		}

		weaponsByBase[baseName] = append(weaponsByBase[baseName], weapon)
	}

	// Filter out already mastered base weapons
	availableBaseWeapons := make([]*models.Weapon, 0)
	for _, weapon := range baseWeapons {
		if !masteredBaseNames[weapon.Name] {
			availableBaseWeapons = append(availableBaseWeapons, weapon)
		}
	}

	// Calculate available masteries and slots
	availableSlots := models.GetAvailableMasterySlots(character.Class, character.Level)
	grandMasteryAvailable := character.Level >= 4
	grandMasteryUsed := false

	for _, mastery := range masteries {
		if mastery.MasteryLevel == "grand_mastery" {
			grandMasteryUsed = true
			break
		}
	}

	response := map[string]interface{}{
		"available_weapons": availableBaseWeapons,
		"current_masteries": masteries,
		"total_slots":       availableSlots,
		"used_slots":        len(masteries),
		"can_grand_master":  grandMasteryAvailable && !grandMasteryUsed,
		"character_level":   character.Level,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to extract base weapon name
func extractBaseWeaponName(name string) string {
	// Remove "+X" suffix
	if idx := strings.Index(name, " +"); idx != -1 {
		name = name[:idx]
	}

	// Remove magical affixes
	suffixes := []string{
		" of Slaying",
		" of Fire",
		" of Frost",
		" of Lightning",
		" of Venom",
		" of Speed",
		" of Accuracy",
		" of Power",
	}

	for _, suffix := range suffixes {
		if idx := strings.Index(name, suffix); idx != -1 {
			name = name[:idx]
			break
		}
	}

	return strings.TrimSpace(name)
}
