package main

import (
	"context"
	"database/sql"
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
		fmt.Printf("* %v | %v\n%v\n", post.Title, post.FeedName, post.Url)
	}

	return nil
}

// Display recent posts from saved feeds
func handlerBrowseFeed(s *state, cmd command) error {
	// Args: url, max number of feeds (optional)
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <saved rss url> <num of posts, optional>", cmd.name)
	}
	feedUrl := cmd.args[0]

	maxNumPosts := 10
	if len(cmd.args) > 1 {
		i, err := strconv.Atoi(cmd.args[1])
		if err != nil {
			return fmt.Errorf("usage: %v <saved rss url> <num of posts, optional>", cmd.name)
		}
		maxNumPosts = i
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return fmt.Errorf("feed url %v has not been added yet", feedUrl)
	}
	if err != nil {
		return err
	}

	posts, err := s.db.GetRecentPostsFromFeed(context.Background(), database.GetRecentPostsFromFeedParams{
		FeedID: feed.ID,
		Limit:  int32(maxNumPosts),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("\n* %v | %v\n", post.Title, post.FeedName)
	}

	return nil
}
