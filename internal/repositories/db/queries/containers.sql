-- name: GetContainer :one
SELECT * FROM containers
WHERE id = ? LIMIT 1;

-- name: GetContainerByName :one
SELECT * FROM containers
WHERE name = ? LIMIT 1;

-- name: ListContainers :many
SELECT * FROM containers
ORDER BY name;

-- name: CreateContainer :execresult
INSERT INTO containers (
  name, max_weight, allowed_items, cost, weight
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: UpdateContainer :execresult
UPDATE containers
SET name = ?,
    max_weight = ?,
    allowed_items = ?,
    cost = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteContainer :execresult
DELETE FROM containers
WHERE id = ?;