-- +goose Up
-- Create class-specific table for pyromancer abilities
CREATE TABLE pyromancer_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all pyromancer abilities directly into the class-specific table
INSERT INTO pyromancer_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Alchemy', 'To practice the sorcery-science of alchemy. Apprentice pyromancers learn how to identify potions by taste alone; albeit the practice is not always safe. At 7th level, a pyromancer may concoct potions with the assistance of an alchemist. By 9th level, a pyromancer is taught the secret formula to incendiary oil. By 11th level, the assistance of an alchemist is no longer required. For more information, refer to Chapter 7: Sorcery, alchemy.', 1),
('Candle', 'Once per day per level of experience, evoke a heatless, candle-like flame to rise from the palm. Candle sheds a 15-foot radius of light and can be placed on an object. This effect lasts 6 turns (1 hour). Multiple candles can be placed concurrently, restricted only by the pyromancer''s daily limit. With a gesture, a pyromancer can cause one or more candles to singe, each causing 1 hp of damage to any creature on which it is directly placed, and possibly enkindling dry, combustible materials. Pyromancers are immune to the damaging effects of this cantrap.', 1),
('Fire / Heat Affinity', '+2 bonus to saving throws versus fire- and heat-related effects, cumulative with the fire resistance spell.', 1),
('Ice / Cold Vulnerability', '−2 penalty to saving throws versus ice- and cold-related effects.', 1),
('Read Magic', 'The ability to decipher unintelligible magical inscriptions or symbols placed on weapons, armour, items, doors, walls, and other media by means of the sorcerer mark spell or other like methods.', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Pyromancer Spell List (see Chapter 7: Sorcery, Table 97), unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer''s ink, such as sepia. This involved process requires one week per spell level.', 1),
('Sorcery', 'Pyromancers cast spells that they memorize from arcane tomes; they also gain favour from elemental forces and otherworldly beings associated with fire. The number and levels of spells cast per day are charted above (see Table 38), though pyromancers of high intelligence gain bonus spells cast per day (see Chapter 3: Statistics, intelligence). For example, a 4th-level pyromancer with 13 intelligence can cast four level 1 spells and two level 2 spells per day. The pyromancer begins with a spell book that contains three level 1 spells drawn from the Pyromancer Spell List (see Chapter 7: Sorcery, Table 97). Through personal research and the patronage of elemental powers, pyromancers develop one new spell at each level gain; this spell is learnt automatically, with no need of a qualification roll, but it must be of a castable level (see Table 38 above). Initial no. of spells known: ×3 No. of spells gained per level: ×1 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a pyromancer may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, a pyromancer who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up pyromancer abilities table if reverting
DROP TABLE pyromancer_abilities;