-- +goose Up
-- Create class-specific table for warlock abilities
CREATE TABLE warlock_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all warlock abilities directly into the class-specific table
INSERT INTO warlock_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Extraordinary', '+8% chance to perform extraordinary feats of strength (see Chapter 3: Statistics, strength).', 1),
('Magic Item Use', 'Can utilize magic items normally restricted to magicians.', 1),
('Read Magic', 'The ability to decipher unintelligible magical inscriptions or symbols placed on weapons, armour, items, doors, walls, and other media by means of the sorcerer mark spell or other like methods.', 1),
('Scroll Use', 'To decipher and invoke scroll scrolls with spells that are included in the warlock''s chosen school of sorcery: magician, cryomancer, necromancer, or pyromancer (see Chapter 7: Sorcery, Table 93, 94, 96, or 97), unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer''s ink, such as sepia. This involved process requires one week per spell level.', 1),
('Sorcery', 'Warlocks cast spells that they memorize from arcane tomes. At character creation, the player must select which school of magic his or her warlock practices: that of the magician, cryomancer, necromancer, or pyromancer. This decision is irrevocable. The number and levels of spells cast per day are charted above (see Table 30), though warlocks of high intelligence gain bonus spells cast per day (see Chapter 3: Statistics, intelligence). For example, a 4th-level warlock with 13 intelligence can cast two level 1 spells and one level 2 spell per day. A warlock begins with a spell book that contains one level 1 spell selected from the Magician-, Cryomancer-, Necromancer-, or Pyromancer Spell List (see Chapter 7, Table 93, 94, 96, or 97), depending on which school of sorcery was selected at character creation. Through personal research, the warlock develops one new spell at each level gain; each is learnt automatically, with no need of a qualification roll, but it must be of a castable level (see Table 30 above). Initial no. of spells known: ×1 No. of spells gained per level: ×1 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a warlock may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 7 abilities
('Melee Expert', 'At 7th level, a warlock''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7),

-- Level 9 abilities
('Lordship', 'At 9th level, a warlock who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up warlock abilities table if reverting
DROP TABLE warlock_abilities;