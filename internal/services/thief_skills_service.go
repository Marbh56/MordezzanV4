package services

import (
	"context"
	"fmt"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	"mordezzanV4/internal/repositories"
)

// ThiefSkillsService handles business logic for thief skills
type ThiefSkillsService struct {
	thiefSkillsRepo repositories.ThiefSkillsRepository
}

// NewThiefSkillsService creates a new thief skills service
func NewThiefSkillsService(thiefSkillsRepo repositories.ThiefSkillsRepository) *ThiefSkillsService {
	return &ThiefSkillsService{
		thiefSkillsRepo: thiefSkillsRepo,
	}
}

// GetThiefSkillsForCharacter returns thief skills for a character based on class, level, and attributes
func (s *ThiefSkillsService) GetThiefSkillsForCharacter(
	ctx context.Context,
	charClass string,
	level int64,
	attributes map[string]int,
) ([]*models.ThiefSkillWithChance, error) {
	logger.Debug("Getting thief skills for character - Class: %s, Level: %d", charClass, level)

	// Get the effective thief level for this character
	effectiveLevel, err := s.thiefSkillsRepo.GetEffectiveThiefLevel(ctx, charClass, level)
	if err != nil {
		logger.Error("Error determining effective thief level: %v", err)
		return nil, fmt.Errorf("error determining effective thief level: %w", err)
	}

	// If the effective level is 0, this class doesn't have thief skills
	if effectiveLevel == 0 {
		logger.Debug("Character class %s has no thief skills", charClass)
		return []*models.ThiefSkillWithChance{}, nil
	}
	logger.Debug("Effective thief level for %s (level %d): %d", charClass, level, effectiveLevel)

	// Get thief skills for the effective level
	skills, err := s.thiefSkillsRepo.GetThiefSkillsByLevel(ctx, effectiveLevel)
	if err != nil {
		logger.Error("Error fetching thief skills: %v", err)
		return nil, fmt.Errorf("error fetching thief skills: %w", err)
	}
	logger.Debug("Found %d thief skills for level %d", len(skills), effectiveLevel)

	// Apply attribute bonuses
	skillsWithBonuses := s.thiefSkillsRepo.ApplyAttributeBonus(skills, attributes)
	logger.Debug("Applied attribute bonuses to thief skills")

	return skillsWithBonuses, nil
}
