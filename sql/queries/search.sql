-- name: SearchPostTitles :many
SELECT feeds.name AS feed_name, posts.title, posts.published_at FROM posts
INNER JOIN feeds on feeds.id = posts.feed_id
WHERE posts.title ILIKE $1
ORDER BY feeds.name ASC, posts.published_at DESC;

-- name: SearchFeeds :many
SELECT feeds.name, feeds.description
FROM feeds
WHERE (feeds.description ILIKE $1) OR (feeds.name ILIKE $1)
ORDER BY feeds.name ASC;