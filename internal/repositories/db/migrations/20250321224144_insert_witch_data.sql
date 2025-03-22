-- +goose Up
INSERT INTO class_data (
    class_name, 
    level, 
    experience_points, 
    hit_dice, 
    saving_throw, 
    fighting_ability, 
    casting_ability, 
    spell_slots_level1, 
    spell_slots_level2, 
    spell_slots_level3, 
    spell_slots_level4, 
    spell_slots_level5, 
    spell_slots_level6
) VALUES
    ('Witch', 1, 0, '1d4', 16, 0, 1, 1, 0, 0, 0, 0, 0),
    ('Witch', 2, 3000, '2d4', 16, 0, 2, 2, 0, 0, 0, 0, 0),
    ('Witch', 3, 6000, '3d4', 15, 1, 3, 2, 1, 0, 0, 0, 0),
    ('Witch', 4, 12000, '4d4', 15, 1, 4, 3, 2, 0, 0, 0, 0),
    ('Witch', 5, 24000, '5d4', 14, 2, 5, 3, 2, 1, 0, 0, 0),
    ('Witch', 6, 48000, '6d4', 14, 2, 6, 4, 3, 2, 0, 0, 0),
    ('Witch', 7, 96000, '7d4', 13, 3, 7, 4, 3, 2, 1, 0, 0),
    ('Witch', 8, 192000, '8d4', 13, 3, 8, 4, 4, 3, 2, 0, 0),
    ('Witch', 9, 384000, '9d4', 12, 4, 9, 5, 4, 3, 2, 1, 0),
    ('Witch', 10, 576000, '9d4+1', 12, 4, 10, 5, 4, 4, 3, 2, 0),
    ('Witch', 11, 768000, '9d4+2', 11, 5, 11, 5, 5, 4, 3, 2, 1),
    ('Witch', 12, 960000, '9d4+3', 11, 5, 12, 5, 5, 4, 4, 3, 2);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Witch';