-- +goose Up
ALTER TABLE feeds
ADD COLUMN description TEXT NOT NULL
                            DEFAULT '';

-- +goose Down
ALTER TABLE feeds
DROP COLUMN description;
