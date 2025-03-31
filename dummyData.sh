#!/bin/zsh

go run . reset

go run . register Ant
go run . register Bird

go run . login Ant
go run . addFeed "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "GW2 Dev Tracker" https://en-forum.guildwars2.com/discover/6.xml
go run . addFeed "MassivelyOP" https://massivelyop.com/feed/


go run . login Bird
go run . follow "Pivot To AI" https://pivot-to-ai.com/feed/
go run . addFeed "Guild Wars 2" https://www.guildwars2.com/en/feed/
go run . follow "GW2 Dev Tracker" https://en-forum.guildwars2.com/discover/6.xml

go run . login Ant
go run . follow https://www.guildwars2.com/en/feed/

go run . users
go run . feeds
