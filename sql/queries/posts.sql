-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetPostsFromFollowedFeeds :many
WITH user_followed_feeds AS (
-- Get feeds you follow
    SELECT feeds.name AS feed_name, feeds.id AS feed_id
    FROM feed_follows
    INNER JOIN users
    ON feed_follows.user_id = users.id
    INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
    WHERE feed_follows.user_id = $1
)
-- Get all posts from those feeds
SELECT posts.title, posts.published_at, user_followed_feeds.feed_name
FROM posts
INNER JOIN user_followed_feeds
ON user_followed_feeds.feed_id = posts.feed_id
ORDER BY posts.published_at DESC
LIMIT $2;