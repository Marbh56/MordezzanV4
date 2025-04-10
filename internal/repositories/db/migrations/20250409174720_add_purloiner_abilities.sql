-- +goose Up
-- Create class-specific table for purloiner abilities
CREATE TABLE purloiner_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all purloiner abilities directly into the class-specific table
INSERT INTO purloiner_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Agile', '+1 AC bonus when unarmoured and unencumbered (small shield allowed).', 1),
('Backstab', 'A backstab attempt with a class 1 or 2 melee weapon. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible "back" (e.g., green slime, purple worm exempt). If the requirements are met, the following benefits are derived: The attack roll is made at +4 "to hit." Additional weapon damage dice are rolled according to the purloiner''s level of experience: 1st to 4th levels = ×2, 5th to 8th levels = ×3, 9th to 12th levels = ×4. Other damage modifiers (strength, sorcery, etc.) are added afterwards (e.g., a 5th-level purloiner with 13 strength and a +1 short sword rolls 3d6+2).', 1),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of dexterity (see Chapter 3: Statistics, dexterity).', 1),
('Magic Item Use', 'Can utilize magic items normally restricted to clerics.', 1),
('Scroll Use', 'To decipher and invoke scrolls with spells that are included in the Cleric Spell List (see Chapter 7: Sorcery, Table 99), unless the scroll was created by a thaumaturgical sorcerer (one who casts the spells of a magician or magician subclass).', 1),
('Scroll Writing', 'To scribe a known spell onto a scroll, creating a single-use magical device at a cost of 500 gp + 100 gp per spell level. Materials vary: Some purloiners engrave thin tablets of stone, whereas others use vellum or parchment, a fine quill, and sorcerer''s ink, such as sepia. This involved process requires one week per spell level and must be completed on consecrated ground, such as a shrine, fane, or temple.', 1),
('Sorcery', 'Purloiners memorize and cast cleric spells, but they do not maintain spell books; rather, they might bear the scriptures of their faiths in prayer books, on sacred scrolls, or on graven tablets. The number and levels of spells cast per day are charted above (see Table 62), though purloiners of high wisdom gain bonus spells cast per day (see Chapter 3: Statistics, wisdom). For example, a 4th-level purloiner with 13 wisdom can cast two level 1 spells and one level 2 spell per day. Purloiners begin with knowledge of two level 1 spells, sacred mysteries revealed through initiation into a sect or cult devoted to an otherworldly power, deific being, or ethos. These spells are drawn from the Cleric Spell List (see Chapter 7: Sorcery, Table 99). Purloiners develop two new spells at each level gain. Typically, they are acquired via spiritual revelation, otherworldly favour, or the piecing together of abstract theologies. Such spells are learnt automatically, with no need of qualification rolls, but they must be of castable levels (see Table 62 above). Initial no. of spells known: ×2 No. of spells gained per level: ×2 Additional spells may be learnt outside of level training, but the process is more arduous (see Chapter 7: Sorcery, acquiring new spells).', 1),
('Thieves'' Cant', 'The secret language of thieves, a strange pidgin in which some words may be unintelligible to an ignorant listener, whereas others might be common yet of alternative meaning. This covert tongue is used in conjunction with specific body language, hand gestures, and facial expressions. Two major dialects of Thieves'' Cant are used in Hyperborea: one by city thieves, the other by pirates; commonalities exist betwixt the two.', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.', 1),
('Decipher Script', 'To translate texts otherwise not understood. Maps can be interpreted, instructions decoded, and so forth. Ancient or alien languages, however, may remain unintelligible, lacking any basis for comparison.', 1),
('Discern Noise', 'To hearken at a door and detect the faintest of noises on the other side, perceive the distant footfalls of a wandering monster, or distinguish a single voice in a crowd. Six rounds (one minute) of concentrated listening are required.', 1),
('Hide', 'To vanish into shadows, camouflage oneself, or flatten one''s body to a seemingly impossible degree—all whilst remaining still as a statue. Only the slightest movement is permissible (e.g., unsheathing a blade, opening a pouch). Hiding is impossible in direct sunlight, or if the purloiner is observed.', 1),
('Manipulate Traps', 'To find, remove, and reset traps both magical and mundane. Separate checks must be made to accomplish each facet of this skill: find, remove, reset. Failure by more than two, or if a natural 12 is rolled, may cause the trap to detonate on the purloiner. Also, a new trap may be built if the mechanism is simple and the parts available; anything more complex requires the assistance of an engineer. Thieves'' tools are required when practicing this ability.', 1),
('Move Silently', 'To move with preternatural quiet, even across squeaky floorboards, dry leaves, loose debris, and the like. This skill is executed at half the purloiner''s normal movement rate.', 1),
('Open Locks', 'To pick locks or disable latching mechanisms both magical and mundane. Thieves'' tools are required. Picking or dismantling a lock may be attempted but once; if the attempt fails, the purloiner cannot try again until he has gained a level of experience. Most locks require 1d4 minutes to pick; complex locks might necessitate 3d6 minutes.', 1),
('Pick Pockets', 'To filch items from a pocket, pouch, backpack, or garment using nimble fingers and distraction. Failure by a margin of 3 or greater indicates the attempt has been observed (though not necessarily by the victim). If the roll is successful, the referee must determine what has been procured. If a purloiner attempts to pick the pocket of a higher-level purloiner, thief, or legerdemainist, a penalty equal to the difference in levels must be applied to the check. This skill also covers the gamut of sleight-of-hand trickery a thief might employ to deceive onlookers.', 1),

-- Level 3 abilities
('Turn Undead', 'At 3rd level, exert control over the undead, causing them to flee and/or cower. Refer to Table 13 at the cleric class entry. At 3rd level the purloiner has 1st-level turning ability (TA 1); at 4th level, TA 2; and so on. The purloiner must stand before the undead and speak boldly a commandment of faith, displaying a holy symbol whilst so doing. This ability can be used a number of times per day equal to the character''s TA; however, the purloiner can make but one attempt per encounter. Evil purloiners may opt to command undead on a successful turn undead check. For more information, refer to evil command of undead in the turn undead cleric class ability.', 3),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a purloiner may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, a purloiner who builds or assumes control of suitable headquarters becomes a lord and is eligible to attract a band of reverent thieves. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up purloiner abilities table if reverting
DROP TABLE purloiner_abilities;