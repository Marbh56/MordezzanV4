-- name: GetAmmo :one
SELECT * FROM ammo
WHERE id = ? LIMIT 1;

-- name: GetAmmoByName :one
SELECT * FROM ammo
WHERE name = ? LIMIT 1;

-- name: ListAmmo :many
SELECT * FROM ammo
ORDER BY name;

-- name: CreateAmmo :execresult
INSERT INTO ammo (
  name, cost, weight
) VALUES (
  ?, ?, ?
);

-- name: UpdateAmmo :execresult
UPDATE ammo
SET name = ?,
    cost = ?,
    weight = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteAmmo :execresult
DELETE FROM ammo
WHERE id = ?;