-- +goose Up
CREATE TABLE magic_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    item_type TEXT NOT NULL, -- 'Rod', 'Wand', or 'Staff'
    description TEXT NOT NULL,
    charges INTEGER,
    cost REAL NOT NULL,
    weight INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE magic_items;
