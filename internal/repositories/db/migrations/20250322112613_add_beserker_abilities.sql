-- +goose Up
-- Create class-specific table for berserker abilities
CREATE TABLE berserker_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all berserker abilities directly into the class-specific table
INSERT INTO berserker_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Berserk Rage', 'A furious battle lust and feral madness unleashed. A rage may be entered any time during battle; though it is most often witnessed when the berserker''s blood has been drawn. Frequency: 1st to 4th levels = ×1 per day, 5th to 8th levels = ×2 per day, 9th to 12th levels = ×3 per day. Duration: Rounds equal to the berserker''s constitution score. Benefits and drawbacks: Berserk Attack Rate: Melee attack rate of 2/1 (or 5/2 with weapon mastery). Enhanced Combat: +2 "to hit" and damage on all melee attacks. Fire Immunity: Impervious to normal fire; saves vs. magical fire always successful. Frightening Aspect: Fearsome to behold; enemy morale checks at −2 penalty. Hit Point Burst: Temporary hit points equal to one-half of constitution score, rounded up. Refusal to Fall: Can fight to as low as −3 hp. Refusal to Surrender: Cannot yield, retreat, or withdraw from melee once the rage is begun. Unbreakable Willpower: Immunity to fear, charm, and like sorcery. Uncontrollable: Once all enemies are defeated, on a 1-in-8 chance, the berserker attacks any living creature within 30 feet for 1d6 rounds. Exhaustion: When the rage ends, the berserker is exhausted for 1d3 turns with −2 "to hit" and damage, reduced attack rate, and no running.', 1),
('Thick Skin', 'Flesh not unlike the hide of a bull, which toughens over time. The berserker has natural AC 8 at 1st level, AC 7 at 3rd level, and so on. Body armour does not "stack" with this ability, but thick skin does provide a +1 AC bonus (from 1st to 6th levels) or +2 AC bonus (from 7th to 12th levels) to berserkers clad in light armour. Also, the berserker can function in subfreezing temperatures (as low as −15°F) with little need of protection.', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear, if lightly armoured or unarmoured, as a thief of equal level (see Table 16). Chance-in-twelve reduced by 4 if wearing medium armour; impossible in heavy armour. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength and dexterity (see Chapter 3: Statistics, strength and dexterity).', 1),
('Hardy', 'Physical resilience and an indomitable will to prevail; +2 bonus to all saving throws.', 1),
('Leap', 'Mighty thews enable leaps of 25 feet or greater (if unencumbered), bridging pits, chasms, and the like. Vertical leaps of up to 5 feet can also be accomplished.', 1),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),

-- Higher level abilities
('Enlist Henchmen', 'At 6th level, a berserker may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Bestial Form', 'At 7th level, the berserker is blessed by a deity or spirit with the ability to transmogrify into a semi-human, bipedal shape whilst in berserk rage. The bestial form assumed is one typically associated with the berserker''s culture or ancestry: bear, lion, tiger, or wolf. Benefits: Enlargement: ×1.5 height and ×2 weight. Recovery: Half of any lost hit points are recovered. Melee Weapon Use: Can wield melee weapons and attack as normal. Bestial Attack: Can opt to claw/claw/bite for a base 1d6/1d6/1d8 damage. If both claws strike a single opponent of Small or Medium size, the berserker can hug the victim automatically for an additional 2d6 hp damage. Restrictions: Limited Usage: Bestial form ends when berserk rage ends. Armour: No armour allowed; the transformation rips clothes. Rage Connexion: All benefits and detriments associated with berserk rage remain.', 7),
('Melee Expert', 'At 7th level, a berserker''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7),
('Lordship', 'At 9th level, a berserker who builds or assumes control of a wilderness fortress becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up berserker tables if reverting
DROP TABLE berserker_natural_ac;
DROP TABLE berserker_abilities;