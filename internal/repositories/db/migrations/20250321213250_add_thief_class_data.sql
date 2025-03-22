-- +goose Up
CREATE TABLE thief_class_data (
    level INTEGER PRIMARY KEY,
    experience_points INTEGER NOT NULL,
    hit_dice TEXT NOT NULL,
    saving_throw INTEGER NOT NULL,
    fighting_ability INTEGER NOT NULL
);

-- Insert the thief class level data
INSERT INTO thief_class_data (level, experience_points, hit_dice, saving_throw, fighting_ability)
VALUES 
    (1, 0, '1d6', 16, 1),
    (2, 1500, '2d6', 16, 1),
    (3, 3000, '3d6', 15, 2),
    (4, 6000, '4d6', 15, 3),
    (5, 12000, '5d6', 14, 3),
    (6, 24000, '6d6', 14, 4),
    (7, 48000, '7d6', 13, 5),
    (8, 96000, '8d6', 13, 5),
    (9, 192000, '9d6', 12, 6),
    (10, 288000, '9d6+2', 12, 7),
    (11, 384000, '9d6+4', 11, 7),
    (12, 480000, '9d6+6', 11, 8);

-- +goose Down
DROP TABLE IF EXISTS thief_class_data;