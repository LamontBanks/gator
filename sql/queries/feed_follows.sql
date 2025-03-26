-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- name: CreateFeedFollow :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, feeds.name AS feed_name, users.name AS feed_creator_name FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE feed_follows.user_id = $1;