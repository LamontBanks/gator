-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html

-- Creates a Go function named CreateUser, with one row returned
-- Parameters:
--  id UUID
--  created_at  timestamp
--  updated_at  timestamp
--  name string
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- DEV/TESTING ONLY
-- name: Reset :exec
DELETE FROM USERS *;