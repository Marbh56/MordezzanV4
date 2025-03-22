-- +goose Up

INSERT INTO class_data (class_name, level, experience_points, hit_dice, saving_throw, fighting_ability)
VALUES 
    ('Fighter', 1, 0, '1d10', 16, 1),
    ('Fighter', 2, 2000, '2d10', 16, 2),
    ('Fighter', 3, 4000, '3d10', 15, 3),
    ('Fighter', 4, 8000, '4d10', 15, 4),
    ('Fighter', 5, 16000, '5d10', 14, 5),
    ('Fighter', 6, 32000, '6d10', 14, 6),
    ('Fighter', 7, 64000, '7d10', 13, 7),
    ('Fighter', 8, 128000, '8d10', 13, 8),
    ('Fighter', 9, 256000, '9d10', 12, 9),
    ('Fighter', 10, 384000, '9d10+3', 12, 10),
    ('Fighter', 11, 512000, '9d10+6', 11, 11),
    ('Fighter', 12, 640000, '9d10+9', 11, 12);

-- +goose Down

DELETE FROM class_data WHERE class_name = 'Fighter';