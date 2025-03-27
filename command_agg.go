package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strconv"
)

// Downloads and prints an RSS
func handlerAggregator(s *state, cmd command) error {
	// Args: RSS feed url
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}
	feedUrl := cmd.args[0]

	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)

	return nil
}

// Continuously fetch updates for all feeds
func getAllFeedUpdates(s *state, cmd command) error {
	// Optional args: number of recent posts
	maxNumPosts := math.MaxInt
	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		maxNumPosts = i
		if err != nil {
			return fmt.Errorf("usage: %v <max number of posts> (optional)", cmd.name)
		}
	}

	// Get the oldest feed info
	oldestFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err == sql.ErrNoRows {
		return fmt.Errorf("no feeds to update")
	}
	if err != nil {
		return err
	}

	// Download the latest feed
	rssFeed, err := fetchFeed(context.Background(), oldestFeed.Url)
	if err != nil {
		return err
	}

	// Mark the feed as updated
	err = s.db.MarkFeedAsFetched(context.Background(), oldestFeed.ID)
	if err != nil {
		return fmt.Errorf("failed marking %v as updated", oldestFeed.Name)
	}

	// Print the feed
	maxNumPosts = min(len(rssFeed.Channel.Item), maxNumPosts)

	fmt.Printf("- %v -\n", rssFeed.Channel.Title)
	for i := range maxNumPosts {
		fmt.Println(rssFeed.Channel.Item[i].Title)
		// fmt.Println(rssFeed.Channel.Item[i].Description)
	}
	return nil
}
