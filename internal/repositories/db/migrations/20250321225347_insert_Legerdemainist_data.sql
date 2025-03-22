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
    ('Legerdemainist', 1, 0, '1d6', 16, 1, 1, 1, 0, 0),
    ('Legerdemainist', 2, 2750, '2d6', 16, 1, 2, 1, 0, 0),
    ('Legerdemainist', 3, 5500, '3d6', 15, 2, 3, 1, 1, 0),
    ('Legerdemainist', 4, 11000, '4d6', 15, 3, 4, 1, 1, 0),
    ('Legerdemainist', 5, 22000, '5d6', 14, 3, 5, 1, 1, 1),
    ('Legerdemainist', 6, 44000, '6d6', 14, 4, 6, 1, 1, 1),
    ('Legerdemainist', 7, 88000, '7d6', 13, 5, 7, 2, 1, 1),
    ('Legerdemainist', 8, 176000, '8d6', 13, 5, 8, 2, 2, 1),
    ('Legerdemainist', 9, 352000, '9d6', 12, 6, 9, 2, 2, 2),
    ('Legerdemainist', 10, 528000, '9d6+2', 12, 7, 10, 3, 2, 2),
    ('Legerdemainist', 11, 704000, '9d6+4', 11, 7, 11, 3, 3, 2),
    ('Legerdemainist', 12, 880000, '9d6+6', 11, 8, 12, 3, 3, 3);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Legerdemainist';