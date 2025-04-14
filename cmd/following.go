/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// followingCmd represents the following command
var followingCmd = &cobra.Command{
	Use:   "following",
	Short: "Lists all feeds a user if following",
	Long:  `Lists all feeds a user if following`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(followedFeeds)(appState)
	},
}

func init() {
	rootCmd.AddCommand(followingCmd)
}

func followedFeeds(s *state, user database.User) error {
	feedsAlreadyFollowed, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feedsAlreadyFollowed {
		fmt.Printf("* %v\n", feed.FeedName)
	}
	fmt.Println()

	return nil
}
