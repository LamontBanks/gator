#!/bin/sh

go run . reset

go run . register Alice
go run . register Bob
go run . register Cat

go run . login Alice
go run . addFeed "Nasa Image of the Day" https://www.nasa.gov/feeds/iotd-feed/

go run . login Bob
go run . addFeed "Phys.org | Space News" https://phys.org/rss-feed/space-news/

go run . login Alice
go run . follow https://phys.org/rss-feed/space-news/

go run . feeds