-- name: GetWeaponMasteriesByCharacter :many
SELECT wm.id, wm.character_id, wm.weapon_base_name, wm.mastery_level, wm.created_at, wm.updated_at
FROM weapon_masteries wm
WHERE wm.character_id = ? 
ORDER BY wm.mastery_level DESC, wm.weapon_base_name ASC;

-- name: AddWeaponMastery :exec
INSERT INTO weapon_masteries (character_id, weapon_base_name, mastery_level)
VALUES (?, ?, ?);

-- name: UpdateWeaponMasteryLevel :exec
UPDATE weapon_masteries
SET mastery_level = ?, updated_at = CURRENT_TIMESTAMP
WHERE character_id = ? AND weapon_base_name = ?;

-- name: DeleteWeaponMastery :exec
DELETE FROM weapon_masteries
WHERE character_id = ? AND weapon_base_name = ?;

-- name: CountWeaponMasteries :one
SELECT COUNT(*) as count 
FROM weapon_masteries 
WHERE character_id = ? AND mastery_level = ?;

-- name: GetWeaponMasteryByID :one
SELECT id, character_id, weapon_base_name, mastery_level, created_at, updated_at
FROM weapon_masteries
WHERE id = ?;

-- name: GetWeaponMasteryByBaseName :one
SELECT id, character_id, weapon_base_name, mastery_level, created_at, updated_at
FROM weapon_masteries
WHERE character_id = ? AND weapon_base_name = ?;