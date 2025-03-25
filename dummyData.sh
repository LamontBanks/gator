#!/bin/zsh

go run . reset

go run . register alice
go run . register bob
go run . register cat

go run . login alice
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/

go run . login bob
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss

go run . login cat
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/
