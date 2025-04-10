-- +goose Up
-- Create class-specific table for ranger abilities
CREATE TABLE ranger_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all ranger abilities directly into the class-specific table
INSERT INTO ranger_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Alertness', 'Reduces by one (−1) on a d6 roll the party''s chance to be surprized.', 1),
('Ambusher', 'When traversing the wilds alone or with others of like ability, the ranger''s base surprize chance increases by one (+1) on a d6 roll. Furthermore, when outdoors, even an untrained party''s chance to surprize increases by one (+1) if the ranger positions and prepares them accordingly.', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear, if lightly armoured or unarmoured, as a thief of equal level (see Table 16). Chance in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.', 1),
('Discern Noise', 'Unusually perceptive, detecting the faintest sounds. The ranger can discern noise as a thief of equal level (see Table 16). Six rounds (one minute) of concentrated listening are required.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength (see Chapter 3: Statistics, strength).', 1),
('Hide', 'If lightly armoured or unarmoured, able to hide outdoors (wilderness) as a thief of equal level (see Table 16), lurking behind bushes, rocks, trees, and the like. Chance-in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. A ranger can employ camouflage or remain perfectly still whilst flattened to the ground. Only the slightest movement is permissible (e.g., unsheathing a blade, opening a pouch). Hiding is impossible in direct sunlight, or if the ranger is observed. In non-wilderness areas (e.g., cities, dungeons), the chance-in-twelve is decreased by 2.', 1),
('Move Silently', 'To stalk like a panther, moving with preternatural quiet, comparable to a thief of equal level (see Table 16), if the ranger is lightly armoured or unarmoured. Chance-in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. This skill is executed at half the ranger''s normal movement rate.', 1),
('Otherworldly Enemies', 'From their earliest training (typically aged 10–12), rangers are awakened to the terrible knowledge of malevolent alien beings and the nameless horrors they represent. Through painstaking instruction and supernatural insight rangers can cognize the most effective means to harm the otherworldly. Rangers become uncannily perspicacious to these creatures; consequently, they inflict +1 hp damage per level of experience versus the abominations that are categorized as "otherworldly."', 1),
('Track', 'To stalk prey, tracing physical signs and discerning subtle clues. A ranger can track at the below suggested probabilities: Wilderness: A base 10-in-12 chance to find, identify, and follow fresh tracks outdoors or in natural caverns. Non-Wilderness: A base 3-in-12 chance to discern tracks in a dungeon, castle, city street, or like setting. Furthermore, the ranger can identify in general terms the species tracked if it is a known animal type (e.g., a large feline, a heavy bovine, a small canine). N.B.: The referee may adjust the chance-in-twelve to track based on prevailing circumstances.', 1),
('Track Concealment', 'In the wilderness, obscure the tracks of a number of companions equal to the ranger''s level of experience; however, maximum speed is restricted to half the ranger''s normal movement rate (MV).', 1),
('Traverse Overgrowth', 'Negotiate natural overgrowth (e.g., briars, brush, tangles, thorns, vines) at normal movement rate (MV), without leaving a discernible trail (if so desired). The ranger cannot perform this skill wearing heavy armour.', 1),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),
('Wilderness Survival', 'Hunting, trapping, fishing, boating, shelter building, fire building (including tribal smoke signals), logging, woodworking, raft building, and so on. These tasks are performed without need of a check; they are simply the ranger''s province. Under adverse conditions, the referee may assign a reasonable probability of success. Whether a chance of failure applies is at the discretion of the referee, as reflected by the prevailing conditions and abilities of the ranger.', 1),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a ranger may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 7 abilities
('Melee Expert', 'At 7th level, a ranger''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7),
('Scroll Use', 'At 7th level, decipher and invoke ecclesiastical scroll spells from the Druid Spell List, and at 8th level, decipher and invoke thaumaturgical scroll spells from the Magician Spell List (see Chapter 7: Sorcery, Table 100 and Table 93).', 7),
('Sorcery', 'At 7th level, the ranger develops the ability to cast spells as a shaman, practicing the ecclesiastical sorcery of druids and the thaumaturgical sorcery of magicians. The number and levels of spells cast per day are charted above (see Table 28), though rangers of high intelligence and/or wisdom gain bonus spells cast per day (see Chapter 3: Statistics, intelligence and wisdom). For example, a 9th-level ranger with 13 wisdom and 10 intelligence can cast two level 1 druid spells, one level 2 druid spell, and one level 1 magician spell per day. The ranger''s spells are drawn from the Magician Spell List and the Druid Spell List (see Chapter 7: Sorcery, Table 93 and Table 100). Druid spells are granted by animistic and elemental spirits, and magician spells are memorized from rune-etched stone tablets, bark sheets, or animal skins; such media functioning as the ranger''s spell book, as it were. At 7th level, the ranger cultivates a level 1 druid spell; at 8th level, a level 1 magician spell. The ranger develops one new spell each level gain thereafter. The schedule is as follows: 7th level: ×1 level 1 druid spell 8th level: ×1 level 1 magician spell 9th level: ×1 level 2 druid spell 10th level: ×1 level 2 magician spell 11th level: ×1 level 3 druid spell 12th level: ×1 level 3 magician spell. Spells are gained automatically, with no need of qualification rolls. The ranger cannot learn additional spells beyond those developed during level training. Note that a 7th-level ranger has 1st-level casting ability (CA 1) and progresses accordingly.', 7),

-- Level 9 abilities
('Lordship', 'At 9th level, a ranger who builds or assumes control of a wilderness fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up ranger abilities table if reverting
DROP TABLE ranger_abilities;