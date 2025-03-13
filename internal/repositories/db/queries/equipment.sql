-- name: GetEquipment :one
SELECT * FROM equipment
WHERE id = ? LIMIT 1;

-- name: GetEquipmentByName :one
SELECT * FROM equipment
WHERE name = ? LIMIT 1;

-- name: ListEquipment :many
SELECT * FROM equipment
ORDER BY name;

-- name: CreateEquipment :execresult
INSERT INTO equipment (
  name, description, cost, weight
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateEquipment :execresult
UPDATE equipment
SET name = ?,
    description = ?,
    cost = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteEquipment :execresult
DELETE FROM equipment
WHERE id = ?;