-- name: GetSpellbook :one
SELECT * FROM spellbooks WHERE id = ? LIMIT 1;

-- name: GetSpellbookByName :one
SELECT * FROM spellbooks WHERE name = ? LIMIT 1;

-- name: ListSpellbooks :many
SELECT * FROM spellbooks ORDER BY name;

-- name: CreateSpellbook :execresult
INSERT INTO spellbooks (
    name, description, total_pages, used_pages, value, weight
) VALUES (
    ?, ?, ?, 0, ?, ?
);

-- name: UpdateSpellbook :exec
UPDATE spellbooks
SET name = ?, description = ?, total_pages = ?, used_pages = ?, value = ?, weight = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateSpellbookUsedPages :exec
UPDATE spellbooks
SET used_pages = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteSpellbook :exec
DELETE FROM spellbooks WHERE id = ?;

-- name: AddSpellToSpellbook :exec
INSERT INTO spellbook_spells (spellbook_id, spell_id, character_class, pages_used)
VALUES (?, ?, ?, ?);

-- name: RemoveSpellFromSpellbook :exec
DELETE FROM spellbook_spells
WHERE spellbook_id = ? AND spell_id = ?;

-- name: DeleteAllSpellsFromSpellbook :exec
DELETE FROM spellbook_spells
WHERE spellbook_id = ?;

-- name: GetSpellsInSpellbook :many
SELECT spell_id FROM spellbook_spells
WHERE spellbook_id = ?
ORDER BY spell_id;

-- name: GetSpellFromSpellbook :one
SELECT * FROM spellbook_spells
WHERE spellbook_id = ? AND spell_id = ? LIMIT 1;
