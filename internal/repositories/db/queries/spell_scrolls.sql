-- name: GetSpellScroll :one
SELECT 
    ss.id, 
    ss.spell_id, 
    s.name as spell_name,
    ss.casting_level, 
    ss.cost, 
    ss.weight, 
    ss.description, 
    ss.created_at, 
    ss.updated_at
FROM spell_scrolls ss
JOIN spells s ON ss.spell_id = s.id
WHERE ss.id = ? LIMIT 1;

-- name: ListSpellScrolls :many
SELECT 
    ss.id, 
    ss.spell_id, 
    s.name as spell_name,
    ss.casting_level, 
    ss.cost, 
    ss.weight, 
    ss.description, 
    ss.created_at, 
    ss.updated_at
FROM spell_scrolls ss
JOIN spells s ON ss.spell_id = s.id
ORDER BY s.name;

-- name: GetSpellScrollsBySpell :many
SELECT 
    ss.id, 
    ss.spell_id, 
    s.name as spell_name,
    ss.casting_level, 
    ss.cost, 
    ss.weight, 
    ss.description, 
    ss.created_at, 
    ss.updated_at
FROM spell_scrolls ss
JOIN spells s ON ss.spell_id = s.id
WHERE ss.spell_id = ?
ORDER BY ss.casting_level;

-- name: CreateSpellScroll :execresult
INSERT INTO spell_scrolls (
  spell_id, casting_level, cost, weight, description
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: UpdateSpellScroll :execresult
UPDATE spell_scrolls
SET spell_id = ?,
    casting_level = ?,
    cost = ?,
    weight = ?,
    description = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteSpellScroll :execresult
DELETE FROM spell_scrolls
WHERE id = ?;