-- +goose Up
-- Insert Thief-specific abilities
INSERT INTO abilities (name, description) VALUES
('Backstab', 'A backstab attempt with a class 1 or 2 melee weapon. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible "back" (e.g., green slime, purple worm exempt). If the requirements are met, the following benefits result: The attack roll is made at +4 "to hit." Additional weapon damage dice are rolled according to the thief''s level of experience: 1st to 4th levels = ×2, 5th to 8th levels = ×3, 9th to 12th levels = ×4. Other damage modififiers (strength, sorcery, etc.) are added afterwards (e.g., a 5th-level thief with 13 strength and a +1 short sword rolls 3d6+2).'),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.'),
('Thieves'' Cant', 'The secret language of thieves, a strange pidgin in which some words may be unintelligible to an ignorant listener, whereas others might be common yet of alternative meaning. This covert tongue is used in conjunction with specific body language, hand gestures, and facial expressions. Two major dialects of Thieves'' Cant are used in Hyperborea: one by city thieves, the other by pirates; commonalities exist betwixt the two.');

-- Link Thief abilities at level 1
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Thief', id, 1 FROM abilities 
WHERE name IN ('Backstab', 'Detect Secret Doors', 'Thieves'' Cant');

-- Link existing abilities to Thief at level 1
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Thief', id, 1 FROM abilities 
WHERE name IN ('Agile', 'Extraordinary');

-- Link New Weapon Skill (already exists) to Thief at level 4
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Thief', id, 4 FROM abilities 
WHERE name = 'New Weapon Skill';

-- Link Enlist Henchmen (already exists) to Thief at level 6
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Thief', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

-- Link Lordship (already exists) to Thief at level 9
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Thief', id, 9 FROM abilities 
WHERE name = 'Lordship';

-- +goose Down
-- Remove Thief ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Thief';

-- Clean up abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Backstab', 'Detect Secret Doors', 'Thieves'' Cant')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);