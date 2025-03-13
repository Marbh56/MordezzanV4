-- name: GetSpell :one
SELECT * FROM spells
WHERE id = ? LIMIT 1;

-- name: GetSpellsByCharacter :many
SELECT * FROM spells
WHERE character_id = ?
ORDER BY name;

-- name: ListSpells :many
SELECT * FROM spells
ORDER BY name;

-- name: CreateSpell :execresult
INSERT INTO spells (
  character_id, name, 
  mag_level, cry_level, ill_level, nec_level, 
  pyr_level, wch_level, clr_level, drd_level,
  range, duration, area_of_effect, components, description
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateSpell :execresult
UPDATE spells
SET name = ?,
    mag_level = ?,
    cry_level = ?,
    ill_level = ?,
    nec_level = ?,
    pyr_level = ?,
    wch_level = ?,
    clr_level = ?,
    drd_level = ?,
    range = ?,
    duration = ?,
    area_of_effect = ?,
    components = ?,
    description = ?,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteSpell :execresult
DELETE FROM spells
WHERE id = ?;