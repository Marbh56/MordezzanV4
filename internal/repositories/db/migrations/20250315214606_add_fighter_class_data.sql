-- +goose Up
CREATE TABLE fighter_class_data (
    level INTEGER PRIMARY KEY,
    experience_points INTEGER NOT NULL,
    hit_dice TEXT NOT NULL,
    saving_throw INTEGER NOT NULL,
    fighting_ability INTEGER NOT NULL
);

-- Insert the fighter class level data
INSERT INTO fighter_class_data (level, experience_points, hit_dice, saving_throw, fighting_ability)
VALUES 
    (1, 0, '1d10', 16, 1),
    (2, 2000, '2d10', 16, 2),
    (3, 4000, '3d10', 15, 3),
    (4, 8000, '4d10', 15, 4),
    (5, 16000, '5d10', 14, 5),
    (6, 32000, '6d10', 14, 6),
    (7, 64000, '7d10', 13, 7),
    (8, 128000, '8d10', 13, 8),
    (9, 256000, '9d10', 12, 9),
    (10, 384000, '9d10+3', 12, 10),
    (11, 512000, '9d10+6', 11, 11),
    (12, 640000, '9d10+9', 11, 12);

-- +goose Down
DROP IF EXISTS fighter_class_data;
