-- +goose Up
CREATE TABLE treasures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER,
    platinum_coins INTEGER NOT NULL DEFAULT 0,
    gold_coins INTEGER NOT NULL DEFAULT 0,
    electrum_coins INTEGER NOT NULL DEFAULT 0,
    silver_coins INTEGER NOT NULL DEFAULT 0,
    copper_coins INTEGER NOT NULL DEFAULT 0,
    gems TEXT,
    art_objects TEXT,
    other_valuables TEXT,
    total_value_gold REAL NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE treasures;
