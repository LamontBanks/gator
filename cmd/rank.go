/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// rankCmd represents the rank command
var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "List feeds by number of followers in descending order",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		follower_count, err := appState.db.GetFeedFollowerCount(context.Background())
		if err != nil {
			return fmt.Errorf("error retrieving feed follower count")
		}

		if len(follower_count) == 0 {
			fmt.Println("- No feeds saved")
			return nil
		}

		for _, count := range follower_count {
			fmt.Printf("%v\n\t%v\n", count.FeedName, count.NumFollowers)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rankCmd)
}
