package main

import (
	"context"
	"fmt"
)

// Prints details of all the feeds the current user is following
func handlerFollowing(s *state, cmd command) error {
	// Get user info from username
	username := s.config.CurrentUserName
	if username == "" {
		return fmt.Errorf("unable to save feed - no user logged in")
	}
	userRow, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("%v not registered", username)
	}

	// Get followed feed details
	feedDetails, err := s.db.GetFeedFollowsForUser(context.Background(), userRow.ID)
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
