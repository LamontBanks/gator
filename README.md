# Blog Aggregator

A CLI tool to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Tools
- PostgresQL - Data storage

        $ brew install postgresql@15

- Goose: Database migrations

        $ go install github.com/pressly/goose/v3/cmd/goose@latest

- SQLC - Generate Go code from SQL queries
    
        $ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest