#!/bin/sh

go run main.go --reset

go run main.go register Alice

go run main.go add -n "Nasa Image of the Day" -u https://www.nasa.gov/feeds/iotd-feed/
go run main.go add -n "Pivot To AI" -u https://pivot-to-ai.com/feed/
go run main.go add -n "Guild Wars 2" -u https://www.guildwars2.com/en/feed/