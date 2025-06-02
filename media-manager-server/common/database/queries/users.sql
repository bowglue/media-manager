-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :exec
INSERT INTO users (id, username)
VALUES (?, ?);