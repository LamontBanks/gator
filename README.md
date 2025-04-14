# Blog Aggregator

A CLI tool for read RSS feeds in the terminal, built using the [cobra-cli](https://github.com/spf13/cobra) CLI library.

- Read RSS feeds
- Supports multple users, each with their own custom RSS feed subscriptions

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

Start `gator`, preferably as a background process, providing an update frequency: `(30s, 1m, 5, 1h, 24h, etc.)`:
        
        $ gator update 15m &       # Update RSS feeds every 15 minutes
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
 `$ gator add -n <custom name of feed> -u <url>`
 
        # The user will automatically follow feeds they add
        Added RSS feed "Nasa Image of the Day" (https://www.nasa.gov/feeds/iotd-feed/)
        Following "Nasa Image of the Day" (https://www.nasa.gov/feeds/iotd-feed/)

### See followed RSS Feeds overview

`$ gator`

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
                - Hotfix: Patch 51.10
                - Update 51.09: Forbidden Coffins now unusable

        Guild Wars 2 | https://www.guildwars2.com/en/feed/
                - Snack on the Job with the Sweet Treat Gathering Tools
                - Stop and Smell the Flowers with the Wisteria Arborscale Skyscale Skin

        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Apollo 13 Launch: 55 Years Ago
                - Linear Sand Dunes in the Great Sandy Desert

        Pivot To AI | https://pivot-to-ai.com/feed/
                - New York mayoral candidate Andrew Cuomo writes his housing plan with ChatGPT
                - How to sell AI slop to the US military — with Vannevar Labs and MIT Tech Review

Provide a number to change the number of posts per feed:

`$ gator -n 10`:

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
                - Hotfix: Patch 51.10
                - Update 51.09: Forbidden Coffins now unusable
                - Hotfix: Patch 51.08
                - Update 51.07: Easier reconstruction and pet protection
                - Lua Beta + Patch 51.06 ⛏ Dwarf Fortress Dev News
                - Next Steps for Dwarf Fortress + Patch 51.05 ⛏ Dwarf Fortress Dev News
                - Hotfix: Patch 51.04

        Guild Wars 2 | https://www.guildwars2.com/en/feed/
                - Snack on the Job with the Sweet Treat Gathering Tools
                - Stop and Smell the Flowers with the Wisteria Arborscale Skyscale Skin
                - Our New Avian Aspect Helm Skin Is a Real Feather in Your Cap
                - Shine on the Tyrian Runway This Week with Our New Emote Tome!
                - Spring into Adventure with the Guild Wars 2 Spring Sale!
                - “Repentance” Is Now Live
                - Ferocious New Mounts and an Otherworldly Black Lion Chest Update!

        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Apollo 13 Launch: 55 Years Ago
                - Linear Sand Dunes in the Great Sandy Desert
                - Expedition 73 Crew Launches to International Space Station
                etc...

### Read posts
1. `$ gator read` to navigate into your feeds:

        Choose a feed:
        1: Dwarf Fortress
        2: Guild Wars 2
        3: Nasa Image of the Day
        4: Pivot To AI

1. Enter the corresponding number to display recent posts:

        3       // user input

        Nasa Image of the Day
        Choose a post:
        1: Apollo 13 Launch: 55 Years Ago
                11:59 AM, Fri, 11 Apr 25
        2: Linear Sand Dunes in the Great Sandy Desert
                11:49 AM, Thu, 10 Apr 25
        3: Expedition 73 Crew Launches to International Space Station
                02:42 PM, Wed, 09 Apr 25

1. Enter the number corresponding to read the post:

        1       // user input

        Apollo 13 Launch: 55 Years Ago
        11:59 AM, Fri, 11 Apr 25

        NASA astronauts Jim Lovell, Fred Haise, and Jack Swigert launch aboard the Apollo 13 spacecraft from NASA’s Kennedy Space Center in Florida on April 11, 1970.

        https://www.nasa.gov/image-detail/apollo-13-launch-ref-msfc-70-msg-2200-13-mix-file-2/

### List feeds being followed
`$ gator following`

        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/


### List all available feeds
`$ gator -a`

        All RSS Feeds
        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

### Follow additional feeds
Users can follow RSS feeds that others have added. 

For example, Alice is following 2 feeds:

        * Nasa Image of the Day
        * Pivot To AI


1. `$ gator follow` to follow additional RSS feeds:

        Saved feeds:
        * Nasa Image of the Day
        * Pivot To AI

        Choose a new RSS feed to follow:
        1: Dwarf Fortress
        2: Guild Wars 2

1. Enter the number of the desired feed:

        1       // user input

        Alice followed Dwarf Fortress (https://store.steampowered.com/feeds/news/app/975370)

Feeds can also be followed directly by providing the feed url. The feed must have already been added to `gator`:

`$ gator follow https://www.nasa.gov/feeds/iotd-feed/`

### Unfollow feeds

`$ gator unfollow` to remove a feed you're following:

        $ gator unfollow

        - Choose an RSS feed to unfollow
        1: Guild Wars 2
        2: Nasa Image of the Day
        3: Pivot To AI

Enter the number of feed:

        1 (user input)

        Unfollowed Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370

### Help

Show all commands:

$ `gator -h` | `gator --help`

        Gator is a terminal-based RSS reader.
        It is best ran as a terminal background process (ex: gator ... &).
        Then, interact with the tool to read and manage RSS feeds.

        Usage:
        gator [flags]
        gator [command]

        Available Commands:
        add         Add a feed
        completion  Generate the autocompletion script for the specified shell
        delete      Delete a feed
        follow      Follow updates from a feed
        following   Lists all feeds a user if following
        help        Help about any command
        ...

## Development

### SQLC

Install [`sqlc`](https://sqlc.dev/), a useful library for generating boilerplate Go code from SQL queries:

        $ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

Find SQLC's config file in the base project directory:

        $ gator/sqlc.yaml

### Debugging

Sample Go `launch.json` config for Visual Studio Code:

```json
launch.json:

{
        "name": "add",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "/<path>/gator/go run main.go",
        "args": "add -n \"Nasa Image of the Day\" -u https://www.nasa.gov/feeds/iotd-feed/",
        "console": "integratedTerminal"
}
```

### Run locally
`$ go run . <command>`

Examples:
- `$ go run . update 15m &`
- `$ go run . login Alice`

### Test data

A shell script can quickly add users and feeds, see `test-data.sh`, etc.

### `gator --reset`

`gator` has a `reset` command that deletes all users from the database. Since all the the tables have `CASCADE`-ing deletes on the `users.id` field, this effectively clears the entire database (see: https://www.postgresql.org/docs/15/sql-createtable.html, search `CASCADE`).

Only use `reset` for development/local-use - it should (obviously) be commented-out/removed otherwise.

### Adding a new command

#### Modifying the database (optional)
1. If needed, create new `goose` up/down migrations (new tables, `ALTER` tables, etc.) in `sql/schema`

#### New/modify SQL statements (optional)
1. Add/edit `.sql` files in `sql/queries`. Reference existing queries for the needed `SQLC` syntax (i.e., `-- name: CreateFeed :one`), etc.
2. Run `$ sqlc generate` from the base project directory to generate Go code. Files within `internal/database` will be generated, including a new function based on your SQL statement. For example, a query named `CreateFeed` will generate a function: `func CreateFeed(...) ...`

#### Create a new command
1. Use cobra-cli to add a new command:

        $ cobra-cli add <name of new command>

If the command is a sub-command (ex: `gator <existing command> <new command>`):

        $ cobra-cli add <name of new command> -p <parent command>

Refer to [cobra-cli's README](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md) for more details.

2. Cobra-cli creates an easy-to-use template for commands with just 2 cobra-specific sections.

Using `add.go` (add a feed) as an example:

- Set the command description and the code the actually run when the command is called:
        ```
        // addCmd represents the add command
        var addCmd = &cobra.Command{
                Use:   "add",
                Short: "Add a feed",
                Long:  `Add a feed directly using the required flags.`,
                RunE: func(cmd *cobra.Command, args []string) error {
                        return userAuthCall(addFeed)(appState)
                },
        }
        ```
- init() - Set command-line flags. Refer to cobra-cli documentation on available functions:

        func init() {
                rootCmd.AddCommand(addCmd)

                addCmd.Flags().StringVarP(&feedNameArg, "name", "n", "", "Name of the  RSS feed (required)")
                addCmd.Flags().StringVarP(&feedUrlArg, "url", "u", "", "Url to the RSS feed (required)")

                addCmd.MarkFlagRequired("name")
                addCmd.MarkFlagRequired("url")

                addCmd.MarkFlagsRequiredTogether("name", "url")
        }

All other code are just plain Go functions.

#### User-authenticated commands
Commands that require a user id (ex: `gator add`) should use the handler function signature:

```go
func (s *state, user database.User) error { ... }
```

Then, in `main.go` the handler function will be wrapped in `userAuthCall(...)` closure. The closure will provide the currently logged-in user to the handler function. 

Commands should *not* need to duplicate manually retrieving the current user from the database.

### Manually test
From the base directory, run:

        $ go run . agg 15m &
        $ go run . add...
        $ go run . etc...

### Build/install the program
From the base directory:

        $ go build .
        $ go install .

Use the program:

        $ gator update 15m &
        [1] 4015
        Updating RSS feeds...

        $ gator
        Nasa Image of the Day | https://www.nasa.gov/feeds/iotd-feed/
                - Artemis II Insignia Honors All
                - X-ray Clues Reveal Destroyed Planet

## Improvements
- Tests, unit/integration
- Process HTML info RSS feed posts
- Retrieve full blog posts from feed URL
