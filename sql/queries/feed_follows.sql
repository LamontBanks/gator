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
INNER JOIN users
ON new_feed_follow_row.user_id = users.id
INNER JOIN feeds
ON new_feed_follow_row.feed_id = feeds.id;

-- Get all feeds the user is following
-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name AS feed_creator_name 
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE feed_follows.user_id = $1;
