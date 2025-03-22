-- +goose Up

-- Create cleric_turning_ability table for cleric-specific data
CREATE TABLE cleric_turning_ability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    UNIQUE(class_name, level),
    FOREIGN KEY (class_name, level) REFERENCES class_data(class_name, level) ON DELETE CASCADE
);

-- Insert data into class_data for Cleric
INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Cleric', 1, 0, '1d8', 16, 1, 1, 1, 0, 0, 0, 0, 0),
    ('Cleric', 2, 2000, '2d8', 16, 1, 2, 2, 0, 0, 0, 0, 0),
    ('Cleric', 3, 4000, '3d8', 15, 2, 3, 2, 1, 0, 0, 0, 0),
    ('Cleric', 4, 8000, '4d8', 15, 3, 4, 2, 2, 0, 0, 0, 0),
    ('Cleric', 5, 16000, '5d8', 14, 3, 5, 3, 2, 1, 0, 0, 0),
    ('Cleric', 6, 32000, '6d8', 14, 4, 6, 3, 2, 2, 0, 0, 0),
    ('Cleric', 7, 64000, '7d8', 13, 5, 7, 3, 3, 2, 1, 0, 0),
    ('Cleric', 8, 128000, '8d8', 13, 5, 8, 3, 3, 2, 2, 0, 0),
    ('Cleric', 9, 256000, '9d8', 12, 6, 9, 4, 3, 3, 2, 1, 0),
    ('Cleric', 10, 384000, '9d8+2', 12, 7, 10, 4, 3, 3, 2, 2, 0),
    ('Cleric', 11, 512000, '9d8+4', 11, 7, 11, 4, 4, 3, 3, 2, 1),
    ('Cleric', 12, 640000, '9d8+6', 11, 8, 12, 4, 4, 3, 3, 2, 2);

-- Insert data into cleric_turning_ability
INSERT INTO cleric_turning_ability (class_name, level, turning_ability) VALUES
    ('Cleric', 1, 1),
    ('Cleric', 2, 2),
    ('Cleric', 3, 3),
    ('Cleric', 4, 4),
    ('Cleric', 5, 5),
    ('Cleric', 6, 6),
    ('Cleric', 7, 7),
    ('Cleric', 8, 8),
    ('Cleric', 9, 9),
    ('Cleric', 10, 10),
    ('Cleric', 11, 11),
    ('Cleric', 12, 12);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Cleric';
DROP TABLE cleric_turning_ability;