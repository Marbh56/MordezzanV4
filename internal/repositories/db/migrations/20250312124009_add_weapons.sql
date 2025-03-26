-- +goose Up
CREATE TABLE weapons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    weapon_class INTEGER NOT NULL DEFAULT 1,
    cost REAL NOT NULL,
    weight INTEGER NOT NULL,
    range_short INTEGER,
    range_medium INTEGER,
    range_long INTEGER,
    rate_of_fire TEXT,
    damage TEXT NOT NULL,
    damage_two_handed TEXT,
    properties TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO weapons (name, category, weapon_class, cost, weight, range_short, range_medium, range_long, rate_of_fire, damage, damage_two_handed, properties)
VALUES
('Hand Axe', 'Melee', 1, 5, 2, 15, 30, 45, '1/1', '1d6', '-', '-'),
('Battle Axe', 'Melee', 2, 10, 5, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Great Axe', 'Melee', 4, 20, 10, 0, 0, 0, '-', '2d6', '-', '+, #'),
('Cæstuses', 'Melee', 0, 1, 1, 0, 0, 0, '-', '+1', '-', '-'),
('Chain Whip', 'Melee', 4, 10, 3, 0, 0, 0, '-', '1d6', '-', '↵'),
('Light Club', 'Melee', 1, 1, 2, 10, 20, 30, '1/1', '1d4', '-', '-'),
('War Club', 'Melee', 2, 3, 4, 0, 0, 0, '-', '1d6', '1d8', '-'),
('Dagger', 'Melee', 1, 4, 1, 10, 20, 30, '3/2', '1d4', '-', '-'),
('Silver Dagger', 'Melee', 1, 25, 1, 10, 20, 30, '3/2', '1d4', '-', '-'),
('Falcata', 'Melee', 1, 10, 3, 0, 0, 0, '-', '1d6', '-', '-'),
('Horseman''s Flail', 'Melee', 1, 5, 3, 0, 0, 0, '-', '1d6', '-', '↵'),
('Footman''s Flail', 'Melee', 3, 10, 10, 0, 0, 0, '-', '1d10', '-', '↵, +'),
('Halberd', 'Melee', 4, 15, 8, 0, 0, 0, '-', '1d10', '-', '+, ^, #'),
('Horseman''s Hammer', 'Melee', 1, 5, 3, 10, 20, 30, '1/1', '1d6', '-', '-'),
('War Hammer', 'Melee', 2, 10, 5, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Great Hammer', 'Melee', 4, 20, 10, 0, 0, 0, '-', '2d6', '-', '+, #'),
('Javelin', 'Melee', 2, 3, 3, 20, 40, 80, '1/1', '1d4', '1d6', '-'),
('Lance', 'Melee', 5, 15, 8, 0, 0, 0, '-', '1d8', '-', '^ ∇ o'),
('Horseman''s Mace', 'Melee', 1, 4, 3, 0, 0, 0, '-', '1d6', '-', '-'),
('Footman''s Mace', 'Melee', 2, 10, 5, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Great Mace', 'Melee', 4, 20, 10, 0, 0, 0, '-', '2d6', '-', '+, #'),
('Morning Star', 'Melee', 2, 15, 5, 0, 0, 0, '-', '1d8', '1d10', 'Ω'),
('Horseman''s Pick', 'Melee', 1, 5, 3, 0, 0, 0, '-', '1d6', '-', 'Ω'),
('War Pick', 'Melee', 2, 15, 5, 0, 0, 0, '-', '1d8', '1d10', 'Ω'),
('Pike', 'Melee', 6, 7, 12, 0, 0, 0, '-', '1d8', '-', '+, ^'),
('Quarterstaff', 'Melee', 3, 5, 5, 0, 0, 0, '-', '1d6', '-', '↔'),
('Short Scimitar', 'Melee', 1, 10, 3, 0, 0, 0, '-', '1d6', '-', '-'),
('Long Scimitar', 'Melee', 2, 20, 4, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Two-Handed Scimitar', 'Melee', 4, 40, 6, 0, 0, 0, '-', '3d4', '-', '+'),
('Sickle', 'Melee', 1, 3, 2, 0, 0, 0, '-', '1d4', '-', '-'),
('Short Spear', 'Melee', 3, 4, 5, 15, 30, 45, '1/1', '1d6', '1d8', '^'),
('Long Spear', 'Melee', 4, 5, 7, 0, 0, 0, '-', '1d6', '1d8', '^'),
('Great Spear', 'Melee', 5, 7, 9, 0, 0, 0, '-', '1d8', '-', '+, ^, ∇'),
('Spiked Staff', 'Melee', 3, 15, 7, 0, 0, 0, '-', '1d10', '-', '+, ^, #'),
('Short Sword', 'Melee', 1, 10, 3, 0, 0, 0, '-', '1d6', '-', '-'),
('Broad Sword', 'Melee', 2, 20, 4, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Long Sword', 'Melee', 2, 20, 4, 0, 0, 0, '-', '1d8', '1d10', '-'),
('Bastard Sword', 'Melee', 3, 30, 5, 0, 0, 0, '-', '1d8', '2d6', '-'),
('Two-Handed Sword', 'Melee', 4, 40, 6, 0, 0, 0, '-', '3d4', '-', '+'),
('Tonfa', 'Melee', 1, 2, 1, 0, 0, 0, '-', '1d4', '-', '↔'),
('Hand Trident', 'Melee', 1, 7, 1, 0, 0, 0, '-', '1d4', '-', '↔'),
('Long Trident', 'Melee', 4, 10, 6, 10, 20, 30, '1/1', '1d6', '1d8', '^'),
('Whip', 'Melee', 5, 1, 2, 0, 0, 0, '-', '1d2', '-', '-'),
('Bola', 'Hurled Missile Type', 0, 3, 2, 15, 30, 45, '1/1', '1d2', '-', '⤢'),
('Bommerang', 'Hurled Missile Type', 0, 1, 1, 50, 100, 150, '1/1', '1d2', '-', '⤢'),
('Dart', 'Hurled Missile Type', 0, 1, 1, 15, 30, 45, '2/1', '1d3', '-', '⤢' ),
('Hooked Throwing Knife', 'Hurled Missile Type', 0, 20, 2, 30, 60, 90, '1/1', '1d6', '-', '↵, ⤢'),
('Lasso', 'Hurled Missile Type', 0, 3, 3, 0, 20, 0, '1/2', '-', '-', '-'),
('Fighting Net', 'Hurled Missile Type', 0, 5, 7, 0, 10, 0, '1/2', '-', '-', '-'),
('Stone', 'Hurled Missile Type', 0, 0, 1, 30, 60, 90, '2/1', '1', '-', '⤢'),
('Blowgun', 'Launched Missile Type', 0, 5, 1, 30, 60, 90, '1/1', '1', '-', '-'),
('Bow, Long-', 'Launched Missile Type', 0, 60, 3, 70, 140, 210, '3/2', '1d6', '-', '⤤'),
('Bow, Long-, Composite', 'Launched Missile Type', 0, 100, 3, 80, 160, 240, '3/2', '1d6', '-', '⤤'),
('Bow, Short', 'Launched Missile Type', 0, 20, 2, 50, 100, 150, '3/2', '1d6', '-', '-'),
('Bow, Short, Composite', 'Launched Missile Type', 0, 50, 2, 60, 120, 180, '3/2', '1d6', '-', '-'),
('Crossbow, Heavy', 'Launched Missile Type', 0, 25, 10, 80, 160, 240, '1/2', '1d6+2', '-', '-'),
('Crossbow, Light', 'Launched Missile Type', 0, 15, 5, 60, 120, 180, '1/1', '1d6+1', '-', '-'),
('Crossbow, Repeating', 'Launched Missile Type', 0, 100, 6, 50, 100, 150, '3/1', '1d6', '-', '-'),
('Sling', 'Launched Missile Type', 0, 2, 1, 50, 100, 150, '1/1', '1d4', '-', '⤢');

-- +goose Down
DROP TABLE weapons;
