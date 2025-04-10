-- +goose Up
-- Create class-specific table for necromancer abilities
CREATE TABLE necromancer_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all necromancer abilities directly into the class-specific table
INSERT INTO necromancer_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Alchemy', 'To practice the sorcery-science of alchemy. Apprentice necromancers learn how to identify potions by taste alone; albeit the practice is not always safe. At 7th level, a necromancer may concoct potions and/or poisons with the assistance of an alchemist. By 11th level, the assistance of an alchemist is no longer required. For more information, refer to Chapter 7: Sorcery, alchemy.', 1),
('Read Magic', 'The ability to decipher unintelligible magical inscriptions or symbols placed on weapons, armour, items, doors, walls, and other media by means of the sorcerer mark spell or other like methods.', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Necromancer Spell List (see Chapter 7: Sorcery, Table 96), unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer''s ink, such as sepia. This involved process requires one week per spell level.', 1),
('Sorcery', 'Necromancers cast spells that they memorize from arcane tomes; they also gain divine favour from dæmons, netherworldly beings, and ineffable powers. The number and levels of spells cast per day are charted above (see Table 36), though necromancers of high intelligence gain bonus spells cast per day (see Chapter 3: Statistics, intelligence). For example, a 4th-level necromancer with 13 intelligence can cast four level 1 spells and two level 2 spells per day. The necromancer begins with a spell book that contains three level 1 spells drawn from the Necromancer Spell List (see Chapter 7: Sorcery, Table 96). Through personal research and unspeakable pacts, necromancers develop one new spell at each level gain; this spell is learnt automatically, with no need of a qualification roll, but it must be of a castable level (see Table 36 above). Initial no. of spells known: ×3 No. of spells gained per level: ×1 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),

-- Level 3 abilities
('Command Undead', 'At 3rd level, cause undead to submit and serve. Refer to the cleric class and Table 13 for rules regarding turn undead, with special attention given to evil command of undead. At 3rd level, a necromancer has 1st-level turning ability (TA 1). Like the cleric, the necromancer must stand before the undead, within 30 feet, and speak boldly a malefic commandment, whilst displaying a holy symbol (or a necromantic symbol/glyph of death). All other rules and restrictions can be found at the cleric''s aforementioned turn undead ability.', 3),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a necromancer may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, a necromancer who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up necromancer abilities table if reverting
DROP TABLE necromancer_abilities;