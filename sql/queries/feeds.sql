-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- Represent a feed and the user who added the feed
-- Nulling `last_fetched_at` timestamp field - that will be set when the feed is downloaded
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at, description)
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

-- name: GetFeeds :many
SELECT feeds.id, feeds.name AS feed_name, feeds.url, feeds.description, users.name AS user_name
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id
ORDER BY feeds.name ASC;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateFeedDescription :exec
UPDATE feeds
SET description = $2
WHERE id = $1;

-- name: GetFeedsEligibleForDeletion :many
-- 1. Get feeds created by a given user
-- 2. Get feed follower counts for feed created by user
-- 3. Select feeds with either: 0 followers, or the user is the only follower
WITH feeds_created_by_user AS (
    SELECT feeds.id as feed_id, feeds.name AS feed_name, feeds.user_id
    FROM feeds
    INNER JOIN users ON feeds.user_id = users.id
    WHERE feeds.user_id = sqlc.arg(user_id)::uuid
),
num_followers_per_feed AS (
    SELECT feed_follows.feed_id, COUNT(*) AS num_followers
    FROM feed_follows
    WHERE feed_follows.feed_id IN (SELECT feeds_created_by_user.feed_id FROM feeds_created_by_user)
    GROUP BY feed_follows.feed_id
)
SELECT feeds.name, num_followers_per_feed.feed_id, num_followers_per_feed.num_followers
    FROM num_followers_per_feed
    INNER JOIN feeds ON feeds.id = num_followers_per_feed.feed_id
    WHERE 
        (num_followers_per_feed.num_followers = 0)
        OR
        (num_followers_per_feed.num_followers = 1 AND sqlc.arg(user_id)::uuid IN (SELECT feed_follows.user_id FROM feed_follows) )
    ORDER BY feeds.name;

-- name: DeleteFeedById :exec
DELETE FROM feeds
WHERE feeds.id = $1;