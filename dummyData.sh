#!/bin/zsh

go run . reset

go run . register Ant
go run . register Bird
go run . register Cat

go run . login Ant
go run . addFeed "Nasa Image of the Day" https://www.nasa.gov/feeds/iotd-feed/
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "GW2 Dev Tracker" https://en-forum.guildwars2.com/discover/6.xml


go run . login Bird
go run . follow https://pivot-to-ai.com/feed/
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/
go run . follow https://en-forum.guildwars2.com/discover/6.xml

go run . login Ant
go run . follow https://www.guildwars2.com/en/feed/

go run . login Cat
go run . addFeed "Dwarf Fortress" https://store.steampowered.com/feeds/news/app/975370
go run . addFeed "Baldur's Gate 3" https://store.steampowered.com/feeds/news/app/1086940

go run . users
go run . feeds
