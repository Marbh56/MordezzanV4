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

type ClassRepository interface {
	GetClassData(ctx context.Context, className string, level int) (*models.ClassData, error)
	GetAllClassData(ctx context.Context, className string) ([]*models.ClassData, error)
	GetNextLevelData(ctx context.Context, className string, currentLevel int) (*models.ClassData, error)
	GetCharacterClassInfo(ctx context.Context, characterID int64) (*models.Character, error)
	GetClassAbilities(ctx context.Context, className string) ([]*models.ClassAbility, error)
	GetClassAbilitiesByLevel(ctx context.Context, className string, level int) ([]*models.ClassAbility, error)
	GetThiefSkillsForClass(ctx context.Context, className string) ([]*models.ThiefSkill, error)
	GetThiefSkillsForCharacter(ctx context.Context, className string, level int) (map[string]string, error)
	GetThiefSkillsByClass(ctx context.Context, className string) ([]models.ThiefSkill, error)
	GetThiefSkillByName(ctx context.Context, skillName string) (*models.ThiefSkill, error)
	GetThiefSkillChance(ctx context.Context, skillID int64, level int) (string, error)
	GetClericTurningAbility(ctx context.Context, level int) (int, error)
	GetPaladinTurningAbility(ctx context.Context, level int) (int, error)
	GetNecromancerTurningAbility(ctx context.Context, level int) (int, error)
	GetMonkACBonus(ctx context.Context, level int) (int, error)
	GetMonkEmptyHandDamage(ctx context.Context, level int) (string, error)
	GetBerserkerNaturalAC(ctx context.Context, level int) (int, error)
	GetSpecialClassSpellSlots(ctx context.Context, className string, level int) (map[string]int, error)
	GetRunegraverRunesPerDay(ctx context.Context, level int) (map[string]int, error)

	// New methods for each class's abilities
	GetBarbarianAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetBerserkerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetBardAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetCataphractAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetClericAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetCryomancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetDruidAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetFighterAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetHuntsmanAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetIllusionistAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetLegerdemainistAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetMagicianAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetMonkAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetNecromancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetPaladinAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetPriestAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetPurloinerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetPyromancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetRangerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetRunegraverAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetScoutAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetShamanAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetThiefAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetWarlockAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
	GetWitchAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error)
}

// SQLCClassRepository implements ClassRepository using SQLC
type SQLCClassRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

func (r *SQLCClassRepository) GetClassLevelData(ctx context.Context, className string, level int) (*models.ClassData, error) {
	return r.GetClassData(ctx, className, level)
}

// NewSQLCClassRepository creates a new class repository
func NewSQLCClassRepository(db *sql.DB) *SQLCClassRepository {
	return &SQLCClassRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

func (r *SQLCClassRepository) GetCharacterClassInfo(ctx context.Context, characterID int64) (*models.Character, error) {
	// Define the SQL query to fetch the character with class-related information
	query := `
		SELECT 
			id, name, class, level, strength, dexterity, constitution, 
			intelligence, wisdom, charisma 
		FROM characters 
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, characterID)

	// Create a character object to hold the result
	character := &models.Character{}

	// Scan the row into the character struct
	err := row.Scan(
		&character.ID,
		&character.Name,
		&character.Class,
		&character.Level,
		&character.Strength,
		&character.Dexterity,
		&character.Constitution,
		&character.Intelligence,
		&character.Wisdom,
		&character.Charisma,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("character", fmt.Sprintf("%d", characterID))
		}
		return nil, fmt.Errorf("error fetching character class info: %w", err)
	}

	return character, nil
}

func (r *SQLCClassRepository) GetThiefSkillByName(ctx context.Context, skillName string) (*models.ThiefSkill, error) {
	// SQL query to get thief skill by name
	query := `
		SELECT id, skill_name, attribute
		FROM thief_skills
		WHERE skill_name = ?
	`

	row := r.db.QueryRowContext(ctx, query, skillName)

	skill := &models.ThiefSkill{}
	err := row.Scan(&skill.ID, &skill.Name, &skill.Attribute)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("thief skill", skillName)
		}
		return nil, fmt.Errorf("error fetching thief skill: %w", err)
	}

	return skill, nil
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

func (r *SQLCClassRepository) GetThiefSkillChance(ctx context.Context, skillID int64, level int) (string, error) {
	// Query to get the skill chance for a specific skill and level range
	query := `
		SELECT success_chance
		FROM thief_skill_progression
		WHERE skill_id = ? AND ? BETWEEN 
			CAST(SUBSTR(level_range, 1, INSTR(level_range, '-') - 1) AS INTEGER) 
			AND 
			CAST(SUBSTR(level_range, INSTR(level_range, '-') + 1) AS INTEGER)
	`

	var successChance string
	err := r.db.QueryRowContext(ctx, query, skillID, level).Scan(&successChance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperrors.NewNotFound("thief skill chance", fmt.Sprintf("skill_id %d, level %d", skillID, level))
		}
		return "", fmt.Errorf("error fetching thief skill chance: %w", err)
	}

	return successChance, nil
}

func (r *SQLCClassRepository) GetThiefSkillsForClass(ctx context.Context, className string) ([]*models.ThiefSkill, error) {
	// Query to get thief skills for a specific class
	query := `
		SELECT ts.id, ts.skill_name, ts.attribute
		FROM thief_skills ts
		JOIN class_thief_skill_mapping ctsm ON ts.id = ctsm.skill_id
		WHERE ctsm.class_name = ?
	`

	rows, err := r.db.QueryContext(ctx, query, className)
	if err != nil {
		return nil, fmt.Errorf("error fetching thief skills: %w", err)
	}
	defer rows.Close()

	var skills []*models.ThiefSkill
	for rows.Next() {
		skill := &models.ThiefSkill{}
		if err := rows.Scan(&skill.ID, &skill.Name, &skill.Attribute); err != nil {
			return nil, fmt.Errorf("error scanning thief skill row: %w", err)
		}
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating thief skill rows: %w", err)
	}

	return skills, nil
}

func (r *SQLCClassRepository) GetThiefSkillsForCharacter(ctx context.Context, className string, level int) (map[string]string, error) {
	// Use a direct SQL query instead of the missing sqlc method
	rows, err := r.db.QueryContext(ctx, `
        SELECT 
            ts.skill_name, 
            tsp.success_chance
        FROM thief_skills ts
        JOIN class_thief_skill_mapping ctsm ON ts.id = ctsm.skill_id
        JOIN thief_skill_progression tsp ON ts.id = tsp.skill_id
        WHERE ctsm.class_name = ?
        AND ? BETWEEN 
            CAST(SUBSTR(tsp.level_range, 1, INSTR(tsp.level_range, '-') - 1) AS INTEGER) 
            AND 
            CAST(SUBSTR(tsp.level_range, INSTR(tsp.level_range, '-') + 1) AS INTEGER)
    `, className, level)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var skillName, successChance string
		if err := rows.Scan(&skillName, &successChance); err != nil {
			return nil, err
		}
		result[skillName] = successChance
	}

	return result, nil
}

func (r *SQLCClassRepository) GetThiefSkillsByClass(ctx context.Context, className string) ([]models.ThiefSkill, error) {
	skills, err := r.q.GetThiefSkillsByClassName(ctx, className)
	if err != nil {
		return nil, err
	}

	result := make([]models.ThiefSkill, len(skills))
	for i, skill := range skills {
		result[i] = models.ThiefSkill{
			ID:        skill.ID,
			Name:      skill.SkillName,
			Attribute: skill.Attribute,
		}
	}

	return result, nil
}

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

// GetClassAbilities method using the class-specific methods
func (r *SQLCClassRepository) GetClassAbilities(ctx context.Context, className string) ([]*models.ClassAbility, error) {
	// Just call the ByLevel version with a very high level to get all abilities
	// This assumes that no class has abilities with an extremely high level
	return r.GetClassAbilitiesByLevel(ctx, className, 99) // Using 99 as a high level to get all abilities
}

// Updated GetClassAbilitiesByLevel method to use the class-specific methods
func (r *SQLCClassRepository) GetClassAbilitiesByLevel(ctx context.Context, className string, level int) ([]*models.ClassAbility, error) {
	// Use the class-specific abilities methods instead of the generic one
	switch className {
	case "Barbarian":
		return r.GetBarbarianAbilities(ctx, level)
	case "Berserker":
		return r.GetBerserkerAbilities(ctx, level)
	case "Bard":
		return r.GetBardAbilities(ctx, level)
	case "Cataphract":
		return r.GetCataphractAbilities(ctx, level)
	case "Cleric":
		return r.GetClericAbilities(ctx, level)
	case "Cryomancer":
		return r.GetCryomancerAbilities(ctx, level)
	case "Druid":
		return r.GetDruidAbilities(ctx, level)
	case "Fighter":
		return r.GetFighterAbilities(ctx, level)
	case "Huntsman":
		return r.GetHuntsmanAbilities(ctx, level)
	case "Illusionist":
		return r.GetIllusionistAbilities(ctx, level)
	case "Legerdemainist":
		return r.GetLegerdemainistAbilities(ctx, level)
	case "Magician":
		return r.GetMagicianAbilities(ctx, level)
	case "Monk":
		return r.GetMonkAbilities(ctx, level)
	case "Necromancer":
		return r.GetNecromancerAbilities(ctx, level)
	case "Paladin":
		return r.GetPaladinAbilities(ctx, level)
	case "Priest":
		return r.GetPriestAbilities(ctx, level)
	case "Purloiner":
		return r.GetPurloinerAbilities(ctx, level)
	case "Pyromancer":
		return r.GetPyromancerAbilities(ctx, level)
	case "Ranger":
		return r.GetRangerAbilities(ctx, level)
	case "Runegraver":
		return r.GetRunegraverAbilities(ctx, level)
	case "Scout":
		return r.GetScoutAbilities(ctx, level)
	case "Shaman":
		return r.GetShamanAbilities(ctx, level)
	case "Thief":
		return r.GetThiefAbilities(ctx, level)
	case "Warlock":
		return r.GetWarlockAbilities(ctx, level)
	case "Witch":
		return r.GetWitchAbilities(ctx, level)
	default:
		// Fallback to the old method (keeping it for backward compatibility)
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
}

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

func (r *SQLCClassRepository) getShamanSpellSlots(ctx context.Context, level int) (map[string]int, error) {
	// Implement based on shaman_divine_spells and shaman_arcane_spells tables
	return map[string]int{}, nil
}

func (r *SQLCClassRepository) getBardSpellSlots(ctx context.Context, level int) (map[string]int, error) {
	// Implement based on bard_druid_spells and bard_illusionist_spells tables
	return map[string]int{}, nil
}

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

// Helper Functions
func getNullInt64Value(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

func getLevelRangeForLevel(level int) string {
	// Implement your level range logic here
	// For example, if levels 1-4 are stored as "1-4", 5-8 as "5-8", etc.
	if level <= 4 {
		return "1-4"
	} else if level <= 8 {
		return "5-8"
	} else if level <= 12 {
		return "9-12"
	}
	return "1-12" // Fallback
}

// GetBarbarianAbilities gets all barbarian abilities for the specified level
func (r *SQLCClassRepository) GetBarbarianAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetBarbarianAbilities(ctx, int64(level))
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

// GetBerserkerAbilities gets all berserker abilities for the specified level
func (r *SQLCClassRepository) GetBerserkerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetBerserkerAbilities(ctx, int64(level))
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

// GetBardAbilities gets all bard abilities for the specified level
func (r *SQLCClassRepository) GetBardAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetBardAbilities(ctx, int64(level))
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

// GetCataphractAbilities gets all cataphract abilities for the specified level
func (r *SQLCClassRepository) GetCataphractAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetCataphractAbilities(ctx, int64(level))
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

// GetClericAbilities gets all cleric abilities for the specified level
func (r *SQLCClassRepository) GetClericAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetClericAbilities(ctx, int64(level))
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

// GetCryomancerAbilities gets all cryomancer abilities for the specified level
func (r *SQLCClassRepository) GetCryomancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetCryomancerAbilities(ctx, int64(level))
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

// GetDruidAbilities gets all druid abilities for the specified level
func (r *SQLCClassRepository) GetDruidAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetDruidAbilities(ctx, int64(level))
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

// GetFighterAbilities gets all fighter abilities for the specified level
func (r *SQLCClassRepository) GetFighterAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetFighterAbilities(ctx, int64(level))
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

// GetHuntsmanAbilities gets all huntsman abilities for the specified level
func (r *SQLCClassRepository) GetHuntsmanAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetHuntsmanAbilities(ctx, int64(level))
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

// GetIllusionistAbilities gets all illusionist abilities for the specified level
func (r *SQLCClassRepository) GetIllusionistAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetIllusionistAbilities(ctx, int64(level))
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

// GetLegerdemainistAbilities gets all legerdemainist abilities for the specified level
func (r *SQLCClassRepository) GetLegerdemainistAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetLegerdemainistAbilities(ctx, int64(level))
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

// GetMagicianAbilities gets all magician abilities for the specified level
func (r *SQLCClassRepository) GetMagicianAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetMagicianAbilities(ctx, int64(level))
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

// GetMonkAbilities gets all monk abilities for the specified level
func (r *SQLCClassRepository) GetMonkAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetMonkAbilities(ctx, int64(level))
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

// GetNecromancerAbilities gets all necromancer abilities for the specified level
func (r *SQLCClassRepository) GetNecromancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetNecromancerAbilities(ctx, int64(level))
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

// GetPaladinAbilities gets all paladin abilities for the specified level
func (r *SQLCClassRepository) GetPaladinAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetPaladinAbilities(ctx, int64(level))
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

// GetPriestAbilities gets all priest abilities for the specified level
func (r *SQLCClassRepository) GetPriestAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetPriestAbilities(ctx, int64(level))
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

// GetPurloinerAbilities gets all purloiner abilities for the specified level
func (r *SQLCClassRepository) GetPurloinerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetPurloinerAbilities(ctx, int64(level))
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

// GetPyromancerAbilities gets all pyromancer abilities for the specified level
func (r *SQLCClassRepository) GetPyromancerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetPyromancerAbilities(ctx, int64(level))
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

// GetRangerAbilities gets all ranger abilities for the specified level
func (r *SQLCClassRepository) GetRangerAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetRangerAbilities(ctx, int64(level))
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

// GetRunegraverAbilities gets all runegraver abilities for the specified level
func (r *SQLCClassRepository) GetRunegraverAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetRunegraverAbilities(ctx, int64(level))
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

// GetScoutAbilities gets all scout abilities for the specified level
func (r *SQLCClassRepository) GetScoutAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetScoutAbilities(ctx, int64(level))
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

// GetShamanAbilities gets all shaman abilities for the specified level
func (r *SQLCClassRepository) GetShamanAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetShamanAbilities(ctx, int64(level))
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

// GetThiefAbilities gets all thief abilities for the specified level
func (r *SQLCClassRepository) GetThiefAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetThiefAbilities(ctx, int64(level))
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

// GetWarlockAbilities gets all warlock abilities for the specified level
func (r *SQLCClassRepository) GetWarlockAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetWarlockAbilities(ctx, int64(level))
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

// GetWitchAbilities gets all witch abilities for the specified level
func (r *SQLCClassRepository) GetWitchAbilities(ctx context.Context, level int) ([]*models.ClassAbility, error) {
	abilities, err := r.q.GetWitchAbilities(ctx, int64(level))
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
