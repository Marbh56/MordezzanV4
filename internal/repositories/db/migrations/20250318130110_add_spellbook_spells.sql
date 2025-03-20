-- +goose Up
CREATE TABLE spellbook_spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    spellbook_id INTEGER NOT NULL,
    spell_id INTEGER NOT NULL,
    character_class TEXT NOT NULL DEFAULT 'Magician',
    pages_used INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (spellbook_id) REFERENCES spellbooks(id) ON DELETE CASCADE,
    FOREIGN KEY (spell_id) REFERENCES spells(id) ON DELETE CASCADE,
    UNIQUE(spellbook_id, spell_id, character_class)
);

-- +goose Down
DROP TABLE spellbook_spells;