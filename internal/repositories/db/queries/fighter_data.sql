-- name: GetFighterClassData :one
SELECT * FROM fighter_class_data
WHERE level = ? LIMIT 1;

-- name: ListFighterClassData :many
SELECT * FROM fighter_class_data
ORDER BY level;

-- name: GetNextFighterLevel :one
SELECT * FROM fighter_class_data
WHERE level > ? 
ORDER BY level ASC
LIMIT 1;