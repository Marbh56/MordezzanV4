-- +goose Up
CREATE TABLE spellbooks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    total_pages INTEGER NOT NULL,
    used_pages INTEGER NOT NULL DEFAULT 0,
    value INTEGER NOT NULL,
    weight REAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE spellbooks;