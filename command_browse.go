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

	posts, err := s.db.GetPostsFromFollowedFeeds(context.Background(), database.GetPostsFromFollowedFeedsParams{
		UserID: user.ID,
		Limit:  int32(maxNumPosts),
	})
	if err != nil {
		return err
	}

	allPosts := groupPosts(posts)

	for feedName, postTitles := range allPosts {
		fmt.Printf("-- %v --\n", feedName)
		for _, title := range postTitles {
			fmt.Printf("* %v\n", title)
		}
	}

	return nil
}

func groupPosts(posts []database.GetPostsFromFollowedFeedsRow) map[string][]string {
	groupedPosts := make(map[string][]string)

	for _, post := range posts {
		_, exists := groupedPosts[post.FeedName]
		if !exists {
			groupedPosts[post.FeedName] = []string{}
		}
		groupedPosts[post.FeedName] = append(groupedPosts[post.FeedName], post.Title)
	}

	return groupedPosts
}
