package main

import (
	"context"
	"fmt"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

func followingCommandInfo() commandInfo {
	return commandInfo{
		description: "List all followed RSS feeds",
		usage:       "following",
		examples:    []string{},
	}
}

// Prints details of all the feeds the current user is following
func handlerFollowing(s *state, cmd command, user database.User) error {
	feedDetails, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	if len(feedDetails) == 0 {
		fmt.Println("you are not following any feeds")
		return nil
	}

	for _, feed := range feedDetails {
		title := fmt.Sprintf("%v | %v", feed.FeedName, feed.Description)
		fmt.Println(formatTitleAndLink(title, feed.FeedUrl))
	}

	return nil
}
