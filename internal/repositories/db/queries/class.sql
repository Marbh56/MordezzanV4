-- name: GetClassData :one
SELECT * FROM class_data
WHERE class_name = ? AND level = ?;

-- name: GetAllClassData :many
SELECT * FROM class_data
WHERE class_name = ?
ORDER BY level;

-- name: GetNextLevelData :one
SELECT * FROM class_data
WHERE class_name = ? AND level > ?
ORDER BY level
LIMIT 1;

-- name: GetClericTurningAbility :one
SELECT turning_ability FROM cleric_turning_ability
WHERE class_name = 'Cleric' AND level = ?;

-- name: GetPaladinTurningAbility :one
SELECT turning_ability FROM paladin_turning_ability
WHERE class_name = 'Paladin' AND level = ?;

-- name: GetNecromancerTurningAbility :one
SELECT turning_ability FROM necromancer_turning_ability
WHERE class_name = 'Necromancer' AND level = ?;

-- name: GetMonkACBonus :one
SELECT ac_bonus FROM monk_ac_bonus
WHERE level = ?;

-- name: GetMonkEmptyHandDamage :one
SELECT damage FROM monk_empty_hand_damage
WHERE level = ?;

-- name: GetBerserkerNaturalAC :one
SELECT natural_ac FROM berserker_natural_ac
WHERE class_name = ? AND level = ?;

-- name: GetClassAbilities :many
SELECT a.id, a.name, a.description, cam.min_level
FROM abilities a
JOIN class_ability_mapping cam ON a.id = cam.ability_id
WHERE cam.class_name = ?
ORDER BY cam.min_level, a.name;

-- name: GetClassAbilitiesByLevel :many
SELECT a.id, a.name, a.description, cam.min_level
FROM abilities a
JOIN class_ability_mapping cam ON a.id = cam.ability_id
WHERE cam.class_name = ? AND cam.min_level <= ?
ORDER BY cam.min_level, a.name;

-- name: GetRangerDruidSpellSlots :many
SELECT spell_level, slots FROM ranger_druid_spell_slots
WHERE class_level <= ?
ORDER BY spell_level;

-- name: GetRangerMagicianSpellSlots :many
SELECT spell_level, slots FROM ranger_magician_spell_slots
WHERE class_level <= ?
ORDER BY spell_level;

-- name: GetShamanDivineSpells :one
SELECT * FROM shaman_divine_spells
WHERE level = ?;

-- name: GetShamanArcaneSpells :one
SELECT * FROM shaman_arcane_spells
WHERE level = ?;

-- name: GetBardDruidSpells :one
SELECT * FROM bard_druid_spells
WHERE level = ?;

-- name: GetBardIllusionistSpells :one
SELECT * FROM bard_illusionist_spells
WHERE level = ?;

-- name: GetRunesPerDay :one
SELECT level1, level2, level3, level4, level5, level6
FROM runes_per_day
WHERE class_name = 'Runegraver' AND level = ?;
