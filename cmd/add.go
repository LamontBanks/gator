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
var feedUrlArg string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a feed",
	Long: `Add a feed by providing the URL:

	gator add "https://phys.org/rss-feed/space-news/"
	
	Added RSS feed "Space News" (https://phys.org/rss-feed/space-news/)
	Following "Space News" (https://phys.org/rss-feed/space-news/)
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		feedUrlArg = args[0]
		return userAuthCall(addFeed)(appState)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
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

	// Attempt to download feed
	rssFeed, err := FetchFeed(context.Background(), feedUrlArg)
	if err != nil {
		return fmt.Errorf("could not download feed %v: %v", feedUrlArg, err)
	}

	// Save the feed to the database
	newFeed, err := saveFeed(s, rssFeed.Channel.Title, feedUrlArg, user)
	if err != nil {
		return err
	}

	// Immediately download feed updates
	err = saveFeedPosts(s, rssFeed, newFeed.ID)
	if err != nil {
		return fmt.Errorf("error saving feed posts, %v", err)
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
		return database.Feed{}, fmt.Errorf("could not save new feed entry, %v", err)
	}

	fmt.Printf("Added RSS feed \"%v\"\n", newFeed.Name)
	return newFeed, nil
}
