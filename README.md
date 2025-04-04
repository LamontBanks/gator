# Blog Aggregator

A CLI tool to:
- View RSS post titles in the CLI
- Support multple users, each having saved, custom RSS feed
- Follow and unfollow RSS feeds that other users have added

## Usage
### Start gator

Start `gator` as a background process, provide an update frequency `(30s, 1m, 5, 1h, 24h, etc.)`:
        
        $ gator agg 15m &
        Updating RSS feeds...

### Register/login a user

Register a new user:

        $ gator register Bob
        Registered user Bob
        Logged in as Bob

Or, login an existing user:

        $ gator login Alice
        Logged in as Alice

### Add RSS feeds
 `$ gator addFeed <url>`
 
 Newly added feeds will immediately download updates:

        $ gator addFeed  "Nasa Image of the Day" https://www.nasa.gov/feeds/iotd-feed/
        Saved "Nasa Image of the Day" (https://www.nasa.gov/feeds/iotd-feed/) for Alice

        $ gator addFeed "Dwarf Fortress" https://store.steampowered.com/feeds/news/app/975370
        Saved "Dwarf Fortress" (https://store.steampowered.com/feeds/news/app/975370) for Alice
        Alice followed Dwarf Fortress


### See RSS Feeds overview
`$ gator browse`

        $ gator browse

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

        $ gator browseFeed

        Choose a feed:
        1: Dwarf Fortress
        2: Nasa Image of the Day

1. Enter the corresponding number to display recent posts:

        2 (user input)

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

        3 (user input)

        Studying Ice for the Future of Flight
        12:52 PM, Tue, 01 Apr 25

        Thomas Ozoroski, a researcher at NASA’s Glenn Research Center in Cleveland, takes icing accretion measurements in October 2024 as part of transonic truss-braced wing concept research. Researchers at NASA Glenn conducted another test campaign in March 2025.

        https://www.nasa.gov/image-detail/grc-2024-c-12100-2/

### List feeds you follow
`$ gator following`

        $ gator following

        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/


### List all available feeds
`$ gator feeds`

        # gator feeds

        All RSS Feeds
        Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

### Follow additional feeds
Users can follow RSS feeds that others have added. For example, Bob is currenetly not following any feeds:

        $ gator login Bob
        Logged in as Bob

        $ gator following
        you are not following any feeds


`$ gator follow` to select additional RSS feeds to follow:

        Already following:
        no feeds

        - Choose a new RSS feed to follow:
        1: Dwarf Fortress
                Events and Announcements for 975370
                https://store.steampowered.com/feeds/news/app/975370
        2: Nasa Image of the Day
                The latest NASA "Image of the Day" image.
                https://www.nasa.gov/feeds/iotd-feed/

Enter the number of the desired option:

        1 (user input)

        Bob followed Dwarf Fortress (https://store.steampowered.com/feeds/news/app/975370)

View updates:

        $ gator browse

        Dwarf Fortress | https://store.steampowered.com/feeds/news/app/975370
	- Hotfix: Patch 51.10
	- Update 51.09: Forbidden Coffins now unusable
	- Hotfix: Patch 51.08

## Tools required
- PostgresQL - Data storage

        $ brew install postgresql@15

- Goose: Database migrations

        $ go install github.com/pressly/goose/v3/cmd/goose@latest

- SQLC - Generate Go code from SQL queries
    
        $ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest



## Adding new commands

1. Add 