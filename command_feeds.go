package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	// Args: feedName, feedUrl
	if len(cmd.args) < 2 {
		return fmt.Errorf("missing args: name, RSS url")
	}
	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	// Get userId from username
	username := s.config.CurrentUserName
	if username == "" {
		return fmt.Errorf("unable to save feed - no user logged in")
	}

	userRow, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("%v not registered", username)
	}

	// Insert feed info
	addFeedResult, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    userRow.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add: %v (%v) for %v - possible duplicate feed?", feedName, feedUrl, username)
	}

	fmt.Println(addFeedResult)
	fmt.Printf("Saved \"%v\" (%v) for user %v\n", feedName, feedUrl, username)

	return nil
}

// Lists all feeds from all users
func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Print(feeds)

	return nil
}
