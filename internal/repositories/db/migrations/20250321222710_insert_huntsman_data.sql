-- +goose Up
INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Huntsman', 1, 0, '1d10', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 2, 2250, '2d10', 16, 2, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 3, 4500, '3d10', 15, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 4, 9000, '4d10', 15, 4, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 5, 18000, '5d10', 14, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 6, 36000, '6d10', 14, 6, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 7, 72000, '7d10', 13, 7, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 8, 144000, '8d10', 13, 8, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 9, 288000, '9d10', 12, 9, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 10, 432000, '9d10+3', 12, 10, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 11, 576000, '9d10+6', 11, 11, 0, 0, 0, 0, 0, 0, 0),
    ('Huntsman', 12, 720000, '9d10+9', 11, 12, 0, 0, 0, 0, 0, 0, 0);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Huntsman';