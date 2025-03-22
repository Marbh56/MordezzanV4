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
    ('Druid', 1, 0, '1d8', 16, 1, 1, 1, 0, 0, 0, 0, 0),
    ('Druid', 2, 2000, '2d8', 16, 1, 2, 2, 0, 0, 0, 0, 0),
    ('Druid', 3, 4000, '3d8', 15, 2, 3, 2, 1, 0, 0, 0, 0),
    ('Druid', 4, 8000, '4d8', 15, 3, 4, 3, 2, 0, 0, 0, 0),
    ('Druid', 5, 16000, '5d8', 14, 3, 5, 3, 2, 1, 0, 0, 0),
    ('Druid', 6, 32000, '6d8', 14, 4, 6, 4, 3, 2, 0, 0, 0),
    ('Druid', 7, 64000, '7d8', 13, 5, 7, 4, 3, 2, 1, 0, 0),
    ('Druid', 8, 128000, '8d8', 13, 5, 8, 4, 4, 3, 2, 0, 0),
    ('Druid', 9, 256000, '9d8', 12, 6, 9, 5, 4, 3, 2, 1, 0),
    ('Druid', 10, 384000, '9d8+2', 12, 7, 10, 5, 4, 4, 3, 2, 0),
    ('Druid', 11, 512000, '9d8+4', 11, 7, 11, 5, 5, 4, 3, 2, 1),
    ('Druid', 12, 640000, '9d8+6', 11, 8, 12, 6, 5, 4, 4, 3, 2);

-- +goose Down
DELETE FROM class_data WHERE class_name = 'Druid';