-- +goose Up
CREATE TABLE berserker_natural_ac (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    natural_ac INTEGER NOT NULL,
    UNIQUE(class_name, level),
    FOREIGN KEY (class_name, level) REFERENCES class_data(class_name, level) ON DELETE CASCADE
);

INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Berserker', 1, 0, '1d12', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 2, 2500, '2d12', 16, 2, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 3, 5000, '3d12', 15, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 4, 10000, '4d12', 15, 4, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 5, 20000, '5d12', 14, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 6, 40000, '6d12', 14, 6, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 7, 80000, '7d12', 13, 7, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 8, 160000, '8d12', 13, 8, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 9, 320000, '9d12', 12, 9, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 10, 480000, '9d12+4', 12, 10, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 11, 640000, '9d12+8', 11, 11, 0, 0, 0, 0, 0, 0, 0),
    ('Berserker', 12, 800000, '9d12+12', 11, 12, 0, 0, 0, 0, 0, 0, 0);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac) VALUES
    ('Berserker', 1, 8),
    ('Berserker', 2, 8),
    ('Berserker', 3, 7),
    ('Berserker', 4, 7),
    ('Berserker', 5, 6),
    ('Berserker', 6, 6),
    ('Berserker', 7, 5),
    ('Berserker', 8, 5),
    ('Berserker', 9, 4),
    ('Berserker', 10, 4),
    ('Berserker', 11, 3),
    ('Berserker', 12, 3);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Berserker';
DROP TABLE berserker_natural_ac;