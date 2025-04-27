package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"mordezzanV4/internal/logger"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
	"strconv"
	"strings"
)

type ThiefSkillsRepository interface {
	GetThiefSkillsByLevel(ctx context.Context, level int64) ([]*models.ThiefSkillWithChance, error)
	GetEffectiveThiefLevel(ctx context.Context, class string, level int64) (int64, error)
	ApplyAttributeBonus(skills []*models.ThiefSkillWithChance, attributes map[string]int) []*models.ThiefSkillWithChance
}

type SQLCThiefSkillsRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func NewSQLCThiefSkillsRepository(db *sql.DB) *SQLCThiefSkillsRepository {
	return &SQLCThiefSkillsRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCThiefSkillsRepository) GetThiefSkillsByLevel(ctx context.Context, level int64) ([]*models.ThiefSkillWithChance, error) {
	// No need to adjust level since we have data for each individual level
	skills, err := r.q.GetThiefSkillsByLevel(ctx, level)
	if err != nil {
		return nil, fmt.Errorf("error fetching thief skills: %w", err)
	}

	result := make([]*models.ThiefSkillWithChance, len(skills))
	for i, skill := range skills {
		result[i] = &models.ThiefSkillWithChance{
			ID:            skill.ID,
			Name:          skill.SkillName,
			Attribute:     skill.Attribute,
			SuccessChance: skill.SuccessChance,
		}
	}
	return result, nil
}

// ApplyAttributeBonus applies attribute bonuses to thief skills
// If the character has 16 or more in the relevant attribute, they get a +1 to their success chance
func (r *SQLCThiefSkillsRepository) ApplyAttributeBonus(skills []*models.ThiefSkillWithChance, attributes map[string]int) []*models.ThiefSkillWithChance {
	result := make([]*models.ThiefSkillWithChance, len(skills))

	for i, skill := range skills {
		result[i] = &models.ThiefSkillWithChance{
			ID:            skill.ID,
			Name:          skill.Name,
			Attribute:     skill.Attribute,
			SuccessChance: skill.SuccessChance,
		}

		// Check if the character has 16+ in the relevant attribute
		attrValue, exists := attributes[skill.Attribute]
		if exists && attrValue >= 16 {
			// Parse the current success chance (format: "X:12")
			parts := strings.Split(skill.SuccessChance, ":")
			if len(parts) == 2 {
				// Skip for N/A values (like early Read Scrolls)
				if parts[0] == "N/A" {
					continue
				}

				currentChance, err := strconv.Atoi(parts[0])
				if err == nil {
					// Apply +1 bonus, but don't exceed 11:12
					newChance := currentChance + 1
					if newChance > 11 {
						newChance = 11
					}
					result[i].SuccessChance = fmt.Sprintf("%d:%s", newChance, parts[1])
				}
			}
		}
	}

	return result
}

func (r *SQLCThiefSkillsRepository) GetEffectiveThiefLevel(ctx context.Context, class string, level int64) (int64, error) {
	logger.Debug("Getting effective thief level for class: '%s', level: %d", class, level)

	// Normalize class name by trimming spaces and converting to lowercase
	normalizedClass := strings.TrimSpace(strings.ToLower(class))
	logger.Debug("Normalized class name: '%s'", normalizedClass)

	switch normalizedClass {
	case "thief":
		logger.Debug("Character is a Thief, using actual level: %d", level)
		return level, nil
	default:
		logger.Debug("Character class '%s' has no thief skills", class)
		return 0, nil
	}
}
