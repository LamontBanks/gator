package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
)

func browseCommandInfo() commandInfo {
	return commandInfo{
		description: "Show recent posts for feeds followed by the current user",
		usage:       "browse <max number of posts per feed, default: 10>",
		examples: []string{
			"browse",
			"browse 5",
		},
	}
}

// Display most recent posts from user's feeds
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Args: <max number of posts, optional, default 3>
	maxNumPosts := 3
	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf(browseCommandInfo().usage)
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
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(maxNumPosts),
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				printPost(post.Title, post.Description, post.Url, post.PublishedAt)
			}
		} else {
			fmt.Println("No posts")
		}
	}

	return nil
}

func browseFeedCommandInfo() commandInfo {
	return commandInfo{
		description: "Show posts for an RSS feed URL",
		usage:       "browseFeed <RSS feed URL> <optional: max number of posts, default: 10>",
		examples: []string{
			"browseFeed http://example.com/rss/feed",
			"browseFeed http://example.com/rss/feed 5",
			"browseFeed http://example.com/rss/feed 10",
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
		return fmt.Errorf("failed to browseFeed %v - not yet added", feedUrl)
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
		printPost(post.Title, post.Url, post.Description, post.PublishedAt)
	}

	return nil
}

func printPost(title, desc, link string, published_at time.Time) {
	s := fmt.Sprintf("- %v\n", title)
	s += fmt.Sprintf("  %v\n", published_at.Format("3:04 PM, Monday 02 Jan 06"))
	s += fmt.Sprintf("  %v\n", link)
	s += fmt.Sprintf("  %v\n", desc)

	fmt.Println(s)
}
