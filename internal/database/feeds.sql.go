// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one

INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at, description)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at, description
`

type CreateFeedParams struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	Description   string
}

// https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html
// Represent a feed and the user who added the feed
// Nulling `last_fetched_at` timestamp field - that will be set when the feed is downloaded
func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
		arg.LastFetchedAt,
		arg.Description,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
		&i.Description,
	)
	return i, err
}

const deleteFeedById = `-- name: DeleteFeedById :exec
DELETE FROM feeds
WHERE feeds.id = $1
`

func (q *Queries) DeleteFeedById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeedById, id)
	return err
}

const getFeedByUrl = `-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at, description FROM feeds
WHERE url = $1
`

func (q *Queries) GetFeedByUrl(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUrl, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
		&i.Description,
	)
	return i, err
}

const getFeedFollowerCount = `-- name: GetFeedFollowerCount :many
SELECT feed_id, feeds.name as feed_name, COUNT(*) as num_followers
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
GROUP BY feed_id, feeds.name
ORDER BY num_followers DESC
`

type GetFeedFollowerCountRow struct {
	FeedID       uuid.UUID
	FeedName     string
	NumFollowers int64
}

func (q *Queries) GetFeedFollowerCount(ctx context.Context) ([]GetFeedFollowerCountRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowerCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowerCountRow
	for rows.Next() {
		var i GetFeedFollowerCountRow
		if err := rows.Scan(&i.FeedID, &i.FeedName, &i.NumFollowers); err != nil {
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

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.id, feeds.name AS feed_name, feeds.url, feeds.description, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
ORDER BY feeds.name ASC
`

type GetFeedsRow struct {
	ID          uuid.UUID
	FeedName    string
	Url         string
	Description string
	UserName    sql.NullString
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.FeedName,
			&i.Url,
			&i.Description,
			&i.UserName,
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

const getFeedsEligibleForDeletion = `-- name: GetFeedsEligibleForDeletion :many
WITH feeds_created_by_user AS (
    SELECT feeds.id as feed_id, feeds.name AS feed_name, feeds.user_id
    FROM feeds
    INNER JOIN users ON feeds.user_id = users.id
    WHERE feeds.user_id = $1::uuid
),
num_followers_per_feed AS (
    SELECT feed_follows.feed_id, COUNT(*) AS num_followers
    FROM feed_follows
    WHERE feed_follows.feed_id IN (SELECT feeds_created_by_user.feed_id FROM feeds_created_by_user)
    GROUP BY feed_follows.feed_id
)
SELECT feeds.name, num_followers_per_feed.feed_id, num_followers_per_feed.num_followers
    FROM num_followers_per_feed
    INNER JOIN feeds ON feeds.id = num_followers_per_feed.feed_id
    WHERE 
        (num_followers_per_feed.num_followers = 0)
        OR
        (num_followers_per_feed.num_followers = 1 AND $1::uuid IN (SELECT feed_follows.user_id FROM feed_follows) )
    ORDER BY feeds.name
`

type GetFeedsEligibleForDeletionRow struct {
	Name         string
	FeedID       uuid.UUID
	NumFollowers int64
}

// 1. Get feeds created by a given user
// 2. Get feed follower counts for feed created by user
// 3. Select feeds with either: 0 followers, or the user is the only follower
func (q *Queries) GetFeedsEligibleForDeletion(ctx context.Context, userID uuid.UUID) ([]GetFeedsEligibleForDeletionRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsEligibleForDeletion, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsEligibleForDeletionRow
	for rows.Next() {
		var i GetFeedsEligibleForDeletionRow
		if err := rows.Scan(&i.Name, &i.FeedID, &i.NumFollowers); err != nil {
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

const markFeedAsFetched = `-- name: MarkFeedAsFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
`

func (q *Queries) MarkFeedAsFetched(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markFeedAsFetched, id)
	return err
}

const updateFeedDescription = `-- name: UpdateFeedDescription :exec
UPDATE feeds
SET description = $2
WHERE id = $1
`

type UpdateFeedDescriptionParams struct {
	ID          uuid.UUID
	Description string
}

func (q *Queries) UpdateFeedDescription(ctx context.Context, arg UpdateFeedDescriptionParams) error {
	_, err := q.db.ExecContext(ctx, updateFeedDescription, arg.ID, arg.Description)
	return err
}
