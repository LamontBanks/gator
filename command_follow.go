package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
)

func followCommandInfo() commandInfo {
	return commandInfo{
		description: "Follow a registered feed",
		usage:       "follow <feed URL, optional>",
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
		// If no URL is provided, choose from list of available feeds
		feedsNotFollowed, err := s.db.GetFeedsNotFollowedByUser(context.Background(), user.ID)
		if err != nil {
			return err
		}
		if len(feedsNotFollowed) == 0 {
			fmt.Println("No feeds to follow")
			return nil
		}

		// Show followed feeds to remind the user of what they already have
		feedsAlreadyFollowed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
		if err != nil {
			return err
		}
		fmt.Println("Already following:")

		if len(feedsAlreadyFollowed) == 0 {
			fmt.Println("Not following any feeds")
		}

		for _, feed := range feedsAlreadyFollowed {
			printFeed(feed.FeedName, feed.Description, feed.FeedUrl)
		}
		fmt.Println()

		// Create label-value 2D array for the option picker, choose feed to follow
		feedOptions := make([][]string, len(feedsNotFollowed))
		for i := range feedsNotFollowed {
			feedOptions[i] = make([]string, 2)
			feedOptions[i][0] = feedsNotFollowed[i].Name + "\n\t" + feedsNotFollowed[i].Description + "\n\t" + feedsNotFollowed[i].Url
			feedOptions[i][1] = feedsNotFollowed[i].Url
		}

		choice, err := listOptionsReadChoice(feedOptions, "- Choose a new RSS feed to follow:")
		if err != nil {
			return err
		}

		feedUrl = feedsNotFollowed[choice].Url
	}

	// Get feed by url, follow it
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
