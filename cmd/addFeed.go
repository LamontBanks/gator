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
		user, err := getCurrentUser(appState)
		if err != nil {
			return err
		}

		// Don't do anything if the feed url has already been added
		existingFeed, err := appState.db.GetFeedByUrl(context.Background(), feedUrlArg)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err != sql.ErrNoRows {
			fmt.Printf("Feed already exists:\n%v (%v)\n", existingFeed.Name, existingFeed.Url)
			return nil
		}

		// Add the feed
		newFeed, err := appState.db.CreateFeed(context.Background(), database.CreateFeedParams{
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
	},
}

func init() {
	feedsCmd.AddCommand(addFeedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")
	addFeedCmd.Flags().StringVarP(&feedNameArg, "name", "n", "", "Name of new feed (required)")
	addFeedCmd.Flags().StringVarP(&feedUrlArg, "url", "u", "", "Url of new feed (required)")

	addFeedCmd.MarkFlagRequired("name")
	addFeedCmd.MarkFlagRequired("url")

	addFeedCmd.MarkFlagsRequiredTogether("name", "url")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
