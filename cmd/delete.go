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

	// Prompt to select feed
	// Create label-value 2D array for the option picker
	feedOptions := make([][]string, len(feedsToDelete))
	for i := range feedsToDelete {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = feedsToDelete[i].Name
		feedOptions[i][1] = feedsToDelete[i].FeedID.URN()
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed to delete:")
	if err != nil {
		return err
	}

	// Delete the feed
	err = s.db.DeleteFeedById(context.Background(), feedsToDelete[choice].FeedID)
	if err != nil {
		return fmt.Errorf("failed to delete %v", feedsToDelete[choice].Name)
	}

	fmt.Printf("Deleted feed %v\n", feedsToDelete[choice].Name)
	return nil
}
