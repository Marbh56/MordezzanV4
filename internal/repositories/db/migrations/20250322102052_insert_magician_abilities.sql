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

('Scroll Use', 'To decipher and invoke scrolls with spells from the Magician Spell List (see Chapter 7: Sorcery, Table 93), unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).'),

('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer''s ink, such as sepia. This involved process requires one week per spell level.'),

('Sorcery', 'Magicians cast spells that they memorize from arcane tomes. The number and levels of spells cast per day are charted in Table 9, though magicians of high intelligence gain bonus spells cast per day (see Chapter 3: Statistics, intelligence); also, magicians who retain a familiar gain bonus spells cast per day. For example, a 4th-level magician with 13 intelligence can cast four level 1 spells and two level 2 spells per day. If the same magician also keeps a familiar, spells cast per day improve to five level 1 spells and three level 2 spells. A magician begins with a spell book that contains three level 1 spells drawn from the Magician Spell List (see Chapter 7: Sorcery, Table 93). Through personal research, magicians develop one new spell at each level gain; this spell is learnt automatically, with no need of a qualification roll, but it must be of a castable level (see Table 9). Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).'),

('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).'),

('Enlist Henchmen', 'At 6th level, a magician may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.'),

('Lordship', 'At 9th level, a magician who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.');

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
WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Scroll Use', 'Scroll Writing', 'Sorcery');

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 4 FROM abilities 
WHERE name = 'New Weapon Skill';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 9 FROM abilities 
WHERE name = 'Lordship';

-- +goose Down
-- Remove Magician ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Magician' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Scroll Use')
);

-- Clean up abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Alchemy', 'Familiar', 'Read Magic', 'Read Scrolls')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);

-- Drop the familiar_results table
DROP TABLE IF EXISTS familiar_results;