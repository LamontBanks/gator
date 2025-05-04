/*
Download the latest posts for RSS feeds
*/
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	fuzzytimestamp "github.com/LamontBanks/gator/internal/fuzzy_timestamp"
	rss "github.com/LamontBanks/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all feeds",
	Long: `Update all feeds.
Run without arguments to immediately update all feeds:

	gator update

	Guild Wars 2
			- 1 unread posts
	Nasa Image of the Day
			- 3 unread posts
	Pivot To AI
			- 2 unread posts

Run with a time frequency format <number><seconds | minutes | hours>

	gator update 15m	# 15 minutes
	gator update 30s	# 30 seconds
	gator update 2h		# 2 hours
`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Update preiodically using provided frequency
		if len(args) == 1 {
			freq, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Updating feeds every %v\n", freq)
			ticker := time.NewTicker(freq)
			for ; ; <-ticker.C {
				fmt.Println("Updating RSS feeds...")
				userAuthCall(updateAllFeeds)(appState)
			}
		}

		// Or, do single update
		fmt.Println("Updating feeds...")
		return userAuthCall(updateAllFeeds)(appState)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateAllFeeds(s *state, user database.User) error {
	allFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Update all feeds at once using goroutines
	// And prints the unread feed count (for followed feeds)
	// Arbitrary limit of 10 requests at a time
	feedUpdatedCh := make(chan struct{}, 10)

	followedFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range allFeeds {
		go func() error {
			fmt.Printf("- %v...\n", feed.FeedName)
			updateSingleFeed(s, feed.Url)

			// Filter by followed feeds, get unread posts count
			if slices.ContainsFunc(followedFeeds, func(followedFeed database.GetFeedsForUserRow) bool {
				return followedFeed.FeedUrl == feed.Url
			}) {
				count, msg, err := getUnreadPostInfo(s, user, feed.FeedName, feed.ID)
				if err != nil {
					return err
				}

				if count > 0 {
					fmt.Println(msg)
				}
			}

			feedUpdatedCh <- struct{}{}
			return nil
		}()
	}

	// Wait for all feeds to update
	for range allFeeds {
		<-feedUpdatedCh
	}

	fmt.Printf("Feeds updated at %v\n", time.Now().Format("3:04PM"))
	return nil
}

// Download and save all current posts for the given feed
// Returns true/false if there are additional posts since the last update
func updateSingleFeed(s *state, feedUrl string) error {
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	// Update feed's "last updated" timestamp
	err = s.db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("failed marking %v as updated", feedUrl)
	}

	return saveFeedPosts(s, rssFeed, feed.ID)
}

// Save posts to the database
// Only saves posts newer than the previous latest post
func saveFeedPosts(s *state, rssFeed *rss.RSSFeed, feedId uuid.UUID) error {
	err := s.db.UpdateFeedDescription(context.Background(), database.UpdateFeedDescriptionParams{
		ID:          feedId,
		Description: rssFeed.Channel.Description,
	})
	if err != nil {
		return fmt.Errorf("error updating feed %v description, %v", rssFeed.Channel.Title, err)
	}

	// Save new posts
	lastPostTimestamp, err := s.db.GetLastPostTimestamp(context.Background(), feedId)
	if err == sql.ErrNoRows {
		var earliestTimestamp int64 = 0
		lastPostTimestamp = time.Unix(earliestTimestamp, 0)
	}

	for i := range len(rssFeed.Channel.Item) {
		pubDate, err := rss.ParseRSSPubDate(rssFeed.Channel.Item[i].PubDate)
		if err != nil {
			return err
		}

		// Assumes RSS posts are received in descending order (latest to oldest)
		// So, skip posts we've seen before
		if pubDate.Before(lastPostTimestamp) || pubDate.Equal(lastPostTimestamp) {
			break
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

func getUnreadPostCount(s *state, user database.User, feedID uuid.UUID) (int, error) {
	unreadPosts, err := s.db.GetUnreadPostsForFeed(context.Background(), database.GetUnreadPostsForFeedParams{
		UserID: user.ID,
		FeedID: feedID,
	})
	if err != nil {
		return 0, fmt.Errorf("error getting unread posts, %v", err)
	}

	return len(unreadPosts), nil
}

/*
Returns the unread post count and message
*/
func getUnreadPostInfo(s *state, user database.User, feedName string, feedId uuid.UUID) (int, string, error) {
	unreadPosts, err := s.db.GetUnreadPostsForFeed(context.Background(), database.GetUnreadPostsForFeedParams{
		UserID: user.ID,
		FeedID: feedId,
	})
	if err != nil {
		return 0, "", fmt.Errorf("error getting unread posts, %v", err)
	}

	msg := ""
	if len(unreadPosts) > 0 {
		msg = fmt.Sprintf("%v\n\t- %v unread posts since %v", feedName, len(unreadPosts), fuzzytimestamp.FuzzyTimestamp(unreadPosts[len(unreadPosts)-1].PublishedAt.Local()))
	}

	return len(unreadPosts), msg, nil
}
