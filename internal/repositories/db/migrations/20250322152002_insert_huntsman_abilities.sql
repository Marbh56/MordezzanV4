-- +goose Up
-- Insert Huntsman-specific abilities
INSERT INTO abilities (name, description) VALUES
('Harvest Venom', 'To extract venom sacs from snakes, spiders, giant insects, and the like, should the opportunity present and appropriate containers be available. Huntsmen learn to dissect venomous creatures. This skill is performed at a base 9-in-12 chance of success. If a 12 is rolled, a mishap occurs, and the huntsman is exposed (e.g., eyes, nose, skin) to the poisonous fluid. For more information on the extraction of venom, see Chapter 9: Combat, poison, venom harvesting.'),
('Predator', 'Trained from earliest youth to hunt and kill animals both fleet and robust. The huntsman inflicts +1 hp damage per level of experience when combatting normal and giant-sized animals: amphibians, birds, crustaceans, dinosaurs, fishes, insects, mammals, and reptiles. Does not apply to constructs, dæmons, elementals, fungi, giants, humanoids, magical beasts, moulds, oozes, otherworldly and alien beings, slimes and jellies, or undead. When used in concert with a successful hide attempt, the initial attack roll is made at a +4 bonus.'),
('Subdue Animal', 'To soothe and tame an animal (normal, not magical) of hit dice equal to or less than the huntsman''s level. To succeed, the following steps must be completed: Through combat, physically reduce the animal to half or less its hit point total (the huntsman can assess this with accuracy). The huntsman must be a prominent aggressor in the beast''s impairment. Restrain the creature. Many a huntsman will use a bola, lasso, fighting net, or whip to make prone the target before attempting to restrain it. On the round following restraint, the huntsman attempts to assert mental and physical dominance. The base chance of success is 4-in-12. This chance-in-twelve may be increased by the following modifiers: +1 if the huntsman''s strength is 16+, +1 if the huntsman''s wisdom is 16+, +1 if the huntsman''s charisma is 16+, +1 if the huntsman''s level is 7+, +1 if the huntsman has dominated a member of this species before. Failure indicates the animal is impossible to tame. It may continue attempting to break free. Success indicates the animal is subdued; the huntsman must continue to restrain the creature for 1 turn (10 minutes), kneeling on it, commanding it, and forcing submission. Thereafter, it will be docile and relatively obedient. A defeated animal can be tamed to complete loyalty (ML 12) after 1d4 months of training. It can be trained to attack, fetch, guard, hunt, track, or perform other tasks. Multiple animals can be trained—even working in flawless synchronization if they are reasonably compatible—but their total hit dice can never exceed the huntsman''s level, and the training time for multiple animals is cumulative.'),
('Wilderness Traps', 'To set an outdoor trap, including pits, deadfalls (falling logs/rocks), snares, and spring traps. The huntsman is also adept at finding and removing such traps. These tasks are performed as a thief of equal level performs the manipulate traps skill (see Table 16), but the huntsman has no facility with mechanical and/or magical traps.'),
('Werewolf Slayer', 'At 4th level, huntsmen develop the aptitude to slay lycanthropes. Indeed, when men and women suffer the curse of the beast, huntsmen rise to stamp them out. When wielding silver or magical weapons versus lycanthropes, huntsmen gain all the benefits of the predator ability, regardless of the advanced intelligence of the afflicted.');

-- Link Huntsman abilities at level 1
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Huntsman', id, 1 FROM abilities 
WHERE name IN ('Alertness', 'Ambusher', 'Climb', 'Extraordinary', 'Harvest Venom', 
               'Hide', 'Move Silently', 'Predator', 'Subdue Animal', 'Track', 
               'Weapon Mastery', 'Wilderness Survival', 'Wilderness Traps');

-- Link Werewolf Slayer at level 4
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Huntsman', id, 4 FROM abilities 
WHERE name = 'Werewolf Slayer';

-- Link existing abilities to Huntsman at appropriate levels
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Huntsman', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Huntsman', id, 7 FROM abilities 
WHERE name = 'Melee Expert';

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Huntsman', id, 9 FROM abilities 
WHERE name = 'Lordship';

-- Add thief skills for Huntsman
INSERT INTO class_thief_skill_mapping (class_name, skill_id)
SELECT 'Huntsman', id FROM thief_skills 
WHERE skill_name IN ('Move Silently', 'Climb', 'Hide');

-- +goose Down
-- Remove Huntsman ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Huntsman';

-- Remove Huntsman thief skill mappings
DELETE FROM class_thief_skill_mapping
WHERE class_name = 'Huntsman';

-- Clean up any Huntsman-specific abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Harvest Venom', 'Predator', 'Subdue Animal', 'Werewolf Slayer', 'Wilderness Traps')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);