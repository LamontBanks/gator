/*
Deletes an RSS feed
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a feed",
	Long: `Deletes a feed interactively.
A feed can only be delete by its creator and must have either no followers other than (optionally) the creator:

	gator delete

	You can only delete feeds you've added and have no other followers
	Choose a feed to delete:
	1: Phys.org | Space News

	1	# User choice

	Deleted feed Phys.org | Space News
`,
	Run: func(cmd *cobra.Command, args []string) {
		userAuthCall(interactiveDelete)(appState)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func interactiveDelete(s *state, user database.User) error {
	fmt.Println("You can only delete feeds you've added and have no other followers")

	feedsToDelete, err := s.db.GetFeedsEligibleForDeletion(context.Background(), user.ID)
	if len(feedsToDelete) == 0 {
		fmt.Println("- No eligible feeds to delete")
		return nil
	}
	if err != nil {
		return err
	}

	// Make option picker from list of feed names
	feedOptions := make([]string, len(feedsToDelete))
	for i := range feedsToDelete {
		feedOptions[i] = feedsToDelete[i].Name
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed to delete:")
	if err != nil {
		return err
	}

	feedToDelete := feedsToDelete[choice]

	// Delete the feed
	err = s.db.DeleteFeedById(context.Background(), feedToDelete.FeedID)
	if err != nil {
		return fmt.Errorf("failed to delete %v", feedToDelete.Name)
	}

	fmt.Printf("Deleted feed %v\n", feedToDelete.Name)
	return nil
}
