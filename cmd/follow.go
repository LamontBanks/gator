/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// followCmd represents the followFeed command
var followCmd = &cobra.Command{
	Use:   "follow",
	Short: "Follow a feed to see its updates",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return interactiveFollowFeed()
	},
}

func init() {
	rootCmd.AddCommand(followCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// followFeedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// followFeedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func interactiveFollowFeed() error {
	user, err := getCurrentUser(appState)
	if err != nil {
		return err
	}

	// Get feeds not followed...
	feedsNotFollowed, err := appState.db.GetFeedsNotFollowedByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feedsNotFollowed) == 0 {
		fmt.Println("No feeds to follow")
		return nil
	}

	// ...and followed feeds to show the user what they already have
	feedsAlreadyFollowed, err := appState.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Saved feeds")
	if len(feedsAlreadyFollowed) == 0 {
		fmt.Println("- No saved feeds")
	}

	for _, feed := range feedsAlreadyFollowed {
		fmt.Printf("- %v", feed.FeedName)
	}
	fmt.Println()

	// Create label-value 2D array for the option picker, choose feed to follow
	feedOptions := make([][]string, len(feedsNotFollowed))
	for i := range feedsNotFollowed {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = feedsNotFollowed[i].Name
		feedOptions[i][1] = feedsNotFollowed[i].Url
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a new RSS feed to follow:")
	if err != nil {
		return err
	}

	feedUrl := feedsNotFollowed[choice].Url

	// Get desired feed by url, follow it
	feed, err := appState.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	newFeedFollow, err := appState.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Following %v (%v)\n", newFeedFollow.FeedName, feed.Url)

	return nil
}
