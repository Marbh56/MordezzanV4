-- name: GetMagicianClassData :one
SELECT * FROM magician_class_data 
WHERE level = ? 
LIMIT 1;

-- name: GetNextMagicianLevel :one
SELECT * FROM magician_class_data 
WHERE level > ? 
ORDER BY level ASC 
LIMIT 1;

-- name: ListMagicianClassData :many
SELECT * FROM magician_class_data 
ORDER BY level ASC;