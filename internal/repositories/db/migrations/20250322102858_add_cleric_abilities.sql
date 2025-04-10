-- +goose Up
-- Create class-specific table for cleric abilities
CREATE TABLE cleric_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all cleric abilities directly into the class-specific table
INSERT INTO cleric_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Turn Undead', 'All clerics can exert control over the undead and some dæmonic beings, causing them to flee and/or cower. Evil clerics can opt instead to compel the submission and service of these foul creatures. In either case, the cleric must do the following: Stand before the undead (within 30 feet); Speak boldly a commandment whilst displaying a holy symbol. The referee must cross-reference the cleric''s turning ability (TA) with the Undead Type to determine the cleric''s chance of success. Clerics of above-average charisma (15+) are more commanding; hence, their chance-in-twelve of success is improved by one (+1).', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells from the Cleric Spell List, unless the scroll was created by a non-thaumaturgical sorcerer.', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 250 gp + 100 gp per spell level. This requires a set of costly pens and inks, typically contained within a portable wooden case. This elaborate process requires one week per spell level.', 1),
('Sorcery', 'Insight into divine matters. Secrets of prayers, rituals, and other symbolic magic are especially interesting to clerics. Clerics have a 4-in-6 chance to understand an unknown nonverbal magical inscription of divine nature. Furthermore, they can discern the general purpose of unidentified holy or unholy items with a 4-in-6 chance.', 1),

-- Higher level abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.', 4),
('Enlist Henchmen', 'At 6th level, a cleric may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Lordship', 'At 9th level, a cleric who builds or assumes control of a fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up cleric abilities table if reverting
DROP TABLE cleric_abilities;