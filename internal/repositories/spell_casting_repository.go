// repositories/spell_repository.go

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

// SpellCastingRepository defines methods for managing character spells
type SpellCastingRepository interface {
	// Known spells management
	GetKnownSpells(ctx context.Context, characterID int64) ([]models.KnownSpell, error)
	GetKnownSpellsByClass(ctx context.Context, characterID int64, spellClass string) ([]models.KnownSpell, error)
	AddKnownSpell(ctx context.Context, input *models.AddKnownSpellInput) (int64, error)
	RemoveKnownSpell(ctx context.Context, characterID, spellID int64) error

	// Prepared spells management
	GetPreparedSpells(ctx context.Context, characterID int64) ([]models.PreparedSpell, error)
	GetPreparedSpellsByClass(ctx context.Context, characterID int64, spellClass string) ([]models.PreparedSpell, error)
	PrepareSpell(ctx context.Context, input *models.PrepareSpellInput) (int64, error)
	UnprepareSpell(ctx context.Context, characterID, spellID int64) error
	ClearPreparedSpells(ctx context.Context, characterID int64) error

	// Spell slot and limit information
	GetSpellSlots(ctx context.Context, characterID int64) (map[string]int, error)
	GetMaxKnownSpells(ctx context.Context, characterID int64, spellClass string) (map[string]int, error)
	CalculateBonusSpells(ctx context.Context, characterID int64) (map[string]map[string]int, error)
}

// SQLCSpellCastingRepository implements SpellCastingRepository using SQLC
type SQLCSpellCastingRepository struct {
	db *sql.DB
	q  *sqlcdb.Queries
}

// NewSQLCSpellCastingRepository creates a new spell casting repository
func NewSQLCSpellCastingRepository(db *sql.DB) *SQLCSpellCastingRepository {
	return &SQLCSpellCastingRepository{
		db: db,
		q:  sqlcdb.New(db),
	}
}

// GetKnownSpells retrieves all spells known by a character
func (r *SQLCSpellCastingRepository) GetKnownSpells(ctx context.Context, characterID int64) ([]models.KnownSpell, error) {
	knownSpells, err := r.q.GetKnownSpells(ctx, characterID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.KnownSpell, len(knownSpells))
	for i, spell := range knownSpells {
		result[i] = models.KnownSpell{
			ID:          spell.ID,
			CharacterID: spell.CharacterID,
			SpellID:     spell.SpellID,
			SpellName:   spell.SpellName,
			SpellLevel:  int(spell.SpellLevel),
			SpellClass:  spell.SpellClass,
			IsMemorized: spell.IsMemorized,
			Notes:       spell.Notes.String,
			CreatedAt:   spell.CreatedAt,
			UpdatedAt:   spell.UpdatedAt,
		}
	}

	return result, nil
}

// GetKnownSpellsByClass retrieves all spells known by a character for a specific spell class
func (r *SQLCSpellCastingRepository) GetKnownSpellsByClass(ctx context.Context, characterID int64, spellClass string) ([]models.KnownSpell, error) {
	knownSpells, err := r.q.GetKnownSpellsByClass(ctx, sqlcdb.GetKnownSpellsByClassParams{
		CharacterID: characterID,
		SpellClass:  spellClass,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.KnownSpell, len(knownSpells))
	for i, spell := range knownSpells {
		result[i] = models.KnownSpell{
			ID:          spell.ID,
			CharacterID: spell.CharacterID,
			SpellID:     spell.SpellID,
			SpellName:   spell.SpellName,
			SpellLevel:  int(spell.SpellLevel),
			SpellClass:  spell.SpellClass,
			IsMemorized: spell.IsMemorized,
			Notes:       spell.Notes.String,
			CreatedAt:   spell.CreatedAt,
			UpdatedAt:   spell.UpdatedAt,
		}
	}

	return result, nil
}

// AddKnownSpell adds a spell to a character's known spells
func (r *SQLCSpellCastingRepository) AddKnownSpell(ctx context.Context, input *models.AddKnownSpellInput) (int64, error) {
	// First, get the spell details
	spell, err := r.q.GetSpell(ctx, input.SpellID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.NewNotFound("spell", input.SpellID)
		}
		return 0, apperrors.NewDatabaseError(err)
	}

	// Check if the spell is already known by the character
	_, err = r.q.GetKnownSpellByCharacterAndSpell(ctx, sqlcdb.GetKnownSpellByCharacterAndSpellParams{
		CharacterID: input.CharacterID,
		SpellID:     input.SpellID,
	})
	if err == nil {
		return 0, apperrors.NewValidationError("spell_id", "Spell is already known by this character")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, apperrors.NewDatabaseError(err)
	}

	spellLevel := 0
	switch input.SpellClass {
	case "Magician":
		spellLevel = int(spell.MagLevel)
	case "Cleric":
		spellLevel = int(spell.ClrLevel)
	case "Druid":
		spellLevel = int(spell.DrdLevel)
	case "Illusionist":
		spellLevel = int(spell.IllLevel)
	case "Necromancer":
		spellLevel = int(spell.NecLevel)
	case "Pyromancer":
		spellLevel = int(spell.PyrLevel)
	case "Cryomancer":
		spellLevel = int(spell.CryLevel)
	case "Witch":
		spellLevel = int(spell.WchLevel)
	}

	result, err := r.q.AddKnownSpell(ctx, sqlcdb.AddKnownSpellParams{
		CharacterID: input.CharacterID,
		SpellID:     input.SpellID,
		SpellName:   spell.Name,
		SpellLevel:  int64(spellLevel),
		SpellClass:  input.SpellClass,
		Notes: sql.NullString{
			String: input.Notes,
			Valid:  input.Notes != "",
		},
	})

	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return id, nil
}

// RemoveKnownSpell removes a spell from a character's known spells
func (r *SQLCSpellCastingRepository) RemoveKnownSpell(ctx context.Context, characterID, spellID int64) error {
	// First check if this spell is prepared - we shouldn't allow removing it if it is
	preparedSpell, err := r.q.GetPreparedSpellByCharacterAndSpell(ctx, sqlcdb.GetPreparedSpellByCharacterAndSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err == nil && preparedSpell.ID > 0 {
		return apperrors.NewValidationError("spell_id", "Cannot remove a prepared spell. Unprepare it first.")
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return apperrors.NewDatabaseError(err)
	}

	// Check if the spell is known by the character
	knownSpell, err := r.q.GetKnownSpellByCharacterAndSpell(ctx, sqlcdb.GetKnownSpellByCharacterAndSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("known spell", fmt.Sprintf("character %d spell %d", characterID, spellID))
		}
		return apperrors.NewDatabaseError(err)
	}

	// Remove the known spell
	err = r.q.RemoveKnownSpell(ctx, knownSpell.ID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

// GetPreparedSpells retrieves all prepared spells for a character
func (r *SQLCSpellCastingRepository) GetPreparedSpells(ctx context.Context, characterID int64) ([]models.PreparedSpell, error) {
	preparedSpells, err := r.q.GetPreparedSpells(ctx, characterID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.PreparedSpell, len(preparedSpells))
	for i, spell := range preparedSpells {
		result[i] = models.PreparedSpell{
			ID:          spell.ID,
			CharacterID: spell.CharacterID,
			SpellID:     spell.SpellID,
			SpellName:   spell.SpellName,
			SpellLevel:  int(spell.SpellLevel),
			SpellClass:  spell.SpellClass,
			SlotIndex:   int(spell.SlotIndex),
			CreatedAt:   spell.CreatedAt,
			UpdatedAt:   spell.UpdatedAt,
		}
	}

	return result, nil
}

// GetPreparedSpellsByClass retrieves all prepared spells for a character of a specific class
func (r *SQLCSpellCastingRepository) GetPreparedSpellsByClass(ctx context.Context, characterID int64, spellClass string) ([]models.PreparedSpell, error) {
	preparedSpells, err := r.q.GetPreparedSpellsByClass(ctx, sqlcdb.GetPreparedSpellsByClassParams{
		CharacterID: characterID,
		SpellClass:  spellClass,
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	result := make([]models.PreparedSpell, len(preparedSpells))
	for i, spell := range preparedSpells {
		result[i] = models.PreparedSpell{
			ID:          spell.ID,
			CharacterID: spell.CharacterID,
			SpellID:     spell.SpellID,
			SpellName:   spell.SpellName,
			SpellLevel:  int(spell.SpellLevel),
			SpellClass:  spell.SpellClass,
			SlotIndex:   int(spell.SlotIndex),
			CreatedAt:   spell.CreatedAt,
			UpdatedAt:   spell.UpdatedAt,
		}
	}

	return result, nil
}

// PrepareSpell prepares a spell for a character
func (r *SQLCSpellCastingRepository) PrepareSpell(ctx context.Context, input *models.PrepareSpellInput) (int64, error) {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Check if the spell is known by the character
	knownSpell, err := qtx.GetKnownSpellByCharacterAndSpell(ctx, sqlcdb.GetKnownSpellByCharacterAndSpellParams{
		CharacterID: input.CharacterID,
		SpellID:     input.SpellID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.NewValidationError("spell_id", "Character does not know this spell")
		}
		return 0, apperrors.NewDatabaseError(err)
	}

	// Check if character has available slots for this spell level and class
	preparedCount, err := qtx.CountPreparedSpellsByLevelAndClass(ctx, sqlcdb.CountPreparedSpellsByLevelAndClassParams{
		CharacterID: input.CharacterID,
		SpellLevel:  int64(input.SpellLevel),
		SpellClass:  input.SpellClass,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	// Get the character's spell slots from class_data or merged from multiple sources
	character, err := qtx.GetCharacter(ctx, input.CharacterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.NewNotFound("character", input.CharacterID)
		}
		return 0, apperrors.NewDatabaseError(err)
	}

	classData, err := qtx.GetClassData(ctx, sqlcdb.GetClassDataParams{
		ClassName: character.Class,
		Level:     int64(character.Level),
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	var availableSlots int64
	switch input.SpellLevel {
	case 1:
		if classData.SpellSlotsLevel1.Valid {
			availableSlots = classData.SpellSlotsLevel1.Int64
		}
	case 2:
		if classData.SpellSlotsLevel2.Valid {
			availableSlots = classData.SpellSlotsLevel2.Int64
		}
	case 3:
		if classData.SpellSlotsLevel3.Valid {
			availableSlots = classData.SpellSlotsLevel3.Int64
		}
	case 4:
		if classData.SpellSlotsLevel4.Valid {
			availableSlots = classData.SpellSlotsLevel4.Int64
		}
	case 5:
		if classData.SpellSlotsLevel5.Valid {
			availableSlots = classData.SpellSlotsLevel5.Int64
		}
	case 6:
		if classData.SpellSlotsLevel6.Valid {
			availableSlots = classData.SpellSlotsLevel6.Int64
		}
	default:
		return 0, apperrors.NewValidationError("spell_level", "Invalid spell level")
	}

	// Add bonus slots based on relevant attribute (Wisdom for clerics, Intelligence for magicians)
	bonusSlots := int64(0)
	if input.SpellClass == "Cleric" || input.SpellClass == "Druid" || input.SpellClass == "Priest" {
		// Apply Wisdom bonus for divine casters
		if character.Wisdom >= 13 && character.Wisdom <= 14 && input.SpellLevel == 1 {
			bonusSlots = 1
		} else if character.Wisdom >= 15 && character.Wisdom <= 16 && input.SpellLevel <= 2 {
			bonusSlots = 1
		} else if character.Wisdom == 17 && input.SpellLevel <= 3 {
			bonusSlots = 1
		} else if character.Wisdom == 18 && input.SpellLevel <= 4 {
			bonusSlots = 1
		}
	} else if input.SpellClass == "Magician" || input.SpellClass == "Illusionist" || input.SpellClass == "Necromancer" {
		// Apply Intelligence bonus for arcane casters
		if character.Intelligence >= 13 && character.Intelligence <= 14 && input.SpellLevel == 1 {
			bonusSlots = 1
		} else if character.Intelligence >= 15 && character.Intelligence <= 16 && input.SpellLevel <= 2 {
			bonusSlots = 1
		} else if character.Intelligence == 17 && input.SpellLevel <= 3 {
			bonusSlots = 1
		} else if character.Intelligence == 18 && input.SpellLevel <= 4 {
			bonusSlots = 1
		}
	}

	// Total available slots is base slots plus bonus slots
	totalAvailableSlots := availableSlots + bonusSlots

	if preparedCount >= totalAvailableSlots {
		return 0, apperrors.NewValidationError("spell_level", fmt.Sprintf("No available slots for level %d %s spells", input.SpellLevel, input.SpellClass))
	}

	// Find the next available slot index
	nextSlotIndex, err := qtx.GetNextAvailableSlotIndex(ctx, sqlcdb.GetNextAvailableSlotIndexParams{
		CharacterID: input.CharacterID,
		SpellLevel:  int64(input.SpellLevel),
		SpellClass:  input.SpellClass,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	// Add the prepared spell
	result, err := qtx.PrepareSpell(ctx, sqlcdb.PrepareSpellParams{
		CharacterID: input.CharacterID,
		SpellID:     input.SpellID,
		SpellName:   knownSpell.SpellName,
		SpellLevel:  knownSpell.SpellLevel,
		SpellClass:  input.SpellClass,
		SlotIndex:   nextSlotIndex,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	// Mark the spell as memorized in known_spells
	err = qtx.MarkSpellAsMemorized(ctx, sqlcdb.MarkSpellAsMemorizedParams{
		IsMemorized: true,
		ID:          knownSpell.ID,
	})
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(err)
	}

	return id, nil
}

// UnprepareSpell removes a prepared spell
func (r *SQLCSpellCastingRepository) UnprepareSpell(ctx context.Context, characterID, spellID int64) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Get the prepared spell
	preparedSpell, err := qtx.GetPreparedSpellByCharacterAndSpell(ctx, sqlcdb.GetPreparedSpellByCharacterAndSpellParams{
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperrors.NewNotFound("prepared spell", fmt.Sprintf("character %d spell %d", characterID, spellID))
		}
		return apperrors.NewDatabaseError(err)
	}

	// Remove the prepared spell
	err = qtx.UnprepareSpell(ctx, preparedSpell.ID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Mark the spell as not memorized in known_spells
	err = qtx.MarkSpellAsMemorizedBySpellID(ctx, sqlcdb.MarkSpellAsMemorizedBySpellIDParams{
		IsMemorized: false,
		CharacterID: characterID,
		SpellID:     spellID,
	})
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

// ClearPreparedSpells removes all prepared spells for a character
func (r *SQLCSpellCastingRepository) ClearPreparedSpells(ctx context.Context, characterID int64) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Clear all prepared spells
	err = qtx.ClearPreparedSpells(ctx, characterID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Mark all known spells as not memorized
	err = qtx.ResetAllMemorizedSpells(ctx, characterID)
	if err != nil {
		return apperrors.NewDatabaseError(err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(err)
	}

	return nil
}

// GetSpellSlots retrieves the available spell slots for a character
func (r *SQLCSpellCastingRepository) GetSpellSlots(ctx context.Context, characterID int64) (map[string]int, error) {
	// Get character class and level
	character, err := r.q.GetCharacter(ctx, characterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("character", characterID)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	// Get class data for character's level
	classData, err := r.q.GetClassData(ctx, sqlcdb.GetClassDataParams{
		ClassName: character.Class,
		Level:     int64(character.Level),
	})
	if err != nil {
		return nil, apperrors.NewDatabaseError(err)
	}

	spellSlots := map[string]int{
		"level1": getIntValue(classData.SpellSlotsLevel1),
		"level2": getIntValue(classData.SpellSlotsLevel2),
		"level3": getIntValue(classData.SpellSlotsLevel3),
		"level4": getIntValue(classData.SpellSlotsLevel4),
		"level5": getIntValue(classData.SpellSlotsLevel5),
		"level6": getIntValue(classData.SpellSlotsLevel6),
	}

	// Add special class-specific slots if needed (e.g., for hybrid classes)
	// For now this is handled in the ClassService

	return spellSlots, nil
}

// Helper function to safely extract int value from sql.NullInt64
func getIntValue(n sql.NullInt64) int {
	if n.Valid {
		return int(n.Int64)
	}
	return 0
}

// GetMaxKnownSpells calculates how many spells a character can know based on class and attributes
func (r *SQLCSpellCastingRepository) GetMaxKnownSpells(ctx context.Context, characterID int64, spellClass string) (map[string]int, error) {
	// Get character attributes
	character, err := r.q.GetCharacter(ctx, characterID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NewNotFound("character", characterID)
		}
		return nil, apperrors.NewDatabaseError(err)
	}

	// Initialize result
	result := make(map[string]int)

	// Calculate based on spell class and relevant attribute
	switch spellClass {
	case "Cleric", "Druid", "Priest", "Paladin":
		// Divine casters - Wisdom based
		// Base values for each level (minimum)
		result["level1"] = 4
		result["level2"] = 3
		result["level3"] = 2
		result["level4"] = 2
		result["level5"] = 1
		result["level6"] = 1

		// Apply wisdom bonuses
		if character.Wisdom >= 13 {
			result["level1"] += 1
		}
		if character.Wisdom >= 15 {
			result["level2"] += 1
		}
		if character.Wisdom >= 17 {
			result["level3"] += 1
			result["level1"] += 1
		}
		if character.Wisdom >= 18 {
			result["level4"] += 1
			result["level2"] += 1
		}

	case "Magician", "Illusionist", "Necromancer", "Pyromancer", "Cryomancer":
		// Arcane casters - Intelligence based
		// Base values for each level (minimum)
		result["level1"] = 6
		result["level2"] = 4
		result["level3"] = 3
		result["level4"] = 2
		result["level5"] = 2
		result["level6"] = 1

		// Apply intelligence bonuses
		if character.Intelligence >= 13 {
			result["level1"] += 2
		}
		if character.Intelligence >= 15 {
			result["level2"] += 1
			result["level1"] += 1
		}
		if character.Intelligence >= 17 {
			result["level3"] += 1
			result["level2"] += 1
		}
		if character.Intelligence >= 18 {
			result["level4"] += 1
			result["level1"] += 2
		}

	default:
		// For classes with limited spellcasting or hybrid classes
		// Use a more modest progression
		result["level1"] = 3
		result["level2"] = 2
		result["level3"] = 1
		result["level4"] = 1
		result["level5"] = 0
		result["level6"] = 0
	}

	return result, nil
}

func calculateHybridBonusSpells(character *sqlcdb.Character, className string) map[string]int {
	bonusSpells := make(map[string]int)

	switch className {
	case "Paladin":
		// Paladins use Wisdom for their limited divine spells
		if character.Wisdom >= 13 && character.Level >= 9 {
			bonusSpells["level1"] = 1
		}
		if character.Wisdom >= 15 && character.Level >= 11 {
			bonusSpells["level2"] = 1
		}
	case "Ranger":
		// Rangers might use Wisdom or Intelligence depending on your game rules
		if character.Wisdom >= 14 && character.Level >= 8 {
			bonusSpells["level1"] = 1
		}
		if character.Wisdom >= 16 && character.Level >= 10 {
			bonusSpells["level2"] = 1
		}
	}

	return bonusSpells
}

func (r *SQLCSpellCastingRepository) CalculateBonusSpells(ctx context.Context, characterID int64) (map[string]map[string]int, error) {
	// Get character data
	character, err := r.q.GetCharacter(ctx, characterID)
	if err != nil {
		return nil, err
	}

	// Initialize result map (class -> level -> count)
	result := make(map[string]map[string]int)

	// Initialize maps for each relevant class
	className := character.Class

	// Determine primary and secondary classes based on character class
	var primaryClass, secondaryClass string

	switch className {
	case "Magician", "Illusionist", "Necromancer", "Pyromancer", "Cryomancer":
		primaryClass = className
		// Intelligence-based casters
		result[primaryClass] = calculateArcaneBonusSpells(int(character.Intelligence))
	case "Cleric", "Druid", "Witch", "Priest":
		primaryClass = className
		// Wisdom-based casters
		result[primaryClass] = calculateDivineBonusSpells(int(character.Wisdom))
	case "Paladin":
		// Paladins get cleric spells starting at level 9
		if character.Level >= 9 {
			primaryClass = "Cleric"
			result[primaryClass] = calculateDivineBonusSpells(int(character.Wisdom))
		}
	case "Ranger":
		// Rangers get druid spells at level 8 and mage spells at level 9
		if character.Level >= 8 {
			primaryClass = "Druid"
			result[primaryClass] = calculateDivineBonusSpells(int(character.Wisdom))
		}
		if character.Level >= 9 {
			secondaryClass = "Magician"
			result[secondaryClass] = calculateArcaneBonusSpells(int(character.Intelligence))
		}
	case "Bard":
		// Bards might use both illusion and druid spells
		if character.Level >= 2 {
			primaryClass = "Illusionist"
			secondaryClass = "Druid"
			result[primaryClass] = calculateArcaneBonusSpells(int(character.Intelligence))
			result[secondaryClass] = calculateDivineBonusSpells(int(character.Wisdom))
		}
	}

	return result, nil
}

// Helper function for Intelligence-based bonus spells
func calculateArcaneBonusSpells(intelligence int) map[string]int {
	bonusSpells := make(map[string]int)

	// Apply intelligence bonuses
	if intelligence >= 13 {
		bonusSpells["level1"] = 1
	}
	if intelligence >= 15 {
		bonusSpells["level2"] = 1
	}
	if intelligence >= 17 {
		bonusSpells["level3"] = 1
	}
	if intelligence >= 18 {
		bonusSpells["level4"] = 1
	}

	return bonusSpells
}

// Helper function for Wisdom-based bonus spells
func calculateDivineBonusSpells(wisdom int) map[string]int {
	bonusSpells := make(map[string]int)

	// Apply wisdom bonuses
	if wisdom >= 13 {
		bonusSpells["level1"] = 1
	}
	if wisdom >= 15 {
		bonusSpells["level2"] = 1
	}
	if wisdom >= 17 {
		bonusSpells["level3"] = 1
	}
	if wisdom >= 18 {
		bonusSpells["level4"] = 1
	}

	return bonusSpells
}
