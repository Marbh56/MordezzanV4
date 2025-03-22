-- name: GetThiefClassData :one
SELECT level, experience_points, hit_dice, saving_throw, fighting_ability
FROM thief_class_data
WHERE level = ? LIMIT 1;

-- name: ListThiefClassData :many
SELECT level, experience_points, hit_dice, saving_throw, fighting_ability
FROM thief_class_data
ORDER BY level;

-- name: GetNextThiefLevel :one
SELECT level, experience_points, hit_dice, saving_throw, fighting_ability
FROM thief_class_data
WHERE level > ? 
ORDER BY level 
LIMIT 1;