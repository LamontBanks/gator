-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL DEFAULT 'Unspecified',
    url TEXT UNIQUE,
    -- PostgresSQL Foreign Keys: https://www.postgresql.org/docs/current/ddl-constraints.html#DDL-CONSTRAINTS-FK
    -- REFERENCE automatically uses the PRIMARY KEY of the referenced table
    user_id uuid REFERENCES users  
);

-- +goose Down
DROP TABLE feeds;
