-- name: GetUser :one
SELECT * FROM users WHERE 
id=$1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY username 
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (username)
VALUES ($1) RETURNING *;
