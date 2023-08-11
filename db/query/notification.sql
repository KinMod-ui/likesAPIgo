-- name: CreateNotification :one
INSERT INTO notifications(
    user_id 
) VALUES (
  $1
) RETURNING *;

-- name : SelectUsersToNotify :many
SELECT user_likes.user_id FROM (
    SELECT user_id, COUNT(*) AS like_count
    FROM likes
    WHERE liked = true
    GROUP BY user_id
) AS user_likes
WHERE user_likes.like_count >= 10 
AND user_id NOT IN (
    SELECT DISTINCT user_id FROM notifications
);
