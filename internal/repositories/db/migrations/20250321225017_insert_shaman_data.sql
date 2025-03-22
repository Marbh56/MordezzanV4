-- +goose Up

CREATE TABLE IF NOT EXISTS shaman_turning_ability (
    level INTEGER NOT NULL,
    turning_ability INTEGER NOT NULL,
    PRIMARY KEY (level)
);

CREATE TABLE IF NOT EXISTS shaman_divine_spells (
    level INTEGER NOT NULL,
    spell_slots_level1 INTEGER NOT NULL,
    spell_slots_level2 INTEGER NOT NULL,
    spell_slots_level3 INTEGER NOT NULL,
    spell_slots_level4 INTEGER NOT NULL,
    spell_slots_level5 INTEGER NOT NULL,
    spell_slots_level6 INTEGER NOT NULL,
    PRIMARY KEY (level)
);

CREATE TABLE IF NOT EXISTS shaman_arcane_spells (
    level INTEGER NOT NULL,
    spell_slots_level1 INTEGER NOT NULL,
    spell_slots_level2 INTEGER NOT NULL,
    spell_slots_level3 INTEGER NOT NULL,
    spell_slots_level4 INTEGER NOT NULL,
    spell_slots_level5 INTEGER NOT NULL,
    spell_slots_level6 INTEGER NOT NULL,
    PRIMARY KEY (level)
);

INSERT INTO class_data (
    class_name, 
    level, 
    experience_points, 
    hit_dice, 
    saving_throw, 
    fighting_ability,
    casting_ability
) VALUES
    ('Shaman', 1, 0, '1d6', 16, 0, 1),
    ('Shaman', 2, 2500, '2d6', 16, 0, 2),
    ('Shaman', 3, 5000, '3d6', 15, 1, 3),
    ('Shaman', 4, 10000, '4d6', 15, 2, 4),
    ('Shaman', 5, 20000, '5d6', 14, 2, 5),
    ('Shaman', 6, 40000, '6d6', 14, 3, 6),
    ('Shaman', 7, 80000, '7d6', 13, 4, 7),
    ('Shaman', 8, 160000, '8d6', 13, 4, 8),
    ('Shaman', 9, 320000, '9d6', 12, 5, 9),
    ('Shaman', 10, 480000, '9d6+2', 12, 6, 10),
    ('Shaman', 11, 640000, '9d6+4', 11, 6, 11),
    ('Shaman', 12, 800000, '9d6+6', 11, 7, 12);

INSERT INTO shaman_turning_ability (level, turning_ability) VALUES
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

INSERT INTO shaman_divine_spells (
    level,
    spell_slots_level1,
    spell_slots_level2,
    spell_slots_level3,
    spell_slots_level4,
    spell_slots_level5,
    spell_slots_level6
) VALUES
    (1, 1, 0, 0, 0, 0, 0),
    (2, 1, 0, 0, 0, 0, 0),
    (3, 1, 1, 0, 0, 0, 0),
    (4, 1, 1, 0, 0, 0, 0),
    (5, 1, 1, 1, 0, 0, 0),
    (6, 1, 1, 1, 0, 0, 0),
    (7, 2, 1, 1, 1, 0, 0),
    (8, 2, 1, 1, 1, 0, 0),
    (9, 2, 2, 1, 1, 1, 0),
    (10, 2, 2, 1, 1, 1, 0),
    (11, 2, 2, 2, 1, 1, 1),
    (12, 2, 2, 2, 1, 1, 1);

INSERT INTO shaman_arcane_spells (
    level,
    spell_slots_level1,
    spell_slots_level2,
    spell_slots_level3,
    spell_slots_level4,
    spell_slots_level5,
    spell_slots_level6
) VALUES
    (1, 0, 0, 0, 0, 0, 0),
    (2, 1, 0, 0, 0, 0, 0),
    (3, 1, 0, 0, 0, 0, 0),
    (4, 1, 1, 0, 0, 0, 0),
    (5, 1, 1, 0, 0, 0, 0),
    (6, 1, 1, 1, 0, 0, 0),
    (7, 1, 1, 1, 0, 0, 0),
    (8, 2, 1, 1, 1, 0, 0),
    (9, 2, 1, 1, 1, 0, 0),
    (10, 2, 2, 1, 1, 1, 0),
    (11, 2, 2, 1, 1, 1, 0),
    (12, 2, 2, 2, 1, 1, 1);

-- +goose Down
DELETE FROM shaman_arcane_spells;
DELETE FROM shaman_divine_spells;
DELETE FROM shaman_turning_ability;
DELETE FROM class_data WHERE class_name = 'Shaman';
DROP TABLE IF EXISTS shaman_arcane_spells;
DROP TABLE IF EXISTS shaman_divine_spells;
DROP TABLE IF EXISTS shaman_turning_ability;