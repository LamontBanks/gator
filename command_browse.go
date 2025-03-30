package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

// Display most recent posts from user's feeds
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Optional arg: max number of posts, default 3
	maxNumPosts := 3

	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("usage: %v <Optional: max number of posts>", cmd.name)
		}
		maxNumPosts = i
	}

	posts, err := s.db.GetFollowedPosts(context.Background(), database.GetFollowedPostsParams{
		UserID: user.ID,
		Limit:  int32(maxNumPosts),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("* %v | %v | %v\n", post.FeedName, post.Title, post.Url)
	}

	return nil
}
