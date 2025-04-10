package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// agg commands
var updateFreq string

// aggCmd represents the agg command
var aggCmd = &cobra.Command{
	Use:   "agg",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		freq, err := time.ParseDuration(updateFreq)
		if err != nil {
			return err
		}

		fmt.Printf("Updating feeds every %v\n", updateFreq)

		ticker := time.NewTicker(freq)
		for ; ; <-ticker.C {
			fmt.Println("Updating RSS feeds...")
			updateAllFeeds(appState)
		}
	},
}

func init() {
	rootCmd.AddCommand(aggCmd)

	aggCmd.Flags().StringVarP(&updateFreq, "freq", "f", "15m", "feed update frequency")
}

func updateAllFeeds(s *state) error {
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Update all feeds at once using goroutines
	feedUpdatedCh := make(chan struct{})

	for _, feed := range allFeeds {
		go func() {
			updateSingleFeed(s, feed.Url)
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

func updateSingleFeed(s *state, feedUrl string) error {
	fmt.Printf("Updating %v\n", feedUrl)

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	rssFeed, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// Update feed's "last updated" timestamp
	err = s.db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("failed marking %v as updated", feedUrl)
	}

	saveFeedPosts(s, rssFeed, feed.ID)

	return nil
}

// Save posts to the database
func saveFeedPosts(s *state, rssFeed *RSSFeed, feedId uuid.UUID) error {
	for i := range len(rssFeed.Channel.Item) {
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
