#!/bin/sh

go run . reset

go run . register Alice
go run . register Bob
go run . register Cat

go run . login Alice
go run . addFeed "Nasa Image of the Day" https://www.nasa.gov/feeds/iotd-feed/
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/


go run . login Bob
go run . addFeed "Phys.org | Space News" https://phys.org/rss-feed/space-news/
go run . addFeed "NotAnRssFeed" "https://google.com"
go run . addFeed "GW2 Dev Tracker" https://en-forum.guildwars2.com/discover/6.xml

go run . login Cat
go run . follow https://phys.org/rss-feed/space-news/
go run . follow https://pivot-to-ai.com/feed/


go run . feeds