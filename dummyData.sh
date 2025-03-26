#!/bin/zsh

go run . reset

go run . register ant
go run . register bird
go run . register cat
go run . register dog

go run . login ant
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/

go run . login bird
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss

go run . login cat
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "MassivelyOP Podcast" https://feeds.libsyn.com/63150/rss
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/

echo "Sample feeds to add"
echo "https://hnrss.org/newest"
echo "https://www.wagslane.dev/index.xml"

go run . users
