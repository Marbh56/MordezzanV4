-- +goose Up
-- Create class-specific table for druid abilities
CREATE TABLE druid_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all druid abilities directly into the class-specific table
INSERT INTO druid_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Druidic Tongue', 'The secret language of the druids. It has its own runic alphabet; to scribe it, however, is forbidden to all save the highest-ranking druids (9th-level or greater).', 1),
('Fire / Heat Affinity', '+2 bonus to saving throws versus fire- and heat-related effects, cumulative with the fire resistance spell.', 1),
('Natural Identifications', 'A spiritual connexion with nature provides recognition of potable water, the general classification of plant types (e.g., edible, poisonous, curative), and the general classification of animal types (e.g., carnivorous, herbivorous, diurnal, nocturnal, aggressive, docile, natural, unnatural).', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Druid Spell List (see Chapter 7: Sorcery, Table 100), unless the scroll was created by a thaumaturgical sorcerer (one who casts the spells of a magician or magician subclass).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials may include stone tablets, bark, or parchment (the latter of which may be inscribed with ink mixed with animal blood). This involved process requires one week per spell level and must be completed in a sacred grove or henge.', 1),
('Sorcery', 'Druids do not carry spell books, but they may grave runes of religious portent on clay tablets, oak bark, parchment, or other like media. The number and levels of spells cast per day are charted above (see Table 42), though druids of high wisdom gain bonus spells cast per day (see Chapter 3: Statistics, wisdom). For example, a 4th-level druid with 13 wisdom can cast four level 1 spells and two level 2 spells per day. Druids begin with knowledge of three level 1 spells granted upon initiation into the druidic society. These spells are drawn from the Druid Spell List (see Chapter 7: Sorcery, Table 100). Druids develop three new spells at each level gain. These spells are acquired via spiritual revelations gained through communion with ancestral, animistic, and elemental spirits. Such spells are learnt automatically, with no need of qualification rolls, but they must be of castable levels (see Table 42 above). Initial no. of spells known: ×3 No. of spells gained per level: ×3 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),
('Traverse Overgrowth', 'Negotiate natural overgrowth (e.g., briars, brush, tangles, thorns, vines) at normal movement rate (MV), without leaving a discernible trail (if so desired).', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 5 abilities
('Charm Immunity', 'At 5th level, immune to the supernatural charms of magical beasts that ensorcell and beguile (e.g., greater gorgon, harpy, man of Leng). This immunity does not apply to the charm spells of sorcerers.', 5),
('Shapechange', 'At 5th level, the power to assume the form of a normal animal of Small size, once per day. Choices include amphibians, birds, fishes, mammals, and reptiles (e.g., frog or salamander; crow or eagle; carp or trout; raccoon or squirrel; snake or turtle). At 7th level, a Medium animal form can be adopted. Examples include bear (black), boar, deer (red or reindeer), dog (war), hyæna, mountain lion, snake (python), and wolf. The druid assumes the creature''s armour class, movement, and tactile abilities (except venom, disease transmission, and other special attacks); however, personal hit point maximum and saving throw are retained. Upon shapechange, 50% of any prior hit point loss is regained. All clothing, armour, weapons, and items are transformed during the change; magic item enchantments cannot be accessed during the shapechange period, though the powers of a magic ring can be utilized if worn before transformation. Shapechange lasts indefinitely, though it is said a druid who maintains animal form for more than 28 days risks the loss of humanity.', 5),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a druid may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Druidic Hierarchy', 'Druidic society has a strict and far-reaching structure. The druidic hierarchy in Hyperborea comprises nine 9th-level druids, seven 10th-level druids, five 11th-level druids, and three 12th-level druids (The Druidic Triumvirate). When a druid gains enough experience points (XP) to reach 9th level, he or she must seek out and challenge one who has already achieved that rank (unless a vacancy exists). A challenge is met at a sacred grove or henge during an astronomical or astrological event of significance, when members of the sect assemble. The druidic challenge can be one of matched weapons and/or of sorcery; in any case, rites are performed that bring to witness a deity such as Lunaqqua, Thaumagorga, or Yoon''Deh (or an agent thereof). The duel is not necessarily to the death, but loss of life is possible. A defeated but still living druid is reduced in experience to one experience point short of 9th level (255,999 XP) and must abide one year of waiting before issuing a new challenge or rematch. A victorious challenger is promoted to 9th level, awarded the appropriate abilities. This process is repeated in similar fashion at 10th, 11th, and 12th levels.', 9),
('Lordship', 'At 9th level, a druid who builds or assumes control of a shrine or temple becomes a lord and is eligible to attract followers and troops. More information is presented in Appendix B.', 9),

-- Level 11 abilities
('Longevity', 'At 11th level, ageing process slows. For every 13 years (1 Hyperborean cycle), the druid effectively ages but 1 year.', 11);

-- +goose Down
-- Clean up druid abilities table if reverting
DROP TABLE druid_abilities;