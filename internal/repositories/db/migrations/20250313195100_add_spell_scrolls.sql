-- +goose Up
CREATE TABLE spell_scrolls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    spell_id INTEGER NOT NULL,
    casting_level INTEGER NOT NULL DEFAULT 1,
    cost REAL NOT NULL,
    weight INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (spell_id) REFERENCES spells (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE spell_scrolls;
