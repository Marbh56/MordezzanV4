-- +goose Up
CREATE TABLE cleric_class_data (
    level INTEGER PRIMARY KEY,
    experience_points INTEGER NOT NULL,
    hit_dice TEXT NOT NULL,
    saving_throw INTEGER NOT NULL,
    fighting_ability INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    casting_ability INTEGER NOT NULL,
    spell_slot_level1 INTEGER NOT NULL,
    spell_slot_level2 INTEGER NOT NULL,
    spell_slot_level3 INTEGER NOT NULL,
    spell_slot_level4 INTEGER NOT NULL,
    spell_slot_level5 INTEGER NOT NULL,
    spell_slot_level6 INTEGER NOT NULL
);

-- Insert the cleric class level data
INSERT INTO cleric_class_data (
    level, experience_points, hit_dice, saving_throw, 
    fighting_ability, turning_ability, casting_ability, 
    spell_slot_level1, spell_slot_level2, spell_slot_level3, 
    spell_slot_level4, spell_slot_level5, spell_slot_level6
)
VALUES 
    (1, 0, '1d8', 16, 1, 1, 1, 1, 0, 0, 0, 0, 0),
    (2, 2000, '2d8', 16, 1, 2, 2, 2, 0, 0, 0, 0, 0),
    (3, 4000, '3d8', 15, 2, 3, 3, 2, 1, 0, 0, 0, 0),
    (4, 8000, '4d8', 15, 3, 4, 4, 2, 2, 0, 0, 0, 0),
    (5, 16000, '5d8', 14, 3, 5, 5, 3, 2, 1, 0, 0, 0),
    (6, 32000, '6d8', 14, 4, 6, 6, 3, 2, 2, 0, 0, 0),
    (7, 64000, '7d8', 13, 5, 7, 7, 3, 3, 2, 1, 0, 0),
    (8, 128000, '8d8', 13, 5, 8, 8, 3, 3, 2, 2, 0, 0),
    (9, 256000, '9d8', 12, 6, 9, 9, 4, 3, 3, 2, 1, 0),
    (10, 384000, '9d8+2', 12, 7, 10, 10, 4, 3, 3, 2, 2, 0),
    (11, 512000, '9d8+4', 11, 7, 11, 11, 4, 4, 3, 3, 2, 1),
    (12, 640000, '9d8+6', 11, 8, 12, 12, 4, 4, 3, 3, 2, 2);

-- +goose Down
DROP TABLE IF EXISTS cleric_class_data;