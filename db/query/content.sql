-- name: CreateContent :one
INSERT INTO content(
    title,
    user_id 
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetContent :one
SELECT * FROM content 
WHERE id = $1 LIMIT 1;

-- name: ListContentOfUser :many
SELECT * FROM content 
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

