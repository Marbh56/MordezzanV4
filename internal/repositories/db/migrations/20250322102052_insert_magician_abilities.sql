-- +goose Up
-- Insert Magician abilities
INSERT INTO abilities (name, description) VALUES
('Alchemy', 'To practice the sorcery-science of alchemy. Apprentice magicians learn how to identify potions by taste alone; albeit the practice is not always safe. At 7th level, a magician may concoct potions with the assistance of an alchemist. By 11th level, the assistance of an alchemist is no longer required. For more information, refer to Chapter 7: Sorcery, alchemy.'),

('Familiar', 'To summon a small animal (bat, cat, owl, rat, raven, snake, etc.) of 1d3+1 hp to function as a familiar (singular creature with uncanny connexion to the sorcerer). Retaining a familiar provides the following benefits:
- Within range 120 (feet indoors, yards outdoors), the magician can see and hear through the animal; sight is narrowly focused, sounds reverberate metallically.
- The hit point total of the familiar is added to the magician''s total.
- The magician can memorize one extra spell of each available spell level per day (e.g., a 5th-level magician gains bonus level 1, 2, and 3 spells).
These benefits are lost if the familiar is rendered dead, unconscious, or out of range. The familiar is an extraordinary example of the species, has a perfect morale score (ML 12), and always attends and abides the will of its master. To summon a familiar, the magician must perform a series of rites and rituals for 24 hours. To determine result, roll 2d8 on the Familiars table.

If the familiar dies, the magician also sustains 3d6 hp damage. The magician cannot summon another familiar for 1d4+2 months.'),

('Read Magic', 'To decipher magical writings that otherwise are incomprehensible.'),

('Read Scrolls', 'To read a magician spell from a scroll and cast it once. This can be a spell not previously known.'),

('Enlist Henchmen', 'At 6th level, a magician may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.'),

('Craft Magic Items', 'At 11th level, a magician is capable of enchanting magic items for their own use. Even if a magician attains the prowess to craft such items, their knowledge of how to do so must be sought out or discovered, typically in grimoires, scrolls, or tomes of ancient lore; through tutelage by an archmage of like or greater level; or via divine inspiration.'),

('Lordship', 'At 11th level, a magician who builds or occupies a tower became a lord and is eligible to attract followers. More information is presented in Appendix B.'),

('Craft Magic Items for Others', 'At 15th level, a magician is capable of enchanting magic items for use by others. Even if the magician attains the prowess to craft such items, their knowledge of how to do so must be sought out or discovered, typically in grimoires, scrolls, or tomes of ancient lore; through tutelage by an archmage of like or greater level; or via divine inspiration.');

-- Create the Familiar Results table
CREATE TABLE IF NOT EXISTS familiar_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    roll INTEGER NOT NULL,
    animal_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert Familiar Results
INSERT INTO familiar_results (roll, animal_type) VALUES
(2, 'Archæopteryx'),
(3, 'Ice Toad'),
(4, 'Falcon/Hawk'),
(5, 'Squirrel'),
(6, 'Hare'),
(7, 'Gull'),
(8, 'Owl'),
(9, 'Cat'),
(10, 'Rat'),
(11, 'Bat'),
(12, 'Raven'),
(13, 'Weasel'),
(14, 'Fox'),
(15, 'Viper'),
(16, 'Pegomastax');

-- Link Magician abilities at appropriate levels
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 1 FROM abilities 
WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Read Scrolls');

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 11 FROM abilities 
WHERE name IN ('Craft Magic Items', 'Lordship');

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 15 FROM abilities 
WHERE name = 'Craft Magic Items for Others';

-- +goose Down
-- Remove Magician ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Magician' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Read Scrolls', 'Enlist Henchmen', 'Craft Magic Items', 'Lordship', 'Craft Magic Items for Others')
);

-- Clean up abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Read Scrolls', 'Craft Magic Items', 'Craft Magic Items for Others')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);

-- Drop the familiar_results table
DROP TABLE IF EXISTS familiar_results;