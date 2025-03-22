-- +goose Up
CREATE TABLE IF NOT EXISTS bard_druid_spells (
    level INTEGER NOT NULL,
    spell_slots_level1 INTEGER NOT NULL,
    spell_slots_level2 INTEGER NOT NULL,
    spell_slots_level3 INTEGER NOT NULL,
    spell_slots_level4 INTEGER NOT NULL,
    PRIMARY KEY (level)
);

CREATE TABLE IF NOT EXISTS bard_illusionist_spells (
    level INTEGER NOT NULL,
    spell_slots_level1 INTEGER NOT NULL,
    spell_slots_level2 INTEGER NOT NULL,
    spell_slots_level3 INTEGER NOT NULL,
    spell_slots_level4 INTEGER NOT NULL,
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
    ('Bard', 1, 0, '1d8', 16, 1, 1),
    ('Bard', 2, 2500, '2d8', 16, 1, 2),
    ('Bard', 3, 5000, '3d8', 15, 2, 3),
    ('Bard', 4, 10000, '4d8', 15, 3, 4),
    ('Bard', 5, 20000, '5d8', 14, 3, 5),
    ('Bard', 6, 40000, '6d8', 14, 4, 6),
    ('Bard', 7, 80000, '7d8', 13, 5, 7),
    ('Bard', 8, 160000, '8d8', 13, 5, 8),
    ('Bard', 9, 320000, '9d8', 12, 6, 9),
    ('Bard', 10, 480000, '9d8+2', 12, 7, 10),
    ('Bard', 11, 640000, '9d8+4', 11, 7, 11),
    ('Bard', 12, 800000, '9d8+6', 11, 8, 12);

INSERT INTO bard_druid_spells (
    level,
    spell_slots_level1,
    spell_slots_level2,
    spell_slots_level3,
    spell_slots_level4
) VALUES
    (1, 1, 0, 0, 0),
    (2, 1, 0, 0, 0),
    (3, 1, 1, 0, 0),
    (4, 1, 1, 0, 0),
    (5, 1, 1, 1, 0),
    (6, 1, 1, 1, 0),
    (7, 1, 1, 1, 1),
    (8, 1, 1, 1, 1),
    (9, 2, 2, 1, 1),
    (10, 2, 2, 1, 1),
    (11, 2, 2, 2, 2),
    (12, 2, 2, 2, 2);

INSERT INTO bard_illusionist_spells (
    level,
    spell_slots_level1,
    spell_slots_level2,
    spell_slots_level3,
    spell_slots_level4
) VALUES
    (1, 0, 0, 0, 0),
    (2, 1, 0, 0, 0),
    (3, 1, 0, 0, 0),
    (4, 1, 1, 0, 0),
    (5, 1, 1, 0, 0),
    (6, 1, 1, 1, 0),
    (7, 1, 1, 1, 0),
    (8, 1, 1, 1, 1),
    (9, 1, 1, 1, 1),
    (10, 2, 2, 1, 1),
    (11, 2, 2, 1, 1),
    (12, 2, 2, 2, 2);

-- +goose Down
DELETE FROM bard_illusionist_spells;
DELETE FROM bard_druid_spells;
DELETE FROM class_data WHERE class_name = 'Bard';
DROP TABLE IF EXISTS bard_illusionist_spells;
DROP TABLE IF EXISTS bard_druid_spells;