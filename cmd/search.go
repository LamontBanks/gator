/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var postSearchStr string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := appState.db.SearchPostTitles(context.Background(), fmt.Sprintf("%v%v%v", "%", postSearchStr, "%"))
		if err != nil {
			return fmt.Errorf("error searching post title for %v", postSearchStr)
		}

		if len(results) == 0 {
			fmt.Println("- No results found")
			return nil
		}

		fmt.Println("Posts found:")
		fmt.Println("Feed\t\t\t\tTitle\n---\t\t\t\t---")
		for _, result := range results {
			fmt.Printf("%v\t\t\t%v\n", result.FeedName, result.Title)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&postSearchStr, "posts", "p", "", "Search all post titles")
	searchCmd.MarkFlagRequired("posts")
}
