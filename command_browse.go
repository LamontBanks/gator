package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

func browseHelp() commandInfo {
	return commandInfo{
		description: "Show recent posts from current user's followed feed",
		usage:       "browse <max number of posts per feed, default: 3>",
		examples: []string{
			"browse",
			"browse 5",
			"browse 10",
		},
	}
}

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

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	// Pull posts for each feed, but only within the time window
	// and no more posts per feed than specified
	for _, feed := range feeds {
		// Calculate time limit
		timeLimit := "18h"
		t, err := time.ParseDuration(timeLimit)
		if err != nil {
			return err
		}
		withinTimeLimit := time.Now().Add(-1 * t)

		posts, err := s.db.GetRecentPostsFromFeed(context.Background(), database.GetRecentPostsFromFeedParams{
			FeedID:      feed.FeedID,
			PublishedAt: withinTimeLimit,
			Limit:       int32(maxNumPosts),
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				fmt.Println(formatTitleAndLink(post.Title, post.Url))
			}
		} else {
			emptyPost := fmt.Sprintf("Nothing in the last %v", timeLimit)
			fmt.Println(formatTitleAndLink(emptyPost, ""))
		}
	}

	return nil
}

func browseFeedHelp() commandInfo {
	return commandInfo{
		description: "Lists ",
		usage:       "browseFeed <registered RSS feed URL> <optional: >",
		examples: []string{
			"browseFeed http://example.com/rss/feed",
		},
	}
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
		return fmt.Errorf("failed to browser %v - not yet added", feedUrl)
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
		fmt.Printf("%v | %v\n", post.FeedName, post.Title)
	}

	return nil
}
