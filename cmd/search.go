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
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
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
