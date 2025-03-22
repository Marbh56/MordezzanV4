-- +goose Up
-- Insert Turn Undead ability if it doesn't exist
INSERT INTO abilities (name, description)
SELECT 'Turn Undead', 'All clerics can exert control over the undead and some dæmonic beings, causing them to flee and/or cower. Evil clerics can opt instead to compel the submission and service of these foul creatures. In either case, the cleric must do the following: Stand before the undead (within 30 feet); Speak boldly a commandment whilst displaying a holy symbol. The referee must cross-reference the cleric''s turning ability (TA) with the Undead Type to determine the cleric''s chance of success. Clerics of above-average charisma (15+) are more commanding; hence, their chance-in-twelve of success is improved by one (+1).'
WHERE NOT EXISTS (SELECT 1 FROM abilities WHERE name = 'Turn Undead');

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