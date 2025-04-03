package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func addFeedCommandInfo() commandInfo {
	return commandInfo{
		description: "Add a new feed for all users to follow",
		usage:       "addFeed <feed name> <unique RSS feed>",
		examples: []string{
			"addFeed \"Example Feed Name\" http://example.com/rss/feed",
			"addFeed \"Example Feed Name\" http://example.com/rss.xml",
		},
	}
}

// Create a feed, available for all users to follow
// Fails if the feed already exists
func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Args: feedName, feedUrl
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %v <feed name> <url>", cmd.name)
	}
	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	// Attempt to download the feed before saving it
	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("failed to add feed: %v", err)
	}

	// If feed is missing any required elements, also don't save it
	// https://www.rssboard.org/rss-specification#requiredChannelElements
	if feed.Channel.Title == "" {
		return fmt.Errorf("did not save feed %v, missing required RSS fields:\nTitle: %v\nSee https://www.rssboard.org/rss-specification#requiredChannelElements", feedUrl, feed.Channel.Title)
	}

	// Save feed info
	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        feedName,
		Url:         feedUrl,
		UserID:      user.ID,
		Description: feed.Channel.Description,
	})
	if err != nil {
		return fmt.Errorf("could not add: %v (%v) for %v - possible duplicate feed?", feedName, feedUrl, user.Name)
	}
	fmt.Printf("Saved \"%v\" (%v) for user %v\n", newFeed.Name, newFeed.Url, user.Name)

	// Make user follow the new feed
	newFeedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed %v", newFeedFollow.FeedName)
	}
	fmt.Printf("%v followed %v\n", newFeedFollow.UserName, newFeedFollow.FeedName)

	// Download updates
	getFeedUpdates(s, newFeed.Url)

	return nil
}

func feedsCommandInfo() commandInfo {
	return commandInfo{
		description: "List all feeds",
		usage:       "feeds",
		examples:    []string{},
	}
}

// Prints all feeds
func handlerGetFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Sort by Feed name
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].FeedName < feeds[j].FeedName
	})

	fmt.Println("Available RSS Feeds to follow:")
	for _, feed := range feeds {
		printFeed(feed.FeedName, feed.Description, feed.Url)
	}

	return nil
}

func printFeed(title, desc, link string) {
	s := fmt.Sprintf("%v\n\t%v\n\t%v", title, desc, link)

	fmt.Println(s)
}
