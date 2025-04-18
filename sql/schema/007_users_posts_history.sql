-- +goose Up
CREATE TABLE users_posts_history (
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
    post_id uuid    NOT NULL
                    REFERENCES posts
                    -- DELETES this row if a posts in `posts` is deleted
                    ON DELETE CASCADE,
    has_viewed      boolean     NOT NULL
                                DEFAULT false,
    is_bookmarked   boolean     NOT NULL
                                DEFAULT false,
    UNIQUE (user_id, feed_id, post_id)
);

-- +goose Down
DROP TABLE users_posts_history;