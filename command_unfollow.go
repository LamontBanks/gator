package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

// Unfollows the given RSS feel URL
func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Args: url
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	feedUrl := cmd.args[0]

	// Get feed info from the feedUrl
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("feed url %v has not been added yet", feedUrl)
	}
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow %v", feedUrl)
	}

	fmt.Printf("Unfollowed feed %v | %v\n", feed.Name, feedUrl)
	return nil
}
