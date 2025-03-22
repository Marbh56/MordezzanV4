-- +goose Up
INSERT INTO abilities (name, description) VALUES
('Enlist Henchmen', 'At 6th level, a barbarian may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.'),
('Melee Expert', 'At 7th level, a barbarian''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.'),
('Lordship', 'At 9th level, a barbarian who builds or assumes control of a wilderness fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.');

INSERT INTO abilities (name, description) VALUES
('Agile', '+1 AC bonus when unarmoured and unencumbered (shield allowed).'),
('Alertness', 'Reduces by one (−1) on a d6 roll the party''s chance to be surprized.'),
('Ambusher', 'When traversing the wilds alone or with others of like ability, the barbarian''s base surprize chance increases by one (+1) on a d6 roll. Furthermore, when outdoors, even an untrained party''s chance to surprize increases by one (+1) if the barbarian positions and prepares them accordingly.'),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear, if lightly armoured or unarmoured, as a thief of equal level (see Table 16). Chance-in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.'),
('Draw Poison', 'To draw and spit poison from a snakebite or another venomous wound, such as a scorpion sting or spider bite. The attempt must be made within 2 rounds of affliction for a 3-in-6 chance of success, within 4 rounds for a 2-in-6 chance of success, or within 6 rounds for a 1-in-6 chance of success. Success may revive one who has expired from poison, so long as a successful trauma survival check is made (see Chapter 3: Statistics, constitution). The deceased poison victim is restored to 0 hp, albeit at a price: permanent loss of 1 constitution point. N.B.: Victims of envenomed blades or ingested poison are beyond the barbarian''s aid.'),
('Extraordinary', '+8% chance to perform extraordinary feats of strength and dexterity (see Chapter 3: Statistics, strength and dexterity).'),
('Hardy', 'Physical resilience and an indomitable will to prevail; +2 bonus to all saving throws.'),
('Horsemanship', 'Many barbarians are exceptional horsemen, hailing from nomadic tribes that rely on their steeds in times of peace and war. Even the most stubborn of mounts submit to the barbarian''s will. From the saddle of a tamed mount a barbarian can fight with melee weapons and discharge missiles. Depending on geography and background, this skill may apply to camels.'),
('Leap', 'Mighty thews enable leaps of 25 feet or greater (if unencumbered), bridging pits, chasms, and the like. Vertical leaps of up to 5 feet can also be accomplished.'),
('Move Silently', 'To stalk as a panther, moving with preternatural quiet, comparable to a thief of equal level (see Table 16), if the barbarian is lightly armoured or unarmoured. Chance-in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. This skill is executed at half the barbarian''s normal movement rate.'),
('Run', 'To move as swiftly as a tiger; base 50 MV when lightly armoured or unarmoured.'),
('Sense Magic', 'Sorcery raises the hackles of animal fear and superstition; 4-in-12 chance to cognize the presence of magic if the barbarian noses for it. Discerning the precise source is not always possible; merely that it is close at work. This ability does not function as the detect magic spell and usually does not apply to minor magical items and like dweomers.'),
('Sorcerous Distrust', 'Suspicious of sorcery and those who wield it. Some barbarians may not tolerate the company of magicians, but they might esteem tribal shamans, druids, and the like. A barbarian may wield a magical weapon or be girded with a magical belt, but they are unlikely to be bedecked with all manner of dweomered amulets, cloaks, rings, and other trinkets; such behaviour is contrary to their nature. The extent of the barbarian''s sorcerous distrust is best established through individual role-play.'),
('Track', 'To stalk prey, tracing physical signs and scenting as a predator. A barbarian can track at the below suggested probabilities: Wilderness: A base 10-in-12 chance to find, identify, and follow fresh tracks outdoors or in natural caverns. Non-Wilderness: A base 3-in-12 chance to discern tracks in a dungeon, castle, city street, or like setting. Furthermore, the barbarian can identify in general terms the species tracked if it is a known animal type (e.g., a large feline, a heavy bovine, a small canine). N.B.: The referee may decrease or improve the chance-in-twelve to track based on prevailing circumstances.'),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.'),
('Wilderness Survival', 'Hunting, trapping, fishing, boating, shelter building, fire building (including tribal smoke signals), logging, woodworking, raft building, and so on. These tasks are performed without need of a check; they are simply the barbarian''s province. Under adverse conditions, the referee may assign a reasonable probability of success. Whether a chance of failure applies is at the discretion of the referee, as reflected by the prevailing conditions and abilities of the barbarian.');


INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Barbarian', id, 1 FROM abilities 
WHERE name IN ('Agile', 'Alertness', 'Ambusher', 'Climb', 'Draw Poison', 
              'Extraordinary', 'Hardy', 'Horsemanship', 'Leap', 
              'Move Silently', 'Run', 'Sense Magic', 'Sorcerous Distrust',
              'Track', 'Weapon Mastery', 'Wilderness Survival');

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Barbarian', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Barbarian', id, 7 FROM abilities 
WHERE name = 'Melee Expert';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Barbarian', id, 9 FROM abilities 
WHERE name = 'Lordship';

-- +goose Down
DELETE FROM class_ability_mapping 
WHERE class_name = 'Barbarian';

DELETE FROM abilities 
WHERE name IN ('Agile', 'Alertness', 'Ambusher', 'Climb', 'Draw Poison', 
              'Extraordinary', 'Hardy', 'Horsemanship', 'Leap', 'Move Silently', 
              'Run', 'Sense Magic', 'Sorcerous Distrust', 'Track', 'Weapon Mastery', 
              'Wilderness Survival', 'Enlist Henchmen', 'Melee Expert', 'Lordship');