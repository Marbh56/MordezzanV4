package services

import (
	"context"
	"fmt"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
	"strconv"
	"strings"
)

type ClassService struct {
	classRepo          repositories.ClassRepository
	inventoryRepo      repositories.InventoryRepository
	armorRepo          repositories.ArmorRepository
	encumbranceService *EncumbranceService
}

func NewClassService(
	classRepo repositories.ClassRepository,
	inventoryRepo repositories.InventoryRepository,
	armorRepo repositories.ArmorRepository,
) *ClassService {
	return &ClassService{
		classRepo:     classRepo,
		inventoryRepo: inventoryRepo,
		armorRepo:     armorRepo,
	}
}

func (s *ClassService) GetAllClassLevelData(ctx context.Context, className string) ([]*models.ClassData, error) {
	// Simply call through to the repository method
	return s.classRepo.GetAllClassData(ctx, className)
}

func (s *ClassService) SetEncumbranceService(encumbranceService *EncumbranceService) {
	s.encumbranceService = encumbranceService
}

func (s *ClassService) applyAgileBonus(ctx context.Context, character *models.Character) error {
	fmt.Printf("DEBUG: Applying agile bonus for %s (ID: %d)\n", character.Name, character.ID)

	if s.encumbranceService == nil || s.inventoryRepo == nil {
		fmt.Printf("DEBUG: Encumbrance service or inventory repo is nil\n")
		return nil
	}

	encumbranceDetails, err := s.encumbranceService.GetCharacterEncumbrance(ctx, character.ID)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get encumbrance details: %v\n", err)
		return fmt.Errorf("failed to get encumbrance details: %v", err)
	}

	fmt.Printf("DEBUG: Encumbrance status - HeavyEncumbered: %v\n",
		encumbranceDetails.Status.HeavyEncumbered)

	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, character.ID)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get character inventory: %v\n", err)
		return fmt.Errorf("failed to get character inventory: %v", err)
	}

	wearingArmor := false
	for _, item := range inventory.Items {
		if item.ItemType == "armor" && item.IsEquipped {
			wearingArmor = true
			fmt.Printf("DEBUG: Character is wearing armor: %s\n", item.ItemID)
			break
		}
	}

	fmt.Printf("DEBUG: Character is wearing armor: %v\n", wearingArmor)

	if !wearingArmor && !encumbranceDetails.Status.HeavyEncumbered {
		fmt.Printf("DEBUG: Adding +1 DefenceAdjustment for agile bonus\n")
		character.DefenceAdjustment += 1
	} else {
		fmt.Printf("DEBUG: Not applying agile bonus due to armor or encumbrance\n")
	}

	fmt.Printf("DEBUG: Final DefenceAdjustment: %d\n", character.DefenceAdjustment)
	return nil
}

func (s *ClassService) applyRunAbility(ctx context.Context, character *models.Character) error {
	// Only proceed if we have the necessary services
	if s.inventoryRepo == nil {
		return nil
	}

	// Get inventory to check armor
	inventory, err := s.inventoryRepo.GetInventoryByCharacter(ctx, character.ID)
	if err != nil {
		return fmt.Errorf("failed to get character inventory: %v", err)
	}

	// Check if wearing heavy/medium armor
	wearingHeavyArmor := false
	for _, item := range inventory.Items {
		if item.ItemType == "armor" && item.IsEquipped {
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

	// Apply RUN ability - base 50 MV when lightly armored or unarmored
	if !wearingHeavyArmor {
		character.MovementRate = 50
	}

	return nil
}

func (s *ClassService) applyExtraStr(ctx context.Context, character *models.Character) error {
	// Parse the current percentage value
	var currentPercent int
	n, err := fmt.Sscanf(character.ExtraStrengthFeat, "%d%%", &currentPercent)

	// Check for scanning errors
	if err != nil {
		return fmt.Errorf("failed to parse strength feat percentage: %v", err)
	}

	// Check if we scanned the expected number of items
	if n != 1 {
		return fmt.Errorf("unexpected format for strength feat: %s", character.ExtraStrengthFeat)
	}

	// Apply the 8% bonus
	character.ExtraStrengthFeat = fmt.Sprintf("%d%%", currentPercent+8)

	return nil
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

	// Get class abilities using the class-specific method
	var classAbilities []*models.ClassAbility

	// Use a direct call to the appropriate class-specific ability method based on character's class
	switch character.Class {
	case "Barbarian":
		classAbilities, err = s.classRepo.GetBarbarianAbilities(ctx, character.Level)
	case "Berserker":
		classAbilities, err = s.classRepo.GetBerserkerAbilities(ctx, character.Level)
	case "Bard":
		classAbilities, err = s.classRepo.GetBardAbilities(ctx, character.Level)
	case "Cataphract":
		classAbilities, err = s.classRepo.GetCataphractAbilities(ctx, character.Level)
	case "Cleric":
		classAbilities, err = s.classRepo.GetClericAbilities(ctx, character.Level)
	case "Cryomancer":
		classAbilities, err = s.classRepo.GetCryomancerAbilities(ctx, character.Level)
	case "Druid":
		classAbilities, err = s.classRepo.GetDruidAbilities(ctx, character.Level)
	case "Fighter":
		classAbilities, err = s.classRepo.GetFighterAbilities(ctx, character.Level)
	case "Huntsman":
		classAbilities, err = s.classRepo.GetHuntsmanAbilities(ctx, character.Level)
	case "Illusionist":
		classAbilities, err = s.classRepo.GetIllusionistAbilities(ctx, character.Level)
	case "Legerdemainist":
		classAbilities, err = s.classRepo.GetLegerdemainistAbilities(ctx, character.Level)
	case "Magician":
		classAbilities, err = s.classRepo.GetMagicianAbilities(ctx, character.Level)
	case "Monk":
		classAbilities, err = s.classRepo.GetMonkAbilities(ctx, character.Level)
	case "Necromancer":
		classAbilities, err = s.classRepo.GetNecromancerAbilities(ctx, character.Level)
	case "Paladin":
		classAbilities, err = s.classRepo.GetPaladinAbilities(ctx, character.Level)
	case "Priest":
		classAbilities, err = s.classRepo.GetPriestAbilities(ctx, character.Level)
	case "Purloiner":
		classAbilities, err = s.classRepo.GetPurloinerAbilities(ctx, character.Level)
	case "Pyromancer":
		classAbilities, err = s.classRepo.GetPyromancerAbilities(ctx, character.Level)
	case "Ranger":
		classAbilities, err = s.classRepo.GetRangerAbilities(ctx, character.Level)
	case "Runegraver":
		classAbilities, err = s.classRepo.GetRunegraverAbilities(ctx, character.Level)
	case "Scout":
		classAbilities, err = s.classRepo.GetScoutAbilities(ctx, character.Level)
	case "Shaman":
		classAbilities, err = s.classRepo.GetShamanAbilities(ctx, character.Level)
	case "Thief":
		classAbilities, err = s.classRepo.GetThiefAbilities(ctx, character.Level)
	case "Warlock":
		classAbilities, err = s.classRepo.GetWarlockAbilities(ctx, character.Level)
	case "Witch":
		classAbilities, err = s.classRepo.GetWitchAbilities(ctx, character.Level)
	default:
		// Fallback for any classes not specifically handled
		classAbilities, err = s.classRepo.GetClassAbilitiesByLevel(ctx, character.Class, character.Level)
	}

	// Check for errors in ability retrieval, but don't fail the entire function
	if err != nil {
		// Log the error but continue with the rest of the function
		fmt.Printf("failed to fetch %s abilities: %v\n", character.Class, err)
	}

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
		if err := s.applyAgileBonus(ctx, character); err != nil {
			// Log error but continue
			fmt.Printf("failed to apply agile bonus: %v\n", err)
		}

	case "Barbarian":
		// Barbarian-specific save bonuses
		character.AvoidanceSaveBonus = 2
		character.DeviceSaveBonus = 2
		character.SorcerySaveBonus = 2
		character.DeathSaveBonus = 2
		character.TransformationSaveBonus = 2
		character.SurpriseChance = 1

		if err := s.applyExtraStr(ctx, character); err != nil {
			// Log error but continue
			fmt.Printf("failed to apply extra strength feat: %v\n", err)
		}

		if err := s.applyAgileBonus(ctx, character); err != nil {
			// Log error but continue
			fmt.Printf("failed to apply agile bonus: %v\n", err)
		}

		// Apply RUN ability - base 50 MV when lightly armored or unarmored
		if err := s.applyRunAbility(ctx, character); err != nil {
			// Log error but continue
			fmt.Printf("failed to apply run ability: %v\n", err)
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

		if err := s.applyAgileBonus(ctx, character); err != nil {
			fmt.Printf("failed to apply agile bonus: %v\n", err)
		}

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

	// Add class abilities to character.Abilities if we successfully retrieved them
	if classAbilities != nil && len(classAbilities) > 0 {
		// If character.Abilities is nil, initialize it
		if character.Abilities == nil {
			character.Abilities = map[string]interface{}{
				"class_abilities": classAbilities,
			}
		} else if abilitiesMap, ok := character.Abilities.(map[string]interface{}); ok {
			// If it's already a map, add the class abilities to it
			abilitiesMap["class_abilities"] = classAbilities
			character.Abilities = abilitiesMap
		} else {
			// If it's already set but not a map, create a new map with both existing abilities and class abilities
			character.Abilities = map[string]interface{}{
				"class_abilities":    classAbilities,
				"original_abilities": character.Abilities,
			}
		}
	}

	return nil
}

// GetExperienceForNextLevel returns the XP needed for the next level
func (s *ClassService) GetExperienceForNextLevel(ctx context.Context, class string, currentLevel int) (int, error) {
	// Get all class level data
	levelData, err := s.classRepo.GetAllClassData(ctx, class)
	if err != nil {
		return 0, err
	}

	// Find the XP for the next level
	for _, data := range levelData {
		if data.Level == currentLevel+1 {
			return data.ExperiencePoints, nil
		}
	}

	// If no next level found (max level), return -1
	return -1, nil
}

func (s *ClassService) GetClassAbilitiesByLevel(ctx context.Context, class string, level int) ([]*models.ClassAbility, error) {
	return s.classRepo.GetClassAbilitiesByLevel(ctx, class, level)
}

// CalculateMaxSpellsPerLevel determines how many spells a character can know per level
func (s *ClassService) CalculateMaxSpellsPerLevel(ctx context.Context, characterID int64) (map[int]int, error) {
	// Get character
	character, err := s.classRepo.GetCharacterClassInfo(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Get casting stat modifier based on class
	var modifier int
	switch character.Class {
	case "Magician", "Illusionist", "Necromancer", "Pyromancer", "Cryo-mancer":
		// Intelligence-based casters
		modifier = character.RangedModifier // Use the pre-calculated modifier
	case "Cleric", "Druid", "Witch":
		// Wisdom-based casters
		modifier = character.WillpowerModifier // Use the pre-calculated modifier
	default:
		// Non-casting class
		return make(map[int]int), nil
	}

	// Get level data to determine base number of spells
	levelData, err := s.classRepo.GetClassData(ctx, character.Class, character.Level)
	if err != nil {
		return nil, err
	}

	// Calculate max spell level based on available spell slots
	maxSpellLevel := 0
	for level, count := range levelData.SpellSlots {
		if count > 0 {
			// Extract level number from key "level1", "level2", etc.
			levelNum, err := strconv.Atoi(level[5:])
			if err == nil && levelNum > maxSpellLevel {
				maxSpellLevel = levelNum
			}
		}
	}

	// Calculate max spells per level
	result := make(map[int]int)
	for level := 1; level <= maxSpellLevel; level++ {
		// Base formula: level + modifier with a minimum of 1 per level
		maxSpells := level + modifier
		if maxSpells < 1 {
			maxSpells = 1
		}
		result[level] = maxSpells
	}

	return result, nil
}

// ParseSpellSlots converts the spell slots string to a map
func (s *ClassService) ParseSpellSlots(slotsStr string) (map[string]int, error) {
	result := make(map[string]int)
	if slotsStr == "" {
		return result, nil
	}

	slots := strings.Split(slotsStr, ",")
	for i, slot := range slots {
		count, err := strconv.Atoi(strings.TrimSpace(slot))
		if err != nil {
			return nil, fmt.Errorf("invalid spell slot count: %s", slot)
		}
		result[fmt.Sprintf("level_%d", i+1)] = count
	}

	return result, nil
}
