-- name: GetPreparedSpell :one
SELECT * FROM prepared_spells
WHERE id = ? LIMIT 1;

-- name: GetPreparedSpellsByCharacter :many
SELECT * FROM prepared_spells
WHERE character_id = ?
ORDER BY slot_level, prepared_at;

-- name: CountPreparedSpell :one
SELECT COUNT(*) FROM prepared_spells
WHERE character_id = ? AND spell_id = ?;

-- name: CountPreparedSpellsByLevel :one
SELECT COUNT(*) FROM prepared_spells
WHERE character_id = ? AND slot_level = ?;

-- name: PrepareSpell :execresult
INSERT INTO prepared_spells (
    character_id, spell_id, slot_level
) VALUES (
    ?, ?, ?
);

-- name: UnprepareSpell :exec
DELETE FROM prepared_spells
WHERE character_id = ? AND spell_id = ?;

-- name: ClearPreparedSpells :exec
DELETE FROM prepared_spells
WHERE character_id = ?;