// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: search.sql

package database

import (
	"context"
	"time"
)

const searchFeeds = `-- name: SearchFeeds :many
SELECT feeds.name, feeds.description
FROM feeds
WHERE (feeds.description ILIKE $1) OR (feeds.name ILIKE $1)
ORDER BY feeds.name ASC
`

type SearchFeedsRow struct {
	Name        string
	Description string
}

func (q *Queries) SearchFeeds(ctx context.Context, description string) ([]SearchFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, searchFeeds, description)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchFeedsRow
	for rows.Next() {
		var i SearchFeedsRow
		if err := rows.Scan(&i.Name, &i.Description); err != nil {
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

const searchPostTitles = `-- name: SearchPostTitles :many
SELECT feeds.name AS feed_name, posts.title, posts.published_at FROM posts
INNER JOIN feeds on feeds.id = posts.feed_id
WHERE posts.title ILIKE $1
ORDER BY feeds.name ASC, posts.published_at DESC
`

type SearchPostTitlesRow struct {
	FeedName    string
	Title       string
	PublishedAt time.Time
}

func (q *Queries) SearchPostTitles(ctx context.Context, title string) ([]SearchPostTitlesRow, error) {
	rows, err := q.db.QueryContext(ctx, searchPostTitles, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchPostTitlesRow
	for rows.Next() {
		var i SearchPostTitlesRow
		if err := rows.Scan(&i.FeedName, &i.Title, &i.PublishedAt); err != nil {
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
