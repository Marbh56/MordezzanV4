-- +goose Up
INSERT INTO class_data (class_name, level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability)
VALUES 
    ('Ranger', 1, 0, '1d10', 16, 1, 0),
    ('Ranger', 2, 2250, '2d10', 16, 2, 0),
    ('Ranger', 3, 4500, '3d10', 15, 3, 0),
    ('Ranger', 4, 9000, '4d10', 15, 4, 0),
    ('Ranger', 5, 18000, '5d10', 14, 5, 0),
    ('Ranger', 6, 36000, '6d10', 14, 6, 0),
    ('Ranger', 7, 72000, '7d10', 13, 7, 1),
    ('Ranger', 8, 144000, '8d10', 13, 8, 2),
    ('Ranger', 9, 288000, '9d10', 12, 9, 3),
    ('Ranger', 10, 432000, '9d10+3', 12, 10, 4),
    ('Ranger', 11, 576000, '9d10+6', 11, 11, 5),
    ('Ranger', 12, 720000, '9d10+9', 11, 12, 6);

CREATE TABLE ranger_druid_spell_slots (
    class_level INTEGER NOT NULL,
    spell_level INTEGER NOT NULL,
    slots INTEGER NOT NULL,
    PRIMARY KEY (class_level, spell_level)
);

CREATE TABLE ranger_magician_spell_slots (
    class_level INTEGER NOT NULL,
    spell_level INTEGER NOT NULL,
    slots INTEGER NOT NULL,
    PRIMARY KEY (class_level, spell_level)
);

-- +goose Down
DROP TABLE ranger_magician_spell_slots;
DROP TABLE ranger_druid_spell_slots;
DELETE FROM class_data WHERE class_name = 'Ranger';