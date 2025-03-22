-- +goose Up
CREATE TABLE IF NOT EXISTS purloiner_turning_ability (
    level INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    PRIMARY KEY (level)
);

INSERT INTO class_data (
    class_name, 
    level, 
    experience_points, 
    hit_dice, 
    saving_throw, 
    fighting_ability,
    casting_ability, 
    spell_slots_level1, 
    spell_slots_level2, 
    spell_slots_level3
) VALUES
    ('Purloiner', 1, 0, '1d6', 16, 1, 1, 1, 0, 0),
    ('Purloiner', 2, 2500, '2d6', 16, 1, 2, 1, 0, 0),
    ('Purloiner', 3, 5000, '3d6', 15, 2, 3, 1, 1, 0),
    ('Purloiner', 4, 10000, '4d6', 15, 3, 4, 1, 1, 0),
    ('Purloiner', 5, 20000, '5d6', 14, 3, 5, 1, 1, 1),
    ('Purloiner', 6, 40000, '6d6', 14, 4, 6, 1, 1, 1),
    ('Purloiner', 7, 80000, '7d6', 13, 5, 7, 2, 1, 1),
    ('Purloiner', 8, 160000, '8d6', 13, 5, 8, 2, 2, 1),
    ('Purloiner', 9, 320000, '9d6', 12, 6, 9, 2, 2, 2),
    ('Purloiner', 10, 480000, '9d6+2', 12, 7, 10, 3, 2, 2),
    ('Purloiner', 11, 640000, '9d6+4', 11, 7, 11, 3, 3, 2),
    ('Purloiner', 12, 800000, '9d6+6', 11, 8, 12, 3, 3, 3);

INSERT INTO purloiner_turning_ability (level, turning_ability) VALUES
    (1, 0),
    (2, 0),
    (3, 1),
    (4, 2),
    (5, 3),
    (6, 4),
    (7, 5),
    (8, 6),
    (9, 7),
    (10, 8),
    (11, 9),
    (12, 10);

-- +goose Down
DELETE FROM purloiner_turning_ability;
DELETE FROM class_data WHERE class_name = 'Purloiner';
DROP TABLE IF EXISTS purloiner_turning_ability;