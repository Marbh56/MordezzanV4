-- +goose Up
-- Create unified class data table
CREATE TABLE class_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    class_name TEXT NOT NULL,
    level INTEGER NOT NULL,
    experience_points INTEGER NOT NULL,
    hit_dice TEXT NOT NULL,
    saving_throw INTEGER NOT NULL,
    fighting_ability INTEGER NOT NULL,
    casting_ability INTEGER DEFAULT 0,
    spell_slots_level1 INTEGER DEFAULT 0,
    spell_slots_level2 INTEGER DEFAULT 0,
    spell_slots_level3 INTEGER DEFAULT 0,
    spell_slots_level4 INTEGER DEFAULT 0,
    spell_slots_level5 INTEGER DEFAULT 0,
    spell_slots_level6 INTEGER DEFAULT 0,
    UNIQUE(class_name, level)
);

-- Create abilities table
CREATE TABLE abilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    UNIQUE(name)
);

-- Create junction table for class-ability mapping
CREATE TABLE class_ability_mapping (
    class_name TEXT NOT NULL,
    ability_id INTEGER NOT NULL,
    min_level INTEGER NOT NULL,
    PRIMARY KEY (class_name, ability_id),
    FOREIGN KEY (class_name) REFERENCES class_data(class_name),
    FOREIGN KEY (ability_id) REFERENCES abilities(id)
);

-- +goose Down
DROP TABLE class_ability_mapping;
DROP TABLE abilities;
DROP TABLE class_data;