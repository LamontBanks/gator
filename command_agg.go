package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// Periodically Downloads all saved RSSFeeds
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

	// Periodic updates
	ticker := time.NewTicker(freq)
	for ; ; <-ticker.C {
		getAllFeedUpdates(s)
	}
}

// Save updates for the most out-of-date feed
func getAllFeedUpdates(s *state) error {
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for range len(allFeeds) {
		oldestFeed, err := s.db.GetNextFeedToFetch(context.Background())
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feeds to update")
		}
		if err != nil {
			return err
		}

		rssFeed, err := fetchFeed(context.Background(), oldestFeed.Url)
		if err != nil {
			return err
		}

		err = s.db.MarkFeedAsFetched(context.Background(), oldestFeed.ID)
		if err != nil {
			return fmt.Errorf("failed marking %v as updated", oldestFeed.Name)
		}

		saveFeedPosts(s, rssFeed, oldestFeed.ID)
	}

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
