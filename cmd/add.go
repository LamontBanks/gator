/*
Add an RSS feed
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// addCmd flags
var feedNameArg string
var feedUrlArg string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a feed",
	Long: `Add a feed directly using the required flags:

	gator add -n "Space News" -u "https://phys.org/rss-feed/space-news/"
	
	Added RSS feed "Space News" (https://phys.org/rss-feed/space-news/)
	Following "Space News" (https://phys.org/rss-feed/space-news/)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(addFeed)(appState)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&feedNameArg, "name", "n", "", "Name of the  RSS feed (required)")
	addCmd.Flags().StringVarP(&feedUrlArg, "url", "u", "", "Url to the RSS feed (required)")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("url")

	addCmd.MarkFlagsRequiredTogether("name", "url")
}

// Adds a feed record to gator
// Other user's can follow the feed to see updates when they're logged in.
// Does NOT download the posts
func addFeed(s *state, user database.User) error {
	// Don't do anything if the feed url has already been added, or there's any other error
	existingFeed, err := s.db.GetFeedByUrl(context.Background(), feedUrlArg)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != sql.ErrNoRows {
		fmt.Printf("Feed already exists:\n%v (%v)\n", existingFeed.Name, existingFeed.Url)
		return nil
	}

	// Attempt to download feed before adding entry
	rssFeed, err := FetchFeed(context.Background(), feedUrlArg)
	if err != nil {
		return fmt.Errorf("could not add feed %v, %v", feedUrlArg, err)
	}

	// Save the feed to the database
	newFeed, err := saveFeed(s, feedNameArg, feedUrlArg, user)
	if err != nil {
		return err
	}

	// Immediately download feed updates
	err = saveFeedPosts(s, rssFeed, newFeed.ID)
	if err != nil {
		return fmt.Errorf("error updating feed %v, %v", rssFeed.Channel.Title, err)
	}

	// Make user follow feed they just added
	return followFeedByUrl(s, user, newFeed.Url)
}

// Add the feed, attributed to the user
func saveFeed(s *state, feedName, feedUrl string, user database.User) (database.Feed, error) {
	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return database.Feed{}, fmt.Errorf("could not add RSS feed %v (%v)", feedNameArg, feedUrlArg)
	}

	fmt.Printf("Added RSS feed \"%v\" (%v)\n", newFeed.Name, newFeed.Url)
	return newFeed, nil
}
