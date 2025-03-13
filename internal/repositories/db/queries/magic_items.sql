-- name: GetMagicItem :one
SELECT * FROM magic_items
WHERE id = ? LIMIT 1;

-- name: GetMagicItemByName :one
SELECT * FROM magic_items
WHERE name = ? LIMIT 1;

-- name: ListMagicItems :many
SELECT * FROM magic_items
ORDER BY name;

-- name: ListMagicItemsByType :many
SELECT * FROM magic_items
WHERE item_type = ?
ORDER BY name;

-- name: CreateMagicItem :execresult
INSERT INTO magic_items (
  name, item_type, description, charges, cost, weight
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: UpdateMagicItem :execresult
UPDATE magic_items
SET name = ?,
    item_type = ?,
    description = ?,
    charges = ?,
    cost = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteMagicItem :execresult
DELETE FROM magic_items
WHERE id = ?;