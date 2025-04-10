/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read feed posts",
	Long: `Read feed posts.
A interactive menu will help navigate your followed feeds, then to the posts within a feed.
	
Currently only a plaintext <description> is readable in the terminal.
Images will not render, HTML will be raw, etc.
The full-text of the post, if any, will have to be viewed in a web browser.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(readPosts)(appState)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func readPosts(s *state, user database.User) error {
	userFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		return err
	}

	// Copy feedNames, feedUrl into a label-value 2D slive, pass to the option picker, select the feed
	feedOptions := make([][]string, len(userFeeds))
	for i := range userFeeds {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = userFeeds[i].FeedName
		feedOptions[i][1] = userFeeds[i].FeedUrl
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}
	fmt.Println(userFeeds[choice].FeedName)

	// Get posts for the selected feed
	posts, err := s.db.GetRecentPostsFromFeed(context.Background(), database.GetRecentPostsFromFeedParams{
		FeedID: userFeeds[choice].FeedID,
		Limit:  int32(3),
	})
	if err != nil {
		return err
	}

	// Copy postTitle, postId into a label-value 2D slice, pass to the option picker, select the post
	postOptions := make([][]string, len(posts))
	for i := range posts {
		postOptions[i] = make([]string, 2)
		postOptions[i][0] = posts[i].Title + "\n\t" + posts[i].PublishedAt.In(time.Local).Format("03:04 PM, Mon, 02 Jan 06")
		postOptions[i][1] = posts[i].ID.String()
	}

	choice, err = listOptionsReadChoice(postOptions, "Choose a post:")
	if err != nil {
		return err
	}

	printPost(posts[choice].Title, posts[choice].Description, posts[choice].Url, posts[choice].PublishedAt)

	return nil
}

func printPost(title, desc, link string, published_at time.Time) {
	s := fmt.Sprintf("%v\n", title)
	s += fmt.Sprintf("%v\n\n", published_at.In(time.Local).Format("03:04 PM, Mon, 02 Jan 06"))
	s += fmt.Sprintf("%v\n\n", desc)
	s += fmt.Sprintf("%v\n", link)

	fmt.Println(s)
}
