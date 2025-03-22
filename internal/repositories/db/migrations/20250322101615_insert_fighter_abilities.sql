-- +goose Up
-- Insert Fighter abilities
INSERT INTO abilities (name, description) VALUES
('Heroic Fighting', 'To smite multiple foes. When combatting opponents of 1 HD or less, double normal melee attacks per round (2/1, or 3/1 if wielding a mastered weapon). This dramatic attack could be effected as a single, devastating swing or lunge that bursts through multiple foes. At 7th level, when combating foes of 2 HD or less, double normal melee attacks per round (3/1, or 4/1 if wielding a mastered weapon).'),
('Weapon Mastery', 'Mastery of two weapons (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels; however, see grand mastery below, for another option. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.'),
('Grand Mastery', 'At 4th, 8th, or 12th level (player''s choice), when a new weapon mastery is gained, fighters may elect to intensify their training with an already mastered weapon. With this weapon the fighter becomes a grand master (+2 "to hit" and +2 damage, increased attack rate, etc.). A fighter may achieve grand mastery with but one weapon. For more information, see Chapter 9: Combat, weapon skill.'),
('Enlist Henchmen', 'At 6th level, a fighter may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.'),
('Melee Expert', 'At 7th level, a fighter''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.');

-- Link Fighter abilities at appropriate levels
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Fighter', id, 1 FROM abilities 
WHERE name IN ('Heroic Fighting', 'Weapon Mastery');

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Fighter', id, 4 FROM abilities 
WHERE name = 'Grand Mastery';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Fighter', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Fighter', id, 7 FROM abilities 
WHERE name = 'Melee Expert';

-- Link existing Extraordinary ability to Fighter class
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Fighter', id, 1 FROM abilities 
WHERE name = 'Extraordinary';

-- +goose Down
-- Remove Fighter ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Fighter' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Heroic Fighting', 'Weapon Mastery', 'Grand Mastery', 'Enlist Henchmen', 'Melee Expert', 'Extraordinary')
);

-- Clean up any abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Heroic Fighting', 'Grand Mastery')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);