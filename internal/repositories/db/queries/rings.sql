-- name: GetRing :one
SELECT * FROM rings
WHERE id = ? LIMIT 1;

-- name: GetRingByName :one
SELECT * FROM rings
WHERE name = ? LIMIT 1;

-- name: ListRings :many
SELECT * FROM rings
ORDER BY name;

-- name: CreateRing :execresult
INSERT INTO rings (
  name, description, cost, weight
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdateRing :execresult
UPDATE rings
SET name = ?,
    description = ?,
    cost = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteRing :execresult
DELETE FROM rings
WHERE id = ?;