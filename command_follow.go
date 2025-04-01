package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func followCommandInfo() commandInfo {
	return commandInfo{
		description: "Follow a registered feed",
		usage:       "follow <RSS feed URL>",
		examples: []string{
			"follow http://example.com/rss/feed",
		},
	}
}

// Sets the current user as a follower of the given RSS feed.
func handlerFollow(s *state, cmd command, user database.User) error {
	// Args: url
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	feedUrl := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("failed to follow %v - not yet added", feedUrl)
	}
	if err != nil {
		return err
	}

	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow %v; you may already be following it?", feedUrl)
	}
	fmt.Printf("%v followed %v\n", newFeedFollow.UserName, newFeedFollow.FeedName)

	return nil
}
