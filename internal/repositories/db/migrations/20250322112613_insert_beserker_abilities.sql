-- +goose Up
-- Insert Berserker abilities
CREATE TEMPORARY TABLE IF NOT EXISTS temp_abilities_to_insert (
    name TEXT,
    description TEXT
);

INSERT INTO temp_abilities_to_insert (name, description) VALUES
('Berserk Rage', 'A furious battle lust and feral madness unleashed. A rage may be entered any time during battle; though it is most often witnessed when the berserker''s blood has been drawn. Frequency: 1st to 4th levels = ×1 per day, 5th to 8th levels = ×2 per day, 9th to 12th levels = ×3 per day. Duration: Rounds equal to the berserker''s constitution score. Benefits and drawbacks: Berserk Attack Rate: Melee attack rate of 2/1 (or 5/2 with weapon mastery). Enhanced Combat: +2 "to hit" and damage on all melee attacks. Fire Immunity: Impervious to normal fire; saves vs. magical fire always successful. Frightening Aspect: Fearsome to behold; enemy morale checks at −2 penalty. Hit Point Burst: Temporary hit points equal to one-half of constitution score, rounded up. Refusal to Fall: Can fight to as low as −3 hp. Refusal to Surrender: Cannot yield, retreat, or withdraw from melee once the rage is begun. Unbreakable Willpower: Immunity to fear, charm, and like sorcery. Uncontrollable: Once all enemies are defeated, on a 1-in-8 chance, the berserker attacks any living creature within 30 feet for 1d6 rounds. Exhaustion: When the rage ends, the berserker is exhausted for 1d3 turns with −2 "to hit" and damage, reduced attack rate, and no running.'),
('Thick Skin', 'Flesh not unlike the hide of a bull, which toughens over time. The berserker has natural AC 8 at 1st level, AC 7 at 3rd level, and so on. Body armour does not "stack" with this ability, but thick skin does provide a +1 AC bonus (from 1st to 6th levels) or +2 AC bonus (from 7th to 12th levels) to berserkers clad in light armour. Also, the berserker can function in subfreezing temperatures (as low as −15°F) with little need of protection.'),
('Bestial Form', 'At 7th level, the berserker is blessed by a deity or spirit with the ability to transmogrify into a semi-human, bipedal shape whilst in berserk rage. The bestial form assumed is one typically associated with the berserker''s culture or ancestry: bear, lion, tiger, or wolf. Benefits: Enlargement: ×1.5 height and ×2 weight. Recovery: Half of any lost hit points are recovered. Melee Weapon Use: Can wield melee weapons and attack as normal. Bestial Attack: Can opt to claw/claw/bite for a base 1d6/1d6/1d8 damage. If both claws strike a single opponent of Small or Medium size, the berserker can hug the victim automatically for an additional 2d6 hp damage. Restrictions: Limited Usage: Bestial form ends when berserk rage ends. Armour: No armour allowed; the transformation rips clothes. Rage Connexion: All benefits and detriments associated with berserk rage remain.');

-- Insert Berserker abilities that don't already exist
INSERT INTO abilities (name, description)
SELECT t.name, t.description
FROM temp_abilities_to_insert t
WHERE NOT EXISTS (SELECT 1 FROM abilities WHERE name = t.name);

DROP TABLE temp_abilities_to_insert;

-- Link Berserker abilities at level 1
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Berserker', id, 1 FROM abilities 
WHERE name IN ('Berserk Rage', 'Thick Skin', 'Climb', 'Extraordinary', 'Hardy', 'Leap', 'Weapon Mastery')
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Berserker' AND ability_id = abilities.id AND min_level = 1
);

-- Link Bestial Form at level 7
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Berserker', id, 7 FROM abilities 
WHERE name = 'Bestial Form'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Berserker' AND ability_id = abilities.id AND min_level = 7
);

-- Link existing abilities to Berserker at appropriate levels
INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Berserker', id, 6 FROM abilities 
WHERE name = 'Enlist Henchmen'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Berserker' AND ability_id = abilities.id AND min_level = 6
);

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Berserker', id, 7 FROM abilities 
WHERE name = 'Melee Expert'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Berserker' AND ability_id = abilities.id AND min_level = 7
);

INSERT INTO class_ability_mapping (class_name, ability_id, min_level)
SELECT 'Berserker', id, 9 FROM abilities 
WHERE name = 'Lordship'
AND NOT EXISTS (
    SELECT 1 FROM class_ability_mapping
    WHERE class_name = 'Berserker' AND ability_id = abilities.id AND min_level = 9
);

-- Create the berserker_natural_ac table if it doesn't exist
CREATE TABLE IF NOT EXISTS berserker_natural_ac (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    natural_ac INTEGER NOT NULL,
    UNIQUE(class_name, level)
);

-- Insert or update natural AC values using a merge operation
INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 1, 8 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 1);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 2, 8 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 2);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 3, 7 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 3);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 4, 7 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 4);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 5, 6 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 5);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 6, 6 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 6);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 7, 5 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 7);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 8, 5 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 8);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 9, 4 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 9);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 10, 4 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 10);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 11, 3 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 11);

INSERT INTO berserker_natural_ac (class_name, level, natural_ac)
SELECT 'Berserker', 12, 3 WHERE NOT EXISTS (SELECT 1 FROM berserker_natural_ac WHERE class_name = 'Berserker' AND level = 12);

-- +goose Down
-- Remove Berserker ability mappings
DELETE FROM class_ability_mapping 
WHERE class_name = 'Berserker' 
AND ability_id IN (
    SELECT id FROM abilities 
    WHERE name IN ('Berserk Rage', 'Thick Skin', 'Bestial Form', 'Climb', 'Extraordinary', 'Hardy', 'Leap', 'Weapon Mastery', 'Enlist Henchmen', 'Melee Expert', 'Lordship')
);

-- Remove berserker-specific abilities that aren't used by other classes
DELETE FROM abilities 
WHERE name IN ('Berserk Rage', 'Thick Skin', 'Bestial Form')
AND id NOT IN (SELECT ability_id FROM class_ability_mapping);

-- We don't drop the berserker_natural_ac table, just clean up our entries
DELETE FROM berserker_natural_ac WHERE class_name = 'Berserker';