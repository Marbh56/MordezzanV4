-- +goose Up
CREATE TABLE IF NOT EXISTS priest_turning_ability (
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
    spell_slots_level3, 
    spell_slots_level4, 
    spell_slots_level5, 
    spell_slots_level6
) VALUES
    ('Priest', 1, 0, '1d4', 16, 0, 1, 2, 0, 0, 0, 0, 0),
    ('Priest', 2, 2000, '2d4', 16, 0, 2, 3, 0, 0, 0, 0, 0),
    ('Priest', 3, 4000, '3d4', 15, 1, 3, 3, 2, 0, 0, 0, 0),
    ('Priest', 4, 8000, '4d4', 15, 1, 4, 4, 3, 0, 0, 0, 0),
    ('Priest', 5, 16000, '5d4', 14, 2, 5, 4, 3, 2, 0, 0, 0),
    ('Priest', 6, 32000, '6d4', 14, 2, 6, 4, 4, 3, 0, 0, 0),
    ('Priest', 7, 64000, '7d4', 13, 3, 7, 5, 4, 3, 2, 0, 0),
    ('Priest', 8, 128000, '8d4', 13, 3, 8, 5, 4, 4, 3, 0, 0),
    ('Priest', 9, 256000, '9d4', 12, 4, 9, 5, 5, 4, 3, 2, 0),
    ('Priest', 10, 384000, '9d4+1', 12, 4, 10, 6, 5, 4, 4, 3, 0),
    ('Priest', 11, 512000, '9d4+2', 11, 5, 11, 6, 5, 5, 4, 3, 2),
    ('Priest', 12, 640000, '9d4+3', 11, 5, 12, 6, 6, 5, 4, 4, 3);

INSERT INTO priest_turning_ability (level, turning_ability) VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 4),
    (5, 5),
    (6, 6),
    (7, 7),
    (8, 8),
    (9, 9),
    (10, 10),
    (11, 11),
    (12, 12);

-- +goose Down
DELETE FROM priest_turning_ability;
DELETE FROM class_data WHERE class_name = 'Priest';
DROP TABLE IF EXISTS priest_turning_ability;