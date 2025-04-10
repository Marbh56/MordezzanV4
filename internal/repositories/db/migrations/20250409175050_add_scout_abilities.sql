-- +goose Up
-- Create class-specific table for scout abilities
CREATE TABLE scout_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all scout abilities directly into the class-specific table
INSERT INTO scout_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Agile', '+1 AC bonus when unarmoured and unencumbered (small shield allowed).', 1),
('Alertness', 'Reduces by one (−1) on a d6 roll the party''s chance to be surprized.', 1),
('Backstab', 'A backstab attempt with a class 1 or 2 melee weapon. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible "back" (e.g., green slime, purple worm exempt). If the requirements are met, the following benefits are derived: The attack roll is made at +4 "to hit." Additional weapon damage dice are rolled according to the scout''s level of experience: 1st to 4th levels = ×2, 5th to 8th levels = ×3, 9th to 12th levels = ×4. Other damage modifiers (strength, sorcery, etc.) are added afterwards (e.g., a 5th-level scout with 13 strength and a +1 short sword rolls 3d6+2).', 1),
('Controlled Fall', 'To retard descent of precipitous drops. For every level of experience, the scout can fall 10 feet and sustain no damage, so long as a wall or other stable surface is within a five-foot reach throughout the descent. For falls beyond the scout''s limit, normal rules are in force starting at the point at which the controlled fall no longer applies; e.g., a 5th-level scout plummets down a 90-foot pit and thus sustains 4d6 hp damage.', 1),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.', 1),
('Determine Depth and Grade', 'To ascertain the extent of a pit, chasm, or shaft by dropping a coin or pebble and listening; to determine the slope of dungeon passages, detecting even the shallowest of slants. The chance of success for determine depth and grade is as follows: 1st to 4th level = 2-in-6, 5th to 8th level = 3-in-6, 9th to 12th level = 4-in-6. If the roll is off by one or two, the estimate is off by 20%. For example, if a 2nd-level scout rolls a 4 when attempting to determine the depth of a 50-foot pit, he or she will believe the pit to be either 60 feet deep or 40 feet deep, as best decided by the referee. If the roll is off by more than 2, then the result is failure.', 1),
('Disguise', 'To fashion a façade that simulates a particular race/culture/social class, possibly making one appear a few inches taller or shorter and/or several pounds heavier or thinner (cf. the disguise self spell). The scout also can appear as the opposite sex. This ruse may be accomplished through a combination of acting, makeup, apparel, and perhaps even subtle sorcery. The base chance of the disguise being discerned is 2-in-12, adjusted as the referee deems appropriate. If the scout has a 16+ charisma, the base chance for being discerned is reduced to 1-in-12.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of dexterity (see Chapter 3: Statistics, dexterity).', 1),
('Run', 'To move as swiftly as a hare; base 50 MV when lightly armoured or unarmoured.', 1),
('Track', 'To stalk prey, tracing physical signs and discerning subtle clues. A scout can track at the below suggested probabilities: Wilderness: A base 7-in-12 chance to find, identify, and follow fresh tracks outdoors or in natural caverns. Non-Wilderness: A base 9-in-12 chance to discern tracks in a dungeon, castle, city street, or like setting. Furthermore, the scout can identify in general terms the species tracked if it is a known animal type (e.g., a large feline, a heavy bovine, a small canine). N.B.: The referee may adjust the chance-in-twelve to track based on prevailing circumstances.', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.', 1),
('Discern Noise', 'To hearken at a door and detect the faintest of noises on the other side, perceive the distant footfalls of a wandering monster, or distinguish a single voice in a crowd. Six rounds (one minute) of concentrated listening are required.', 1),
('Hide', 'To vanish into shadows, camouflage oneself, or flatten one''s body to a seemingly impossible degree—all whilst remaining still as a statue. Only the slightest movement is permissible (e.g., unsheathing a blade, opening a pouch). Hiding is impossible in direct sunlight, or if the scout is observed.', 1),
('Manipulate Traps', 'To find, remove, and reset traps both magical and mundane. Separate checks must be made to accomplish each facet of this skill: find, remove, reset. Failure by more than two, or if a natural 12 is rolled, may cause the trap to detonate on the scout. Also, a new trap may be built if the mechanism is simple and the parts available; anything more complex requires the assistance of an engineer. Thieves'' tools are required when practicing this ability.', 1),
('Move Silently', 'To move with preternatural quiet, even across squeaky floorboards, dry leaves, loose debris, and the like. This skill is executed at half the scout''s normal movement rate.', 1),
('Open Locks', 'To pick locks or disable latching mechanisms both magical and mundane. Thieves'' tools are required. Picking or dismantling a lock may be attempted but once; if the attempt fails, the thief cannot try again until he has gained a level of experience. Most locks require 1d4 minutes to pick; complex locks might necessitate 3d6 minutes.', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a scout may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, a scout who builds or assumes control of suitable headquarters becomes a lord and is eligible to attract a band of thieves. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up scout abilities table if reverting
DROP TABLE scout_abilities;