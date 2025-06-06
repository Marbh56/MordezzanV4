// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: treasure.sql

package db

import (
	"context"
	"database/sql"
)

const createTreasure = `-- name: CreateTreasure :execresult
INSERT INTO treasures (
  character_id, platinum_coins, gold_coins, electrum_coins,
  silver_coins, copper_coins, gems, art_objects, 
  other_valuables, total_value_gold
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type CreateTreasureParams struct {
	CharacterID    sql.NullInt64
	PlatinumCoins  int64
	GoldCoins      int64
	ElectrumCoins  int64
	SilverCoins    int64
	CopperCoins    int64
	Gems           sql.NullString
	ArtObjects     sql.NullString
	OtherValuables sql.NullString
	TotalValueGold float64
}

func (q *Queries) CreateTreasure(ctx context.Context, arg CreateTreasureParams) (sql.Result, error) {
	return q.exec(ctx, q.createTreasureStmt, createTreasure,
		arg.CharacterID,
		arg.PlatinumCoins,
		arg.GoldCoins,
		arg.ElectrumCoins,
		arg.SilverCoins,
		arg.CopperCoins,
		arg.Gems,
		arg.ArtObjects,
		arg.OtherValuables,
		arg.TotalValueGold,
	)
}

const deleteTreasure = `-- name: DeleteTreasure :execresult
DELETE FROM treasures
WHERE id = ?
`

func (q *Queries) DeleteTreasure(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteTreasureStmt, deleteTreasure, id)
}

const getTreasure = `-- name: GetTreasure :one
SELECT id, character_id, platinum_coins, gold_coins, electrum_coins, silver_coins, copper_coins, gems, art_objects, other_valuables, total_value_gold, created_at, updated_at FROM treasures
WHERE id = ? LIMIT 1
`

func (q *Queries) GetTreasure(ctx context.Context, id int64) (Treasure, error) {
	row := q.queryRow(ctx, q.getTreasureStmt, getTreasure, id)
	var i Treasure
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.PlatinumCoins,
		&i.GoldCoins,
		&i.ElectrumCoins,
		&i.SilverCoins,
		&i.CopperCoins,
		&i.Gems,
		&i.ArtObjects,
		&i.OtherValuables,
		&i.TotalValueGold,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTreasureByCharacter = `-- name: GetTreasureByCharacter :one
SELECT id, character_id, platinum_coins, gold_coins, electrum_coins, silver_coins, copper_coins, gems, art_objects, other_valuables, total_value_gold, created_at, updated_at FROM treasures
WHERE character_id = ? LIMIT 1
`

func (q *Queries) GetTreasureByCharacter(ctx context.Context, characterID sql.NullInt64) (Treasure, error) {
	row := q.queryRow(ctx, q.getTreasureByCharacterStmt, getTreasureByCharacter, characterID)
	var i Treasure
	err := row.Scan(
		&i.ID,
		&i.CharacterID,
		&i.PlatinumCoins,
		&i.GoldCoins,
		&i.ElectrumCoins,
		&i.SilverCoins,
		&i.CopperCoins,
		&i.Gems,
		&i.ArtObjects,
		&i.OtherValuables,
		&i.TotalValueGold,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTreasures = `-- name: ListTreasures :many
SELECT id, character_id, platinum_coins, gold_coins, electrum_coins, silver_coins, copper_coins, gems, art_objects, other_valuables, total_value_gold, created_at, updated_at FROM treasures
ORDER BY character_id
`

func (q *Queries) ListTreasures(ctx context.Context) ([]Treasure, error) {
	rows, err := q.query(ctx, q.listTreasuresStmt, listTreasures)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Treasure{}
	for rows.Next() {
		var i Treasure
		if err := rows.Scan(
			&i.ID,
			&i.CharacterID,
			&i.PlatinumCoins,
			&i.GoldCoins,
			&i.ElectrumCoins,
			&i.SilverCoins,
			&i.CopperCoins,
			&i.Gems,
			&i.ArtObjects,
			&i.OtherValuables,
			&i.TotalValueGold,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTreasure = `-- name: UpdateTreasure :execresult
UPDATE treasures
SET platinum_coins = ?,
    gold_coins = ?,
    electrum_coins = ?,
    silver_coins = ?,
    copper_coins = ?,
    gems = ?,
    art_objects = ?,
    other_valuables = ?,
    total_value_gold = ?,
    updated_at = datetime('now')
WHERE id = ?
`

type UpdateTreasureParams struct {
	PlatinumCoins  int64
	GoldCoins      int64
	ElectrumCoins  int64
	SilverCoins    int64
	CopperCoins    int64
	Gems           sql.NullString
	ArtObjects     sql.NullString
	OtherValuables sql.NullString
	TotalValueGold float64
	ID             int64
}

func (q *Queries) UpdateTreasure(ctx context.Context, arg UpdateTreasureParams) (sql.Result, error) {
	return q.exec(ctx, q.updateTreasureStmt, updateTreasure,
		arg.PlatinumCoins,
		arg.GoldCoins,
		arg.ElectrumCoins,
		arg.SilverCoins,
		arg.CopperCoins,
		arg.Gems,
		arg.ArtObjects,
		arg.OtherValuables,
		arg.TotalValueGold,
		arg.ID,
	)
}
