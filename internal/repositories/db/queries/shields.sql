-- name: GetShield :one
SELECT * FROM shields
WHERE id = ? LIMIT 1;

-- name: GetShieldByName :one
SELECT * FROM shields
WHERE name = ? LIMIT 1;

-- name: ListShields :many
SELECT * FROM shields
ORDER BY name;

-- name: CreateShield :execresult
INSERT INTO shields (
  name, cost, weight, defense_modifier
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateShield :execresult
UPDATE shields
SET name = ?,
    cost = ?,
    weight = ?,
    defense_modifier = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteShield :execresult
DELETE FROM shields
WHERE id = ?;