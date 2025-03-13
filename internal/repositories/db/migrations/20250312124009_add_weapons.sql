-- +goose Up
CREATE TABLE weapons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    weapon_class INTEGER NOT NULL DEFAULT 1,
    cost REAL NOT NULL,
    weight INTEGER NOT NULL,
    range_short INTEGER,
    range_medium INTEGER,
    range_long INTEGER,
    rate_of_fire TEXT,
    damage TEXT NOT NULL,
    damage_two_handed TEXT,
    properties TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE weapons;
