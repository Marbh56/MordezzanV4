-- +goose Up
-- SQL in this section is executed when the migration is applied
INSERT INTO shields (name, cost, weight, defense_modifier) VALUES ('Small', 5.0, 5, 1);
INSERT INTO shields (name, cost, weight, defense_modifier) VALUES ('Large', 10.0, 10, 2);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back
DELETE FROM shields WHERE name IN ('Small', 'Large');