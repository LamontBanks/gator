package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func followCommandInfo() commandInfo {
	return commandInfo{
		description: "Follow a registered feed",
		usage:       "follow <RSS feed URL, optional>",
		examples: []string{
			"follow\n<choose from list of feeds>",
			"follow http://example.com/rss/feed",
		},
	}
}

// Sets the current user as a follower of the given RSS feed.
func handlerFollow(s *state, cmd command, user database.User) error {
	// Args: url, optional
	var feedUrl string
	if len(cmd.args) > 0 {
		feedUrl = cmd.args[0]
	} else {
		// If no URL is provided, make user choose from feed they're not following
		feedsNotFollowed, err := s.db.GetFeedsNotFollowedByUser(context.Background(), user.ID)
		if err != nil {
			return err
		}

		// Show followed feeds to remind the user of what they already have
		feedsAlreadyFollowed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
		if err != nil {
			return err
		}

		fmt.Println("\nAlready following:")
		for _, feed := range feedsAlreadyFollowed {
			printFeed(feed.FeedName, feed.Description, feed.FeedUrl)
		}
		fmt.Println()

		// Create label-value 2D array for the menu generator
		feedOptions := make([][]string, len(feedsNotFollowed))
		for i := range feedsNotFollowed {
			feedOptions[i] = make([]string, 2)
			feedOptions[i][0] = feedsNotFollowed[i].Name + "\n\t" + feedsNotFollowed[i].Description + "\n\t" + feedsNotFollowed[i].Url
			feedOptions[i][1] = feedsNotFollowed[i].Url
		}

		// Choose feed
		_, feedUrl, err = listOptionsReadChoice(feedOptions, "- Choose a new RSS feed to follow:")
		if err != nil {
			return err
		}
	}

	// Get feed, make user follow
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
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
		return err
	}

	fmt.Printf("%v followed %v (%v)\n", newFeedFollow.UserName, newFeedFollow.FeedName, feed.Url)

	return nil
}
