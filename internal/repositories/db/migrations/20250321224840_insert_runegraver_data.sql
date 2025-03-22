-- +goose Up
CREATE TABLE IF NOT EXISTS runes_per_day (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    level1 INTEGER DEFAULT 0,
    level2 INTEGER DEFAULT 0,
    level3 INTEGER DEFAULT 0,
    level4 INTEGER DEFAULT 0,
    level5 INTEGER DEFAULT 0,
    level6 INTEGER DEFAULT 0,
    UNIQUE(class_name, level),
    FOREIGN KEY (class_name, level) REFERENCES class_data(class_name, level)
);

INSERT INTO class_data (
    class_name, 
    level, 
    experience_points,
    hit_dice, 
    saving_throw, 
    fighting_ability, 
    casting_ability
) VALUES
    ('Runegraver', 1, 0, '1d8', 16, 1, 1),
    ('Runegraver', 2, 3000, '2d8', 16, 2, 2),
    ('Runegraver', 3, 6000, '3d8', 15, 3, 3),
    ('Runegraver', 4, 12000, '4d8', 15, 4, 4),
    ('Runegraver', 5, 24000, '5d8', 14, 5, 5),
    ('Runegraver', 6, 48000, '6d8', 14, 6, 6),
    ('Runegraver', 7, 96000, '7d8', 13, 7, 7),
    ('Runegraver', 8, 192000, '8d8', 13, 8, 8),
    ('Runegraver', 9, 384000, '9d8', 12, 9, 9),
    ('Runegraver', 10, 576000, '9d8+2', 12, 10, 10),
    ('Runegraver', 11, 768000, '9d8+4', 11, 11, 11),
    ('Runegraver', 12, 960000, '9d8+6', 11, 12, 12);

INSERT INTO runes_per_day (
    class_name,
    level,
    level1,
    level2,
    level3,
    level4,
    level5,
    level6
) VALUES
    ('Runegraver', 1, 1, 0, 0, 0, 0, 0),
    ('Runegraver', 2, 2, 0, 0, 0, 0, 0),
    ('Runegraver', 3, 3, 1, 0, 0, 0, 0),
    ('Runegraver', 4, 3, 2, 0, 0, 0, 0),
    ('Runegraver', 5, 3, 3, 1, 0, 0, 0),
    ('Runegraver', 6, 3, 3, 2, 0, 0, 0),
    ('Runegraver', 7, 3, 3, 3, 1, 0, 0),
    ('Runegraver', 8, 3, 3, 3, 2, 0, 0),
    ('Runegraver', 9, 3, 3, 3, 3, 1, 0),
    ('Runegraver', 10, 3, 3, 3, 3, 2, 0),
    ('Runegraver', 11, 3, 3, 3, 3, 3, 0),
    ('Runegraver', 12, 3, 3, 3, 3, 3, 1);

-- +goose Down
DELETE FROM runes_per_day WHERE class_name = 'Runegraver';
DELETE FROM class_data WHERE class_name = 'Runegraver';