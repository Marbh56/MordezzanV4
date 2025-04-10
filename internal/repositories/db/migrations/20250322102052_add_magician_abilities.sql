-- +goose Up
-- Create class-specific table for magician abilities
CREATE TABLE magician_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all magician abilities directly into the class-specific table
INSERT INTO magician_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Sorcery', 'Insight into arcane matters. Secrets of runes, glyphs, sigils, and other symbolic magic are especially interesting to magicians. Magicians have a 4-in-6 chance to understand an unknown nonverbal magical inscription. Furthermore, they can discern the general purpose of unidentified potions and scrolls with a 4-in-6 chance.', 1),
('Spell Preparation', 'Prepare magical formulae in accord with the strictures of the class, as described in Chapter 7: Sorcery. Note that a magician''s intelligence statistic score affects extra spell capacity.', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells from the Magician Spell List, unless the scroll was created by a thaumaturgical sorcerer (one who casts cleric or cleric subclass spells).', 1),
('Spell Book', 'To scribe a tome of magical formulae, allowing the magician to memorize spells for eventual casting, based on the contents of his spellbook. In the case of nonpreparation, the referee may allow a 15% random chance per spell level for the magician to prepare a commonly used spell without recourse to his spellbook. For more information, see Chapter 7: Sorcery.', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 250 gp + 100 gp per spell level. This requires a set of costly pens and inks (e.g., sepia, distilled from ink-devil secretions), typically contained within a portable wooden case. This elaborate process requires one week per spell level, and it is time consuming in nature, requiring delicate, precise penmanhip and notation.', 1),

-- Higher level abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.', 4),
('Enlist Henchmen', 'At 6th level, a magician may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Lordship', 'At 9th level, a magician who builds or assumes control of a fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up magician abilities table if reverting
DROP TABLE magician_abilities;