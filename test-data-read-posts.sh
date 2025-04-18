#!/bin/sh

go run main.go --reset

go run main.go register Alice

go run main.go add -n "Nasa Image of the Day" -u https://www.nasa.gov/feeds/iotd-feed/
go run main.go add -n "Pivot To AI" -u https://pivot-to-ai.com/feed/
go run main.go add -n "Guild Wars 2" -u https://www.guildwars2.com/en/feed/

# Send choice to program: echo "<feed number>\n<post number>"
# Ex:
# echo "2\n3" | go run main.go read 
#
# Choose a feed:
# 1: Guild Wars 2
# 2: Nasa Image of the Day  <---
# 3: Pivot To AI

# Nasa Image of the Day
# Choose a post:
# 1: Hubble Spies Cosmic Pillar in Eagle Nebula
#         03:32 PM, Fri, 18 Apr 25
# 2: Space Shuttle Discovery Lifts Off
#         04:27 PM, Thu, 17 Apr 25
# 3: Scrub Jay at the Vehicle Assembly Building     <---
#         05:08 PM, Wed, 16 Apr 25
#
# Scrub Jay at the Vehicle Assembly Building
# 05:08 PM, Wednesday, 16 Apr

# A scrub jay perches on a branch near the Vehicle Assembly Building at NASA’s Kennedy Space Center in Florida on June 22, 2020.

# https://www.nasa.gov/image-detail/afs-8-101-1017/
echo "2\n3" | go run main.go read 

echo "3\n1" | go run main.go read 
echo "3\n2" | go run main.go read 
echo "2\n1" | go run main.go read 