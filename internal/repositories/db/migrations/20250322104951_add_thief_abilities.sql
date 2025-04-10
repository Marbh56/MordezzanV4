-- +goose Up
-- Create class-specific table for thief abilities
CREATE TABLE thief_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all thief abilities directly into the class-specific table
INSERT INTO thief_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Backstab', 'A backstab attempt with a class 1 or 2 melee weapon. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible "back" (e.g., green slime, purple worm exempt). If the requirements are met, the following benefits result: The attack roll is made at +4 "to hit." Additional weapon damage dice are rolled according to the thief''s level of experience: 1st to 4th levels = ×2, 5th to 8th levels = ×3, 9th to 12th levels = ×4. Other damage modififiers (strength, sorcery, etc.) are added afterwards (e.g., a 5th-level thief with 13 strength and a +1 short sword rolls 3d6+2).', 1),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.', 1),
('Thieves'' Cant', 'The secret language of thieves, a strange pidgin in which some words may be unintelligible to an ignorant listener, whereas others might be common yet of alternative meaning. This covert tongue is used in conjunction with specific body language, hand gestures, and facial expressions. Two major dialects of Thieves'' Cant are used in Hyperborea: one by city thieves, the other by pirates; commonalities exist betwixt the two.', 1),
('Agile', '+1 AC bonus when unarmoured and unencumbered (shield allowed).', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength and dexterity (see Chapter 3: Statistics, strength and dexterity).', 1),

-- Higher level abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.', 4),
('Enlist Henchmen', 'At 6th level, a thief may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Lordship', 'At 9th level, a thief who builds or assumes control of a fortress or guild becomes a lord and is eligible to attract followers. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up thief abilities table if reverting
DROP TABLE thief_abilities;