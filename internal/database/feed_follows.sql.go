// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one

INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

// https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
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

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT feed_follows.id, feed_follows.created_at, feed_follows.updated_at, feed_follows.user_id, feed_follows.feed_id, feeds.name AS feed_name, users.name AS feed_creator_name FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserID          uuid.UUID
	FeedID          uuid.UUID
	FeedName        string
	FeedCreatorName string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.FeedName,
			&i.FeedCreatorName,
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
