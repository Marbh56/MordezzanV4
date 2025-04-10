-- +goose Up
-- Create class-specific table for monk abilities
CREATE TABLE monk_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all monk abilities directly into the class-specific table
INSERT INTO monk_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Block Missile', 'Deflect a fired arrow, bolt, or bullet; likewise, a hurled axe, boomerang, dagger, dart, javelin, spear, or like weapon. Even spells such as flaming missile, magic ice dart, magic missile or acid arrow can be blocked. Siege missiles, giant-hurled boulders, and the like do not apply. To repel a missile, the monk must make an avoidance saving throw. The monk must drop anything held to use this ability, which can be attempted as many times per round as the monk has levels of experience.', 1),
('Cellular Adjustment', 'Once per day, supernaturally alter cells to heal physical damage equal to 2 hp per level of experience. Alternatively, the monk may confer this healing onto an injured ally. Also, once per week the monk can purge him- or herself or another of disease or poison, per the spells cure disease or neutralize poison, though not lycanthropy.', 1),
('Controlled Fall', 'To retard descent of precipitous falls. For every level of experience, the monk can fall 10 feet and sustain no damage, so long as a wall or other stable surface is within a five-foot reach throughout the descent. For falls beyond the monk''s limit, normal rules are in force starting at the point at which the controlled fall no longer applies; e.g., a 5th-level monk plummets down a 90-foot pit and thus sustains 4d6 hp damage.', 1),
('Defensive Ability', 'To avoid and deflect blows and damage through physical, mental, and spiritual (ka/qi) mastery. The monk gains an AC bonus that increases as levels of experience are gained (see Table 44).', 1),
('Detect Secret Doors', 'Find a secret door on a base 3-in-6 chance.', 1),
('Empty Hand', 'Master of the unarmed attack (hand, foot, knee, elbow, etc.). The monk enjoys the following benefits when fighting sans weapons: 2/1 attack rate; requires two free hands, attacks not exclusively by hand. Weapon-like damage (see Table 44). The monk may use cæstuses (leather thongs wrapped around the hands and weighted with iron or lead plates or spikes) for a +1 damage bonus. On a natural 20 attack roll, a stunning blow can be delivered. A Small or Medium creature must make a transformation saving throw or be stunned for 2d4 rounds. At 7th level, Large creatures can be stunned, but they are afforded a +4 bonus to the save. Does not affect undead, constructs, oozes, slimes, and the like. If the optional critical hits rule is used (see Chapter 9: Combat, critical hits and misses), the target is stunned in addition to any bonus damage inflicted. At 5th level, the empty hand attack (due to heightened ka/qi) is equivalent to a magical weapon; the monk gains a +1 "to hit" bonus. At 12th level, once per day, deliver a quivering palm death blow. The monk mystically vibrates the empty hand to match the rhythm of the target''s heart or other vital organ. If hit, the victim must make a death saving throw or die instantly; otherwise, normal damage applies. If the attack misses, subsequent attempts may be made, so long as the monk does nothing else but focus on the quivering palm attacks. Quivering palm has no effect on the undead, constructs, oozes, slimes, and the like.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of dexterity (see Chapter 3: Statistics, dexterity).', 1),
('Run', 'To move as swiftly as a tiger. If unarmoured, achieve a movement rate of 50; at 7th level, this speed increases to a superhuman MV 60 for short bursts.', 1),
('Superior Willpower', '+2 bonus to saving throws versus any sorcery that would influence the mind, including illusions, charms, etc. This bonus is cumulative with willpower adjustment, if applicable (see Chapter 3: Statistics, wisdom).', 1),
('Climb', 'To ascend or descend sheer cliffs or walls without need of climbing gear, as a thief of equal level. If vertical, the surface must be rough or cracked. At least one check must be made per 100 feet of climbing. Failure indicates the climber has slipped and fallen at about the midpoint of the check (however, see controlled fall ability).', 1),
('Discern Noise', 'To hearken at a door and detect the faintest of noises on the other side, perceive the distant footfalls of a wandering monster, or distinguish a single voice in a crowd, as a thief of equal level. Six rounds (one minute) of concentrated listening are required.', 1),
('Hide', 'To vanish into shadows, camouflage oneself, or flatten one''s body to a seemingly impossible degree—all whilst remaining still as a statue. This ability is performed as a thief of equal level. Only the slightest movement is permissible (e.g., unsheathing a blade, opening a pouch). Hiding is impossible in direct sunlight, or if the monk is observed.', 1),
('Move Silently', 'To move with preternatural quiet, even across squeaky floorboards, dry leaves, loose debris, and the like, as a thief of equal level. This skill is executed at half the monk''s normal movement rate.', 1),

-- Level 4 abilities
('New Weapon Skill', 'At 4th, 8th, and 12th levels, become skilled in a new weapon that is not included in the favoured weapons list (noted above). This new proficiency is dependent upon training and practice (see Chapter 9: Combat, weapon skill).', 4),
('Speak With Nature', 'At 4th level, speak with animals (as the spell) once per day; at 8th level, also speak with plants (as the spell) once per day.', 4),

-- Level 5 abilities
('Simulate Death', 'At 5th level of experience, enter a deep trance in which the monk can feign a deathlike condition, as per the cataleptic state spell (q.v.).', 5),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a monk may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 9 abilities
('Lordship', 'At 9th level, a monk who builds or assumes control of a monastery becomes a lord and is eligible to attract followers. More information is presented in Appendix B.', 9),

-- Level 11 abilities
('Longevity', 'At 11th level, ageing process slows. For every 13 years (1 cycle), the monk effectively ages but 1 year.', 11);

-- +goose Down
-- Clean up monk abilities table if reverting
DROP TABLE monk_abilities;