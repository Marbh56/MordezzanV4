-- +goose Up
CREATE TABLE armors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    armor_type TEXT NOT NULL,
    ac INTEGER NOT NULL,
    cost REAL NOT NULL,
    damage_reduction INTEGER NOT NULL DEFAULT 0,
    weight INTEGER NOT NULL,
    weight_class TEXT NOT NULL,
    movement_rate INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE armors;
