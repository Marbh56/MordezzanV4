// internal/services/weapon_mastery_service.go

package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"strings"
)

// WeaponMasteryService handles business logic related to weapon masteries
type WeaponMasteryService struct {
	weaponMasteryRepo repositories.WeaponMasteryRepository
	characterRepo     repositories.CharacterRepository
	weaponRepo        repositories.WeaponRepository
}

// NewWeaponMasteryService creates a new weapon mastery service
func NewWeaponMasteryService(
	weaponMasteryRepo repositories.WeaponMasteryRepository,
	characterRepo repositories.CharacterRepository,
	weaponRepo repositories.WeaponRepository,
) *WeaponMasteryService {
	return &WeaponMasteryService{
		weaponMasteryRepo: weaponMasteryRepo,
		characterRepo:     characterRepo,
		weaponRepo:        weaponRepo,
	}
}

// GetAvailableWeaponMasteriesData retrieves weapon mastery data for a character
func (s *WeaponMasteryService) GetAvailableWeaponMasteriesData(ctx context.Context, characterID int64) (map[string]interface{}, error) {
	// Get character data to check class and level
	character, err := s.characterRepo.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Check if character class has weapon mastery capabilities
	if !s.hasWeaponMasteryAbility(character.Class) {
		return map[string]interface{}{
			"available_weapons": []interface{}{},
			"current_masteries": []interface{}{},
			"total_slots":       0,
			"used_slots":        0,
			"can_grand_master":  false,
			"character_level":   character.Level,
		}, nil
	}

	// Calculate the total mastery slots based on character level and class
	totalSlots := s.calculateMasterySlots(character.Class, character.Level)

	// Fetch the character's current weapon masteries
	currentMasteries, err := s.weaponMasteryRepo.GetWeaponMasteriesByCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Count used slots and check for grand mastery
	usedSlots := 0
	hasGrandMastery := false
	for _, mastery := range currentMasteries {
		if mastery.MasteryLevel == "grand_mastery" {
			usedSlots += 2 // Grand mastery takes two slots
			hasGrandMastery = true
		} else {
			usedSlots++ // Regular mastery takes one slot
		}
	}

	// Determine if the character can have grand mastery
	// Usually requires level 4+ and appropriate class
	canGrandMaster := character.Level >= 4 && s.canHaveGrandMastery(character.Class)

	// Get all available weapons the character could master

	weapons, err := s.weaponRepo.ListWeapons(ctx)
	if err != nil {
		return nil, err
	}
	availableWeapons := s.extractBaseWeapons(weapons)

	// Filter out weapons the character has already mastered
	filteredWeapons := s.filterMasteredWeapons(availableWeapons, currentMasteries)

	// Prepare the response
	response := map[string]interface{}{
		"available_weapons": filteredWeapons,
		"current_masteries": currentMasteries,
		"total_slots":       totalSlots,
		"used_slots":        usedSlots,
		"can_grand_master":  canGrandMaster && !hasGrandMastery, // Can only have one grand mastery
		"character_level":   character.Level,
	}

	return response, nil
}

// Helper function to determine if a character class has weapon mastery
func (s *WeaponMasteryService) hasWeaponMasteryAbility(class string) bool {
	// Classes that don't have weapon mastery in Hyperborea 3E
	noWeaponMasteryClasses := map[string]bool{
		"Magician":    true,
		"Cryomancer":  true,
		"Illusionist": true,
		"Necromancer": true,
		"Pyromancer":  true,
		"Witch":       true,
		// Add other spellcasting classes that don't get weapon mastery
	}

	return !noWeaponMasteryClasses[class]
}

// Helper function to calculate the number of mastery slots based on character level and class
func (s *WeaponMasteryService) calculateMasterySlots(class string, level int) int {
	// Base slots depend on class
	baseSlots := 1

	// Fighter-type classes might get more starting slots
	if class == "Fighter" || class == "Ranger" || class == "Paladin" ||
		class == "Barbarian" || class == "Berserker" || class == "Cataphract" {
		baseSlots = 2
	}

	// Add slots at certain level thresholds
	additionalSlots := 0
	if level >= 4 {
		additionalSlots++
	}
	if level >= 8 {
		additionalSlots++
	}
	if level >= 12 {
		additionalSlots++
	}

	return baseSlots + additionalSlots
}

// Helper function to determine if a character class can have grand mastery
func (s *WeaponMasteryService) canHaveGrandMastery(class string) bool {
	// Only certain classes can achieve grand mastery
	grandMasteryClasses := map[string]bool{
		"Fighter":    true,
		"Ranger":     true,
		"Paladin":    true,
		"Barbarian":  true,
		"Berserker":  true,
		"Cataphract": true,
		// Add other classes that can achieve grand mastery
	}

	return grandMasteryClasses[class]
}

// Helper function to filter out weapons the character has already mastered
func (s *WeaponMasteryService) filterMasteredWeapons(available []models.WeaponBase, mastered []*models.WeaponMastery) []models.WeaponBase {
	// Create a map of mastered weapon base names for quick lookup
	masteredMap := make(map[string]bool)
	for _, mastery := range mastered {
		masteredMap[mastery.WeaponBaseName] = true
	}

	// Filter the available weapons
	var filtered []models.WeaponBase
	for _, weapon := range available {
		if !masteredMap[weapon.Name] {
			filtered = append(filtered, weapon)
		}
	}

	return filtered
}

func (s *WeaponMasteryService) extractBaseWeapons(weapons []*models.Weapon) []models.WeaponBase {
	// Keep track of base names to avoid duplicates
	baseNames := make(map[string]bool)
	result := make([]models.WeaponBase, 0)

	for _, weapon := range weapons {
		baseName := s.extractBaseWeaponName(weapon.Name)

		// Skip if we've already added this base weapon
		if baseNames[baseName] {
			continue
		}

		baseNames[baseName] = true
		result = append(result, models.WeaponBase{
			Name:     baseName,
			Category: weapon.Category,
		})
	}

	return result
}

func (s *WeaponMasteryService) extractBaseWeaponName(name string) string {
	if idx := strings.Index(name, " +"); idx != -1 {
		name = name[:idx]
	}

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
