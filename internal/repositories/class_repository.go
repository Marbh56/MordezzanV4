package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	apperrors "mordezzanV4/internal/errors"
	"mordezzanV4/internal/models"
	sqlcdb "mordezzanV4/internal/repositories/db/sqlc"
)

// ClassRepository defines methods for accessing class data
type ClassRepository interface {
	GetClassData(ctx context.Context, className string, level int) (*models.ClassData, error)
	GetAllClassData(ctx context.Context, className string) ([]*models.ClassData, error)
	GetNextLevelData(ctx context.Context, className string, currentLevel int) (*models.ClassData, error)
	GetClassAbilities(ctx context.Context, className string) ([]*models.ClassAbility, error)
	GetClassAbilitiesByLevel(ctx context.Context, className string, level int) ([]*models.ClassAbility, error)
	// Thief skills
	GetThiefSkillsForClass(ctx context.Context, className string) ([]*models.ThiefSkill, error)
	GetThiefSkillsForCharacter(ctx context.Context, className string, level int) (map[string]string, error)
	// Class-specific turning abilities
	GetClericTurningAbility(ctx context.Context, level int) (int, error)
	GetPaladinTurningAbility(ctx context.Context, level int) (int, error)
	GetNecromancerTurningAbility(ctx context.Context, level int) (int, error)
	// Monk abilities
	GetMonkACBonus(ctx context.Context, level int) (int, error)
	GetMonkEmptyHandDamage(ctx context.Context, level int) (string, error)
	// Berserker abilities
	GetBerserkerNaturalAC(ctx context.Context, level int) (int, error)
	// Special spell slot handling
	GetSpecialClassSpellSlots(ctx context.Context, className string, level int) (map[string]int, error)
	// Runegraver abilities
	GetRunegraverRunesPerDay(ctx context.Context, level int) (map[string]int, error)
}

// SQLCClassRepository implements ClassRepository using SQLC
type SQLCClassRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCClassRepository creates a new class repository
func NewSQLCClassRepository(db *sql.DB) *SQLCClassRepository {
	return &SQLCClassRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetClassData retrieves class data for a specific class and level
func (r *SQLCClassRepository) GetClassData(ctx context.Context, className string, level int) (*models.ClassData, error) {
	data, err := r.q.GetClassData(ctx, sqlcdb.GetClassDataParams{
		ClassName: className,
		Level:     int64(level),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("class data", fmt.Sprintf("%s level %d", className, level))
		}
		return nil, err
	}

	return mapDbClassDataToModel(data), nil
}

// GetAllClassData retrieves all level data for a specific class
func (r *SQLCClassRepository) GetAllClassData(ctx context.Context, className string) ([]*models.ClassData, error) {
	data, err := r.q.GetAllClassData(ctx, className)
	if err != nil {
		return nil, err
	}

	result := make([]*models.ClassData, len(data))
	for i, d := range data {
		result[i] = mapDbClassDataToModel(d)
	}

	return result, nil
}

// GetNextLevelData retrieves the next level data for a character
func (r *SQLCClassRepository) GetNextLevelData(ctx context.Context, className string, currentLevel int) (*models.ClassData, error) {
	data, err := r.q.GetNextLevelData(ctx, sqlcdb.GetNextLevelDataParams{
		ClassName: className,
		Level:     int64(currentLevel),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("next level data", fmt.Sprintf("%s level > %d", className, currentLevel))
		}
		return nil, err
	}

	return mapDbClassDataToModel(data), nil
}

// GetClassAbilities retrieves all abilities for a specific class
func (r *SQLCClassRepository) GetClassAbilities(ctx context.Context, className string) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetClassAbilities(ctx, className)
	if err != nil {
		return nil, err
	}

	result := make([]*models.ClassAbility, len(abilities))
	for i, a := range abilities {
		result[i] = &models.ClassAbility{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			MinLevel:    int(a.MinLevel),
		}
	}

	return result, nil
}

// GetClassAbilitiesByLevel retrieves abilities for a class at a specific level
func (r *SQLCClassRepository) GetClassAbilitiesByLevel(ctx context.Context, className string, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetClassAbilitiesByLevel(ctx, sqlcdb.GetClassAbilitiesByLevelParams{
		ClassName: className,
		MinLevel:  int64(level),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*models.ClassAbility, len(abilities))
	for i, a := range abilities {
		result[i] = &models.ClassAbility{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			MinLevel:    int(a.MinLevel),
		}
	}

	return result, nil
}

// GetClericTurningAbility retrieves turning ability for clerics
func (r *SQLCClassRepository) GetClericTurningAbility(ctx context.Context, level int) (int, error) {
	ability, err := r.q.GetClericTurningAbility(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Return 0 if no turning ability found
		}
		return 0, err
	}
	return int(ability), nil
}

// GetPaladinTurningAbility retrieves turning ability for paladins
func (r *SQLCClassRepository) GetPaladinTurningAbility(ctx context.Context, level int) (int, error) {
	ability, err := r.q.GetPaladinTurningAbility(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Return 0 if no turning ability found
		}
		return 0, err
	}
	return int(ability), nil
}

// GetNecromancerTurningAbility retrieves turning ability for necromancers
func (r *SQLCClassRepository) GetNecromancerTurningAbility(ctx context.Context, level int) (int, error) {
	ability, err := r.q.GetNecromancerTurningAbility(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // Return 0 if no turning ability found
		}
		return 0, err
	}
	return int(ability), nil
}

// GetMonkACBonus retrieves AC bonus for monks
func (r *SQLCClassRepository) GetMonkACBonus(ctx context.Context, level int) (int, error) {
	bonus, err := r.q.GetMonkACBonus(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return int(bonus), nil
}

// GetMonkEmptyHandDamage retrieves empty hand damage for monks
func (r *SQLCClassRepository) GetMonkEmptyHandDamage(ctx context.Context, level int) (string, error) {
	damage, err := r.q.GetMonkEmptyHandDamage(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "1d4", nil // Default damage
		}
		return "", err
	}
	return damage, nil
}

// GetBerserkerNaturalAC retrieves natural AC for berserkers
func (r *SQLCClassRepository) GetBerserkerNaturalAC(ctx context.Context, level int) (int, error) {
	ac, err := r.q.GetBerserkerNaturalAC(ctx, sqlcdb.GetBerserkerNaturalACParams{
		ClassName: "Berserker",
		Level:     int64(level),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return int(ac), nil
}

// GetSpecialClassSpellSlots handles complex spell slot logic for multiclass spellcasters
func (r *SQLCClassRepository) GetSpecialClassSpellSlots(ctx context.Context, className string, level int) (map[string]int, error) {
	switch className {
	case "Ranger":
		return r.getRangerSpellSlots(ctx, level)
	case "Shaman":
		return r.getShamanSpellSlots(ctx, level)
	case "Bard":
		return r.getBardSpellSlots(ctx, level)
	default:
		return nil, nil
	}
}

// getRangerSpellSlots handles Ranger's special spellcasting
func (r *SQLCClassRepository) getRangerSpellSlots(ctx context.Context, level int) (map[string]int, error) {
	// Rangers use both druid and magician spells, implement this based on your game rules
	// This is a placeholder implementation
	if level < 7 {
		return map[string]int{}, nil // No spells before level 7
	}

	// Implement ranger spell slot logic here
	return map[string]int{
		"druid_level1":    1,
		"magician_level1": 1,
	}, nil
}

// getShamanSpellSlots handles Shaman's special spellcasting
func (r *SQLCClassRepository) getShamanSpellSlots(ctx context.Context, level int) (map[string]int, error) {
	// Implement based on shaman_divine_spells and shaman_arcane_spells tables
	return map[string]int{}, nil
}

// getBardSpellSlots handles Bard's special spellcasting
func (r *SQLCClassRepository) getBardSpellSlots(ctx context.Context, level int) (map[string]int, error) {
	// Implement based on bard_druid_spells and bard_illusionist_spells tables
	return map[string]int{}, nil
}

// GetRunegraverRunesPerDay retrieves runes per day for runegravers
func (r *SQLCClassRepository) GetRunegraverRunesPerDay(ctx context.Context, level int) (map[string]int, error) {
	runes, err := r.q.GetRunesPerDay(ctx, int64(level))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]int{}, nil
		}
		return nil, err
	}

	result := map[string]int{
		"level1": 0,
		"level2": 0,
		"level3": 0,
		"level4": 0,
		"level5": 0,
		"level6": 0,
	}

	if runes.Level1.Valid {
		result["level1"] = int(runes.Level1.Int64)
	}
	if runes.Level2.Valid {
		result["level2"] = int(runes.Level2.Int64)
	}
	if runes.Level3.Valid {
		result["level3"] = int(runes.Level3.Int64)
	}
	if runes.Level4.Valid {
		result["level4"] = int(runes.Level4.Int64)
	}
	if runes.Level5.Valid {
		result["level5"] = int(runes.Level5.Int64)
	}
	if runes.Level6.Valid {
		result["level6"] = int(runes.Level6.Int64)
	}

	return result, nil
}

func mapDbClassDataToModel(data sqlcdb.ClassDatum) *models.ClassData {
	return &models.ClassData{
		ID:               data.ID,
		ClassName:        data.ClassName,
		Level:            int(data.Level),
		ExperiencePoints: int(data.ExperiencePoints),
		HitDice:          data.HitDice,
		SavingThrow:      int(data.SavingThrow),
		FightingAbility:  int(data.FightingAbility),
		CastingAbility:   int(getNullInt64Value(data.CastingAbility)),
		SpellSlots: map[string]int{
			"level1": int(getNullInt64Value(data.SpellSlotsLevel1)),
			"level2": int(getNullInt64Value(data.SpellSlotsLevel2)),
			"level3": int(getNullInt64Value(data.SpellSlotsLevel3)),
			"level4": int(getNullInt64Value(data.SpellSlotsLevel4)),
			"level5": int(getNullInt64Value(data.SpellSlotsLevel5)),
			"level6": int(getNullInt64Value(data.SpellSlotsLevel6)),
		},
	}
}

// Helper function to safely get value from sql.NullInt64
func getNullInt64Value(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}
