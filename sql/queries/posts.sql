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

-- name: GetFollowedPosts :many
WITH user_followed_feeds AS (
-- Get feeds being followed...
    SELECT feeds.name AS feed_name, feeds.id AS feed_id
    FROM feed_follows
    INNER JOIN users ON feed_follows.user_id = users.id
    INNER JOIN feeds ON feed_follows.feed_id = feeds.id
    WHERE feed_follows.user_id = $1
)
-- ...get all posts from those feeds
SELECT posts.*, user_followed_feeds.feed_name
FROM posts
INNER JOIN user_followed_feeds ON user_followed_feeds.feed_id = posts.feed_id
ORDER BY posts.published_at DESC
LIMIT $2;

-- name: GetRecentPostsFromFeed :many
SELECT feeds.name AS feed_name, posts.title, posts.description, posts.published_at, posts.Url
FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
WHERE posts.feed_id = $1 AND posts.published_at >= $2 
ORDER BY posts.published_at DESC
LIMIT $3;

-- name: GetPostsFromFeed :many
SELECT feeds.name AS feed_name, posts.title, posts.description, posts.published_at, posts.Url
FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
WHERE posts.feed_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;