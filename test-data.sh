#!/bin/sh

go run main.go --reset

go run main.go register Alice
go run main.go register Bob
go run main.go register Cat

go run main.go login Alice
go run main.go add https://www.nasa.gov/feeds/iotd-feed/
go run main.go add https://pivot-to-ai.com/feed/
go run main.go add https://www.guildwars2.com/en/feed/

go run main.go login Bob
go run main.go add https://phys.org/rss-feed/space-news/
go run main.go add "https://example.com" # purposefully attempt to add a non RSS url
go run main.go add https://en-forum.guildwars2.com/discover/6.xml

go run main.go login Cat
# Make choice for interactive commands: echo "<feed number>\n<post number>"
# Ex:
# Choose option 2, then 3:
# echo "2\n3" | go run main.go read

# Follow first 2 feeds
echo "1" | go run main.go follow
echo "1" | go run main.go follow

go run main.go login Alice
go run main.go