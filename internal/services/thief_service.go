package services

import (
	"context"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ThiefService handles business logic for thief class characters
type ThiefService struct {
	thiefDataRepo repositories.ThiefDataRepository
}

// NewThiefService creates a new thief service instance
func NewThiefService(thiefDataRepo repositories.ThiefDataRepository) *ThiefService {
	return &ThiefService{
		thiefDataRepo: thiefDataRepo,
	}
}

// EnrichCharacterWithThiefData adds thief class specific data to a character
func (s *ThiefService) EnrichCharacterWithThiefData(ctx context.Context, character *models.Character) error {
	if character.Class != "Thief" {
		return nil
	}

	thiefData, err := s.thiefDataRepo.GetThiefClassData(ctx, character.Level)
	if err != nil {
		return err
	}

	// First calculate derived stats from ability scores
	character.CalculateDerivedStats()

	// THEN set class data that should override basic calculations
	character.HitDice = thiefData.HitDice
	character.SavingThrow = thiefData.SavingThrow
	character.FightingAbility = thiefData.FightingAbility

	// Set Thief-specific save bonuses AFTER calculating derived stats
	character.DeathSaveBonus = 0
	character.TransformationSaveBonus = 0
	character.DeviceSaveBonus = 2
	character.SorcerySaveBonus = 0
	character.AvoidanceSaveBonus = 2

	// Get thief skills based on level
	thiefSkills := models.GetThiefSkillsByLevel(character.Level)

	// Add thief abilities and skills
	character.Abilities = map[string]interface{}{
		"thief_skills": thiefSkills,
		"abilities":    s.GetAvailableThiefAbilities(character.Level),
	}

	return nil
}

// GetAvailableThiefAbilities returns thief abilities available at the given level
func (s *ThiefService) GetAvailableThiefAbilities(level int) []*models.ThiefAbility {
	allAbilities := models.GetThiefAbilities()
	available := make([]*models.ThiefAbility, 0)

	for _, ability := range allAbilities {
		if ability.MinLevel <= level {
			available = append(available, ability)
		}
	}

	return available
}

// GetExperienceForNextLevel returns the experience needed for the next level
func (s *ThiefService) GetExperienceForNextLevel(ctx context.Context, currentLevel int) (int, error) {
	// If at max level, return the current level's experience
	if currentLevel >= 12 {
		thiefData, err := s.thiefDataRepo.GetThiefClassData(ctx, 12)
		if err != nil {
			return 0, err
		}
		return thiefData.ExperiencePoints, nil
	}

	// Otherwise, get the next level's experience requirement
	nextLevel, err := s.thiefDataRepo.GetNextThiefLevel(ctx, currentLevel)
	if err != nil {
		return 0, err
	}

	return nextLevel.ExperiencePoints, nil
}

// GetAllThiefLevelData returns all thief level data
func (s *ThiefService) GetAllThiefLevelData(ctx context.Context) ([]*models.ThiefClassData, error) {
	return s.thiefDataRepo.ListThiefClassData(ctx)
}
