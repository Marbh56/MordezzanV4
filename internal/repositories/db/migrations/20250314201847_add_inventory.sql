-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE inventories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    character_id INTEGER NOT NULL,
    max_weight REAL NOT NULL DEFAULT 0,
    current_weight REAL NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE
);

CREATE TABLE inventory_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inventory_id INTEGER NOT NULL,
    item_type TEXT NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    is_equipped BOOLEAN NOT NULL DEFAULT 0,
    slot TEXT,
    notes TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (inventory_id) REFERENCES inventories(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_inventory_character_id ON inventories(character_id);
CREATE INDEX idx_inventory_items_inventory_id ON inventory_items(inventory_id);
CREATE INDEX idx_inventory_items_type_and_id ON inventory_items(item_type, item_id);
CREATE INDEX idx_inventory_items_slot ON inventory_items(slot);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP INDEX IF EXISTS idx_inventory_items_slot;
DROP INDEX IF EXISTS idx_inventory_items_type_and_id;
DROP INDEX IF EXISTS idx_inventory_items_inventory_id;
DROP INDEX IF EXISTS idx_inventory_character_id;
DROP TRIGGER IF EXISTS update_inventory_items_timestamp;
DROP TRIGGER IF EXISTS update_inventories_timestamp;
DROP TABLE IF EXISTS inventory_items;
DROP TABLE IF EXISTS inventories;