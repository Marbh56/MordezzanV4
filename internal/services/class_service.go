package services

import (
	"context"
	"fmt"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ClassService handles business logic for character classes
type ClassService struct {
	classRepo          repositories.ClassRepository
	inventoryRepo      repositories.InventoryRepository
	encumbranceService *EncumbranceService
}

// NewClassService creates a new class service
func NewClassService(
	classRepo repositories.ClassRepository,
	inventoryRepo repositories.InventoryRepository,
	encumbranceService *EncumbranceService,
) *ClassService {
	return &ClassService{
		classRepo:          classRepo,
		inventoryRepo:      inventoryRepo,
		encumbranceService: encumbranceService,
	}
}

// EnrichCharacterWithClassData applies class-specific data to a character
func (s *ClassService) EnrichCharacterWithClassData(ctx context.Context, character *models.Character) error {
	// Get class data for this character's class and level
	classData, err := s.classRepo.GetClassData(ctx, character.Class, character.Level)
	if err != nil {
		return fmt.Errorf("failed to get class data: %v", err)
	}

	// First calculate basic derived stats
	character.CalculateDerivedStats()

	// Apply common class data
	character.HitDice = classData.HitDice
	character.SavingThrow = classData.SavingThrow
	character.FightingAbility = classData.FightingAbility
	character.CastingAbility = classData.CastingAbility
	character.SpellSlots = classData.SpellSlots

	// Set default save bonuses
	character.DeathSaveBonus = 0
	character.TransformationSaveBonus = 0
	character.DeviceSaveBonus = 0
	character.SorcerySaveBonus = 0
	character.AvoidanceSaveBonus = 0

	// Apply class-specific data based on class type
	switch character.Class {
	// Warrior types
	case "Fighter":
		// Fighter-specific save bonuses
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2
		character.DeviceSaveBonus = 0
		character.SorcerySaveBonus = 0
		character.AvoidanceSaveBonus = 0

		// Parse the current percentage value
		var currentPercent int
		fmt.Sscanf(character.ExtraStrengthFeat, "%d%%", &currentPercent)

		// Apply the 8% bonus
		character.ExtraStrengthFeat = fmt.Sprintf("%d%%", currentPercent+8)

	case "Barbarian":
		// Barbarian-specific save bonuses
		character.AvoidanceSaveBonus = 2
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2
		character.SurpriseChance = 1

		// Parse the current percentage value
		var currentPercent int
		fmt.Sscanf(character.ExtraStrengthFeat, "%d%%", &currentPercent)

		// Apply the 8% bonus
		character.ExtraStrengthFeat = fmt.Sprintf("%d%%", currentPercent+8)

		// Check for Agile ability and RUN ability
		if s.encumbranceService != nil && s.inventoryRepo != nil {
			// Get encumbrance details
			encumbranceDetails, err := s.encumbranceService.GetCharacterEncumbrance(ctx, character.ID)
			if err == nil {
				// Get inventory to check armor
				inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, character.ID)
				if err == nil {
					// Check if wearing armor and what type
					wearingArmor := false
					wearingHeavyArmor := false

					for _, item := range inventory.Items {
						if item.ItemType == "armor" && item.IsEquipped {
							wearingArmor = true

							// Check armor weight class
							if s.armorRepo != nil {
								armor, armorErr := s.armorRepo.GetArmor(ctx, item.ItemID)
								if armorErr == nil {
									if armor.WeightClass == "Heavy" || armor.WeightClass == "Medium" {
										wearingHeavyArmor = true
										break
									}
								}
							}
						}
					}

					// Apply Agile bonus if unarmored and not heavily encumbered
					// Shield is allowed per the requirement
					if !wearingArmor && !encumbranceDetails.Status.HeavyEncumbered {
						character.DefenceAdjustment += 1
					}

					// Apply RUN ability - base 50 MV when lightly armored or unarmored
					if !wearingHeavyArmor {
						character.MovementRate = 50
					}
				}
			}
		}

	case "Berserker":
		naturalAC, _ := s.classRepo.GetBerserkerNaturalAC(ctx, character.Level)
		// Add to character struct or create a special field
		character.Abilities = map[string]interface{}{
			"natural_ac": naturalAC,
		}
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2

	case "Cataphract":
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2

	case "Huntsman":
		character.DeathSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	// Divine casters
	case "Cleric":
		turningAbility, _ := s.classRepo.GetClericTurningAbility(ctx, character.Level)
		character.TurningAbility = turningAbility
		character.TransformationSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Druid":
		character.TransformationSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Paladin":
		turningAbility, _ := s.classRepo.GetPaladinTurningAbility(ctx, character.Level)
		character.TurningAbility = turningAbility
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2

	case "Priest":
		character.TransformationSaveBonus = 2
		character.SorcerySaveBonus = 2

	// Arcane casters
	case "Magician":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Cryomancer":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Illusionist":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Necromancer":
		turningAbility, _ := s.classRepo.GetNecromancerTurningAbility(ctx, character.Level)
		character.TurningAbility = turningAbility
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Pyromancer":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Witch":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Warlock":
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2

	// Hybrid and special classes
	case "Ranger":
		specialSlots, _ := s.classRepo.GetSpecialClassSpellSlots(ctx, character.Class, character.Level)
		if specialSlots != nil {
			// Add special slots to character
			if character.SpellSlots == nil {
				character.SpellSlots = specialSlots
			} else {
				// Merge the maps
				for k, v := range specialSlots {
					character.SpellSlots[k] = v
				}
			}
		}
		character.DeathSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Runegraver":
		runesPerDay, _ := s.classRepo.GetRunegraverRunesPerDay(ctx, character.Level)
		if character.Abilities == nil {
			character.Abilities = map[string]interface{}{
				"runes_per_day": runesPerDay,
			}
		} else if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			abilities["runes_per_day"] = runesPerDay
			character.Abilities = abilities
		}
		character.DeathSaveBonus = 2
		character.DeviceSaveBonus = 2

	case "Shaman":
		specialSlots, _ := s.classRepo.GetSpecialClassSpellSlots(ctx, character.Class, character.Level)
		if specialSlots != nil {
			// Merge with existing spell slots
			if character.SpellSlots == nil {
				character.SpellSlots = specialSlots
			} else {
				for k, v := range specialSlots {
					character.SpellSlots[k] = v
				}
			}
		}
		character.TransformationSaveBonus = 2
		character.SorcerySaveBonus = 2

	case "Bard":
		specialSlots, _ := s.classRepo.GetSpecialClassSpellSlots(ctx, character.Class, character.Level)
		if specialSlots != nil {
			if character.SpellSlots == nil {
				character.SpellSlots = specialSlots
			} else {
				for k, v := range specialSlots {
					character.SpellSlots[k] = v
				}
			}
		}
		character.DeviceSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Monk":
		acBonus, _ := s.classRepo.GetMonkACBonus(ctx, character.Level)
		emptyHandDamage, _ := s.classRepo.GetMonkEmptyHandDamage(ctx, character.Level)

		if character.Abilities == nil {
			character.Abilities = map[string]interface{}{
				"ac_bonus":          acBonus,
				"empty_hand_damage": emptyHandDamage,
			}
		} else if abilities, ok := character.Abilities.(map[string]interface{}); ok {
			abilities["ac_bonus"] = acBonus
			abilities["empty_hand_damage"] = emptyHandDamage
			character.Abilities = abilities
		}
		character.DeathSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	// Rogue-types
	case "Thief":
		character.DeviceSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Assassin":
		character.DeathSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Legerdemainist":
		character.DeviceSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Purloiner":
		character.DeviceSaveBonus = 2
		character.AvoidanceSaveBonus = 2

	case "Scout":
		character.DeathSaveBonus = 2
		character.AvoidanceSaveBonus = 2
	}

	// Get class abilities for this level
	abilities, err := s.classRepo.GetClassAbilitiesByLevel(ctx, character.Class, character.Level)
	if err == nil && len(abilities) > 0 {
		// Store class abilities in character object
		// If character.Abilities is already a map, add to it, otherwise replace
		if character.Abilities == nil {
			character.Abilities = abilities
		} else if abilitiesMap, ok := character.Abilities.(map[string]interface{}); ok {
			abilitiesMap["class_abilities"] = abilities
			character.Abilities = abilitiesMap
		}
	}

	return nil
}

// GetExperienceForNextLevel returns the experience needed for the next level
func (s *ClassService) GetExperienceForNextLevel(ctx context.Context, className string, currentLevel int) (int, error) {
	// Get the next level's data
	nextLevel, err := s.classRepo.GetNextLevelData(ctx, className, currentLevel)
	if err != nil {
		return 0, err
	}

	return nextLevel.ExperiencePoints, nil
}

// GetAllClassLevelData returns all level data for a specific class
func (s *ClassService) GetAllClassLevelData(ctx context.Context, className string) ([]*models.ClassData, error) {
	return s.classRepo.GetAllClassData(ctx, className)
}

// GetClassDataByLevel returns class data for a specific class and level
func (s *ClassService) GetClassDataByLevel(ctx context.Context, className string, level int) (*models.ClassData, error) {
	return s.classRepo.GetClassData(ctx, className, level)
}

// GetClassAbilitiesByLevel returns abilities for a class at a specific level
func (s *ClassService) GetClassAbilitiesByLevel(ctx context.Context, className string, level int) ([]*models.ClassAbility, error) {
	return s.classRepo.GetClassAbilitiesByLevel(ctx, className, level)
}

// GetAllClassAbilities returns all abilities for a class
func (s *ClassService) GetAllClassAbilities(ctx context.Context, className string) ([]*models.ClassAbility, error) {
	return s.classRepo.GetClassAbilities(ctx, className)
}
