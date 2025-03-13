-- name: GetArmor :one
SELECT * FROM armors
WHERE id = ? LIMIT 1;

-- name: ListArmors :many
SELECT * FROM armors
ORDER BY name;

-- name: CreateArmor :execresult
INSERT INTO armors (
  name, armor_type, ac, cost,
  damage_reduction, weight, weight_class, movement_rate
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateArmor :execresult
UPDATE armors
SET name = ?,
    armor_type = ?,
    ac = ?,
    cost = ?,
    damage_reduction = ?,
    weight = ?,
    weight_class = ?,
    movement_rate = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteArmor :execresult
DELETE FROM armors
WHERE id = ?;

-- name: GetArmorByName :one
SELECT * FROM armors
WHERE name = ? LIMIT 1;