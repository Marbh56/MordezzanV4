// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: spells.sql

package db

import (
	"context"
	"database/sql"
)

const createSpell = `-- name: CreateSpell :execresult
INSERT INTO spells (
  name, 
  mag_level, cry_level, ill_level, nec_level, 
  pyr_level, wch_level, clr_level, drd_level,
  range, duration, area_of_effect, components, description
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
`

type CreateSpellParams struct {
	Name         string
	MagLevel     int64
	CryLevel     int64
	IllLevel     int64
	NecLevel     int64
	PyrLevel     int64
	WchLevel     int64
	ClrLevel     int64
	DrdLevel     int64
	Range        string
	Duration     string
	AreaOfEffect sql.NullString
	Components   sql.NullString
	Description  string
}

func (q *Queries) CreateSpell(ctx context.Context, arg CreateSpellParams) (sql.Result, error) {
	return q.exec(ctx, q.createSpellStmt, createSpell,
		arg.Name,
		arg.MagLevel,
		arg.CryLevel,
		arg.IllLevel,
		arg.NecLevel,
		arg.PyrLevel,
		arg.WchLevel,
		arg.ClrLevel,
		arg.DrdLevel,
		arg.Range,
		arg.Duration,
		arg.AreaOfEffect,
		arg.Components,
		arg.Description,
	)
}

const deleteSpell = `-- name: DeleteSpell :execresult
DELETE FROM spells
WHERE id = ?
`

func (q *Queries) DeleteSpell(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteSpellStmt, deleteSpell, id)
}

const getSpell = `-- name: GetSpell :one
SELECT id, name, mag_level, cry_level, ill_level, nec_level, pyr_level, wch_level, clr_level, drd_level, "range", duration, area_of_effect, components, description, created_at, updated_at FROM spells
WHERE id = ? LIMIT 1
`

func (q *Queries) GetSpell(ctx context.Context, id int64) (Spell, error) {
	row := q.queryRow(ctx, q.getSpellStmt, getSpell, id)
	var i Spell
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MagLevel,
		&i.CryLevel,
		&i.IllLevel,
		&i.NecLevel,
		&i.PyrLevel,
		&i.WchLevel,
		&i.ClrLevel,
		&i.DrdLevel,
		&i.Range,
		&i.Duration,
		&i.AreaOfEffect,
		&i.Components,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listSpells = `-- name: ListSpells :many
SELECT id, name, mag_level, cry_level, ill_level, nec_level, pyr_level, wch_level, clr_level, drd_level, "range", duration, area_of_effect, components, description, created_at, updated_at FROM spells
ORDER BY name
`

func (q *Queries) ListSpells(ctx context.Context) ([]Spell, error) {
	rows, err := q.query(ctx, q.listSpellsStmt, listSpells)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Spell{}
	for rows.Next() {
		var i Spell
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MagLevel,
			&i.CryLevel,
			&i.IllLevel,
			&i.NecLevel,
			&i.PyrLevel,
			&i.WchLevel,
			&i.ClrLevel,
			&i.DrdLevel,
			&i.Range,
			&i.Duration,
			&i.AreaOfEffect,
			&i.Components,
			&i.Description,
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

const updateSpell = `-- name: UpdateSpell :execresult
UPDATE spells
SET name = ?,
    mag_level = ?,
    cry_level = ?,
    ill_level = ?,
    nec_level = ?,
    pyr_level = ?,
    wch_level = ?,
    clr_level = ?,
    drd_level = ?,
    range = ?,
    duration = ?,
    area_of_effect = ?,
    components = ?,
    description = ?,
    updated_at = datetime('now')
WHERE id = ?
`

type UpdateSpellParams struct {
	Name         string
	MagLevel     int64
	CryLevel     int64
	IllLevel     int64
	NecLevel     int64
	PyrLevel     int64
	WchLevel     int64
	ClrLevel     int64
	DrdLevel     int64
	Range        string
	Duration     string
	AreaOfEffect sql.NullString
	Components   sql.NullString
	Description  string
	ID           int64
}

func (q *Queries) UpdateSpell(ctx context.Context, arg UpdateSpellParams) (sql.Result, error) {
	return q.exec(ctx, q.updateSpellStmt, updateSpell,
		arg.Name,
		arg.MagLevel,
		arg.CryLevel,
		arg.IllLevel,
		arg.NecLevel,
		arg.PyrLevel,
		arg.WchLevel,
		arg.ClrLevel,
		arg.DrdLevel,
		arg.Range,
		arg.Duration,
		arg.AreaOfEffect,
		arg.Components,
		arg.Description,
		arg.ID,
	)
}
