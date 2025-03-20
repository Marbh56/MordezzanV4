package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// MagicianService handles business logic for magician class characters
type MagicianService struct {
	magicianDataRepo repositories.MagicianRepository
}

// NewMagicianService creates a new magician service instance
func NewMagicianService(magicianDataRepo repositories.MagicianRepository) *MagicianService {
	return &MagicianService{
		magicianDataRepo: magicianDataRepo,
	}
}

func (s *MagicianService) EnrichCharacterWithMagicianData(ctx context.Context, character *models.Character) error {
	if character.Class != "Magician" {
		return nil
	}

	magicianData, err := s.magicianDataRepo.GetMagicianClassData(ctx, character.Level)
	if err != nil {
		return err
	}

	// First calculate derived stats from ability scores
	character.CalculateDerivedStats()

	// THEN set class data that should override calculations
	character.HitDice = magicianData.HitDice
	character.SavingThrow = magicianData.SavingThrow
	character.FightingAbility = magicianData.FightingAbility
	character.CastingAbility = magicianData.CastingAbility

	// Set Magician-specific save bonuses AFTER calculating derived stats
	character.DeathSaveBonus = 0
	character.TransformationSaveBonus = 0
	character.DeviceSaveBonus = 2
	character.SorcerySaveBonus = 2
	character.AvoidanceSaveBonus = 0

	character.SpellSlots = map[string]int{
		"level1": magicianData.SpellSlotsLevel1,
		"level2": magicianData.SpellSlotsLevel2,
		"level3": magicianData.SpellSlotsLevel3,
		"level4": magicianData.SpellSlotsLevel4,
		"level5": magicianData.SpellSlotsLevel5,
		"level6": magicianData.SpellSlotsLevel6,
	}

	character.Abilities = s.GetAvailableMagicianAbilities(character.Level)

	return nil
}

// GetAvailableMagicianAbilities returns magician abilities available at the given level
func (s *MagicianService) GetAvailableMagicianAbilities(level int) []*models.MagicianAbility {
	allAbilities := models.GetMagicianAbilities()
	available := make([]*models.MagicianAbility, 0)

	for _, ability := range allAbilities {
		if ability.MinLevel <= level {
			available = append(available, ability)
		}
	}

	return available
}

// GetExperienceForNextLevel returns the experience needed for the next level
func (s *MagicianService) GetExperienceForNextLevel(ctx context.Context, currentLevel int) (int, error) {
	// If at max level, return the current level's experience
	if currentLevel >= 12 {
		magicianData, err := s.magicianDataRepo.GetMagicianClassData(ctx, 12)
		if err != nil {
			return 0, err
		}
		return magicianData.ExperiencePoints, nil
	}

	// Otherwise, get the next level's experience requirement
	nextLevel, err := s.magicianDataRepo.GetNextMagicianLevel(ctx, currentLevel)
	if err != nil {
		return 0, err
	}

	return nextLevel.ExperiencePoints, nil
}

// GetAllMagicianLevelData returns all magician level data
func (s *MagicianService) GetAllMagicianLevelData(ctx context.Context) ([]*models.MagicianClassData, error) {
	return s.magicianDataRepo.ListMagicianClassData(ctx)
}
