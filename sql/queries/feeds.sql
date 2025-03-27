-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- Represent a feed and the user who added the feed
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
LEFT JOIN users
ON feeds.user_id = users.id
ORDER BY feeds.name ASC;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Gets oldest feed to update
-- Feeds that have never been updated (last_fetched_at = NULL) are prioritized (though no specific order can be gauranteed)
-- name: GetNextFeedToFetch :one
SELECT id, name, url, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC
NULLS FIRST;
