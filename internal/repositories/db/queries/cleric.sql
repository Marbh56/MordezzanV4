-- name: GetClericClassData :one
SELECT * FROM cleric_class_data 
WHERE level = ? LIMIT 1;

-- name: ListClericClassData :many
SELECT * FROM cleric_class_data 
ORDER BY level ASC;

-- name: GetNextClericLevel :one
SELECT * FROM cleric_class_data 
WHERE level > ? 
ORDER BY level ASC 
LIMIT 1;