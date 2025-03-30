package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

// Display posts from feeds the logged in use is following
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Optional arg: max number of posts, default 2
	maxNumPosts := 3

	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <Optional: max number of posts>", cmd.name)
		}
		maxNumPosts = i
	}

	feeds, err := s.db.GetFollowedFeeds(context.Background(), database.GetFollowedFeedsParams{
		UserID: user.ID,
		Limit:  int32(maxNumPosts),
	})
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%v | %v\n", feed.FeedName, feed.Title)
	}

	return nil
}
