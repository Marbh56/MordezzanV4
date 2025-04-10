-- +goose Up
-- Create class-specific table for cataphract abilities
CREATE TABLE cataphract_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all cataphract abilities directly into the class-specific table
INSERT INTO cataphract_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Honour', 'Cataphracts who operate as knights (must be Lawful Good, Lawful Evil, or Neutral [Lawful]), serving a lord, monarch, etc., enjoy all the political and social benefits derived therefrom. In fact, people of similar affiliation are expected to accommodate the knight to the best of their abilities. To enjoy these benefits, knights must comport themselves to a strict code of honour, abiding the following precepts: duty, integrity, justice, loyalty, respect, and valour. Failure to do so may result in disgrace, banishment, and in some cases, execution. If honour is comported to with competence and distinction, a cataphract may be knighted at 5th level or greater.', 1),
('Horsemanship', 'Many cataphracts are exceptional horsemen, hailing from mounted warrior traditions that rely on their steeds in times of peace and war. Even the most stubborn of mounts submit to the cataphract''s will. From the saddle of a tamed mount a cataphract can fight with melee weapons and discharge missiles. Depending on geography and background, this skill may apply to camels.', 1),
('Mounted Charge', 'A thunderous mounted onset both feared and renowned. The cataphract''s lance charge from horseback or camelback is at +2 to the attack roll (+3 versus footmen) and treble damage dice (other modifiers added afterwards, such as strength, weapon mastery, etc.).', 1),
('Shield Sacrifice', 'To sacrifice a shield and escape harm from a single melee blow. When wielding a shield in combat, if the cataphract is struck by a melee blow, the player may opt to announce a shield sacrifice to avoid damage; however, the shield is destroyed by the blow. If the shield is magical, it has a chance-in-eight to survive destruction equal to the shield''s bonus (e.g., a +1 small shield has a 1-in-8 chance of surviving destruction). This ability cannot be used after damage is rolled, and is usable but once per day, regardless of results.', 1),
('Skilful Defender', 'To avail armour to its utmost. When clad in medium or heavy armour, the cataphract gains a +1 AC bonus from 1st to 6th levels, and a +2 AC bonus from 7th to 12th levels.', 1),
('Unbreakable Willpower', 'Immune to the effects of magically induced fear.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength and dexterity (see Chapter 3: Statistics, strength and dexterity).', 1),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),

-- Higher level abilities
('Enlist Henchmen', 'At 6th level, a cataphract may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Melee Expert', 'At 7th level, a cataphract''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7),
('Lordship', 'At 9th level, a cataphract who builds or assumes control of a fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up cataphract abilities table if reverting
DROP TABLE cataphract_abilities;