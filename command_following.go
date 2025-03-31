package main

import (
	"context"
	"fmt"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

func followingHelp() commandHelp {
	return commandHelp{
		description: "List all followed feeds",
		usage:       "gator following",
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
		fmt.Println(formatTitleAndLink(feed.FeedName, feed.FeedUrl))
	}

	return nil
}
