-- name: GetKnownSpells :many
SELECT * FROM known_spells
WHERE character_id = ?
ORDER BY spell_level, spell_name;

-- name: GetKnownSpellsByClass :many
SELECT * FROM known_spells
WHERE character_id = ? AND spell_class = ?
ORDER BY spell_level, spell_name;

-- name: GetKnownSpellByCharacterAndSpell :one
SELECT * FROM known_spells
WHERE character_id = ? AND spell_id = ?;

-- name: AddKnownSpell :execresult
INSERT INTO known_spells (
    character_id, spell_id, spell_name, spell_level, spell_class, notes
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: MarkSpellAsMemorized :exec
UPDATE known_spells
SET is_memorized = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: MarkSpellAsMemorizedBySpellID :exec
UPDATE known_spells
SET is_memorized = ?, updated_at = CURRENT_TIMESTAMP
WHERE character_id = ? AND spell_id = ?;

-- name: RemoveKnownSpell :exec
DELETE FROM known_spells
WHERE id = ?;

-- name: ResetAllMemorizedSpells :exec
UPDATE known_spells
SET is_memorized = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE character_id = ?;

-- name: GetPreparedSpells :many
SELECT * FROM prepared_spells
WHERE character_id = ?
ORDER BY spell_level, slot_index;

-- name: GetPreparedSpellsByClass :many
SELECT * FROM prepared_spells
WHERE character_id = ? AND spell_class = ?
ORDER BY spell_level, slot_index;

-- name: GetPreparedSpellByCharacterAndSpell :one
SELECT * FROM prepared_spells
WHERE character_id = ? AND spell_id = ?;

-- name: CountPreparedSpellsByLevelAndClass :one
SELECT COUNT(*) as count FROM prepared_spells
WHERE character_id = ? AND spell_level = ? AND spell_class = ?;

-- name: GetNextAvailableSlotIndex :one
SELECT COALESCE(MAX(slot_index), 0) + 1 as next_slot_index
FROM prepared_spells
WHERE character_id = ? AND spell_level = ? AND spell_class = ?;

-- name: PrepareSpell :execresult
INSERT INTO prepared_spells (
    character_id, spell_id, spell_name, spell_level, spell_class, slot_index
) VALUES (
    ?, ?, ?, ?, ?, ?
);

-- name: UnprepareSpell :exec
DELETE FROM prepared_spells
WHERE id = ?;

-- name: ClearPreparedSpells :exec
DELETE FROM prepared_spells
WHERE character_id = ?;

-- name: GetClassDataForSpellcasting :one
SELECT * FROM class_data
WHERE class_name = ? AND level = ?;

-- name: GetCharacterForSpellcasting :one
SELECT * FROM characters
WHERE id = ?;

-- name: GetSpellForSpellcasting :one
SELECT * FROM spells
WHERE id = ?;

-- name: GetSpellsByClassLevel :many
SELECT * FROM spells
WHERE
    (? = 'Magician' AND mag_level = ?) OR
    (? = 'Cleric' AND clr_level = ?) OR
    (? = 'Druid' AND drd_level = ?) OR
    (? = 'Illusionist' AND ill_level = ?) OR
    (? = 'Necromancer' AND nec_level = ?) OR
    (? = 'Pyromancer' AND pyr_level = ?) OR
    (? = 'Cryomancer' AND cry_level = ?) OR
    (? = 'Witch' AND wch_level = ?)
ORDER BY name;