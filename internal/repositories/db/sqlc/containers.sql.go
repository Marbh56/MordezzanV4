// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: containers.sql

package db

import (
	"context"
	"database/sql"
)

const createContainer = `-- name: CreateContainer :execresult
INSERT INTO containers (
  name, max_weight, allowed_items, cost
) VALUES (
  ?, ?, ?, ?
)
`

type CreateContainerParams struct {
	Name         string
	MaxWeight    int64
	AllowedItems string
	Cost         float64
}

func (q *Queries) CreateContainer(ctx context.Context, arg CreateContainerParams) (sql.Result, error) {
	return q.exec(ctx, q.createContainerStmt, createContainer,
		arg.Name,
		arg.MaxWeight,
		arg.AllowedItems,
		arg.Cost,
	)
}

const deleteContainer = `-- name: DeleteContainer :execresult
DELETE FROM containers
WHERE id = ?
`

func (q *Queries) DeleteContainer(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteContainerStmt, deleteContainer, id)
}

const getContainer = `-- name: GetContainer :one
SELECT id, name, max_weight, allowed_items, cost, created_at, updated_at FROM containers
WHERE id = ? LIMIT 1
`

func (q *Queries) GetContainer(ctx context.Context, id int64) (Container, error) {
	row := q.queryRow(ctx, q.getContainerStmt, getContainer, id)
	var i Container
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MaxWeight,
		&i.AllowedItems,
		&i.Cost,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getContainerByName = `-- name: GetContainerByName :one
SELECT id, name, max_weight, allowed_items, cost, created_at, updated_at FROM containers
WHERE name = ? LIMIT 1
`

func (q *Queries) GetContainerByName(ctx context.Context, name string) (Container, error) {
	row := q.queryRow(ctx, q.getContainerByNameStmt, getContainerByName, name)
	var i Container
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MaxWeight,
		&i.AllowedItems,
		&i.Cost,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listContainers = `-- name: ListContainers :many
SELECT id, name, max_weight, allowed_items, cost, created_at, updated_at FROM containers
ORDER BY name
`

func (q *Queries) ListContainers(ctx context.Context) ([]Container, error) {
	rows, err := q.query(ctx, q.listContainersStmt, listContainers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Container{}
	for rows.Next() {
		var i Container
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MaxWeight,
			&i.AllowedItems,
			&i.Cost,
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

const updateContainer = `-- name: UpdateContainer :execresult
UPDATE containers
SET name = ?,
    max_weight = ?,
    allowed_items = ?,
    cost = ?,
    updated_at = datetime('now')
WHERE id = ?
`

type UpdateContainerParams struct {
	Name         string
	MaxWeight    int64
	AllowedItems string
	Cost         float64
	ID           int64
}

func (q *Queries) UpdateContainer(ctx context.Context, arg UpdateContainerParams) (sql.Result, error) {
	return q.exec(ctx, q.updateContainerStmt, updateContainer,
		arg.Name,
		arg.MaxWeight,
		arg.AllowedItems,
		arg.Cost,
		arg.ID,
	)
}
