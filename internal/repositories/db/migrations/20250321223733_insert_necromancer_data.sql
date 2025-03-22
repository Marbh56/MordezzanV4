-- +goose Up
CREATE TABLE IF NOT EXISTS necromancer_turning_ability (
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    PRIMARY KEY (class_name, level)
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
    spell_slots_level3, 
    spell_slots_level4, 
    spell_slots_level5, 
    spell_slots_level6
) VALUES
    ('Necromancer', 1, 0, '1d4', 16, 0, 1, 1, 0, 0, 0, 0, 0),
    ('Necromancer', 2, 2500, '2d4', 16, 0, 2, 2, 0, 0, 0, 0, 0),
    ('Necromancer', 3, 5000, '3d4', 15, 1, 3, 2, 1, 0, 0, 0, 0),
    ('Necromancer', 4, 10000, '4d4', 15, 1, 4, 3, 2, 0, 0, 0, 0),
    ('Necromancer', 5, 20000, '5d4', 14, 2, 5, 3, 2, 1, 0, 0, 0),
    ('Necromancer', 6, 40000, '6d4', 14, 2, 6, 4, 3, 2, 0, 0, 0),
    ('Necromancer', 7, 80000, '7d4', 13, 3, 7, 4, 3, 2, 1, 0, 0),
    ('Necromancer', 8, 160000, '8d4', 13, 3, 8, 4, 4, 3, 2, 0, 0),
    ('Necromancer', 9, 320000, '9d4', 12, 4, 9, 5, 4, 3, 2, 1, 0),
    ('Necromancer', 10, 480000, '9d4+1', 12, 4, 10, 5, 4, 4, 3, 2, 0),
    ('Necromancer', 11, 640000, '9d4+2', 11, 5, 11, 5, 5, 4, 3, 2, 1),
    ('Necromancer', 12, 800000, '9d4+3', 11, 5, 12, 5, 5, 4, 4, 3, 2);

-- Insert Necromancer turning ability data
INSERT INTO necromancer_turning_ability (class_name, level, turning_ability) VALUES
    ('Necromancer', 1, 0),
    ('Necromancer', 2, 0),
    ('Necromancer', 3, 1),
    ('Necromancer', 4, 2),
    ('Necromancer', 5, 3),
    ('Necromancer', 6, 4),
    ('Necromancer', 7, 5),
    ('Necromancer', 8, 6),
    ('Necromancer', 9, 7),
    ('Necromancer', 10, 8),
    ('Necromancer', 11, 9),
    ('Necromancer', 12, 10);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

-- Remove turning ability data
DELETE FROM necromancer_turning_ability WHERE class_name = 'Necromancer';

-- Remove class data
DELETE FROM class_data WHERE class_name = 'Necromancer';