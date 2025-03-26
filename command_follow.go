package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// Sets the current user as a follower of the given RSS feed.
// If the feed does not exist, it will be created.
func handlerFollow(s *state, cmd command) error {
	// Args: url
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <RSS Feed URL>", cmd.name)
	}
	feedUrl := cmd.args[0]

	// Get user info from username
	username := s.config.CurrentUserName
	if username == "" {
		return fmt.Errorf("unable to save feed - no user logged in")
	}
	userRow, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user %v not registered", username)
	}

	// Get feed info from the feedUrl
	feedRow, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("feed url %v has not been added yet", feedUrl)
	}
	if err != nil {
		return err
	}

	// Save userId -> feedId mapping
	queryResult, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userRow.ID,
		FeedID:    feedRow.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %v", queryResult.FeedName)
	}
	fmt.Printf("%v followed %v\n", queryResult.UserName, queryResult.FeedName)

	return nil
}
