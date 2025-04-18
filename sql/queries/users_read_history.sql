-- name: GetPostFromUserReadHisory :one
SELECT * FROM users_posts_history
WHERE user_id = $1 
        AND feed_id = $2
        AND post_id = $3;

-- name: CreatePostInUserReadHistory :one
INSERT INTO users_posts_history (id, created_at, updated_at, user_id, feed_id, post_id, has_viewed, is_bookmarked)
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
RETURNING *;

-- name: MarkPostAsViewed :exec
UPDATE users_posts_history
SET has_viewed = true
WHERE user_id = $1 
    AND feed_id = $2
    AND post_id = $3;