-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- Common Table Expressions (CTE)
-- Use result of query as temporary table for other parts of the query
-- https://www.postgresql.org/docs/13/queries-with.html#QUERIES-WITH-CTE

-- Add entry showing user follwed a feed, then returns the name of the user and feed
-- name: CreateFeedFollow :one
WITH new_feed_follow_row AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT new_feed_follow_row.*, users.name AS user_name, feeds.name AS feed_name
FROM new_feed_follow_row
INNER JOIN users ON new_feed_follow_row.user_id = users.id
INNER JOIN feeds ON new_feed_follow_row.feed_id = feeds.id;

-- Get all feeds the user is following
-- name: GetFeedFollowsForUser :many
SELECT feed_follows.feed_id, feeds.description, feeds.name AS feed_name, feeds.url as feed_url 
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY feeds.name ASC;

-- Get all feeds the user is NOT following
-- name: GetFeedsNotFollowedByUser :many
SELECT feeds.name, feeds.description, feeds.id, feeds.url
FROM feeds
WHERE feeds.id NOT IN (
	SELECT feed_follows.feed_id
	FROM feed_follows
	WHERE feed_follows.user_id = $1
)
ORDER BY feeds.name ASC;

-- name: DeleteFeedFollowForUser :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;