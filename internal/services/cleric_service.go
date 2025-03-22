package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ClericService handles business logic for cleric class characters
type ClericService struct {
	clericDataRepo repositories.ClericRepository
}

// NewClericService creates a new cleric service instance
func NewClericService(clericDataRepo repositories.ClericRepository) *ClericService {
	return &ClericService{
		clericDataRepo: clericDataRepo,
	}
}

func (s *ClericService) EnrichCharacterWithClericData(ctx context.Context, character *models.Character) error {
	if character.Class != "Cleric" {
		return nil
	}

	clericData, err := s.clericDataRepo.GetClericClassData(ctx, character.Level)
	if err != nil {
		return err
	}

	// First calculate derived stats from ability scores
	character.CalculateDerivedStats()

	// THEN set class data that should override calculations
	character.HitDice = clericData.HitDice
	character.SavingThrow = clericData.SavingThrow
	character.FightingAbility = clericData.FightingAbility
	character.TurningAbility = clericData.TurningAbility

	// Set Cleric-specific save bonuses AFTER calculating derived stats
	character.DeathSaveBonus = 2
	character.TransformationSaveBonus = 0
	character.DeviceSaveBonus = 0
	character.SorcerySaveBonus = 2
	character.AvoidanceSaveBonus = 0

	// Hardcode spell slots based on level until sqlc is regenerated
	character.SpellSlots = s.getClericSpellSlots(character.Level)

	character.Abilities = s.GetAvailableClericAbilities(character.Level)

	return nil
}

// GetAvailableClericAbilities returns cleric abilities available at the given level
func (s *ClericService) GetAvailableClericAbilities(level int) []*models.ClericAbility {
	allAbilities := models.GetClericAbilities()
	available := make([]*models.ClericAbility, 0)

	for _, ability := range allAbilities {
		if ability.MinLevel <= level {
			available = append(available, ability)
		}
	}

	return available
}

// getClericSpellSlots returns spell slots based on level
func (s *ClericService) getClericSpellSlots(level int) map[string]int {
	slots := make(map[string]int)

	switch level {
	case 1:
		slots["level1"] = 0
	case 2, 3:
		slots["level1"] = level - 1
	case 4, 5:
		slots["level1"] = 2
		slots["level2"] = level - 3
	case 6, 7:
		slots["level1"] = 3
		slots["level2"] = level - 4
		slots["level3"] = 1
	case 8:
		slots["level1"] = 3
		slots["level2"] = 3
		slots["level3"] = 2
	case 9:
		slots["level1"] = 4
		slots["level2"] = 3
		slots["level3"] = 2
		slots["level4"] = 1
	case 10:
		slots["level1"] = 4
		slots["level2"] = 3
		slots["level3"] = 3
		slots["level4"] = 2
	case 11:
		slots["level1"] = 4
		slots["level2"] = 4
		slots["level3"] = 3
		slots["level4"] = 2
		slots["level5"] = 1
	case 12:
		slots["level1"] = 5
		slots["level2"] = 4
		slots["level3"] = 3
		slots["level4"] = 3
		slots["level5"] = 1
	}

	return slots
}

// GetExperienceForNextLevel returns the experience needed for the next level
func (s *ClericService) GetExperienceForNextLevel(ctx context.Context, currentLevel int) (int, error) {
	// If at max level, return the current level's experience
	if currentLevel >= 12 {
		clericData, err := s.clericDataRepo.GetClericClassData(ctx, 12)
		if err != nil {
			return 0, err
		}
		return clericData.ExperiencePoints, nil
	}

	// Otherwise, get the next level's experience requirement
	nextLevel, err := s.clericDataRepo.GetNextClericLevel(ctx, currentLevel)
	if err != nil {
		return 0, err
	}

	return nextLevel.ExperiencePoints, nil
}

// GetAllClericLevelData returns all cleric level data
func (s *ClericService) GetAllClericLevelData(ctx context.Context) ([]*models.ClericClassData, error) {
	return s.clericDataRepo.ListClericClassData(ctx)
}
