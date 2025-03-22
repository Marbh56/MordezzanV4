-- +goose Up
-- Create a temporary table to track abilities we need to insert
CREATE TEMPORARY TABLE IF NOT EXISTS temp_abilities_to_insert (
    name TEXT,
    description TEXT
);

-- Insert abilities into the temporary table
INSERT INTO temp_abilities_to_insert (name, description) VALUES
('Turn Undead', 'All clerics can exert control over the undead and some dæmonic beings, causing them to flee and/or cower. Evil clerics can opt instead to compel the submission and service of these foul creatures. In either case, the cleric must do the following: Stand before the undead (within 30 feet); Speak boldly a commandment whilst displaying a holy symbol. The referee must cross-reference the cleric''s turning ability (TA) with the Undead Type to determine the cleric''s chance of success. Clerics of above-average charisma (15+) are more commanding; hence, their chance-in-twelve of success is improved by one (+1).'),
('Scroll Use', 'To decipher and invoke scrolls with spells from the Cleric Spell List, unless the scroll was created by a thaumaturgical sorcerer (one who casts magician or magician subclass spells).'),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials vary: Some clerics engrave thin tablets of stone, whereas others use vellum or parchment, a fine quill, and sorcerer''s ink, such as sepia. This involved process requires one week per spell level and must be completed on consecrated ground, such as a shrine, fane, or temple.'),
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.');

-- Insert Cleric abilities that don't already exist in the database
INSERT INTO abilities (name, description)
SELECT t.name, t.description
FROM temp_abilities_to_insert t
WHERE NOT EXISTS (SELECT 1 FROM abilities WHERE name = t.name);

-- Drop the temporary table
DROP TABLE temp_abilities_to_insert;

-- Link Cleric abilities at appropriate levels (1st level abilities)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cleric', id, 1 FROM abilities 
WHERE name IN ('Turn Undead', 'Scroll Use', 'Scroll Writing')
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Cleric' AND ability_id = abilities.id AND min_level = 1
);

-- Link New Weapon Skill at level 4
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cleric', id, 4 FROM abilities 
WHERE name = 'New Weapon Skill'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Cleric' AND ability_id = abilities.id AND min_level = 4
);

-- Link Enlist Henchmen at level 6 (reuse existing ability)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cleric', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Cleric' AND ability_id = abilities.id AND min_level = 6
);

-- Link Lordship at level 9 (reuse existing ability)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cleric', id, 9 FROM abilities 
WHERE name = 'Lordship'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Cleric' AND ability_id = abilities.id AND min_level = 9
);

-- Link Sorcery at level 1 (reuse from Magician if exists)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cleric', id, 1 FROM abilities 
WHERE name = 'Sorcery'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Cleric' AND ability_id = abilities.id AND min_level = 1
);

-- +goose Down
-- Remove Cleric ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Cleric' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Turn Undead', 'Scroll Use', 'Scroll Writing', 'Sorcery', 'New Weapon Skill', 'Enlist Henchmen', 'Lordship')
);

-- Clean up any abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Turn Undead', 'Scroll Use', 'Scroll Writing', 'New Weapon Skill')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);