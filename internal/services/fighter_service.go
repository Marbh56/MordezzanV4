package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// FighterService handles business logic for fighter class characters
type FighterService struct {
	fighterDataRepo repositories.FighterDataRepository
}

// NewFighterService creates a new fighter service instance
func NewFighterService(fighterDataRepo repositories.FighterDataRepository) *FighterService {
	return &FighterService{
		fighterDataRepo: fighterDataRepo,
	}
}

func (s *FighterService) EnrichCharacterWithFighterData(ctx context.Context, character *models.Character) error {
	if character.Class != "Fighter" {
		return nil
	}

	fighterData, err := s.fighterDataRepo.GetFighterClassData(ctx, character.Level)
	if err != nil {
		return err
	}

	// First calculate derived stats from ability scores
	character.CalculateDerivedStats()

	// THEN set class data that should override basic calculations
	character.HitDice = fighterData.HitDice
	character.SavingThrow = fighterData.SavingThrow
	character.FightingAbility = fighterData.FightingAbility

	// Set Fighter-specific save bonuses AFTER calculating derived stats
	character.DeathSaveBonus = 2
	character.TransformationSaveBonus = 2
	character.DeviceSaveBonus = 0
	character.SorcerySaveBonus = 0
	character.AvoidanceSaveBonus = 0

	character.Abilities = s.GetAvailableFighterAbilities(character.Level)

	return nil
}

func (s *FighterService) GetAvailableFighterAbilities(level int) []*models.FighterAbility {
	allAbilities := models.GetFighterAbilities()
	available := make([]*models.FighterAbility, 0)

	for _, ability := range allAbilities {
		if ability.MinLevel <= level {
			available = append(available, ability)
		}
	}

	return available
}

// GetExperienceForNextLevel returns the experience needed for the next level
func (s *FighterService) GetExperienceForNextLevel(ctx context.Context, currentLevel int) (int, error) {
	// If at max level, return the current level's experience
	if currentLevel >= 12 {
		fighterData, err := s.fighterDataRepo.GetFighterClassData(ctx, 12)
		if err != nil {
			return 0, err
		}
		return fighterData.ExperiencePoints, nil
	}

	// Otherwise, get the next level's experience requirement
	nextLevel, err := s.fighterDataRepo.GetNextFighterLevel(ctx, currentLevel)
	if err != nil {
		return 0, err
	}

	return nextLevel.ExperiencePoints, nil
}

func (s *FighterService) GetAllFighterLevelData(ctx context.Context) ([]*models.FighterClassData, error) {
	return s.fighterDataRepo.ListFighterClassData(ctx)
}
