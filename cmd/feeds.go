package cmd

import (
	"context"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// feedsCmd represents the feeds command
var (
	// Feed flags
	showAllFeeds    bool
	numPostsPerFeed int

	feedsCmd = &cobra.Command{
		Use:   "feeds",
		Short: "View and manage your feeds",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	feeds			Show recent posts from all feeds you're following
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showAllFeeds {
				return printAllFeeds(appState)
			} else {
				return userAuthCall(printFollowedFeeds)(appState)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(feedsCmd)

	feedsCmd.Flags().BoolVarP(&showAllFeeds, "all", "a", false, "Show all feeds added to gator")
	feedsCmd.Flags().IntVarP(&numPostsPerFeed, "numPosts", "n", 2, "maximum number of posts per feed")
}

func printFollowedFeeds(s *state, user database.User) error {
	if numPostsPerFeed < 0 {
		return fmt.Errorf("number of posts must be >= 0")
	}

	// Get feeds followed by user
	feeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("- Not following any feeds")
		return nil
	}

	// Pull posts for each feed

	fmt.Println("Your RSS Feeds:")
	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(numPostsPerFeed),
		})
		if err != nil {
			return err
		}

		// TODO: Print function for feeds, posts
		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.Description)
		if len(posts) > 0 {
			for _, post := range posts {
				fmt.Printf("\t- %v\n", post.Title)
			}
		} else {
			fmt.Println("- No posts")
		}
	}

	return nil
}

func printAllFeeds(s *state) error {
	if numPostsPerFeed < 0 {
		return fmt.Errorf("number of posts must be >= 0")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Print posts for each feed
	// TODO: Print function for feeds, posts
	fmt.Println("All RSS Feeds:")
	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.ID,
			Limit:  int32(numPostsPerFeed),
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.Description)
		if len(posts) > 0 {
			for _, post := range posts {
				fmt.Printf("\t- %v\n", post.Title)
			}
		} else {
			fmt.Println("- No posts")
		}
	}

	return nil
}
