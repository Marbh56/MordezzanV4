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
    spell_slots_level3
) VALUES
    ('Warlock', 1, 0, '1d8', 16, 1, 1, 1, 0, 0),
    ('Warlock', 2, 3000, '2d8', 16, 2, 2, 1, 0, 0),
    ('Warlock', 3, 6000, '3d8', 15, 3, 3, 1, 1, 0),
    ('Warlock', 4, 12000, '4d8', 15, 4, 4, 1, 1, 0),
    ('Warlock', 5, 24000, '5d8', 14, 5, 5, 1, 1, 1),
    ('Warlock', 6, 48000, '6d8', 14, 6, 6, 1, 1, 1),
    ('Warlock', 7, 96000, '7d8', 13, 7, 7, 2, 1, 1),
    ('Warlock', 8, 192000, '8d8', 13, 8, 8, 2, 2, 1),
    ('Warlock', 9, 384000, '9d8', 12, 9, 9, 2, 2, 2),
    ('Warlock', 10, 576000, '9d8+2', 12, 10, 10, 3, 2, 2),
    ('Warlock', 11, 768000, '9d8+4', 11, 11, 11, 3, 3, 2),
    ('Warlock', 12, 960000, '9d8+6', 11, 12, 12, 3, 3, 3);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Warlock';