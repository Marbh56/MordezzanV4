-- +goose Up
-- Create a temporary table to track abilities we need to insert
CREATE TEMPORARY TABLE IF NOT EXISTS temp_abilities_to_insert (
    name TEXT,
    description TEXT
);

-- Insert abilities into the temporary table
INSERT INTO temp_abilities_to_insert (name, description) VALUES
('Sorcery', 'Insight into arcane matters. Secrets of runes, glyphs, sigils, and other symbolic magic are especially interesting to magicians. Magicians have a 4-in-6 chance to understand an unknown nonverbal magical inscription. Furthermore, they can discern the general purpose of unidentified potions and scrolls with a 4-in-6 chance.'),
('Spell Preparation', 'Prepare magical formulae in accord with the strictures of the class, as described in Chapter 7: Sorcery. Note that a magician''s intelligence statistic score affects extra spell capacity.'),
('Scroll Use', 'To decipher and invoke scrolls with spells from the Magician Spell List, unless the scroll was created by a thaumaturgical sorcerer (one who casts cleric or cleric subclass spells).'),
('Spell Book', 'To scribe a tome of magical formulae, allowing the magician to memorize spells for eventual casting, based on the contents of his spellbook. In the case of nonpreparation, the referee may allow a 15% random chance per spell level for the magician to prepare a commonly used spell without recourse to his spellbook. For more information, see Chapter 7: Sorcery.'),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 250 gp + 100 gp per spell level. This requires a set of costly pens and inks (e.g., sepia, distilled from ink-devil secretions), typically contained within a portable wooden case. This elaborate process requires one week per spell level, and it is time consuming in nature, requiring delicate, precise penmanhip and notation.'),
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list. This new proficiency is dependent upon training and practice.');

-- Insert Magician abilities that don't already exist in the database
INSERT INTO abilities (name, description)
SELECT t.name, t.description
FROM temp_abilities_to_insert t
WHERE NOT EXISTS (SELECT 1 FROM abilities WHERE name = t.name);

-- Drop the temporary table
DROP TABLE temp_abilities_to_insert;

-- Link Magician abilities at appropriate levels (1st level abilities)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 1 FROM abilities 
WHERE name IN ('Sorcery', 'Spell Preparation', 'Scroll Use', 'Spell Book', 'Scroll Writing')
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Magician' AND ability_id = abilities.id AND min_level = 1
);

-- Link New Weapon Skill at level 4
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 4 FROM abilities 
WHERE name = 'New Weapon Skill'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Magician' AND ability_id = abilities.id AND min_level = 4
);

-- Link Enlist Henchmen at level 6 (reuse existing ability)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Magician' AND ability_id = abilities.id AND min_level = 6
);

-- Link Lordship at level 9 (reuse existing ability)
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Magician', id, 9 FROM abilities 
WHERE name = 'Lordship'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Magician' AND ability_id = abilities.id AND min_level = 9
);

-- +goose Down
-- Remove Magician ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Magician' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Sorcery', 'Spell Preparation', 'Scroll Use', 'Spell Book', 'Scroll Writing', 'New Weapon Skill', 'Enlist Henchmen', 'Lordship')
);

-- Clean up any abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Sorcery', 'Spell Preparation', 'Spell Book')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);