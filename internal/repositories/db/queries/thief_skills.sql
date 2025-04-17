-- name: GetThiefSkillsByLevel :many
SELECT 
    id,
    skill_name,
    attribute,
    level,
    success_chance
FROM thief_skills
WHERE level = ?
ORDER BY skill_name;