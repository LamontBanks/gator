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

var postSearchStr string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		postSearchStr = args[0]

		if postSearchStr == "" {
			return fmt.Errorf("no search string provided")
		}

		// Search for posts
		postsFound, err := appState.db.SearchPostTitles(context.Background(), fmt.Sprintf("%v%v%v", "%", postSearchStr, "%"))
		if err != nil {
			return fmt.Errorf("error searching post title for %v", postSearchStr)
		}

		// Format posts output
		var postsFoundOutput string

		if len(postsFound) == 0 {
			postsFoundOutput = "- No posts found"
		} else {
			foundPostsSlice := []string{"Feed | Date | Post"}
			for _, result := range postsFound {
				foundPostsSlice = append(foundPostsSlice, fmt.Sprintf("%v | %v | %v", result.FeedName, fuzzytimestamp.FuzzyTimestamp(result.PublishedAt), result.Title))
			}
			postsFoundOutput = columnize.SimpleFormat(foundPostsSlice)
		}

		fmt.Println(postsFoundOutput)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
