
-- +goose Up
CREATE TABLE magician_class_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level INTEGER NOT NULL,
    experience_points INTEGER NOT NULL,
    hit_dice VARCHAR(10) NOT NULL,
    saving_throw INTEGER NOT NULL,
    fighting_ability INTEGER NOT NULL,
    casting_ability INTEGER NOT NULL,
    spell_slots_level_1 INTEGER NOT NULL DEFAULT 0,
    spell_slots_level_2 INTEGER NOT NULL DEFAULT 0,
    spell_slots_level_3 INTEGER NOT NULL DEFAULT 0,
    spell_slots_level_4 INTEGER NOT NULL DEFAULT 0,
    spell_slots_level_5 INTEGER NOT NULL DEFAULT 0,
    spell_slots_level_6 INTEGER NOT NULL DEFAULT 0
);

INSERT INTO magician_class_data (level, experience_points, hit_dice, saving_throw, fighting_ability, casting_ability, spell_slots_level_1, spell_slots_level_2, spell_slots_level_3, spell_slots_level_4, spell_slots_level_5, spell_slots_level_6) VALUES
(1, 0, '1d4', 16, 0, 1, 1, 0, 0, 0, 0, 0),
(2, 2500, '2d4', 16, 0, 2, 2, 0, 0, 0, 0, 0),
(3, 5000, '3d4', 15, 1, 3, 2, 1, 0, 0, 0, 0),
(4, 10000, '4d4', 15, 1, 4, 3, 2, 0, 0, 0, 0),
(5, 20000, '5d4', 14, 2, 5, 3, 2, 1, 0, 0, 0),
(6, 40000, '6d4', 14, 2, 6, 4, 3, 2, 0, 0, 0),
(7, 80000, '7d4', 13, 3, 7, 4, 3, 2, 1, 0, 0),
(8, 160000, '8d4', 13, 3, 8, 4, 4, 3, 2, 0, 0),
(9, 320000, '9d4', 12, 4, 9, 5, 4, 3, 2, 1, 0),
(10, 480000, '9d4+1', 12, 4, 10, 5, 4, 4, 3, 2, 0),
(11, 640000, '9d4+2', 11, 5, 11, 5, 5, 4, 3, 2, 1),
(12, 800000, '9d4+3', 11, 5, 12, 5, 5, 4, 4, 3, 2);

-- +goose Down
DROP TABLE magician_class_data;