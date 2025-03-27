package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Downloads and prints an RSS
func handlerAggregator(s *state, cmd command) error {
	// Args: <update freq string, ex: 1s, 30s, 1m, 5m, 1h>
	freqFormat := "1s, 30s, 1m, 5m, 1h, 1d"
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <update freq string, ex: %v, etc.>", cmd.name, freqFormat)
	}

	freq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not parse frequency - format is %v", freqFormat)
	}

	ticker := time.NewTicker(freq)
	for ; ; <-ticker.C {
		getAllFeedUpdates(s)
	}
}

// Continuously fetch updates for all feeds
func getAllFeedUpdates(s *state) error {
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
	fmt.Printf("- %v -\n", rssFeed.Channel.Title)

	maxItems := 3
	for i := range maxItems {
		fmt.Println(rssFeed.Channel.Item[i].Title)
		// fmt.Println(rssFeed.Channel.Item[i].Description)
	}
	return nil
}
