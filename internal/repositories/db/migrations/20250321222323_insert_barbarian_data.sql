-- +goose Up
INSERT INTO class_data (
    class_name, level, experience_points, hit_dice, saving_throw, 
    fighting_ability, casting_ability, spell_slots_level1, 
    spell_slots_level2, spell_slots_level3, spell_slots_level4, 
    spell_slots_level5, spell_slots_level6
) VALUES
    ('Barbarian', 1, 0, '1d12', 16, 1, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 2, 3000, '2d12', 16, 2, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 3, 6000, '3d12', 15, 3, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 4, 12000, '4d12', 15, 4, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 5, 24000, '5d12', 14, 5, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 6, 48000, '6d12', 14, 6, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 7, 96000, '7d12', 13, 7, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 8, 192000, '8d12', 13, 8, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 9, 384000, '9d12', 12, 9, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 10, 576000, '9d12+4', 12, 10, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 11, 768000, '9d12+8', 11, 11, 0, 0, 0, 0, 0, 0, 0),
    ('Barbarian', 12, 960000, '9d12+12', 11, 12, 0, 0, 0, 0, 0, 0, 0);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Barbarian';