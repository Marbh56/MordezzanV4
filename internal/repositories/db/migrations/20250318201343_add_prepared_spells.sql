-- +goose Up
CREATE TABLE prepared_spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    spell_id INTEGER NOT NULL,
    slot_level INTEGER NOT NULL,
    prepared_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (spell_id) REFERENCES spells(id) ON DELETE CASCADE,
    UNIQUE(character_id, spell_id)
);

CREATE INDEX idx_prepared_spells_character_id ON prepared_spells(character_id);

-- +goose Down
DROP TABLE prepared_spells;