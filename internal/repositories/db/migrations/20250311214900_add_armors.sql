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

INSERT INTO armors (name, armor_type, ac, cost, damage_reduction, weight, weight_class, movement_rate, created_at, updated_at)
VALUES 
    ('None', '-', 9, 0, 0, 0, '-', 40, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Padded', 'Padded', 8, 10, 0, 10, 'Light', 40, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Leather', 'Leather', 7, 15, 0, 15, 'Light', 40, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Studded', 'Studded', 6, 25, 0, 20, 'Light', 40, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Scale Mail', 'Scale Mail', 6, 50, 1, 25, 'Medium', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Chain Mail', 'Chain Mail', 5, 75, 1, 30, 'Medium', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Laminated', 'Laminated', 5, 75, 1, 30, 'Medium', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Banded Mail', 'Banded Mail', 4, 150, 1, 35, 'Medium', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Splint Mail', 'Splint Mail', 4, 150, 1, 35, 'Medium', 30, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Plate Mail', 'Plate Mail', 3, 350, 2, 40, 'Heavy', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Field Plate', 'Field Plate', 2, 1000, 2, 50, 'Heavy', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('Full Plate', 'Full Plate', 1, 2000, 2, 60, 'Heavy', 20, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);


-- +goose Down
DROP TABLE armors;
