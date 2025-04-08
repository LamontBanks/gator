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

// addFeedCmd represents the add command
var addFeedCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds an RSS feed to gator",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(addFeed)(appState)
	},
}

func init() {
	feedsCmd.AddCommand(addFeedCmd)

	addFeedCmd.Flags().StringVarP(&feedNameArg, "name", "n", "", "Name of new feed (required)")
	addFeedCmd.Flags().StringVarP(&feedUrlArg, "url", "u", "", "Url of new feed (required)")

	addFeedCmd.MarkFlagRequired("name")
	addFeedCmd.MarkFlagRequired("url")

	addFeedCmd.MarkFlagsRequiredTogether("name", "url")
}

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

	// Add the feed, attributed to the user
	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedNameArg,
		Url:       feedUrlArg,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add RSS feed %v (%v)", feedNameArg, feedUrlArg)
	}
	fmt.Printf("Added RSS feed \"%v\" (%v)\n", newFeed.Name, newFeed.Url)

	return nil
}
