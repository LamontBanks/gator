package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LamontBanks/gator/internal/database"
)

func unfollowCommandInfo() commandInfo {
	return commandInfo{
		description: "Stop following a feed",
		usage:       "unfollow <feed url, optional>",
		examples: []string{
			"unfollow\n<choose from list of feeds>",
			"unfollow http://example.com/rss/feed",
		},
	}
}

// Unfollows the given RSS feel URL
func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Args: url, optional
	var feedUrl string
	if len(cmd.args) > 0 {
		feedUrl = cmd.args[0]
	} else {
		// If no URL is provided, choose from list of followed feeds
		followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
		if err != nil {
			return err
		}

		if len(followedFeeds) == 0 {
			fmt.Println("not following any feeds")
			return nil
		}

		// Create label-value 2D array for the option picker
		feedOptions := make([][]string, len(followedFeeds))
		for i := range followedFeeds {
			feedOptions[i] = make([]string, 2)
			feedOptions[i][0] = followedFeeds[i].FeedName + "\n\t" + followedFeeds[i].Description + "\n\t" + followedFeeds[i].FeedUrl
			feedOptions[i][1] = followedFeeds[i].FeedUrl
		}

		// Choose feed to unfollow
		choice, err := listOptionsReadChoice(feedOptions, "- Choose an RSS feed to unfollow")
		if err != nil {
			return err
		}

		feedUrl = followedFeeds[choice].FeedUrl
	}

	// Get all feed info from the url
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("failed to unfollow %v - not yet added", feedUrl)
	}
	if err != nil {
		return err
	}

	// Unfollow the feed
	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow %v", feedUrl)
	}

	fmt.Printf("Unfollowed %v | %v\n", feed.Name, feedUrl)
	return nil
}
