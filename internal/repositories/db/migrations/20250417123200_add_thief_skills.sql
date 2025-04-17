-- +goose Up
-- SQL in this section is executed when the migration is applied

-- Create the thief skills table
CREATE TABLE thief_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    skill_name TEXT NOT NULL,
    attribute TEXT NOT NULL,
    level INTEGER NOT NULL,
    success_chance TEXT NOT NULL
);

-- Insert thief skills data for all levels
INSERT INTO thief_skills (skill_name, attribute, level, success_chance) VALUES
    -- Climb
    ('Climb', 'DX', 1, '8:12'),
    ('Climb', 'DX', 2, '8:12'),
    ('Climb', 'DX', 3, '8:12'),
    ('Climb', 'DX', 4, '8:12'),
    ('Climb', 'DX', 5, '9:12'),
    ('Climb', 'DX', 6, '9:12'),
    ('Climb', 'DX', 7, '9:12'),
    ('Climb', 'DX', 8, '9:12'),
    ('Climb', 'DX', 9, '10:12'),
    ('Climb', 'DX', 10, '10:12'),
    ('Climb', 'DX', 11, '10:12'),
    ('Climb', 'DX', 12, '10:12'),
    
    -- Decipher Script
    ('Decipher Script', 'IN', 1, '0:12'),
    ('Decipher Script', 'IN', 2, '0:12'),
    ('Decipher Script', 'IN', 3, '1:12'),
    ('Decipher Script', 'IN', 4, '1:12'),
    ('Decipher Script', 'IN', 5, '2:12'),
    ('Decipher Script', 'IN', 6, '2:12'),
    ('Decipher Script', 'IN', 7, '3:12'),
    ('Decipher Script', 'IN', 8, '3:12'),
    ('Decipher Script', 'IN', 9, '4:12'),
    ('Decipher Script', 'IN', 10, '4:12'),
    ('Decipher Script', 'IN', 11, '5:12'),
    ('Decipher Script', 'IN', 12, '5:12'),
    
    -- Discern Noise
    ('Discern Noise', 'WS', 1, '4:12'),
    ('Discern Noise', 'WS', 2, '4:12'),
    ('Discern Noise', 'WS', 3, '5:12'),
    ('Discern Noise', 'WS', 4, '5:12'),
    ('Discern Noise', 'WS', 5, '6:12'),
    ('Discern Noise', 'WS', 6, '6:12'),
    ('Discern Noise', 'WS', 7, '7:12'),
    ('Discern Noise', 'WS', 8, '7:12'),
    ('Discern Noise', 'WS', 9, '8:12'),
    ('Discern Noise', 'WS', 10, '8:12'),
    ('Discern Noise', 'WS', 11, '9:12'),
    ('Discern Noise', 'WS', 12, '9:12'),
    
    -- Hide
    ('Hide', 'DX', 1, '5:12'),
    ('Hide', 'DX', 2, '5:12'),
    ('Hide', 'DX', 3, '6:12'),
    ('Hide', 'DX', 4, '6:12'),
    ('Hide', 'DX', 5, '7:12'),
    ('Hide', 'DX', 6, '7:12'),
    ('Hide', 'DX', 7, '8:12'),
    ('Hide', 'DX', 8, '8:12'),
    ('Hide', 'DX', 9, '9:12'),
    ('Hide', 'DX', 10, '9:12'),
    ('Hide', 'DX', 11, '10:12'),
    ('Hide', 'DX', 12, '10:12'),
    
    -- Manipulate Traps
    ('Manipulate Traps', 'DX', 1, '3:12'),
    ('Manipulate Traps', 'DX', 2, '3:12'),
    ('Manipulate Traps', 'DX', 3, '4:12'),
    ('Manipulate Traps', 'DX', 4, '4:12'),
    ('Manipulate Traps', 'DX', 5, '5:12'),
    ('Manipulate Traps', 'DX', 6, '5:12'),
    ('Manipulate Traps', 'DX', 7, '6:12'),
    ('Manipulate Traps', 'DX', 8, '6:12'),
    ('Manipulate Traps', 'DX', 9, '7:12'),
    ('Manipulate Traps', 'DX', 10, '7:12'),
    ('Manipulate Traps', 'DX', 11, '8:12'),
    ('Manipulate Traps', 'DX', 12, '8:12'),
    
    -- Move Silently
    ('Move Silently', 'DX', 1, '5:12'),
    ('Move Silently', 'DX', 2, '5:12'),
    ('Move Silently', 'DX', 3, '6:12'),
    ('Move Silently', 'DX', 4, '6:12'),
    ('Move Silently', 'DX', 5, '7:12'),
    ('Move Silently', 'DX', 6, '7:12'),
    ('Move Silently', 'DX', 7, '8:12'),
    ('Move Silently', 'DX', 8, '8:12'),
    ('Move Silently', 'DX', 9, '9:12'),
    ('Move Silently', 'DX', 10, '9:12'),
    ('Move Silently', 'DX', 11, '10:12'),
    ('Move Silently', 'DX', 12, '10:12'),
    
    -- Open Locks
    ('Open Locks', 'DX', 1, '3:12'),
    ('Open Locks', 'DX', 2, '3:12'),
    ('Open Locks', 'DX', 3, '4:12'),
    ('Open Locks', 'DX', 4, '4:12'),
    ('Open Locks', 'DX', 5, '5:12'),
    ('Open Locks', 'DX', 6, '5:12'),
    ('Open Locks', 'DX', 7, '6:12'),
    ('Open Locks', 'DX', 8, '6:12'),
    ('Open Locks', 'DX', 9, '7:12'),
    ('Open Locks', 'DX', 10, '7:12'),
    ('Open Locks', 'DX', 11, '8:12'),
    ('Open Locks', 'DX', 12, '8:12'),
    
    -- Pick Pockets
    ('Pick Pockets', 'DX', 1, '4:12'),
    ('Pick Pockets', 'DX', 2, '4:12'),
    ('Pick Pockets', 'DX', 3, '5:12'),
    ('Pick Pockets', 'DX', 4, '5:12'),
    ('Pick Pockets', 'DX', 5, '6:12'),
    ('Pick Pockets', 'DX', 6, '6:12'),
    ('Pick Pockets', 'DX', 7, '7:12'),
    ('Pick Pockets', 'DX', 8, '7:12'),
    ('Pick Pockets', 'DX', 9, '8:12'),
    ('Pick Pockets', 'DX', 10, '8:12'),
    ('Pick Pockets', 'DX', 11, '9:12'),
    ('Pick Pockets', 'DX', 12, '9:12'),
    
    -- Read Scrolls
    ('Read Scrolls', 'IN', 1, 'N/A'),
    ('Read Scrolls', 'IN', 2, 'N/A'),
    ('Read Scrolls', 'IN', 3, 'N/A'),
    ('Read Scrolls', 'IN', 4, 'N/A'),
    ('Read Scrolls', 'IN', 5, '0:12'),
    ('Read Scrolls', 'IN', 6, '0:12'),
    ('Read Scrolls', 'IN', 7, '3:12'),
    ('Read Scrolls', 'IN', 8, '3:12'),
    ('Read Scrolls', 'IN', 9, '4:12'),
    ('Read Scrolls', 'IN', 10, '4:12'),
    ('Read Scrolls', 'IN', 11, '5:12'),
    ('Read Scrolls', 'IN', 12, '5:12');

-- Create index for faster lookups
CREATE INDEX idx_thief_skills_lookup ON thief_skills(skill_name, level);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DROP TABLE IF EXISTS thief_skills;