-- name: GetCharacter :one
SELECT * FROM characters
WHERE id = ? LIMIT 1;

-- name: GetCharactersByUser :many
SELECT * FROM characters
WHERE user_id = ?
ORDER BY name;

-- name: ListCharacters :many
SELECT * FROM characters
ORDER BY name;

-- name: CreateCharacter :execresult
INSERT INTO characters (
  user_id, name, strength, dexterity, constitution, 
  wisdom, intelligence, charisma, hit_points
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateCharacter :execresult
UPDATE characters
SET name = ?, 
    strength = ?, 
    dexterity = ?, 
    constitution = ?, 
    wisdom = ?, 
    intelligence = ?, 
    charisma = ?, 
    hit_points = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteCharacter :execresult
DELETE FROM characters
WHERE id = ?;