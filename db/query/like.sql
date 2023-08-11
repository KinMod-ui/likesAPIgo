-- name: CreateLike :one
INSERT INTO likes(
    user_id,
    content_id
) values (
$1, $2 
) RETURNING *;

-- name: GetLike :one
SELECT liked from likes
WHERE content_id = $1 AND user_id = $2; 

-- name: TotalLikesForContent :one
SELECT COUNT(*) AS Total_likes FROM likes
WHERE content_id = $1 AND liked= true;
