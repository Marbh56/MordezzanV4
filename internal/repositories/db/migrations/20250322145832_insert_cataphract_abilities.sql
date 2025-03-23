-- +goose Up
-- Insert Cataphract abilities
INSERT INTO abilities (name, description) VALUES
('Honour', 'Cataphracts who operate as knights (must be Lawful Good, Lawful Evil, or Neutral [Lawful]), serving a lord, monarch, etc., enjoy all the political and social benefits derived therefrom. In fact, people of similar affiliation are expected to accommodate the knight to the best of their abilities. To enjoy these benefits, knights must comport themselves to a strict code of honour, abiding the following precepts: duty, integrity, justice, loyalty, respect, and valour. Failure to do so may result in disgrace, banishment, and in some cases, execution. If honour is comported to with competence and distinction, a cataphract may be knighted at 5th level or greater.'),
('Mounted Charge', 'A thunderous mounted onset both feared and renowned. The cataphract''s lance charge from horseback or camelback is at +2 to the attack roll (+3 versus footmen) and treble damage dice (other modifiers added afterwards, such as strength, weapon mastery, etc.).'),
('Shield Sacrifice', 'To sacrifice a shield and escape harm from a single melee blow. When wielding a shield in combat, if the cataphract is struck by a melee blow, the player may opt to announce a shield sacrifice to avoid damage; however, the shield is destroyed by the blow. If the shield is magical, it has a chance-in-eight to survive destruction equal to the shield''s bonus (e.g., a +1 small shield has a 1-in-8 chance of surviving destruction). This ability cannot be used after damage is rolled, and is usable but once per day, regardless of results.'),
('Skilful Defender', 'To avail armour to its utmost. When clad in medium or heavy armour, the cataphract gains a +1 AC bonus from 1st to 6th levels, and a +2 AC bonus from 7th to 12th levels.'),
('Unbreakable Willpower', 'Immune to the effects of magically induced fear.');

-- Link Cataphract abilities at level 1
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cataphract', id, 1 FROM abilities 
WHERE name IN ('Honour', 'Horsemanship', 'Mounted Charge', 'Shield Sacrifice', 'Skilful Defender', 'Unbreakable Willpower', 'Extraordinary', 'Weapon Mastery');

-- Link Enlist Henchmen at level 6
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cataphract', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

-- Link Melee Expert at level 7
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cataphract', id, 7 FROM abilities 
WHERE name = 'Melee Expert';

-- Link Lordship at level 9
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Cataphract', id, 9 FROM abilities 
WHERE name = 'Lordship';

-- +goose Down
-- Remove Cataphract ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Cataphract';

-- Clean up any Cataphract-specific abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Honour', 'Mounted Charge', 'Shield Sacrifice', 'Skilful Defender', 'Unbreakable Willpower')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);