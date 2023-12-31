// Code generated by sqlc. DO NOT EDIT.
// source: notification.sql

package db

import (
	"context"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications(
    user_id 
) VALUES (
  $1
) RETURNING user_id
`

func (q *Queries) CreateNotification(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, createNotification, userID)
	var user_id int64
	err := row.Scan(&user_id)
	return user_id, err
}


const selectUsersTN= `
SELECT user_likes.user_id FROM (
    SELECT user_id, COUNT(*) AS like_count
    FROM likes
    WHERE liked = true
    GROUP BY user_id
) AS user_likes
WHERE user_likes.like_count >= 2 
AND user_id NOT IN (
    SELECT DISTINCT user_id FROM notifications
);
`

func (q *Queries) SelectUsersToNotify(ctx context.Context) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, selectUsersTN) 
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var i int64 
		if err := rows.Scan(
			&i,
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
