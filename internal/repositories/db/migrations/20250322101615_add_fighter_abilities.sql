-- +goose Up
-- Create class-specific table for fighter abilities
CREATE TABLE fighter_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all fighter abilities directly into the class-specific table
INSERT INTO fighter_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Heroic Fighting', 'To smite multiple foes. When combatting opponents of 1 HD or less, double normal melee attacks per round (2/1, or 3/1 if wielding a mastered weapon). This dramatic attack could be effected as a single, devastating swing or lunge that bursts through multiple foes. At 7th level, when combating foes of 2 HD or less, double normal melee attacks per round (3/1, or 4/1 if wielding a mastered weapon).', 1),
('Weapon Mastery', 'Mastery of two weapons (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels; however, see grand mastery below, for another option. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength and dexterity (see Chapter 3: Statistics, strength and dexterity).', 1),

-- Higher level abilities
('Grand Mastery', 'At 4th, 8th, or 12th level (player''s choice), when a new weapon mastery is gained, fighters may elect to intensify their training with an already mastered weapon. With this weapon the fighter becomes a grand master (+2 "to hit" and +2 damage, increased attack rate, etc.). A fighter may achieve grand mastery with but one weapon. For more information, see Chapter 9: Combat, weapon skill.', 4),
('Enlist Henchmen', 'At 6th level, a fighter may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),
('Melee Expert', 'At 7th level, a fighter''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7);

-- +goose Down
-- Clean up fighter abilities table if reverting
DROP TABLE fighter_abilities;