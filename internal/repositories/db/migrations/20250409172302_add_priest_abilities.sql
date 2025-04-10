-- +goose Up
-- Create class-specific table for priest abilities
CREATE TABLE priest_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all priest abilities directly into the class-specific table
INSERT INTO priest_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Cleric Spell List (see Chapter 7: Sorcery, Table 99), unless the scroll was created by a thaumaturgical sorcerer (one who casts the spells of a magician or magician subclass).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials vary: Some priests engrave thin tablets of stone, whereas others use vellum or parchment, a fine quill, and sorcerer''s ink, such as sepia. This involved process requires one week per spell level and must be completed on consecrated ground, such as a shrine, fane, or temple.', 1),
('Sorcery', 'Priests memorize and cast cleric spells, but they do not maintain spell books; rather, they might bear the scriptures of their faiths in prayer books, on sacred scrolls, or on graven tablets. The number and levels of spells cast per day are charted above (see Table 46), though priests of high wisdom gain bonus spells cast per day (see Chapter 3: Statistics, wisdom). For example, a 4th-level priest with 13 wisdom can cast five level 1 spells and three level 2 spells per day. Priests begin with knowledge of four level 1 spells, sacred mysteries revealed through initiation into a sect or cult devoted to an otherworldly power, deific being, or ethos. These spells are drawn from the Cleric Spell List (see Chapter 7: Sorcery, Table 99). Priests develop four new spells at each level gain. Typically, they are acquired via spiritual revelation, otherworldly favour, or the piecing together of abstract theologies. Such spells are learnt automatically, with no need of qualification rolls, but they must be of castable levels (see Table 46 above). Initial no. of spells known: ×4 No. of spells gained per level: ×4 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),
('Turn Undead', 'Like clerics, priests can exert control over the undead and some dæmonic beings, causing them to flee and/or cower. Evil priests can opt instead to compel the submission and service of these foul creatures. In either case, the cleric''s turn undead ability and Table 13 should be referenced for more information.', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a priest may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Dæmonwrack', 'At 9th level, once per week, conduct a powerful ritual that dismisses or beckons a dæmon or other netherworldly being. This ability functions as the dismissal spell, except that it applies only to dæmons and their ilk.', 9),
('Lordship', 'At 9th level, a priest who builds or assumes control of a shrine or temple becomes a lord and is eligible to attract followers and troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up priest abilities table if reverting
DROP TABLE priest_abilities;