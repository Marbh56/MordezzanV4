-- +goose Up
-- Create class-specific table for runegraver abilities
CREATE TABLE runegraver_abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    min_level INTEGER NOT NULL DEFAULT 1,
    UNIQUE(name)
);

-- Insert all runegraver abilities directly into the class-specific table
INSERT INTO runegraver_abilities (name, description, min_level) VALUES
-- Level 1 abilities
('Extraordinary', '+8% chance to perform extraordinary feats of strength (see Chapter 3: Statistics, strength).', 1),
('Grave Runes', 'Runegravers grave runes on specific materials that are carved, smoothed, and shaped by the runegraver. Each rune emulates a specific spell, with no other material components required. The number and levels of rune spells cast per day are charted above (see Table 49). Runes are enchanted when the runegraver cuts his or her palm (no damage at this time), drips blood on the rune, and then recites the appropriate poem. Much like spells memorized from books, when a rune spell is cast, the rune and its medium remain, but their sorcery is drained until the process is started anew. Runes are invoked (cast) much like a spell, except each invocation costs the runegraver a sacrifice of 1 hit point per spell level, pain coursing up the runegraver''s arms. So, casting a level 3 rune spell entails a loss of 3 hp, which can be recovered in the usual manner. Should the runegraver drop below 1 hp from evoking a rune, unconsciousness is staved off for 1 turn, pending an extraordinary feat of constitution. Each of the 16 runes may be invoked but once per day. The runes and the spells they mimic are as follows: Level 1: giant–enlargement, man–command, riding–mount. Level 2: constraint–hold person, shower–black cloud, wealth–fool''s gold. Level 3: plenty–create food and water, ulcer–inflict disease (reverse of cure disease), yew–twofold missile. Level 4: hail–ice storm (hail), ice–freeze surface, Tyr–dweomered weapon. Level 5: god–true seeing, sun–flame strike, water–control water. Level 6: birch–reincarnation. A runegraver begins with a single rune, a sacred mystery revealed by a master. It is either selected or randomly determined upon character creation (consult your referee). Runegravers unlock the secret of a new rune or runes at each level gain, at the following schedule: 1st level: ×1 level 1 rune. 2nd level: ×1 level 1 rune. 3rd level: ×1 level 1 rune; ×1 level 2 rune. 4th level: ×1 level 2 rune. 5th level: ×1 level 2 rune; ×1 level 3 rune. 6th level: ×1 level 3 rune. 7th level: ×1 level 3 rune; ×1 level 4 rune. 8th level: ×1 level 4 rune. 9th level: ×1 level 4 rune; ×1 level 5 rune. 10th level: ×1 level 5 rune. 11th level: ×1 level 5 rune. 12th level: ×1 level 6 rune. New runes are acquired via spiritual revelation, the uncloaking of runic lore, or deific favour from Ullr. Each new rune is learnt automatically, with no need of a qualification roll. Runegravers do not learn new runes outside of level gains, but by 12th level, all 16 runes have been mastered.', 1),
('Ale Horn', 'An aurochs (or other bovine) drinking horn painstakingly etched in Nordic designs. Once per day, if the horn is filled with 12 ounces of fresh water, the runegraver can turn the water to enchanted ale that restores 2 hp per CA level when drunk. The ale can be shared, but once uncapped, the entire contents must be imbibed within 6 rounds. Once per week, three drops of honey added to the water-filled drinking horn produces magical mead that can cure disease (as the spell).', 1),

-- Level 5 abilities
('Casting of Lots', 'At 5th level, a runegraver can gather 16 fresh twigs. The twigs must be smoothed and engraved, each incised with one of the 16 runes listed earlier. They are then placed in a leather pouch filled with powdered bone. To cast lots, the runegraver must engage in a 1-turn ritual, chanting, shaking the pouch, and finally dumping the twigs from the pouch. The runegraver then articulates a question, selects three bone dust–covered twigs, wipes them clean, and interprets their meaning. Casting of lots is otherwise equivalent to an augury spell. It can be performed once per day.', 5),

-- Level 6 abilities
('Enlist Henchmen', 'At 6th level, a runegraver may seek or be sought out by one or more henchmen, classed individuals (typically of similar class, race, and/or culture) who become loyal followers. For more information, see Chapter 8: Adventure, hirelings and henchmen.', 6);

-- +goose Down
-- Clean up runegraver abilities table if reverting
DROP TABLE runegraver_abilities;