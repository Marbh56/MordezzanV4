// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: class.sql

package db

import (
	"context"
	"database/sql"
)

const getAllClassData = `-- name: GetAllClassData :many
SELECT id, class_name, level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM class_data
WHERE class_name = ?
ORDER BY level
`

func (q *Queries) GetAllClassData(ctx context.Context, className string) ([]ClassDatum, error) {
	rows, err := q.query(ctx, q.getAllClassDataStmt, getAllClassData, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ClassDatum{}
	for rows.Next() {
		var i ClassDatum
		if err := rows.Scan(
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

const getBardDruidSpells = `-- name: GetBardDruidSpells :one
SELECT level, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4 FROM bard_druid_spells
WHERE level = ?
`

func (q *Queries) GetBardDruidSpells(ctx context.Context, level int64) (BardDruidSpell, error) {
	row := q.queryRow(ctx, q.getBardDruidSpellsStmt, getBardDruidSpells, level)
	var i BardDruidSpell
	err := row.Scan(
		&i.Level,
		&i.SpellSlotsLevel1,
		&i.SpellSlotsLevel2,
		&i.SpellSlotsLevel3,
		&i.SpellSlotsLevel4,
	)
	return i, err
}

const getBardIllusionistSpells = `-- name: GetBardIllusionistSpells :one
SELECT level, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4 FROM bard_illusionist_spells
WHERE level = ?
`

func (q *Queries) GetBardIllusionistSpells(ctx context.Context, level int64) (BardIllusionistSpell, error) {
	row := q.queryRow(ctx, q.getBardIllusionistSpellsStmt, getBardIllusionistSpells, level)
	var i BardIllusionistSpell
	err := row.Scan(
		&i.Level,
		&i.SpellSlotsLevel1,
		&i.SpellSlotsLevel2,
		&i.SpellSlotsLevel3,
		&i.SpellSlotsLevel4,
	)
	return i, err
}

const getBerserkerNaturalAC = `-- name: GetBerserkerNaturalAC :one
SELECT natural_ac FROM berserker_natural_ac
WHERE class_name = ? AND level = ?
`

type GetBerserkerNaturalACParams struct {
	ClassName string
	Level     int64
}

func (q *Queries) GetBerserkerNaturalAC(ctx context.Context, arg GetBerserkerNaturalACParams) (int64, error) {
	row := q.queryRow(ctx, q.getBerserkerNaturalACStmt, getBerserkerNaturalAC, arg.ClassName, arg.Level)
	var natural_ac int64
	err := row.Scan(&natural_ac)
	return natural_ac, err
}

const getClassAbilities = `-- name: GetClassAbilities :many
SELECT a.id, a.name, a.description, cam.min_level
FROM abilities a
JOIN class_ability_mapping cam ON a.id = cam.ability_id
WHERE cam.class_name = ?
ORDER BY cam.min_level, a.name
`

type GetClassAbilitiesRow struct {
	ID          int64
	Name        string
	Description string
	MinLevel    int64
}

func (q *Queries) GetClassAbilities(ctx context.Context, className string) ([]GetClassAbilitiesRow, error) {
	rows, err := q.query(ctx, q.getClassAbilitiesStmt, getClassAbilities, className)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetClassAbilitiesRow{}
	for rows.Next() {
		var i GetClassAbilitiesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinLevel,
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

const getClassAbilitiesByLevel = `-- name: GetClassAbilitiesByLevel :many
SELECT a.id, a.name, a.description, cam.min_level
FROM abilities a
JOIN class_ability_mapping cam ON a.id = cam.ability_id
WHERE cam.class_name = ? AND cam.min_level <= ?
ORDER BY cam.min_level, a.name
`

type GetClassAbilitiesByLevelParams struct {
	ClassName string
	MinLevel  int64
}

type GetClassAbilitiesByLevelRow struct {
	ID          int64
	Name        string
	Description string
	MinLevel    int64
}

func (q *Queries) GetClassAbilitiesByLevel(ctx context.Context, arg GetClassAbilitiesByLevelParams) ([]GetClassAbilitiesByLevelRow, error) {
	rows, err := q.query(ctx, q.getClassAbilitiesByLevelStmt, getClassAbilitiesByLevel, arg.ClassName, arg.MinLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetClassAbilitiesByLevelRow{}
	for rows.Next() {
		var i GetClassAbilitiesByLevelRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.MinLevel,
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

const getClassData = `-- name: GetClassData :one
SELECT id, class_name, level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM class_data
WHERE class_name = ? AND level = ?
`

type GetClassDataParams struct {
	ClassName string
	Level     int64
}

func (q *Queries) GetClassData(ctx context.Context, arg GetClassDataParams) (ClassDatum, error) {
	row := q.queryRow(ctx, q.getClassDataStmt, getClassData, arg.ClassName, arg.Level)
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

const getClericTurningAbility = `-- name: GetClericTurningAbility :one
SELECT turning_ability FROM cleric_turning_ability
WHERE class_name = 'Cleric' AND level = ?
`

func (q *Queries) GetClericTurningAbility(ctx context.Context, level int64) (int64, error) {
	row := q.queryRow(ctx, q.getClericTurningAbilityStmt, getClericTurningAbility, level)
	var turning_ability int64
	err := row.Scan(&turning_ability)
	return turning_ability, err
}

const getMonkACBonus = `-- name: GetMonkACBonus :one
SELECT ac_bonus FROM monk_ac_bonus
WHERE level = ?
`

func (q *Queries) GetMonkACBonus(ctx context.Context, level int64) (int64, error) {
	row := q.queryRow(ctx, q.getMonkACBonusStmt, getMonkACBonus, level)
	var ac_bonus int64
	err := row.Scan(&ac_bonus)
	return ac_bonus, err
}

const getMonkEmptyHandDamage = `-- name: GetMonkEmptyHandDamage :one
SELECT damage FROM monk_empty_hand_damage
WHERE level = ?
`

func (q *Queries) GetMonkEmptyHandDamage(ctx context.Context, level int64) (string, error) {
	row := q.queryRow(ctx, q.getMonkEmptyHandDamageStmt, getMonkEmptyHandDamage, level)
	var damage string
	err := row.Scan(&damage)
	return damage, err
}

const getNecromancerTurningAbility = `-- name: GetNecromancerTurningAbility :one
SELECT turning_ability FROM necromancer_turning_ability
WHERE class_name = 'Necromancer' AND level = ?
`

func (q *Queries) GetNecromancerTurningAbility(ctx context.Context, level int64) (int64, error) {
	row := q.queryRow(ctx, q.getNecromancerTurningAbilityStmt, getNecromancerTurningAbility, level)
	var turning_ability int64
	err := row.Scan(&turning_ability)
	return turning_ability, err
}

const getNextLevelData = `-- name: GetNextLevelData :one
SELECT id, class_name, level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM class_data
WHERE class_name = ? AND level > ?
ORDER BY level
LIMIT 1
`

type GetNextLevelDataParams struct {
	ClassName string
	Level     int64
}

func (q *Queries) GetNextLevelData(ctx context.Context, arg GetNextLevelDataParams) (ClassDatum, error) {
	row := q.queryRow(ctx, q.getNextLevelDataStmt, getNextLevelData, arg.ClassName, arg.Level)
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

const getPaladinTurningAbility = `-- name: GetPaladinTurningAbility :one
SELECT turning_ability FROM paladin_turning_ability
WHERE class_name = 'Paladin' AND level = ?
`

func (q *Queries) GetPaladinTurningAbility(ctx context.Context, level int64) (int64, error) {
	row := q.queryRow(ctx, q.getPaladinTurningAbilityStmt, getPaladinTurningAbility, level)
	var turning_ability int64
	err := row.Scan(&turning_ability)
	return turning_ability, err
}

const getRangerDruidSpellSlots = `-- name: GetRangerDruidSpellSlots :many
SELECT spell_level, slots FROM ranger_druid_spell_slots
WHERE class_level <= ?
ORDER BY spell_level
`

type GetRangerDruidSpellSlotsRow struct {
	SpellLevel int64
	Slots      int64
}

func (q *Queries) GetRangerDruidSpellSlots(ctx context.Context, classLevel int64) ([]GetRangerDruidSpellSlotsRow, error) {
	rows, err := q.query(ctx, q.getRangerDruidSpellSlotsStmt, getRangerDruidSpellSlots, classLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRangerDruidSpellSlotsRow{}
	for rows.Next() {
		var i GetRangerDruidSpellSlotsRow
		if err := rows.Scan(&i.SpellLevel, &i.Slots); err != nil {
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

const getRangerMagicianSpellSlots = `-- name: GetRangerMagicianSpellSlots :many
SELECT spell_level, slots FROM ranger_magician_spell_slots
WHERE class_level <= ?
ORDER BY spell_level
`

type GetRangerMagicianSpellSlotsRow struct {
	SpellLevel int64
	Slots      int64
}

func (q *Queries) GetRangerMagicianSpellSlots(ctx context.Context, classLevel int64) ([]GetRangerMagicianSpellSlotsRow, error) {
	rows, err := q.query(ctx, q.getRangerMagicianSpellSlotsStmt, getRangerMagicianSpellSlots, classLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRangerMagicianSpellSlotsRow{}
	for rows.Next() {
		var i GetRangerMagicianSpellSlotsRow
		if err := rows.Scan(&i.SpellLevel, &i.Slots); err != nil {
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

const getRunesPerDay = `-- name: GetRunesPerDay :one
SELECT level1, level2, level3, level4, level5, level6
FROM runes_per_day
WHERE class_name = 'Runegraver' AND level = ?
`

type GetRunesPerDayRow struct {
	Level1 sql.NullInt64
	Level2 sql.NullInt64
	Level3 sql.NullInt64
	Level4 sql.NullInt64
	Level5 sql.NullInt64
	Level6 sql.NullInt64
}

func (q *Queries) GetRunesPerDay(ctx context.Context, level int64) (GetRunesPerDayRow, error) {
	row := q.queryRow(ctx, q.getRunesPerDayStmt, getRunesPerDay, level)
	var i GetRunesPerDayRow
	err := row.Scan(
		&i.Level1,
		&i.Level2,
		&i.Level3,
		&i.Level4,
		&i.Level5,
		&i.Level6,
	)
	return i, err
}

const getShamanArcaneSpells = `-- name: GetShamanArcaneSpells :one
SELECT level, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM shaman_arcane_spells
WHERE level = ?
`

func (q *Queries) GetShamanArcaneSpells(ctx context.Context, level int64) (ShamanArcaneSpell, error) {
	row := q.queryRow(ctx, q.getShamanArcaneSpellsStmt, getShamanArcaneSpells, level)
	var i ShamanArcaneSpell
	err := row.Scan(
		&i.Level,
		&i.SpellSlotsLevel1,
		&i.SpellSlotsLevel2,
		&i.SpellSlotsLevel3,
		&i.SpellSlotsLevel4,
		&i.SpellSlotsLevel5,
		&i.SpellSlotsLevel6,
	)
	return i, err
}

const getShamanDivineSpells = `-- name: GetShamanDivineSpells :one
SELECT level, spell_slots_level1, spell_slots_level2, spell_slots_level3, spell_slots_level4, spell_slots_level5, spell_slots_level6 FROM shaman_divine_spells
WHERE level = ?
`

func (q *Queries) GetShamanDivineSpells(ctx context.Context, level int64) (ShamanDivineSpell, error) {
	row := q.queryRow(ctx, q.getShamanDivineSpellsStmt, getShamanDivineSpells, level)
	var i ShamanDivineSpell
	err := row.Scan(
		&i.Level,
		&i.SpellSlotsLevel1,
		&i.SpellSlotsLevel2,
		&i.SpellSlotsLevel3,
		&i.SpellSlotsLevel4,
		&i.SpellSlotsLevel5,
		&i.SpellSlotsLevel6,
	)
	return i, err
}
