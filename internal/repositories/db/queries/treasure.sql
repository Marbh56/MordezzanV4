-- name: GetTreasure :one
SELECT * FROM treasures
WHERE id = ? LIMIT 1;

-- name: GetTreasureByCharacter :one
SELECT * FROM treasures
WHERE character_id = ? LIMIT 1;

-- name: ListTreasures :many
SELECT * FROM treasures
ORDER BY character_id;

-- name: CreateTreasure :execresult
INSERT INTO treasures (
  character_id, platinum_coins, gold_coins, electrum_coins,
  silver_coins, copper_coins, gems, art_objects, 
  other_valuables, total_value_gold
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateTreasure :execresult
UPDATE treasures
SET platinum_coins = ?,
    gold_coins = ?,
    electrum_coins = ?,
    silver_coins = ?,
    copper_coins = ?,
    gems = ?,
    art_objects = ?,
    other_valuables = ?,
    total_value_gold = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteTreasure :execresult
DELETE FROM treasures
WHERE id = ?;