-- name: GetPotion :one
SELECT * FROM potions
WHERE id = ? LIMIT 1;

-- name: GetPotionByName :one
SELECT * FROM potions
WHERE name = ? LIMIT 1;

-- name: ListPotions :many
SELECT * FROM potions
ORDER BY name;

-- name: CreatePotion :execresult
INSERT INTO potions (
  name, description, uses, weight
) VALUES (
  ?, ?, ?, ?
);

-- name: UpdatePotion :execresult
UPDATE potions
SET name = ?,
    description = ?,
    uses = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeletePotion :execresult
DELETE FROM potions
WHERE id = ?;