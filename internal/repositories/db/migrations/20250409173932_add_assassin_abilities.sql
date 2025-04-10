-- +goose Up
-- Create class-specific table for assassin abilities
CREATE TABLE assassin_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all assassin abilities directly into the class-specific table
INSERT INTO assassin_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Agile', '+1 AC bonus when unarmoured and unencumbered (shield allowed).', 1),
('Assassinate (Backstab)', 'A backstab (cf. the thief ability) attempt made with a class 1 or 2 melee weapon, with the intent to assassinate. The target must be unaware of the attack, which may be the result of hiding or moving silently. Also, the target must have vital organs (e.g., skeleton, zombie exempt) and a discernible "back" (e.g., green slime, purple worm exempt). If the requirements are met, the following benefits are derived: The attack roll is made at +4 "to hit." If target is hit, a second d20 is rolled to verify assassination (see Table 56 below). If the second d20 roll meets the required target number or less, the target must make a death saving throw or die; however, if the original d20 attack roll was a natural 19 or 20, then no saving throw is allowed. Normal backstab damage rules per the thief class apply if the result is a hit but not an automatic assassination. Additional damage dice are rolled according to the assassin''s level of experience: 1st to 4th levels = ×2, 5th to 8th levels = ×3, 9th to 12th levels = ×4. Other damage modifiers (strength, sorcery, etc.) are added afterwards (e.g., a 5th-level assassin with 13 strength and a +1 short sword rolls 3d6+2). Note that if an assassination attempt is made against an assassin of higher level, the chance in-twenty of success is reduced by one for every level of difference. Sniper Attack: Lastly, unlike the thief''s backstab ability, the assassin also can make an assassination attempt with a missile weapon (such as a bow or thrown dagger) versus a human, humanoid, quasi-man, or giant. The attempt requires short range; the mark must be completely unaware of danger and not otherwise engaged in combat. The assassin''s comprehensive knowledge of anthropoid anatomy allows for this specialized termination attempt.', 1),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of dexterity (see Chapter 3: Statistics, dexterity).', 1),
('Disguise', 'To fashion a façade that simulates a particular race/culture/social class, possibly making one appear a few inches taller or shorter and/or several pounds heavier or thinner (cf. the disguise self spell). The assassin also can appear as the opposite sex. This ruse may be accomplished through a combination of acting, makeup, apparel, and perhaps even subtle sorcery. The base chance of the disguise being discerned is 2-in-12, adjusted as the referee deems appropriate. If the assassin has a 16+ charisma, the base chance for being discerned is reduced to 1-in-12.', 1),
('Harvest Venom', 'To extract venom sacs from snakes, spiders, giant insects, and the like, should the opportunity present and appropriate containers be available. Assassins learn to dissect venomous creatures in the field. This skill is performed at a base 9-in-12 chance of success. If a 12 is rolled, a mishap occurs, and the assassin is exposed (e.g., eyes, nose, skin) to the poisonous fluid. For more information on the extraction of venom, see Chapter 9: Combat, poison, venom harvesting.', 1),
('Poison Resistance', 'Toxicological training and exposure to various poisons and toxins provide a +1 bonus on all saving throws versus poison, though not other death saving throws.', 1),
('Poison Use', 'The employment of toxins to kill or assassinate. Some assassins'' guilds have in-house alchemists who concoct poisons and toxins potentially available for purchase. Consult your referee regarding availability.', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.', 1),
('Discern Noise', 'To hearken at a door and detect the faintest of noises on the other side, perceive the distant footfalls of a wandering monster, or distinguish a single voice in a crowd. Six rounds (one minute) of concentrated listening are required.', 1),
('Hide', 'To vanish into shadows, camouflage oneself, or flatten one''s body to a seemingly impossible degree—all whilst remaining still as a statue. Only the slightest movement is permissible (e.g., unsheathing a blade, opening a pouch). Hiding is impossible in direct sunlight, or if the assassin is observed.', 1),
('Manipulate Traps', 'To find, remove, and reset traps both magical and mundane. Separate checks must be made to accomplish each facet of this skill: find, remove, reset. Failure by more than two, or if a natural 12 is rolled, may cause the trap to detonate on the assassin. Also, a new trap may be built if the mechanism is simple and the parts available; anything more complex requires the assistance of an engineer. Thieves'' tools are required when practicing this ability.', 1),
('Move Silently', 'To move with preternatural quiet, even across squeaky floorboards, dry leaves, loose debris, etc. This skill is executed at half the assassin''s normal movement rate.', 1),
('Open Locks', 'To pick locks or disable latching mechanisms both magical and mundane. Thieves'' tools are required. Picking or dismantling a lock may be attempted but once; if the attempt fails, the thief cannot try again until he has gained a level of experience. Most locks require 1d4 minutes to pick; complex locks might necessitate 3d6 minutes.', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, an assassin may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Poison Manufacture', 'At 9th level, eligible to train with a master alchemist/toxicologist and learn how to concoct debilitating and deadly poisons. See Chapter 7: Sorcery, alchemy for more information.', 9),
('Lordship', 'At 9th level, an assassin who builds or assumes control of suitable headquarters becomes a lord and is eligible to attract a murderous band of thieves. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up assassin abilities table if reverting
DROP TABLE assassin_abilities;