-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (id, username, age_restriction, pin_hash, role)
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET username = ?, age_restriction = ?, pin_hash = ?, role = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? RETURNING *;

-- name: DoesUsernameExist :one
SELECT EXISTS (
    SELECT 1 FROM users WHERE username = ?
);