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

-- name: UpdateInventoryWeight :exec
UPDATE inventories
SET current_weight = ?
WHERE id = ?;

-- name: RecalculateInventoryWeight :exec
UPDATE inventories
SET current_weight = (
    SELECT COALESCE(SUM(ii.quantity * COALESCE(
        CASE ii.item_type
            WHEN 'weapon' THEN (SELECT weight FROM weapons WHERE weapons.id = ii.item_id)
            WHEN 'armor' THEN (SELECT weight FROM armors WHERE armors.id = ii.item_id)
            WHEN 'shield' THEN (SELECT weight FROM shields WHERE shields.id = ii.item_id)
            WHEN 'potion' THEN (SELECT weight FROM potions WHERE potions.id = ii.item_id)
            WHEN 'magic_item' THEN (SELECT weight FROM magic_items WHERE magic_items.id = ii.item_id)
            WHEN 'ring' THEN (SELECT weight FROM rings WHERE rings.id = ii.item_id)
            WHEN 'ammo' THEN (SELECT weight FROM ammo WHERE ammo.id = ii.item_id)
            WHEN 'spell_scroll' THEN (SELECT weight FROM spell_scrolls WHERE spell_scrolls.id = ii.item_id)
            WHEN 'container' THEN (SELECT weight FROM containers WHERE containers.id = ii.item_id)
            WHEN 'equipment' THEN (SELECT weight FROM equipment WHERE equipment.id = ii.item_id)
            ELSE 0.1
        END, 0.1)
    ), 0) FROM inventory_items ii WHERE ii.inventory_id = inventories.id)
WHERE inventories.id = ?;