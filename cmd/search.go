/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	fuzzytimestamp "github.com/LamontBanks/gator/internal/fuzzy_timestamp"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

var searchStr string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search all: posts titles, feed titles, and feed descriptions for the given text",
	Long: `Search all: posts titles, feed titles, and feed descriptions for the given text

	gator search word
	gator search "word with spaces"

Examples:

	# Returns feeds and posts that contain the word "guild" in their title or description
	gator search guild

	Feeds: 2 results found
	---
	Feed            Desc
	Dev Tracker     Guild Wars 2 Forums - Dev Tracker
	GuildWars2.com  

	Posts: 56 results found
	---
	Post                                                                                            Date  Feed
	Guild Wars Reforged Launches on Mobile This Summer!                                             2w    Dev Tracker
	Guild Wars Reforged Launches on Mobile This Summer!                                             2w    Dev Tracker
	Announcing Guild Wars 3!                                                                        3w    Dev Tracker
	...

	# Returns feeds and posts that contain the phrase "coming soon" in their title or description
	gator search "coming soon"

	Feeds: 0 results found
	---
	- No feeds found

	Posts: 4 results found
	---
	Post                                                                        Date  Feed
	Coming Soon:  An update to fashion templates to maintain free dye changes.  5mo   Dev Tracker
	Coming Soon:  An update to fashion templates to maintain free dye changes.  5mo   Dev Tracker
	Balance Patch coming soon?                                                  8mo   Dev Tracker
	Coming Soon: Fashion Templates!                                             6mo   GuildWars2.com
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchStr = args[0]

		if searchStr == "" {
			return fmt.Errorf("no search string provided")
		}

		// Search for posts, format into columns
		postsFound, err := appState.db.SearchPostTitles(context.Background(), fmt.Sprintf("%v%v%v", "%", searchStr, "%"))
		if err != nil {
			return fmt.Errorf("error searching post title for %v", searchStr)
		}

		var postsFoundOutput string
		if len(postsFound) == 0 {
			postsFoundOutput = "- No posts found"
		} else {
			foundPostsSlice := []string{"Post | Date | Feed"}
			for _, result := range postsFound {
				foundPostsSlice = append(foundPostsSlice, fmt.Sprintf("%v | %v | %v", result.Title, fuzzytimestamp.FuzzyTimestamp(result.PublishedAt), result.FeedName))
			}
			postsFoundOutput = columnize.SimpleFormat(foundPostsSlice)
		}

		// Search for feeds, format into columns
		feedsFound, err := appState.db.SearchFeeds(context.Background(), fmt.Sprintf("%v%v%v", "%", searchStr, "%"))
		if err != nil {
			return fmt.Errorf("error searching feeds for %v", searchStr)
		}

		var feedsFoundOutput string
		if len(feedsFound) == 0 {
			feedsFoundOutput = "- No feeds found"
		} else {
			foundFeedSlice := []string{"Feed | Desc"}
			for _, result := range feedsFound {
				foundFeedSlice = append(foundFeedSlice, fmt.Sprintf("%v | %v", result.Name, result.Description))
			}
			feedsFoundOutput = columnize.SimpleFormat(foundFeedSlice)
		}

		fmt.Printf("Feeds: %v results found\n---\n", len(feedsFound))
		fmt.Println(feedsFoundOutput)

		fmt.Println()
		fmt.Printf("Posts: %v results found\n---\n", len(postsFound))
		fmt.Println(postsFoundOutput)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
