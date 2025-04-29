/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Marks all posts as read",
	Long: `Marks all posts as read

	gator read clear

	Choose a feed:
	1: Guild Wars 2
	2: Nasa Image of the Day
	3: Pivot To AI

	1	# User choice

	All Guild Wars 2 posts marked as read
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(markAllPostsAsRead)(appState)
	},
}

func init() {
	readCmd.AddCommand(clearCmd)
}

func markAllPostsAsRead(s *state, user database.User) error {
	// Get feed from menu picker
	// TODO: Only get feeds with unread posts
	userFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("not following any feeds")
	}
	if err != nil {
		return err
	}

	// Make option picker from list of feed names
	feedOptions := make([]string, len(userFeeds))
	for i := range userFeeds {
		feedOptions[i] = userFeeds[i].FeedName
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}

	feed := userFeeds[choice]

	// Exit if all posts are already read
	unreadPostCount, err := getUnreadPostCount(s, user, feed.FeedID)
	if err != nil {
		return err
	}
	if unreadPostCount == 0 {
		fmt.Printf("- No unread posts in %v", feed.FeedName)
		return nil
	}

	// Get all posts from chosen feed
	allPosts, err := s.db.GetAllPostsFromFeed(context.Background(), feed.FeedID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error getting all post for feed %v, %v", feed.FeedName, err)
	}
	if err == sql.ErrNoRows {
		fmt.Printf("- No posts for %v\n", feed.FeedName)
	}

	// Mark all posts as read
	for _, post := range allPosts {
		if err := markPostAsRead(s, user, feed.FeedID, post.ID); err != nil {
			continue
		}
	}

	fmt.Printf("All %v posts marked as read\n", feed.FeedName)
	return nil
}
