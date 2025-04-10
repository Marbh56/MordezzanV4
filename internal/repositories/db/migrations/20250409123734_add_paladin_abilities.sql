-- +goose Up
-- Create class-specific table for paladin abilities
CREATE TABLE paladin_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all paladin abilities directly into the class-specific table
INSERT INTO paladin_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Divine Protection', 'The paladin''s indomitable spirit and blameless conduct attract the favour of benign otherworldly forces of Law that provide immunity to natural diseases, a +2 bonus to all saving throws, and a +1 AC bonus versus attacks made by Evil creatures. For fell paladins, the favour of malefic otherworldly forces of Law is enjoyed, and the +1 AC bonus is versus attacks made by Good creatures. N.B.: Divine protection does not provide immunity to magical diseases such as lycanthropy.', 1),
('Extraordinary', '+8% chance to perform extraordinary feats of strength (see "Chapter 3: Statistics, strength).', 1),
('Healing Hands', 'The power to mend wounds and alleviate disease by laying palms on the injured or afflicted. The paladin can restore up to 2 hp per day, per level of experience. The paladin also can cure disease (as the spell) once every seven days. With few exceptions, the paladin will treat allies first. Fell paladins do not have healing hands; rather, they have sapping hands, a touch attack that drains the victim of up to 2 hp per day, per level of experience. Classed individuals (PCs and NPCs) can be drained to –10 hp (death) and monsters to 0 hp (also death). The sapped hit points are transferred to the fell paladin, restoring previously sustained damage. Too, fell paladins can inflict disease (as the reverse spell of cure disease) once every seven days.', 1),
('Honour', 'To comport oneself to a code of honour extolling strength, skill, stoicism, consistency, fidelity, courage in the face of enemies, clemency towards defeated opponents, largess towards dependents, hospitality to associates and superiors, and a willingness to protect the weak. A paladin must never commit murder, perpetrate a felony, or utilize poison. A paladin must oppose tyranny, despotism, cruelty, dæmonism, and other forms of Evil. Failure to adhere to these precepts is grounds for penalization by the referee, possibly including experience point reduction, denial of all supernatural abilities, or level loss. In the worst cases (murder, damning innocents, relations with dæmons, etc.) a paladin can metamorphose, embracing the wickedness of abomination and ultimately transforming into a fell paladin. Such beings become honour-bound paragons of the Lawful Evil alignment.', 1),
('Horsemanship', 'Trained in mounted combat from their earliest youth, paladins can fight from the saddle, urge their mounts to nimble feats on the battlefield, and engage in close-ordered charges. Depending on geography and background, this skill may apply to camels.', 1),
('Sense Evil', 'Perspicacity to Evil most palpable—the nearby presence of a purely Evil sorcerer, undead, dæmons, and other unclean spirits. Particularly strong emanations, such as from a malign artefact or dominion of Evil, may eclipse lesser sensations. Note that this ability will not discern if another character is of Evil alignment unless the subject is about to commit a most vile act or is of a pure and intense Evil (e.g., empowered by dæmons; too, necromancers, witches, and certain priests might qualify, per referee discretion). In any case, the paladin must stop and concentrate, sensing in a 60-foot range (cf. the spell, detect evil). Fell paladins sense similarly, detecting kindred powers.', 1),
('Valiant Resolve', 'Immune to the effects of magically induced fear.', 1),
('Weapon Mastery', 'Mastery of one weapon (+1 "to hit" and +1 damage). Additional weapons may be mastered at 4th, 8th, and 12th levels. As noted in Chapter 6: Equipment, the attack rate for melee weapons and the rates of fire for most missile weapons improve through weapon mastery. For more information on weapon mastery, see Chapter 9: Combat, weapon skill.', 1),

-- Level 3 abilities
('Righteous Wrath', 'At 3rd level, when delivering a charge attack (mounted or afoot) against an Evil foe, paladins gain a damage bonus equal to their level of experience (replacing the +2 damage bonus normally associated with charge attacks; –2 AC penalty applies). When mounted, this bonus is in addition to the double damage inflicted by a lance. For fell paladins, this ability is against Good foes.', 3),

-- Level 5 abilities
('Sacred Mount', 'At 5th level or later, receive a vision and thus learn the location of an extraordinary mount. This thewy, wild stallion is of keen senses and great resolve. The mount must be quested after, lassoed, and trained. This equine is a heavy warhorse of superior health and exceptional wisdom (maximum hit points, 12 morale). This benison can be realized but once per year at most. At 10th level, the paladin may seek the fabled pegasus, and a fell paladin can quest for a nightmare (both maximum hit points, 12 morale).', 5),
('Turn Undead', 'At 5th level, exert control over the undead, causing them to flee and/or cower. Refer to Table 13 at the cleric class entry. At 5th level the paladin has 1st-level turning ability (TA 1); at 6th level, TA 2; and so on. The paladin must stand before the undead and speak boldly a commandment of faith and/or Law, displaying a holy symbol (or Lawful crest) whilst so doing. This ability can be used once per day per TA; other rules and restrictions are noted in the cleric class ability of turn undead. Conversely, fell paladins can command undead on a successful turn undead check. For more information, refer to evil command of undead in the turn undead cleric class ability.', 5),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a paladin may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6),

-- Level 7 abilities
('Melee Expert', 'At 7th level, a paladin''s standard melee attack rate improves from 1/1 to 3/2. Note that weapon mastery can further increase attack rate.', 7),
('Scroll Use', 'At 7th level, decipher and invoke scrolls with spells that are included in the Cleric Spell List (see Chapter 7: Sorcery, Table 99), unless the scroll was created by a thaumaturgical sorcerer (one who casts magician or magician subclass spells).', 7),
('Sorcery', 'At 7th level, paladins gain the ability to cast spells drawn from the Cleric Spell List (see Chapter 7: Sorcery, Table 99). These spells are learnt through prayer, study of scripture, or communion with otherworldly beings associated with Lawful Good; conversely, fell paladins commune with otherworldly beings associated with Lawful Evil. Spell memorization involves prayer, meditation, incantations, and the study and recitation of scriptures dedicated to the tenets of Law. The number and levels of spells cast per day are charted above (see Table 26), though paladins of high wisdom gain bonus spells cast per day (see Chapter 3: Statistics, wisdom). For example, a 9th-level paladin with 13 wisdom can cast three level 1 spells and one level 2 spell per day. The paladin develops a level 1 cleric spell at 7th level, and one new spell at each level gain thereafter. The schedule is as follows: 7th level: ×1 level 1 cleric spell 8th level: ×1 level 1 cleric spell 9th level: ×1 level 2 cleric spell 10th level: ×1 level 2 cleric spell 11th level: ×1 level 3 cleric spell 12th level: ×1 level 3 cleric spell. There is no need of a qualification roll. The paladin cannot learn additional spells beyond those acquired during level training. Note that a 7th-level paladin has 1st-level casting ability (CA 1) and progresses accordingly.', 7),

-- Level 9 abilities
('Lordship', 'At 9th level, a paladin who builds or assumes control of a stronghold becomes a lord and is eligible to attract troops. More information is presented in Appendix B.', 9);

-- +goose Down
-- Clean up paladin abilities table if reverting
DROP TABLE paladin_abilities;