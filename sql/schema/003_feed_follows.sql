-- Map users to the feeds they follow
-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    user_id uuid    NOT NULL
                    REFERENCES users
                    -- DELETES this if a user in `users` is deleted
                    ON DELETE CASCADE,
    feed_id uuid    NOT NULL 
                    REFERENCES feeds
                    -- DELETES this if a feed in `feeds` is deleted
                    ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
