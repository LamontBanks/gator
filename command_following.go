package main

import (
	"context"
	"fmt"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

// Prints details of all the feeds the current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	// Get followed feed details
	feedDetails, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feedDetails) == 0 {
		fmt.Println("you are not following any feeds")
		return nil
	}

	// List feed names
	for _, feed := range feedDetails {
		fmt.Println(feed.FeedName)
	}

	return nil
}
