-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL DEFAULT 'Unspecified',
    url TEXT NOT NULL DEFAULT 'Unspecified',
    -- PostgresSQL Foreign Keys: https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-FK
    -- REFERENCE automatically uses the PRIMARY KEY of the referenced table
    user_id uuid NOT NULL REFERENCES users,
    -- UNIQUE CONSTRAINT for url and user_id:
    -- a user can only add the same feed once, but multiple users can have the same feed
    -- https://www.postgresql.org/docs/17/ddl-constraints.html#DDL-CONSTRAINTS-UNIQUE-CONSTRAINTS
    UNIQUE (url, user_id)
);

-- +goose Down
DROP TABLE feeds;
