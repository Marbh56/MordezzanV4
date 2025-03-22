-- +goose Up
CREATE TABLE containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    max_weight INTEGER NOT NULL,
    allowed_items TEXT NOT NULL,
    cost REAL NOT NULL,
    weight INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE containers;