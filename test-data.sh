#!/bin/sh

go run main.go --reset

go run main.go register Alice
go run main.go register Bob
go run main.go register Cat

go run main.go login Alice
go run main.go feeds add -n "Nasa Image of the Day" -u https://www.nasa.gov/feeds/iotd-feed/
go run main.go feeds add -n "Pivot To AI" -u https://pivot-to-ui.com/feed/
go run main.go feeds add -n "Guild Wars 2" -u https://www.guildwars2.com/en/feed/


go run main.go login Bob
go run main.go feeds add -n "Phys.org | Space News" -u https://phys.org/rss-feed/space-news/
go run main.go feeds add -n "NotAnRssFeed" -u "https://google.com"
go run main.go feeds add -n "GW2 Dev Tracker" -u https://en-forum.guildwars2.com/discover/6.xml

# go run main.go login Cat
# go run main.go follow https://phys.org/rss-feed/space-news/
# go run main.go follow https://pivot-to-ai.com/feed/

go run main.go login Alice
go run main.go feeds