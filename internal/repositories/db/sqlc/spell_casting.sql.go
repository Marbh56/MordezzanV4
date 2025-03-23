// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: spell_casting.sql

package db

import (
	"context"
	"database/sql"
)

const addKnownSpell = `-- name: AddKnownSpell :execresult
INSERT INTO known_spells (
    character_id, spell_id, spell_name, spell_level, spell_class, notes
) VALUES (
    ?, ?, ?, ?, ?, ?
)
`

type AddKnownSpellParams struct {
	CharacterID int64
	SpellID     int64
	SpellName   string
	SpellLevel  int64
	SpellClass  string
	Notes       sql.NullString
}

func (q *Queries) AddKnownSpell(ctx context.Context, arg AddKnownSpellParams) (sql.Result, error) {
	return q.exec(ctx, q.addKnownSpellStmt, addKnownSpell,
		arg.CharacterID,
		arg.SpellID,
		arg.SpellName,
		arg.SpellLevel,
		arg.SpellClass,
		arg.Notes,
	)
}

const clearPreparedSpells = `-- name: ClearPreparedSpells :exec
DELETE FROM prepared_spells
WHERE character_id = ?
`

func (q *Queries) ClearPreparedSpells(ctx context.Context, characterID int64) error {
	_, err := q.exec(ctx, q.clearPreparedSpellsStmt, clearPreparedSpells, characterID)
	return err
}

const countPreparedSpellsByLevelAndClass = `-- name: CountPreparedSpellsByLevelAndClass :one
SELECT COUNT(*) as count FROM prepared_spells
WHERE character_id = ? AND spell_level = ? AND spell_class = ?
`

type CountPreparedSpellsByLevelAndClassParams struct {
	CharacterID int64
	SpellLevel  int64
	SpellClass  string
}

func (q *Queries) CountPreparedSpellsByLevelAndClass(ctx context.Context, arg CountPreparedSpellsByLevelAndClassParams) (int64, error) {
	row := q.queryRow(ctx, q.countPreparedSpellsByLevelAndClassStmt, countPreparedSpellsByLevelAndClass, arg.CharacterID, arg.SpellLevel, arg.SpellClass)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getCharacterForSpellcasting = `-- name: GetCharacterForSpellcasting :one
SELECT id, user_id, name, class, level, experience_points, strength, dexterity, constitution, wisdom, intelligence, charisma, max_hit_points, current_hit_points, temporary_hit_points, created_at, updated_at FROM characters
WHERE id = ?
`

func (q *Queries) GetCharacterForSpellcasting(ctx context.Context, id int64) (Character, error) {
	row := q.queryRow(ctx, q.getCharacterForSpellcastingStmt, getCharacterForSpellcasting, id)
	var i Character
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Class,
		&i.Level,
		&i.ExperiencePoints,
		&i.Strength,
		&i.Dexterity,
		&i.Constitution,
		&i.Wisdom,
		&i.Intelligence,
		&i.Charisma,
		&i.MaxHitPoints,
		&i.CurrentHitPoints,
		&i.TemporaryHitPoints,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassDataForSpellcasting = `-- name: GetClassDataForSpellcasting :one
SELECT id, class_name, level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM class_data
WHERE class_name = ? AND level = ?
`

type GetClassDataForSpellcastingParams struct {
	ClassName string
	Level     int64
}

func (q *Queries) GetClassDataForSpellcasting(ctx context.Context, arg GetClassDataForSpellcastingParams) (ClassDatum, error) {
	row := q.queryRow(ctx, q.getClassDataForSpellcastingStmt, getClassDataForSpellcasting, arg.ClassName, arg.Level)
	var i ClassDatum
	err := row.Scan(
		&i.ID,
		&i.ClassName,
		&i.Level,
		&i.ExperiencePoints,
		&i.HitDice,
		&i.SavingThrow,
		&i.FightingAbility,
		&i.CastingAbility,
		&i.SpellSlotsLevel1,
		&i.SpellSlotsLevel2,
		&i.SpellSlotsLevel3,
		&i.SpellSlotsLevel4,
		&i.SpellSlotsLevel5,
		&i.SpellSlotsLevel6,
	)
	return i, err
}

const getKnownSpellByCharacterAndSpell = `-- name: GetKnownSpellByCharacterAndSpell :one
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, is_memorized, notes, created_at, updated_at FROM known_spells
WHERE character_id = ? AND spell_id = ?
`

type GetKnownSpellByCharacterAndSpellParams struct {
	CharacterID int64
	SpellID     int64
}

func (q *Queries) GetKnownSpellByCharacterAndSpell(ctx context.Context, arg GetKnownSpellByCharacterAndSpellParams) (KnownSpell, error) {
	row := q.queryRow(ctx, q.getKnownSpellByCharacterAndSpellStmt, getKnownSpellByCharacterAndSpell, arg.CharacterID, arg.SpellID)
	var i KnownSpell
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.SpellID,
		&i.SpellName,
		&i.SpellLevel,
		&i.SpellClass,
		&i.IsMemorized,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getKnownSpells = `-- name: GetKnownSpells :many
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, is_memorized, notes, created_at, updated_at FROM known_spells
WHERE character_id = ?
ORDER BY spell_level, spell_name
`

func (q *Queries) GetKnownSpells(ctx context.Context, characterID int64) ([]KnownSpell, error) {
	rows, err := q.query(ctx, q.getKnownSpellsStmt, getKnownSpells, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []KnownSpell{}
	for rows.Next() {
		var i KnownSpell
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.SpellID,
			&i.SpellName,
			&i.SpellLevel,
			&i.SpellClass,
			&i.IsMemorized,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getKnownSpellsByClass = `-- name: GetKnownSpellsByClass :many
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, is_memorized, notes, created_at, updated_at FROM known_spells
WHERE character_id = ? AND spell_class = ?
ORDER BY spell_level, spell_name
`

type GetKnownSpellsByClassParams struct {
	CharacterID int64
	SpellClass  string
}

func (q *Queries) GetKnownSpellsByClass(ctx context.Context, arg GetKnownSpellsByClassParams) ([]KnownSpell, error) {
	rows, err := q.query(ctx, q.getKnownSpellsByClassStmt, getKnownSpellsByClass, arg.CharacterID, arg.SpellClass)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []KnownSpell{}
	for rows.Next() {
		var i KnownSpell
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.SpellID,
			&i.SpellName,
			&i.SpellLevel,
			&i.SpellClass,
			&i.IsMemorized,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextAvailableSlotIndex = `-- name: GetNextAvailableSlotIndex :one
SELECT COALESCE(MAX(slot_index), 0) + 1 as next_slot_index
FROM prepared_spells
WHERE character_id = ? AND spell_level = ? AND spell_class = ?
`

type GetNextAvailableSlotIndexParams struct {
	CharacterID int64
	SpellLevel  int64
	SpellClass  string
}

func (q *Queries) GetNextAvailableSlotIndex(ctx context.Context, arg GetNextAvailableSlotIndexParams) (int64, error) {
	row := q.queryRow(ctx, q.getNextAvailableSlotIndexStmt, getNextAvailableSlotIndex, arg.CharacterID, arg.SpellLevel, arg.SpellClass)
	var next_slot_index int64
	err := row.Scan(&next_slot_index)
	return next_slot_index, err
}

const getPreparedSpellByCharacterAndSpell = `-- name: GetPreparedSpellByCharacterAndSpell :one
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, slot_index, created_at, updated_at FROM prepared_spells
WHERE character_id = ? AND spell_id = ?
`

type GetPreparedSpellByCharacterAndSpellParams struct {
	CharacterID int64
	SpellID     int64
}

func (q *Queries) GetPreparedSpellByCharacterAndSpell(ctx context.Context, arg GetPreparedSpellByCharacterAndSpellParams) (PreparedSpell, error) {
	row := q.queryRow(ctx, q.getPreparedSpellByCharacterAndSpellStmt, getPreparedSpellByCharacterAndSpell, arg.CharacterID, arg.SpellID)
	var i PreparedSpell
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.SpellID,
		&i.SpellName,
		&i.SpellLevel,
		&i.SpellClass,
		&i.SlotIndex,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPreparedSpells = `-- name: GetPreparedSpells :many
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, slot_index, created_at, updated_at FROM prepared_spells
WHERE character_id = ?
ORDER BY spell_level, slot_index
`

func (q *Queries) GetPreparedSpells(ctx context.Context, characterID int64) ([]PreparedSpell, error) {
	rows, err := q.query(ctx, q.getPreparedSpellsStmt, getPreparedSpells, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PreparedSpell{}
	for rows.Next() {
		var i PreparedSpell
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.SpellID,
			&i.SpellName,
			&i.SpellLevel,
			&i.SpellClass,
			&i.SlotIndex,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPreparedSpellsByClass = `-- name: GetPreparedSpellsByClass :many
SELECT id, character_id, spell_id, spell_name, spell_level, spell_class, slot_index, created_at, updated_at FROM prepared_spells
WHERE character_id = ? AND spell_class = ?
ORDER BY spell_level, slot_index
`

type GetPreparedSpellsByClassParams struct {
	CharacterID int64
	SpellClass  string
}

func (q *Queries) GetPreparedSpellsByClass(ctx context.Context, arg GetPreparedSpellsByClassParams) ([]PreparedSpell, error) {
	rows, err := q.query(ctx, q.getPreparedSpellsByClassStmt, getPreparedSpellsByClass, arg.CharacterID, arg.SpellClass)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PreparedSpell{}
	for rows.Next() {
		var i PreparedSpell
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.SpellID,
			&i.SpellName,
			&i.SpellLevel,
			&i.SpellClass,
			&i.SlotIndex,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSpellForSpellcasting = `-- name: GetSpellForSpellcasting :one
SELECT id, name, mag_level, cry_level, ill_level, nec_level, pyr_level, wch_level, clr_level, drd_level, "range", duration, area_of_effect, components, description, created_at, updated_at FROM spells
WHERE id = ?
`

func (q *Queries) GetSpellForSpellcasting(ctx context.Context, id int64) (Spell, error) {
	row := q.queryRow(ctx, q.getSpellForSpellcastingStmt, getSpellForSpellcasting, id)
	var i Spell
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MagLevel,
		&i.CryLevel,
		&i.IllLevel,
		&i.NecLevel,
		&i.PyrLevel,
		&i.WchLevel,
		&i.ClrLevel,
		&i.DrdLevel,
		&i.Range,
		&i.Duration,
		&i.AreaOfEffect,
		&i.Components,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSpellsByClassLevel = `-- name: GetSpellsByClassLevel :many
SELECT id, name, mag_level, cry_level, ill_level, nec_level, pyr_level, wch_level, clr_level, drd_level, "range", duration, area_of_effect, components, description, created_at, updated_at FROM spells
WHERE
    (? = 'Magician' AND mag_level = ?) OR
    (? = 'Cleric' AND clr_level = ?) OR
    (? = 'Druid' AND drd_level = ?) OR
    (? = 'Illusionist' AND ill_level = ?) OR
    (? = 'Necromancer' AND nec_level = ?) OR
    (? = 'Pyromancer' AND pyr_level = ?) OR
    (? = 'Cryomancer' AND cry_level = ?) OR
    (? = 'Witch' AND wch_level = ?)
ORDER BY name
`

type GetSpellsByClassLevelParams struct {
	Column1  interface{}
	MagLevel int64
	Column3  interface{}
	ClrLevel int64
	Column5  interface{}
	DrdLevel int64
	Column7  interface{}
	IllLevel int64
	Column9  interface{}
	NecLevel int64
	Column11 interface{}
	PyrLevel int64
	Column13 interface{}
	CryLevel int64
	Column15 interface{}
	WchLevel int64
}

func (q *Queries) GetSpellsByClassLevel(ctx context.Context, arg GetSpellsByClassLevelParams) ([]Spell, error) {
	rows, err := q.query(ctx, q.getSpellsByClassLevelStmt, getSpellsByClassLevel,
		arg.Column1,
		arg.MagLevel,
		arg.Column3,
		arg.ClrLevel,
		arg.Column5,
		arg.DrdLevel,
		arg.Column7,
		arg.IllLevel,
		arg.Column9,
		arg.NecLevel,
		arg.Column11,
		arg.PyrLevel,
		arg.Column13,
		arg.CryLevel,
		arg.Column15,
		arg.WchLevel,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Spell{}
	for rows.Next() {
		var i Spell
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MagLevel,
			&i.CryLevel,
			&i.IllLevel,
			&i.NecLevel,
			&i.PyrLevel,
			&i.WchLevel,
			&i.ClrLevel,
			&i.DrdLevel,
			&i.Range,
			&i.Duration,
			&i.AreaOfEffect,
			&i.Components,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markSpellAsMemorized = `-- name: MarkSpellAsMemorized :exec
UPDATE known_spells
SET is_memorized = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

type MarkSpellAsMemorizedParams struct {
	IsMemorized bool
	ID          int64
}

func (q *Queries) MarkSpellAsMemorized(ctx context.Context, arg MarkSpellAsMemorizedParams) error {
	_, err := q.exec(ctx, q.markSpellAsMemorizedStmt, markSpellAsMemorized, arg.IsMemorized, arg.ID)
	return err
}

const markSpellAsMemorizedBySpellID = `-- name: MarkSpellAsMemorizedBySpellID :exec
UPDATE known_spells
SET is_memorized = ?, updated_at = CURRENT_TIMESTAMP
WHERE character_id = ? AND spell_id = ?
`

type MarkSpellAsMemorizedBySpellIDParams struct {
	IsMemorized bool
	CharacterID int64
	SpellID     int64
}

func (q *Queries) MarkSpellAsMemorizedBySpellID(ctx context.Context, arg MarkSpellAsMemorizedBySpellIDParams) error {
	_, err := q.exec(ctx, q.markSpellAsMemorizedBySpellIDStmt, markSpellAsMemorizedBySpellID, arg.IsMemorized, arg.CharacterID, arg.SpellID)
	return err
}

const prepareSpell = `-- name: PrepareSpell :execresult
INSERT INTO prepared_spells (
    character_id, spell_id, spell_name, spell_level, spell_class, slot_index
) VALUES (
    ?, ?, ?, ?, ?, ?
)
`

type PrepareSpellParams struct {
	CharacterID int64
	SpellID     int64
	SpellName   string
	SpellLevel  int64
	SpellClass  string
	SlotIndex   int64
}

func (q *Queries) PrepareSpell(ctx context.Context, arg PrepareSpellParams) (sql.Result, error) {
	return q.exec(ctx, q.prepareSpellStmt, prepareSpell,
		arg.CharacterID,
		arg.SpellID,
		arg.SpellName,
		arg.SpellLevel,
		arg.SpellClass,
		arg.SlotIndex,
	)
}

const removeKnownSpell = `-- name: RemoveKnownSpell :exec
DELETE FROM known_spells
WHERE id = ?
`

func (q *Queries) RemoveKnownSpell(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.removeKnownSpellStmt, removeKnownSpell, id)
	return err
}

const resetAllMemorizedSpells = `-- name: ResetAllMemorizedSpells :exec
UPDATE known_spells
SET is_memorized = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE character_id = ?
`

func (q *Queries) ResetAllMemorizedSpells(ctx context.Context, characterID int64) error {
	_, err := q.exec(ctx, q.resetAllMemorizedSpellsStmt, resetAllMemorizedSpells, characterID)
	return err
}

const unprepareSpell = `-- name: UnprepareSpell :exec
DELETE FROM prepared_spells
WHERE id = ?
`

func (q *Queries) UnprepareSpell(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.unprepareSpellStmt, unprepareSpell, id)
	return err
}
