-- +goose Up
-- SQL in this section is executed when the migration is applied
CREATE TABLE weapon_masteries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    weapon_base_name TEXT NOT NULL,  -- Changed from weapon_id
    mastery_level TEXT NOT NULL CHECK (mastery_level IN ('mastered', 'grand_mastery')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    UNIQUE(character_id, weapon_base_name)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE weapon_masteries;