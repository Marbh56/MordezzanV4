-- +goose Up
CREATE TABLE spells (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    mag_level INTEGER NOT NULL DEFAULT 0,  -- Magician level
    cry_level INTEGER NOT NULL DEFAULT 0,  -- Cryomancer level
    ill_level INTEGER NOT NULL DEFAULT 0,  -- Illusionist level
    nec_level INTEGER NOT NULL DEFAULT 0,  -- Necromancer level
    pyr_level INTEGER NOT NULL DEFAULT 0,  -- Pyromancer level
    wch_level INTEGER NOT NULL DEFAULT 0,  -- Witch level
    clr_level INTEGER NOT NULL DEFAULT 0,  -- Cleric level
    drd_level INTEGER NOT NULL DEFAULT 0,  -- Druid level
    range TEXT NOT NULL,
    duration TEXT NOT NULL,
    area_of_effect TEXT,
    components TEXT,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE spells;