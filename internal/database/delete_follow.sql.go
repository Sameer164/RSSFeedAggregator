// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: delete_follow.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const unfollow = `-- name: Unfollow :one
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2
RETURNING id, created_at, updated_at, user_id, feed_id
`

type UnfollowParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) Unfollow(ctx context.Context, arg UnfollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, unfollow, arg.UserID, arg.FeedID)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}
