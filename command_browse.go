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
		description: "Show latest posts for current user's feeds",
		usage:       "browse",
		examples: []string{
			"browse",
			"browse 5",
		},
	}
}

// Display most recent posts from user's feeds
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Args: <max number of posts, optional, default 3> <'-h' to only show titles, optional>
	maxNumPosts := 3
	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf(browseCommandInfo().usage)
		}
		maxNumPosts = i
	}

	// Get feeds
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	// Pull posts for each feed
	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(maxNumPosts),
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("\n%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				printPostTitle(post.Title)
			}
		} else {
			fmt.Println("No posts")
		}

	}

	return nil
}

func browseFeedCommandInfo() commandInfo {
	return commandInfo{
		description: "Read posts for the given feed URL",
		usage:       "browseFeed <feed url> <number of posts, default: 5>",
		examples: []string{
			"browseFeed http://example.com/rss/feed",
			"browseFeed http://example.com/rss/feed 5",
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

	maxNumPosts := 3
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
	s := fmt.Sprintf("- %v | %v\n", title, published_at.Format("Mon, 02 Jan 03:04 PM"))
	s += fmt.Sprintf("  %v\n", desc)
	s += fmt.Sprintf("  %v\n", link)

	fmt.Println(s)
}

func printPostTitle(title string) {
	fmt.Printf("- %v\n", title)
}
