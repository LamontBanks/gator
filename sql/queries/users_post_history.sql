-- name: GetPostFromUserReadHistory :one
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

-- name: GetUnreadPostsForFeed :many
--- Get all posts for a given feed...
SELECT posts.id AS post_id, posts.feed_id, posts.title, posts.published_at, posts.description, posts.url
FROM posts
WHERE posts.feed_id = $2
AND
--- ...but only posts user has not read
posts.id NOT IN
	(SELECT post_id
		FROM users_posts_history
		WHERE users_posts_history.user_id = $1
        AND feed_id = $2
		AND has_viewed = true)
ORDER BY posts.published_at DESC;