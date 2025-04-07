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

// feedsCmd represents the feeds command
var feedsCmd = &cobra.Command{
	Use:   "feeds",
	Short: "View and managed your feeds",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

feeds			Show recent posts from all feeds you're following

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listFollowedFeeds()
	},
}

func init() {
	rootCmd.AddCommand(feedsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// feedsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// feedsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listFollowedFeeds() error {
	user, err := getCurrentUser(appState)
	if err != nil {
		return err
	}

	// Get feeds followed by user
	feeds, err := appState.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("Not following any feeds")
		return nil
	}

	// Pull posts for each feed
	for _, feed := range feeds {
		posts, err := appState.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  3,
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				printPostTitle(post.Title)
			}
		} else {
			fmt.Println("No posts")
		}
	}

	return nil
}

func printPostTitle(title string) {
	fmt.Printf("\t- %v\n", title)
}
