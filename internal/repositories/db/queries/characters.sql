-- name: GetCharacter :one
SELECT id, user_id, name, class, level, strength, dexterity, constitution,
       wisdom, intelligence, charisma, max_hit_points, current_hit_points, temporary_hit_points, experience_points,
       created_at, updated_at
FROM characters
WHERE id = ? LIMIT 1;

-- name: GetCharactersByUser :many
SELECT id, user_id, name, class, level, strength, dexterity, constitution,
       wisdom, intelligence, charisma, max_hit_points, current_hit_points, temporary_hit_points, experience_points,
       created_at, updated_at
FROM characters
WHERE user_id = ?
ORDER BY name;

-- name: ListCharacters :many
SELECT id, user_id, name, class, level, strength, dexterity, constitution,
       wisdom, intelligence, charisma, max_hit_points, current_hit_points, temporary_hit_points, experience_points,
       created_at, updated_at
FROM characters
ORDER BY name;

-- name: CreateCharacter :execresult
INSERT INTO characters (
  user_id, name, class, level, strength, dexterity, constitution,
  wisdom, intelligence, charisma, max_hit_points, current_hit_points, temporary_hit_points, experience_points
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateCharacter :execresult
UPDATE characters
SET name = ?,
    class = ?,
    level = ?,
    strength = ?,
    dexterity = ?,
    constitution = ?,
    wisdom = ?,
    intelligence = ?,
    charisma = ?,
    max_hit_points = ?,
    current_hit_points = ?,
    temporary_hit_points = ?,
    experience_points = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteCharacter :execresult
DELETE FROM characters
WHERE id = ?;