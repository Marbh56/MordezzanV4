-- name: GetWeapon :one
SELECT * FROM weapons
WHERE id = ? LIMIT 1;

-- name: GetWeaponByName :one
SELECT * FROM weapons
WHERE name = ? LIMIT 1;

-- name: ListWeapons :many
SELECT * FROM weapons
ORDER BY name;

-- name: CreateWeapon :execresult
INSERT INTO weapons (
  name, category, weapon_class, cost, weight,  -- Changed from weight_class to weapon_class
  range_short, range_medium, range_long, rate_of_fire, 
  damage, damage_two_handed, properties
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateWeapon :execresult
UPDATE weapons
SET name = ?,
    category = ?,
    weapon_class = ?,
    cost = ?,
    weight = ?,
    range_short = ?,
    range_medium = ?,
    range_long = ?,
    rate_of_fire = ?,
    damage = ?,
    damage_two_handed = ?,
    properties = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteWeapon :execresult
DELETE FROM weapons
WHERE id = ?;