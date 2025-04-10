-- +goose Up
-- Create class-specific table for illusionist abilities
CREATE TABLE illusionist_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all illusionist abilities directly into the class-specific table
INSERT INTO illusionist_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Alchemy', 'To practice the sorcery-science of alchemy. Apprentice illusionists learn how to identify potions by taste alone; albeit the practice is not always safe. At 7th level, an illusionist may concoct potions with the assistance of an alchemist. By 11th level, the assistance of an alchemist is no longer required. For more information, refer to Chapter 7: Sorcery, alchemy.', 1),
('Coloured Globe', 'Once per day per level of experience, evoke a 6-inch diameter coloured globe to rise from the open palm and float within 10 feet, as directed by the illusionist. A coloured globe glows any pastel colour as chosen by the illusionist, shedding light in a 10-foot radius. It lasts 6 turns (1 hour) with no need of concentration. Multiple coloured globes can be simultaneously controlled if so desired.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of dexterity (see Chapter 3: Statistics, dexterity).', 1),
('Perceive Illusion', '+2 bonus to saving throws versus illusions and phantasms. If the illusion is that of a sorcerer 3 or more levels lower than the illusionist, then the saving throw bonus is instead equal to the level difference (e.g., a 7th-level illusionist met by the phantasm spell of a 3rd-level bard gains a +4 saving throw bonus).', 1),
('Read Magic', 'The ability to decipher unintelligible magical inscriptions or symbols placed on weapons, armour, items, doors, walls, and other media by means of the sorcerer mark spell or other like methods.', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Illusionist Spell List (see Chapter 7: Sorcery, Table 95), unless the scroll was created by an ecclesiastical sorcerer (one who casts cleric or druid spells).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include the finest vellum, paper, or papyrus; a fresh quill; and sorcerer''s ink, such as sepia. This involved process requires one week per spell level.', 1),
('Sorcery', 'Illusionists cast spells that they memorize from arcane tomes. They also channel strange energies from darkness, light, and sound. The number and levels of spells cast per day are charted above (see Table 34), though illusionists of high intelligence gain bonus spells cast per day (see Chapter 3: Statistics, intelligence). For example, a 4th-level illusionist with 13 intelligence can cast four level 1 spells and two level 2 spells per day. The illusionist begins with a spell book that contains three level 1 spells drawn from the Illusionist Spell List (see Chapter 7: Sorcery, Table 95). Through personal research, illusionists develop one new spell at each level gain; this spell is learnt automatically, with no need of a qualification roll, but it must be of a castable level (see Table 34 above). Initial no. of spells known: ×3 No. of spells gained per level: ×1 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),
('Vizard', 'A cantrap to alter one''s countenance by means of a simple illusion. Chin, ears, eyes, hair, lips, nose, teeth, and all such facial features can be altered, including shape, colour, thickness, blemishes, and so forth. The illusionist must touch his or her face for one round (10 seconds) and imagine the desired features. The vizard is only detectable if a viewer peers closely (within 12 inches of face) or casts detect phantasm. Otherwise, the effect lasts for 1 turn (10 minutes) and may be executed once per day from 1st to 6th level, twice per day from 7th to 12th level. Regardless of level, any illusionist can see through another illusionist''s vizard.', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, an illusionist may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, an illusionist who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up illusionist abilities table if reverting
DROP TABLE illusionist_abilities;