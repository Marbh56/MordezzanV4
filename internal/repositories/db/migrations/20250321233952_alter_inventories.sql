-- +goose Up
ALTER TABLE inventories ADD COLUMN base_encumbered REAL NOT NULL DEFAULT 75.0;
ALTER TABLE inventories ADD COLUMN base_heavy_encumbered REAL NOT NULL DEFAULT 150.0;
ALTER TABLE inventories ADD COLUMN maximum_capacity REAL NOT NULL DEFAULT 300.0;
ALTER TABLE inventories ADD COLUMN is_encumbered BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE inventories ADD COLUMN is_heavy_encumbered BOOLEAN NOT NULL DEFAULT 0;
ALTER TABLE inventories ADD COLUMN is_overloaded BOOLEAN NOT NULL DEFAULT 0;

CREATE INDEX idx_inventory_encumbrance ON inventories(character_id, is_encumbered, is_heavy_encumbered, is_overloaded);

-- +goose Down
DROP INDEX IF EXISTS idx_inventory_encumbrance;
ALTER TABLE inventories DROP COLUMN is_overloaded;
ALTER TABLE inventories DROP COLUMN is_heavy_encumbered;
ALTER TABLE inventories DROP COLUMN is_encumbered;
ALTER TABLE inventories DROP COLUMN maximum_capacity;
ALTER TABLE inventories DROP COLUMN base_heavy_encumbered;
ALTER TABLE inventories DROP COLUMN base_encumbered;