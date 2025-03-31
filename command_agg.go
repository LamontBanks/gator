package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func aggHelp() commandHelp {
	return commandHelp{
		description: "Aggregate all RRS feeds",
		usage:       "gator agg <update freq> <oldest post time limit>",
		examples: []string{
			"gator agg 30s",
			"gator agg 15m",
			"gator agg 1h",
			"gator agg 48h",
			"gator agg 15m 12h",
			"gator agg 15m 48h",
		},
	}
}

// Aggregate all RSS feeds
func handlerAggregator(s *state, cmd command) error {
	// Args: <update freq string>, optional: <oldest posts to show>
	// Parse update frequency
	freqFormat := "30s, 1m, 5m, 1h, 1d"
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <update freq string, ex: %v, etc.>", cmd.name, freqFormat)
	}
	freq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("%v: could not parse frequency - format is %v", cmd.name, freqFormat)
	}

	// Optional: Parse oldest post time limit
	oldestPostTime, _ := time.ParseDuration("18h") // Skipping error handling
	if len(cmd.args) > 2 {
		t, err := time.ParseDuration(cmd.args[1])
		if err != nil {
			fmt.Printf("error parsing %v, format: 1h, 3h, 24h, 36h, etc.", cmd.args[1])
		} else {
			oldestPostTime = t
		}
	}

	// Periodic updates
	ticker := time.NewTicker(freq)
	for ; ; <-ticker.C {
		fmt.Println("Updating RSS feeds...")
		getAllFeedUpdates(s, oldestPostTime)
	}
}

// Save updates all feeds
func getAllFeedUpdates(s *state, oldestPostTime time.Duration) error {
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for range len(allFeeds) {
		// Start with the most out-of-date feed
		oldestFeed, err := s.db.GetNextFeedToFetch(context.Background())
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

		// Pull posts for each feed, but only within the time window
		// Calculate time limit

		withinTimeLimit := time.Now().Add(-1 * oldestPostTime)

		posts, err := s.db.GetRecentPostsFromFeed(context.Background(), database.GetRecentPostsFromFeedParams{
			FeedID:      oldestFeed.ID,
			PublishedAt: withinTimeLimit,
			Limit:       3,
		})
		if err != nil {
			return err
		}

		fmt.Printf("\n%v | %v\n", oldestFeed.Name, oldestFeed.Url)
		if len(posts) == 0 {
			fmt.Printf("* Nothing in the last %v\n", oldestPostTime)
		}
		for _, post := range posts {
			fmt.Printf("* %v\n", post.Title)
			fmt.Printf("  %v\n", post.Url)
		}
	}

	return nil
}

// Save the posts to the database
func saveFeedPosts(s *state, rssFeed *RSSFeed, feedId uuid.UUID) error {
	for i := range len(rssFeed.Channel.Item) {
		// Convert published data string to time.Time
		pubDate, err := ParseRSSPubDate(rssFeed.Channel.Item[i].PubDate)
		if err != nil {
			fmt.Println(err)
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
