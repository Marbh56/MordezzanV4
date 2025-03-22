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

-- name: GetThiefSkillsForClass :many
SELECT 
    ts.id, 
    ts.skill_name, 
    ts.attribute 
FROM thief_skills ts
JOIN class_thief_skill_mapping ctsm ON ts.id = ctsm.skill_id
WHERE ctsm.class_name = ?;

-- name: GetThiefSkillsForCharacter :many
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
    CAST(SUBSTR(tsp.level_range, INSTR(tsp.level_range, '-') + 1) AS INTEGER);

-- name: AddThiefSkill :execresult
INSERT INTO thief_skills (skill_name, attribute)
VALUES (?, ?);

-- name: AddThiefSkillProgression :exec
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
VALUES (?, ?, ?);

-- name: AssignSkillToClass :exec
INSERT INTO class_thief_skill_mapping (class_name, skill_id)
VALUES (?, ?);

-- name: RemoveSkillFromClass :exec
DELETE FROM class_thief_skill_mapping
WHERE class_name = ? AND skill_id = ?;

-- name: GetThiefSkillByName :one
SELECT id, skill_name, attribute
FROM thief_skills
WHERE skill_name = ?;