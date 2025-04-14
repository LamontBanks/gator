/*
See custom RSS feed updates
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

// Follow params
var feedUrlParam string

// followCmd represents the followFeed command
var followCmd = &cobra.Command{
	Use:   "follow",
	Short: "Follow updates from a feed",
	Long: `Follow updates from a feed
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			feedUrlParam = args[0]
		}
		if feedUrlParam != "" {
			return userAuthCall(followFeedByUrlInternal)(appState)
		} else {
			return userAuthCall(interactiveFollowFeed)(appState)
		}
	},
}

func init() {
	rootCmd.AddCommand(followCmd)
}

func interactiveFollowFeed(s *state, user database.User) error {
	feedsNotFollowed, err := s.db.GetFeedsNotFollowedByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feedsNotFollowed) == 0 {
		fmt.Println("- No feeds to follow")
		return nil
	}

	// Show user what feeds they're already following
	feedsAlreadyFollowed, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Saved feeds")
	if len(feedsAlreadyFollowed) == 0 {
		fmt.Println("- No saved feeds")
	}

	for _, feed := range feedsAlreadyFollowed {
		fmt.Printf("* %v", feed.FeedName)
	}
	fmt.Println()

	// Create label-value 2D array for the option picker
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

	return followFeedByUrl(s, user, feedsNotFollowed[choice].Url)
}

// Wrapper function to use the feed URL supplied with this command
// Only for use within this cobra command - it relies on the url param being set from this command
//
// Done this way to make the actual function - followFeedByUrl - callable outside of this command
// For ex: Use the `add` command to add a feed, then auto follow it
// TODO: Better function name, or better way to handle this?
func followFeedByUrlInternal(s *state, user database.User) error {
	return followFeedByUrl(s, user, feedUrlParam)
}

func followFeedByUrl(s *state, user database.User, feedUrl string) error {
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
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
