-- +goose Up
-- Create tables for the spell system

-- Known spells table
CREATE TABLE IF NOT EXISTS known_spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    spell_id INTEGER NOT NULL,
    spell_name TEXT NOT NULL,
    spell_level INTEGER NOT NULL,
    spell_class TEXT NOT NULL,
    is_memorized BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (spell_id) REFERENCES spells(id),
    UNIQUE(character_id, spell_id)
);

-- Prepared spells table
CREATE TABLE IF NOT EXISTS prepared_spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    spell_id INTEGER NOT NULL,
    spell_name TEXT NOT NULL,
    spell_level INTEGER NOT NULL,
    spell_class TEXT NOT NULL,
    slot_index INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (spell_id) REFERENCES spells(id),
    UNIQUE(character_id, spell_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_known_spells_character ON known_spells (character_id);
CREATE INDEX IF NOT EXISTS idx_prepared_spells_character ON prepared_spells (character_id);
CREATE INDEX IF NOT EXISTS idx_known_spells_spell ON known_spells (spell_id);
CREATE INDEX IF NOT EXISTS idx_prepared_spells_spell ON prepared_spells (spell_id);

-- +goose Down
-- Drop the tables
DROP TABLE IF EXISTS prepared_spells;
DROP TABLE IF EXISTS known_spells;