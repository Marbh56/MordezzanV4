-- +goose Up

CREATE TABLE IF NOT EXISTS thief_skills (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    skill_name TEXT NOT NULL UNIQUE,
    attribute TEXT NOT NULL
);

-- Create level progression table
CREATE TABLE IF NOT EXISTS thief_skill_progression (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    skill_id INTEGER NOT NULL,
    level_range TEXT NOT NULL,
    success_chance TEXT NOT NULL,
    FOREIGN KEY(skill_id) REFERENCES thief_skills(id),
    UNIQUE(skill_id, level_range)
);

-- Create class-skill mapping
CREATE TABLE IF NOT EXISTS class_thief_skill_mapping (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    skill_id INTEGER NOT NULL,
    FOREIGN KEY(skill_id) REFERENCES thief_skills(id),
    UNIQUE(class_name, skill_id)
);

-- Insert thief skills
INSERT INTO thief_skills (skill_name, attribute) VALUES
('Climb', 'DX'),
('Decipher Script', 'IN'),
('Discern Noise', 'WS'),
('Hide', 'DX'),
('Manipulate Traps', 'DX'),
('Move Silently', 'DX'),
('Open Locks', 'DX'),
('Pick Pockets', 'DX'),
('Read Scrolls', 'IN');

-- Insert progression data
-- Climb
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '8:12' FROM thief_skills WHERE skill_name = 'Climb';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '8:12' FROM thief_skills WHERE skill_name = 'Climb';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '9:12' FROM thief_skills WHERE skill_name = 'Climb';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '9:12' FROM thief_skills WHERE skill_name = 'Climb';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '10:12' FROM thief_skills WHERE skill_name = 'Climb';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '10:12' FROM thief_skills WHERE skill_name = 'Climb';

-- Decipher Script
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '0:12' FROM thief_skills WHERE skill_name = 'Decipher Script';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '1:12' FROM thief_skills WHERE skill_name = 'Decipher Script';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '2:12' FROM thief_skills WHERE skill_name = 'Decipher Script';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '3:12' FROM thief_skills WHERE skill_name = 'Decipher Script';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '4:12' FROM thief_skills WHERE skill_name = 'Decipher Script';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '5:12' FROM thief_skills WHERE skill_name = 'Decipher Script';

-- Discern Noise
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '4:12' FROM thief_skills WHERE skill_name = 'Discern Noise';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '5:12' FROM thief_skills WHERE skill_name = 'Discern Noise';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '6:12' FROM thief_skills WHERE skill_name = 'Discern Noise';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '7:12' FROM thief_skills WHERE skill_name = 'Discern Noise';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '8:12' FROM thief_skills WHERE skill_name = 'Discern Noise';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '9:12' FROM thief_skills WHERE skill_name = 'Discern Noise';

-- Hide
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '5:12' FROM thief_skills WHERE skill_name = 'Hide';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '6:12' FROM thief_skills WHERE skill_name = 'Hide';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '7:12' FROM thief_skills WHERE skill_name = 'Hide';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '8:12' FROM thief_skills WHERE skill_name = 'Hide';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '9:12' FROM thief_skills WHERE skill_name = 'Hide';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '10:12' FROM thief_skills WHERE skill_name = 'Hide';

-- Manipulate Traps
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '3:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '4:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '5:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '6:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '7:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '8:12' FROM thief_skills WHERE skill_name = 'Manipulate Traps';

-- Move Silently
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '5:12' FROM thief_skills WHERE skill_name = 'Move Silently';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '6:12' FROM thief_skills WHERE skill_name = 'Move Silently';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '7:12' FROM thief_skills WHERE skill_name = 'Move Silently';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '8:12' FROM thief_skills WHERE skill_name = 'Move Silently';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '9:12' FROM thief_skills WHERE skill_name = 'Move Silently';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '10:12' FROM thief_skills WHERE skill_name = 'Move Silently';

-- Open Locks
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '3:12' FROM thief_skills WHERE skill_name = 'Open Locks';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '4:12' FROM thief_skills WHERE skill_name = 'Open Locks';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '5:12' FROM thief_skills WHERE skill_name = 'Open Locks';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '6:12' FROM thief_skills WHERE skill_name = 'Open Locks';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '7:12' FROM thief_skills WHERE skill_name = 'Open Locks';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '8:12' FROM thief_skills WHERE skill_name = 'Open Locks';

-- Pick Pockets
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '4:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '5:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '6:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '7:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '8:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '9:12' FROM thief_skills WHERE skill_name = 'Pick Pockets';

-- Read Scrolls
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '1-2', '0:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '3-4', '0:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '5-6', '0:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '7-8', '3:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '9-10', '4:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';
INSERT INTO thief_skill_progression (skill_id, level_range, success_chance)
SELECT id, '11-12', '5:12' FROM thief_skills WHERE skill_name = 'Read Scrolls';

-- Assign skills to thief
INSERT INTO class_thief_skill_mapping (class_name, skill_id)
SELECT 'Thief', id FROM thief_skills;

-- Assign skills to barbarian (only Move Silently and Climb)
INSERT INTO class_thief_skill_mapping (class_name, skill_id)
SELECT 'Barbarian', id FROM thief_skills WHERE skill_name IN ('Move Silently', 'Climb');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back

-- Drop tables in reverse order
DROP TABLE IF EXISTS class_thief_skill_mapping;
DROP TABLE IF EXISTS thief_skill_progression;
DROP TABLE IF EXISTS thief_skills;