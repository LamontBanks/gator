# Blog Aggregator

A CLI tool to aggregate RSS feeds:

- Read RSS feeds
- Supports multple users, each with custom RSS feed subscriptions

## Setup

### Database Setup

This project requires a PostgreSQL database with the correct schema.

#### Install tools

1. Install the Goose database migration tool:
        
        $ go install github.com/pressly/goose/v3/cmd/goose@latest

2. Install PostgreSQL:

- Mac:

        $ brew install postgresql@15

- Linux / WSL (Debian), see [Microsoft Docs](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql) for more details:

        $ sudo apt update
        $ sudo apt install postgresql postgresql-contrib

3.  Ensure `psql` is installed with version 15+ (should come with the Postgres installation):

        $ psql --version
        psql (PostgreSQL) 15.12 (Homebrew)

- (Linux only) Update postgres password:

        $ sudo passwd postgres

#### Configure Postgres

1. Start Postgres in the background

- Mac:
        
        $ brew services start postgresql@15

- Linux:
        
        $ sudo service postgresql start

2. Create a database named `gator`:

        postgres=# CREATE DATABASE gator;

- Set the database user password (Linux only)

        postgres=# ALTER USER postgres PASSWORD '<postgres>';

#### Setup the tables

1. Back in the project, in the `sql/schema` directory, run [goose up](https://github.com/pressly/goose?tab=readme-ov-file#up) to create the tables:

        $ cd sql/schema
        $ goose postgres postgres://<postgres username>:<postgres password>@localhost:5432/gator up

        2025/04/04 10:53:17 goose: successfully migrated database to version: 6

2. Verify the schema in Postgres:

        $ psql postgres://<postgres username>:<postgres password>@localhost:5432/gator

        $ gator=# \dt   // Describe tables

        List of relations
        Schema |       Name       | Type  | Owner 
        --------+------------------+-------+-------
        public | feed_follows     | table | lb
        public | feeds            | table | lb
        public | goose_db_version | table | lb
        public | posts            | table | lb
        public | users            | table | lb

## Usage

### Config File (one-time step)

1. Create a `.gatorconfig.json` file in your `$HOME` directory:

        .gatorconfig.json

        {"db_url":"","current_user_name":""}

2. Add your database connection string; Leave `current_user_name` blank - `gator` will update this field.
Append `?sslmode=disable` (local use only) to connection string:

        .gatorconfig.json

        {"db_url":"postgres://<pg username>:<password>@localhost:5432/gator?sslmode=disable","current_user_name":""}

### Run gator

Start `gator` as a background process, providing an update frequency, `(30s, 1m, 5, 1h, 24h, etc.)`:
        
        $ gator agg 15m &       # Update RSS feeds every 15 minutes
        Updating RSS feeds...

### Register/login a user
A user must be created and logged in, with 
Register a new user:

        $ gator register Bob
        Registered user Bob
        Logged in as Bob

Or, login an existing user:

        $ gator login Alice
        Logged in as Alice

### List users

`$ gator users`

        * Alice (current)
        * Bob


### Add RSS feeds
 `$ gator addFeed <url>`
 
 Newly added feeds will immediately download updates:

        $ gator addFeed "Nasa Image of the Day" https://www.nasa.gov/feeds/iotd-feed/
        Saved "Nasa Image of the Day" (https://www.nasa.gov/feeds/iotd-feed/) for Alice

        $ gator addFeed "Dwarf Fortress" https://store.steampowered.com/feeds/news/app/975370
        Saved "Dwarf Fortress" (https://store.steampowered.com/feeds/news/app/975370) for Alice
        Alice followed Dwarf Fortress

### See RSS Feeds overview
`$ gator browse`:

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
                - Hotfix: Patch 51.10
                - Update 51.09: Forbidden Coffins now unusable
                - Hotfix: Patch 51.08
        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Artemis II Insignia Honors All
                - X-ray Clues Reveal Destroyed Planet
                - Studying Ice for the Future of Flight

Provide a number to change the number of posts per feed:

`$ gator browse 10`:

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
                - Hotfix: Patch 51.10
                - Update 51.09: Forbidden Coffins now unusable
                - Hotfix: Patch 51.08
                - Update 51.07: Easier reconstruction and pet protection
                - Lua Beta + Patch 51.06 ⛏ Dwarf Fortress Dev News
                - Next Steps for Dwarf Fortress + Patch 51.05 ⛏ Dwarf Fortress Dev News
                - Hotfix: Patch 51.04
                - Hotfix: Patch 51.03
                - Adventure Mode is Out NOW!
                - Adventure Mode Beta Patch Notes (All of them)
        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Artemis II Insignia Honors All
                - X-ray Clues Reveal Destroyed Planet
                - Studying Ice for the Future of Flight
                - Artemis II Core Stage Integration – Complete!
                - Turning Vanes inside the Altitude Wind Tunnel
                - NEO Surveyor Instrument Enclosure Inside Historic Chamber A
                - Norman Rockwell Commemorates Gemini Program with Grissom and Young
                - NASA’s Spirit Rover Gets Looked Over
                - Like Sands Through the Hourglass…
                - Making Ripples

### Read posts
1. `$ gator browseFeed` to navigate into your feeds:

        Choose a feed:
        1: Dwarf Fortress
        2: Nasa Image of the Day

1. Enter the corresponding number to display recent posts:

        2       // user input

        Nasa Image of the Day
        Choose a post:
        1: Artemis II Insignia Honors All
                01:55 PM, Thu, 03 Apr 25
        2: X-ray Clues Reveal Destroyed Planet
                01:10 PM, Wed, 02 Apr 25
        3: Studying Ice for the Future of Flight
                12:52 PM, Tue, 01 Apr 25
        4: Artemis II Core Stage Integration – Complete!
                02:42 PM, Mon, 31 Mar 25
        ...

1. Enter the number corresponding to read the post:

        3       // user input

        Studying Ice for the Future of Flight
        12:52 PM, Tue, 01 Apr 25

        Thomas Ozoroski, a researcher at NASA’s Glenn Research Center in Cleveland, takes icing accretion measurements in October 2024 as part of transonic truss-braced wing concept research. Researchers at NASA Glenn conducted another test campaign in March 2025.

        https://www.nasa.gov/image-detail/grc-2024-c-12100-2/

### List feeds being followed
`$ gator following`

        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/


### List all available feeds
`$ gator feeds`

        All RSS Feeds
        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

### Follow additional feeds
Users can follow RSS feeds that others have added. 

For example, Bob is currently not following any feeds:

        $ gator login Bob
        Logged in as Bob

        $ gator following
        you are not following any feeds


1. `$ gator follow` to select additional RSS feeds to follow:

        Already following:
        no feeds

        - Choose a new RSS feed to follow:
        1: Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        2: Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

1. Enter the number of the desired feed:

        1       // user input

        Bob followed Dwarf Fortress (https://store.steampowered.com/feeds/news/app/975370)

1. View updates for Bob:

        $ gator browse

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
                - Hotfix: Patch 51.10
                - Update 51.09: Forbidden Coffins now unusable
                - Hotfix: Patch 51.08

Saved feeds can also be followed directly by providing the feed url:

`$ gator follow https://www.nasa.gov/feeds/iotd-feed/`

### Unfollow feeds

`$ gator unfollow` to remove a feed you're following:

        $ gator unfollow

        - Choose an RSS feed to unfollow
        1: Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        2: Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

Enter the number of feed:

        1 (user input)

        Unfollowed Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370

View current feeds:

        $ gator following

        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

### Help

Show all commands:

        $ gator help

        gator is a tool for viewing RSS feeds in the console.

        Usage:

                gator <command> [arguments]

                addFeed		Add a new feed for all users to follow
                agg		Aggregate all feeds, poll for updates, useful when run in the background with '*'
                browse		Show latest posts for current user's feeds
                browseFeed		Read posts from a followed feed
                feeds		List all feeds
                ...


Show help for specific command:

        $ gator help browse

        usage: gator browse <max number of posts per feed>

        Show latest posts for current user's feeds

        Examples:
                gator browse
                gator browse 5




## Development

### SQLC

Install [`sqlc`](https://sqlc.dev/), a useful library for generating boilerplate Go code from SQL queries:

        $ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

Find SQLC's config file in the base project directory:

        $ gator/sqlc.yaml

### Debugging

Visual Studio Code debugger is useful - a plain/standard setup should be sufficient. Here is a sample `launch.json` config:

```json
launch.json:

{
        "name": "addFeed",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "/<path>/gator/.",
        "args": "addFeed \"Nasa Image of the Day\" https://www.nasa.gov/feeds/iotd-feed/",
        "console": "integratedTerminal"
}
```

### Run locally
`$ go run . <command>`

Examples:
- `$ go run . agg 15m &`
- `$ go run . login Alice`

### Test data

A shell script can add users and feeds, see `sampleData.sh`.

### `gator reset`

`gator` has a `reset` command that deletes all users from the database. Since all the the tables have `CASCADE`-ing deletes on the `users.id` field, this effectively clears the entire database (see: https://www.postgresql.org/docs/15/sql-createtable.html, search `CASCADE`).

Only use `reset` for development/local-use - it should (obviously) be commented-out/removed otherwise.

### Adding a new command

#### Modifying the database (optional)
1. If needed, create new `goose` up/down migrations (new tables, `ALTER` tables, etc.) in `sql/schema`

#### New/modify SQL statements (optional)
1. Add/edit `.sql` files in `sql/queries`. Reference existing queries for the needed `SQLC` syntax (i.e., `-- name: CreateFeed :one`), etc.
2. Run `$ sqlc generate` from the base project directory to generate Go code. Files within `internal/database` will be generated, including a new function based on your SQL statement. For example, a query named `CreateFeed` will generate a function: `func CreateFeed(...) ...`

#### Create a new command Go file
1. Command files are named "command_\<name of command\>.go", placed at the base project directory.

#### Handler and info functions
2. Each command needs 2 functions:

- A handler function to read user args, read/write the database, and print output:
                
```go
func (s *state, cmd command) error { ... }
```

- A "command help" function to print usage info. The gator `help` command will call this to display info in a CLI-styled help output.

```go
func () commandInfo { ... }
```

The `handlerGetFeeds()` and `feedsCommandInfo()` functions (mapped to `$ gator feeds`) are good references.

#### User-authenticated commands
Commands that require a user id (ex: `gator browse`) should use the handler function signature:

```go
func (s *state, cmd command, user database.User) error { ... }
```

Then, in `main.go` the handler function will be wrapped in `middlewareLoggedIn(...)` closure. The closure will provide the currently logged-in user to the handler function. 

Commands should *not* need to duplicate manually retrieving the current user from the database.

The `handlerAddFeed()` function (mapped to `$ gator addFeed`) is a good reference of a user-authenticated command. Also, see `middlewareLoggedIn()` in "main.go".

#### Enabling the command

Once the handler and info functions are created, they need to be added to `main.go`.

There are many comands already:

```go
// Register the CLI commands
appCommands := commands{
        cmds: make(map[string]commandDetails),
}
...
appCommands.register("feeds", feedsCommandInfo(), handlerGetFeeds)
appCommands.register("addFeed", addFeedCommandInfo(), middlewareLoggedIn(handlerAddFeed))
...
```

Simply copy-paste the `register` line, adding your command's name, info function *call*, and handler function. Commands that require a user should be wrapped in `middlewareLoggedIn(...)` - see the section "User-authenticated commands":

```go
...
appCommands.register("myCommand", myCommandInfo(), middlewareLoggedIn(handlerMyCommand))
...
```

### Manually test
From the base directory, run:

        $ go run . agg 15m &
        $ go run . addFeed...
        $ go run . etc...

### Build the program
From the base directory:

        $ go build .

Use the program:

        $ ./gator agg 15m &
        [1] 4015
        Updating RSS feeds...

        $ ./gator browse
        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Artemis II Insignia Honors All
                - X-ray Clues Reveal Destroyed Planet
                - Studying Ice for the Future of Flight

## Improvements
- Tests, unit/integration
- Process HTML info RSS feed posts
- Retrieve full blog posts from feed URL
- Formal handling of CLI args
- Improve feed/post navigation/viewing
        - Keep CLI output minimum, but still display sufficient output at a glance
