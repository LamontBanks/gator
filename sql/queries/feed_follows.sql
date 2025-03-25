-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- name: CreateFeedFollow :one
INSERT INTO users (id, created_at, updated_at, user_id, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;