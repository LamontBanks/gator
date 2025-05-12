/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"slices"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := appState.db.SearchPostTitles(context.Background(), fmt.Sprintf("%v%v%v", "%", postSearchStr, "%"))
		if err != nil {
			return fmt.Errorf("error searching post title for %v", postSearchStr)
		}

		if len(results) == 0 {
			fmt.Println("- No results found")
			return nil
		}

		output := []string{"Feed | Date | Post"}

		seenFeedNames := []string{}

		for _, result := range results {
			// Only print feeds names once
			// Assumes feeds are grouped
			feedName := ""
			if !slices.Contains(seenFeedNames, result.FeedName) {
				seenFeedNames = append(seenFeedNames, result.FeedName)
				feedName = result.FeedName
			}

			output = append(output, fmt.Sprintf("%v | %v | %v", feedName, fuzzytimestamp.FuzzyTimestamp(result.PublishedAt), result.Title))
		}

		fmt.Println(columnize.SimpleFormat(output))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&postSearchStr, "posts", "p", "", "Search all post titles")
	searchCmd.MarkFlagRequired("posts")
}
