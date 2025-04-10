-- name: GetBarbarianAbilities :many
-- Gets all barbarian abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM barbarian_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetBerserkerAbilities :many
-- Gets all berserker abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM berserker_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetBardAbilities :many
-- Gets all bard abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM bard_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetCataphractAbilities :many
-- Gets all cataphract abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM cataphract_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetClericAbilities :many
-- Gets all cleric abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM cleric_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetCryomancerAbilities :many
-- Gets all cryomancer abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM cryomancer_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetDruidAbilities :many
-- Gets all druid abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM druid_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetFighterAbilities :many
-- Gets all fighter abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM fighter_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetHuntsmanAbilities :many
-- Gets all huntsman abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM huntsman_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetIllusionistAbilities :many
-- Gets all illusionist abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM illusionist_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetLegerdemainistAbilities :many
-- Gets all legerdemainist abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM legerdemainist_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetMagicianAbilities :many
-- Gets all magician abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM magician_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetMonkAbilities :many
-- Gets all monk abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM monk_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetNecromancerAbilities :many
-- Gets all necromancer abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM necromancer_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetPaladinAbilities :many
-- Gets all paladin abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM paladin_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetPriestAbilities :many
-- Gets all priest abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM priest_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetPurloinerAbilities :many
-- Gets all purloiner abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM purloiner_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetPyromancerAbilities :many
-- Gets all pyromancer abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM pyromancer_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetRangerAbilities :many
-- Gets all ranger abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM ranger_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetRunegraverAbilities :many
-- Gets all runegraver abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM runegraver_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetScoutAbilities :many
-- Gets all scout abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM scout_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetShamanAbilities :many
-- Gets all shaman abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM shaman_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetThiefAbilities :many
-- Gets all thief abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM thief_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetWarlockAbilities :many
-- Gets all warlock abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM warlock_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;

-- name: GetWitchAbilities :many
-- Gets all witch abilities available to a character based on their level
SELECT id, name, description, min_level 
FROM witch_abilities
WHERE min_level <= sqlc.arg(character_level)
ORDER BY min_level, name;


