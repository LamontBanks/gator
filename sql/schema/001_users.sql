-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE users (
    id uuid                 PRIMARY KEY,
    created_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL 
                            DEFAULT CURRENT_TIMESTAMP,
    name TEXT               UNIQUE 
                            NOT NULL DEFAULT 'Unspecified'
);

-- +goose Down
DROP TABLE users;
