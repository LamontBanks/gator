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
	maxNumPosts := 2

	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <Optional: max number of posts>", cmd.name)
		}
		maxNumPosts = i
	}

	posts, err := s.db.GetPostsFromFollowedFeeds(context.Background(), database.GetPostsFromFollowedFeedsParams{
		UserID: user.ID,
		Limit:  int32(maxNumPosts),
	})
	if err != nil {
		return err
	}

	for _, p := range posts {
		fmt.Printf("%v | %v | Published: %v\n", p.FeedName, p.Title, p.PublishedAt)
	}

	return nil
}
