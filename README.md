# Blog Aggregator

A CLI tool to:
- View RSS post titles in the CLI
- Support multple users, each having saved, custom RSS feed
- Follow and unfollow RSS feeds that other users have added


```
GW2 Dev Tracker | https://en-forum.guildwars2.com/discover/6.xml
- Professor Smoll has returned for April Fools 2025!
  https://en-forum.guildwars2.com/topic/157030-professor-smoll-has-returned-for-april-fools-2025/?do=findComment&comment=2268781
- Latest Update Has Shrunk My Male Norn Weapons/Glider
  https://en-forum.guildwars2.com/topic/157011-latest-update-has-shrunk-my-male-norn-weaponsglider/?do=findComment&comment=2268540
- Game Update Notes: March 25, 2025
  https://en-forum.guildwars2.com/topic/156800-game-update-notes-march-25-2025/?do=findComment&comment=2268507

Guild Wars 2 | https://www.guildwars2.com/en/feed/
- Stop and Smell the Flowers with the Wisteria Arborscale Skyscale Skin
  https://www.guildwars2.com/en/news/stop-and-smell-the-flowers-with-the-wisteria-arborscale-skyscale-skin/?utm_source=rss&utm_medium=news&utm_campaign=rss

MassivelyOP | https://massivelyop.com/feed/
- The Stream Team: Reveling in the highest of silly days in EverQuest II
  https://massivelyop.com/2025/04/01/the-stream-team-reveling-in-the-highest-of-silly-days-in-everquest-ii/
- So how is the Lord of the Rings Online 2025 housing rush going so far? Well, it’s going
  https://massivelyop.com/2025/04/01/so-how-is-the-lord-of-the-rings-online-2025-housing-rush-going-so-far-well-its-going/
- DC Universe Online marks the spring season with a pollen-spreading pod and new cosmetics
  https://massivelyop.com/2025/04/01/dc-universe-online-marks-the-spring-season-with-a-pollen-spreading-pod-and-new-cosmetics/

Pivot To AI | https://pivot-to-ai.com/feed/
- Nothing in the last 18h
```

## Tools
- PostgresQL - Data storage

        $ brew install postgresql@15

- Goose: Database migrations

        $ go install github.com/pressly/goose/v3/cmd/goose@latest

- SQLC - Generate Go code from SQL queries
    
        $ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest



## Adding new commands

1. Add 