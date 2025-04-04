// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :execresult
INSERT INTO users (
  username, email, password_hash
) VALUES (
  ?, ?, ?
)
`

type CreateUserParams struct {
	Username     string
	Email        string
	PasswordHash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.exec(ctx, q.createUserStmt, createUser, arg.Username, arg.Email, arg.PasswordHash)
}

const deleteUser = `-- name: DeleteUser :execresult
DELETE FROM users
WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) (sql.Result, error) {
	return q.exec(ctx, q.deleteUserStmt, deleteUser, id)
}

const getFullUserByEmail = `-- name: GetFullUserByEmail :one
SELECT id, username, email, password_hash, created_at, updated_at FROM users
WHERE email = ? LIMIT 1
`

func (q *Queries) GetFullUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.getFullUserByEmailStmt, getFullUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, email, created_at, updated_at FROM users
WHERE id = ? LIMIT 1
`

type GetUserRow struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) GetUser(ctx context.Context, id int64) (GetUserRow, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, id)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at FROM users
ORDER BY username
`

type ListUsersRow struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.query(ctx, q.listUsersStmt, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUsersRow{}
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
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

const updateUser = `-- name: UpdateUser :execresult
UPDATE users
SET username = ?, email = ?, updated_at = datetime('now')
WHERE id = ?
`

type UpdateUserParams struct {
	Username string
	Email    string
	ID       int64
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (sql.Result, error) {
	return q.exec(ctx, q.updateUserStmt, updateUser, arg.Username, arg.Email, arg.ID)
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

type UpdateUserPasswordParams struct {
	PasswordHash string
	ID           int64
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.exec(ctx, q.updateUserPasswordStmt, updateUserPassword, arg.PasswordHash, arg.ID)
	return err
}
