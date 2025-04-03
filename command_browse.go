package main

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func browseCommandInfo() commandInfo {
	return commandInfo{
		description: "Show latest posts for current user's feeds",
		usage:       "browse <max number of posts per feed>",
		examples: []string{
			"browse",
			"browse 5",
		},
	}
}

// Display most recent posts from user's feeds
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Args: <max number of posts per feed, optional, default 3>
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
		description: "Read posts from a followed feed",
		usage:       "browseFeed <max number of posts>",
		examples: []string{
			"browseFeed",
			"browseFeed 10",
		},
	}
}

// Display recent posts from saved feeds
func handlerBrowseFeed(s *state, cmd command, user database.User) error {
	maxNumPosts := 10

	// Get user feeds
	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		fmt.Println("error retrieving your followed feeds")
		return err
	}

	// Copy feedNames+feedUrl into a slice, sort alphabetically, pass to menu generator
	feedOptions := make([][]string, len(userFeeds))
	for i := range userFeeds {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = userFeeds[i].FeedName
		feedOptions[i][1] = userFeeds[i].FeedUrl
	}
	sort.Slice(feedOptions, func(i, j int) bool {
		return feedOptions[i][0] < feedOptions[j][0]
	})

	_, feedUrl, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}

	// Get the feed, then the posts
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

	// Copy postTitle+postId into a 2D slice, pass to menu generator
	postOptions := make([][]string, len(posts))
	for i := range posts {
		postOptions[i] = make([]string, 2)
		postOptions[i][0] = posts[i].Title
		postOptions[i][1] = posts[i].ID.String()
	}

	_, postIdString, err := listOptionsReadChoice(postOptions, "Choose a post:")
	if err != nil {
		return err
	}

	// Get, print post
	postUUID, err := uuid.Parse(postIdString)
	if err != nil {
		return err
	}

	post, err := s.db.GetPostById(context.Background(), postUUID)
	if err != nil {
		return nil
	}
	printPost(post.Title, post.Url, post.Description, post.PublishedAt)

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
