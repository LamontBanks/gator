package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	// Args: feedUrl
	if len(cmd.args) < 1 {
		return fmt.Errorf("RSS feed URL required")
	}
	feedUrl := cmd.args[0]

	// Get userId from username
	username := s.config.CurrentUserName
	if username == "" {
		return fmt.Errorf("unable to save feed - no user logged in")
	}

	userRow, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("%v not registered", username)
	}

	// Get RSSFeed details
	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	// Insert user
	_, err = s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
		Url:       rssFeed.Channel.Link,
		UserID:    userRow.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add: %v (%v) for %v - possible duplicate feed?", rssFeed.Channel.Title, rssFeed.Channel.Link, username)
	}
	fmt.Printf("Saved \"%v\" (%v) for user %v\n", rssFeed.Channel.Title, rssFeed.Channel.Link, username)

	return nil
}
