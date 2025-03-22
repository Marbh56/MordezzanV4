-- +goose Up
INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Thief', 1, 0, '1d6', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 2, 1500, '2d6', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 3, 3000, '3d6', 15, 2, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 4, 6000, '4d6', 15, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 5, 12000, '5d6', 14, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 6, 24000, '6d6', 14, 4, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 7, 48000, '7d6', 13, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 8, 96000, '8d6', 13, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 9, 192000, '9d6', 12, 6, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 10, 288000, '9d6+2', 12, 7, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 11, 384000, '9d6+4', 11, 7, 0, 0, 0, 0, 0, 0, 0),
    ('Thief', 12, 480000, '9d6+6', 11, 8, 0, 0, 0, 0, 0, 0, 0);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Thief';