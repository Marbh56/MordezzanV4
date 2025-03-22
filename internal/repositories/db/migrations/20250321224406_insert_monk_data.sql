-- +goose Up
CREATE TABLE IF NOT EXISTS monk_ac_bonus (
    level INTEGER NOT NULL,
    ac_bonus INTEGER NOT NULL,
    PRIMARY KEY (level)
);

CREATE TABLE IF NOT EXISTS monk_empty_hand_damage (
    level INTEGER NOT NULL,
    damage TEXT NOT NULL,
    PRIMARY KEY (level)
);

INSERT INTO class_data (
    class_name, 
    level, 
    experience_points, 
    hit_dice, 
    saving_throw, 
    fighting_ability
) VALUES
    ('Monk', 1, 0, '1d8', 16, 0),
    ('Monk', 2, 2500, '2d8', 16, 1),
    ('Monk', 3, 5000, '3d8', 15, 2),
    ('Monk', 4, 10000, '4d8', 15, 3),
    ('Monk', 5, 20000, '5d8', 14, 4),
    ('Monk', 6, 40000, '6d8', 14, 5),
    ('Monk', 7, 80000, '7d8', 13, 6),
    ('Monk', 8, 160000, '8d8', 13, 7),
    ('Monk', 9, 320000, '9d8', 12, 8),
    ('Monk', 10, 480000, '9d8+2', 12, 9),
    ('Monk', 11, 640000, '9d8+4', 11, 10),
    ('Monk', 12, 800000, '9d8+6', 11, 11);

INSERT INTO monk_ac_bonus (level, ac_bonus) VALUES
    (1, 1),
    (2, 1),
    (3, 2),
    (4, 2),
    (5, 3),
    (6, 3),
    (7, 4),
    (8, 4),
    (9, 5),
    (10, 5),
    (11, 6),
    (12, 6);

INSERT INTO monk_empty_hand_damage (level, damage) VALUES
    (1, '1d4'),
    (2, '1d4'),
    (3, '1d4'),
    (4, '2d4'),
    (5, '2d4'),
    (6, '2d4'),
    (7, '3d4'),
    (8, '3d4'),
    (9, '3d4'),
    (10, '4d4'),
    (11, '4d4'),
    (12, '4d4');

-- +goose Down
DELETE FROM monk_empty_hand_damage;
DELETE FROM monk_ac_bonus;
DELETE FROM class_data WHERE class_name = 'Monk';
DROP TABLE IF EXISTS monk_empty_hand_damage;
DROP TABLE IF EXISTS monk_ac_bonus;