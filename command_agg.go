package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
)

func aggCommandInfo() commandInfo {
	return commandInfo{
		description: "Aggregate all feeds, poll for updates, useful when run in the background with '*'",
		usage:       "agg <update freq> <oldest post time limit, optional>",
		examples: []string{
			"agg 30s &",
			"agg 15m",
			"agg 1h",
			"agg 48h",
			"agg 15m 12h",
			"agg 15m 48h",
		},
	}
}

// Aggregate all RSS feeds
func handlerAggregator(s *state, cmd command) error {
	// Args: <update freq string>, optional: <oldest posts to show>
	// Parse update frequency
	freqFormat := "30s, 1m, 5m, 1h 5h"
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <update freq string, ex: %v, etc.>", cmd.name, freqFormat)
	}
	freq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("%v: could not parse frequency - format is %v", cmd.name, freqFormat)
	}

	// Periodic updates
	ticker := time.NewTicker(freq)
	for ; ; <-ticker.C {
		// fmt.Println("Updating RSS feeds...")
		getAllFeedUpdates(s)
	}
}

// Update all feeds
func getAllFeedUpdates(s *state) error {
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Goroutine to update all feeds at once
	feedUpdatedCh := make(chan struct{})

	for _, feed := range allFeeds {
		go func() {
			getFeedUpdates(s, feed.Url)
			feedUpdatedCh <- struct{}{}
		}()
	}

	// Wait for all feeds to update
	for range allFeeds {
		<-feedUpdatedCh
	}

	fmt.Printf("Feeds updated at %v\n", time.Now().Format("3:04PM"))
	return nil
}

// Update single feed
func getFeedUpdates(s *state, feedUrl string) error {
	// Get feed from DB
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	// Download updates
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("failed marking %v as updated", feedUrl)
	}

	saveFeedPosts(s, rssFeed, feed.ID)

	return nil
}

// Save the posts to the database
func saveFeedPosts(s *state, rssFeed *RSSFeed, feedId uuid.UUID) error {
	for i := range len(rssFeed.Channel.Item) {
		// Convert published data string to time.Time
		pubDate, err := ParseRSSPubDate(rssFeed.Channel.Item[i].PubDate)
		if err != nil {
			return err
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rssFeed.Channel.Item[i].Title,
			Url:         rssFeed.Channel.Item[i].Link,
			Description: rssFeed.Channel.Item[i].Description,
			PublishedAt: pubDate,
			FeedID:      feedId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
