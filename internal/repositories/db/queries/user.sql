-- name: GetUser :one
SELECT id, username, email, created_at, updated_at FROM users
WHERE id = ? LIMIT 1;

-- name: GetFullUserByEmail :one
SELECT id, username, email, password_hash, created_at, updated_at FROM users
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at FROM users
ORDER BY username;

-- name: CreateUser :execresult
INSERT INTO users (
  username, email, password_hash
) VALUES (
  ?, ?, ?
);

-- name: UpdateUser :execresult
UPDATE users
SET username = ?, email = ?, updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE id = ?;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
