#!/bin/zsh

go run . reset

go run . register alice
go run . register bob
go run . register cal

go run . users

go run . login alice
go run . addFeed "Pivot TO AI" https://pivot-to-ai.com/feed/

go run . login bob
go run . addFeed "Pivot TO AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss

go run . login bob
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss
go run . addFeed "MassivelyOP Podcast" https://www.guildwars2.com/en/feed/

go run . feeds