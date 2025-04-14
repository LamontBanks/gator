/*
 */
package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// flags/params for unfollow
var unfollowfeedUrlParam string

// unfollowCmd represents the unfollow command
var unfollowCmd = &cobra.Command{
	Use:   "unfollow",
	Short: "Stop following updates from a feed",
	Long:  `Stop following updates from a feed`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			unfollowfeedUrlParam = args[0]
			return userAuthCall(unfollowFeedByUrlInternal)(appState)
		}
		return userAuthCall(interactiveUnfollowFeed)(appState)
	},
}

func init() {
	rootCmd.AddCommand(unfollowCmd)
}

// Wrapper function to use the feed URL supplied with this command
// Only for use within this cobra command - it relies on the url param being set from this command
//
// Done this way to make the actual function - followFeedByUrl - callable outside of this command
// For ex: Use the `add` command to add a feed, then auto follow it
// TODO: Better function name, or better way to handle this?
func unfollowFeedByUrlInternal(s *state, user database.User) error {
	return unfollowFeedByUrl(s, user, unfollowfeedUrlParam)
}

func interactiveUnfollowFeed(s *state, user database.User) error {
	followedFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(followedFeeds) == 0 {
		fmt.Println("Not following any feeds")
		return nil
	}

	// Create label-value 2D array for the option picker
	feedOptions := make([][]string, len(followedFeeds))
	for i := range followedFeeds {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = followedFeeds[i].FeedName
		feedOptions[i][1] = followedFeeds[i].FeedUrl
	}

	// Choose feed to unfollow
	choice, err := listOptionsReadChoice(feedOptions, "- Choose an RSS feed to unfollow")
	if err != nil {
		return err
	}

	feedUrl := followedFeeds[choice].FeedUrl

	return unfollowFeedByUrl(s, user, feedUrl)
}

func unfollowFeedByUrl(s *state, user database.User, feedUrl string) error {
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		fmt.Printf("- Can't to unfollow %v - has not been added\n", feedUrl)
		return nil
	}
	if err != nil {
		return err
	}

	// Unfollow the feed
	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow %v", feedUrl)
	}

	fmt.Printf("Unfollowed %v | %v\n", feed.Name, feedUrl)
	return nil
}
