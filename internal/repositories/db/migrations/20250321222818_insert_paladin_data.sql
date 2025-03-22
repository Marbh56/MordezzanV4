-- +goose Up
CREATE TABLE paladin_turning_ability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    UNIQUE(class_name, level),
    FOREIGN KEY (class_name, level) REFERENCES class_data(class_name, level) ON DELETE CASCADE
);

INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Paladin', 1, 0, '1d10', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 2, 2750, '2d10', 16, 2, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 3, 5500, '3d10', 15, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 4, 11000, '4d10', 15, 4, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 5, 22000, '5d10', 14, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 6, 44000, '6d10', 14, 6, 0, 0, 0, 0, 0, 0, 0),
    ('Paladin', 7, 88000, '7d10', 13, 7, 1, 1, 0, 0, 0, 0, 0),
    ('Paladin', 8, 176000, '8d10', 13, 8, 2, 2, 0, 0, 0, 0, 0),
    ('Paladin', 9, 352000, '9d10', 12, 9, 3, 2, 1, 0, 0, 0, 0),
    ('Paladin', 10, 528000, '9d10+3', 12, 10, 4, 2, 2, 0, 0, 0, 0),
    ('Paladin', 11, 704000, '9d10+6', 11, 11, 5, 2, 2, 1, 0, 0, 0),
    ('Paladin', 12, 880000, '9d10+9', 11, 12, 6, 2, 2, 2, 0, 0, 0);

INSERT INTO paladin_turning_ability (class_name, level, turning_ability) VALUES
    ('Paladin', 1, 0),
    ('Paladin', 2, 0),
    ('Paladin', 3, 0),
    ('Paladin', 4, 0),
    ('Paladin', 5, 1),
    ('Paladin', 6, 2),
    ('Paladin', 7, 3),
    ('Paladin', 8, 4),
    ('Paladin', 9, 5),
    ('Paladin', 10, 6),
    ('Paladin', 11, 7),
    ('Paladin', 12, 8);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Paladin';
DROP TABLE paladin_turning_ability;