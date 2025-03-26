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
		return fmt.Errorf("usage: %v <Name> <RSS Feed URL>", cmd.name)
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
	fmt.Printf("Saved \"%v\" (%v) for user %v\n", addFeedResult.Name, addFeedResult.Url, userRow.Name)

	// Also follow the added feed
	// Save userId -> feedId mapping
	queryResult, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userRow.ID,
		FeedID:    addFeedResult.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %v", queryResult.FeedName)
	}
	fmt.Printf("%v followed %v\n", queryResult.UserName, queryResult.FeedName)

	return nil
}

// Lists all feeds from all users
func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds {
		fmt.Printf("%v - %v\n", f.FeedName, f.UserName.String)
	}

	return nil
}
