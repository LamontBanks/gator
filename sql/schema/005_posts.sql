-- Goose for database migrations: https://github.com/pressly/goose

-- Save RSS Feed posts for browsing
-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    title TEXT              NOT NULL
                            DEFAULT 'No Title',
    url TEXT                UNIQUE,
    description TEXT,
    published_at timestamp  NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    feed_id uuid            NOT NULL
                            REFERENCES feeds
                            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;