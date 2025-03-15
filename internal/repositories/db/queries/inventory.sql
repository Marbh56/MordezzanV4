-- name: GetInventory :one
SELECT * FROM inventories
WHERE id = ? LIMIT 1;

-- name: GetInventoryByCharacter :one
SELECT * FROM inventories
WHERE character_id = ? LIMIT 1;

-- name: ListInventories :many
SELECT * FROM inventories
ORDER BY id;

-- name: CreateInventory :execresult
INSERT INTO inventories (
    character_id,
    max_weight,
    current_weight
) VALUES (
    ?,
    ?,
    0
);

-- name: UpdateInventory :execresult
UPDATE inventories
SET 
    max_weight = COALESCE(sqlc.narg(max_weight), max_weight),
    current_weight = COALESCE(sqlc.narg(current_weight), current_weight)
WHERE id = ?;

-- name: DeleteInventory :exec
DELETE FROM inventories
WHERE id = ?;

-- name: GetInventoryItems :many
SELECT * FROM inventory_items
WHERE inventory_id = ?
ORDER BY id;

-- name: GetInventoryItem :one
SELECT * FROM inventory_items
WHERE id = ? LIMIT 1;

-- name: GetInventoryItemsByType :many
SELECT * FROM inventory_items
WHERE inventory_id = ? AND item_type = ?
ORDER BY id;

-- name: GetInventoryItemByTypeAndItemID :one
SELECT * FROM inventory_items
WHERE inventory_id = ? AND item_type = ? AND item_id = ?
LIMIT 1;

-- name: AddInventoryItem :execresult
INSERT INTO inventory_items (
    inventory_id,
    item_type,
    item_id,
    quantity,
    is_equipped,
    notes
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

-- name: UpdateInventoryItem :execresult
UPDATE inventory_items
SET 
    quantity = COALESCE(sqlc.narg(quantity), quantity),
    is_equipped = COALESCE(sqlc.narg(is_equipped), is_equipped),
    notes = COALESCE(sqlc.narg(notes), notes)
WHERE id = ?;

-- name: RemoveInventoryItem :exec
DELETE FROM inventory_items
WHERE id = ?;

-- name: RemoveAllInventoryItems :exec
DELETE FROM inventory_items
WHERE inventory_id = ?;