-- +goose Up
INSERT INTO class_data (
    class_name, 
    level, 
    experience_points, 
    hit_dice, 
    saving_throw, 
    fighting_ability
) VALUES
    ('Assassin', 1, 0, '1d6', 16, 1),
    ('Assassin', 2, 1750, '2d6', 16, 1),
    ('Assassin', 3, 3500, '3d6', 15, 2),
    ('Assassin', 4, 7000, '4d6', 15, 3),
    ('Assassin', 5, 14000, '5d6', 14, 3),
    ('Assassin', 6, 28000, '6d6', 14, 4),
    ('Assassin', 7, 56000, '7d6', 13, 5),
    ('Assassin', 8, 112000, '8d6', 13, 5),
    ('Assassin', 9, 224000, '9d6', 12, 6),
    ('Assassin', 10, 336000, '9d6+2', 12, 7),
    ('Assassin', 11, 448000, '9d6+4', 11, 7),
    ('Assassin', 12, 560000, '9d6+6', 11, 8);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Assassin';