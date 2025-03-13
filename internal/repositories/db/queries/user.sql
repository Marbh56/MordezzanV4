-- name: GetUser :one
SELECT id, username, email, created_at, updated_at FROM users
WHERE id = ? LIMIT 1;

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