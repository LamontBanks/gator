/*
 */
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	fuzzytimestamp "github.com/LamontBanks/gator/internal/fuzzy_timestamp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

/*
Format Post data from the database in a general struct for viewing
*/
type PrintablePost struct {
	PostID         uuid.UUID
	FeedID         uuid.UUID
	Title          string
	PublishedAt    time.Time
	FuzzyTimestamp string
	Description    string
	Url            string
}

var numReadPosts int
var sequentialReadFlag bool
var showOnlyNewPostsFlag bool

const UNREADPOSTMARKER = "new"

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use: "read",

	Short: "Read posts in a feed",
	Long: `Read posts in a feed.
A interactive menu will help navigate through followed feeds, then to the posts within a feed.

	gator read
	gator read <number of posts to display, default: 3>

Examples:

	gator

	Choose a feed:
	1: Nasa Image of the Day
	2: Phys.org | Space News
	3: Pivot To AI

	1	# Choice

	Nasa Image of the Day
	Choose a post:
	1: Testing NASA’s IMAP (Interstellar Mapping and Acceleration Probe)
			12:08 PM, Tue, 15 Apr 25
	2: Sculpted by Luminous Stars
			02:23 PM, Mon, 14 Apr 25
	3: Apollo 13 Launch: 55 Years Ago
			11:59 AM, Fri, 11 Apr 25

	2	# Choice

	Sculpted by Luminous Stars
	02:23 PM, Monday, 14 Apr

	This new image showcases the dazzling young star cluster NGC 346. Although both the James Webb Space Telescope...

Currently only a plaintext <description> is readable in the terminal.
Images will not render, HTML will be raw, etc.
The full-text of the post, if any, will have to be viewed in a web browser.
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Display either default or desired number of posts to display
		numReadPosts = 3
		if len(args) == 1 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			numReadPosts = i
		}

		if err := updateIfOutOfDate(appState); err != nil {
			// If unable to update, show the error, but continue reading the feeds
			fmt.Println(err)
		}

		return userAuthCall(readPosts)(appState)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&showOnlyNewPostsFlag, "new", "n", false, "List only new posts")
	readCmd.Flags().BoolVarP(&sequentialReadFlag, "seq", "s", false, "Read displayed posts from oldest to newest")
}

// Display option picker for user to select a feed, then select the post to read the RSS desc field
// Some RSS feeds contain the full post test in the description, others have only a snippet
// This is dependent on the creators of the feed itself, not a limitation of this program
func readPosts(s *state, user database.User) error {
	userFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		return err
	}

	// Make option picker from list of feed names and unread counts
	feedOptions := []string{}
	haveUnreadPosts := false
	for i := range userFeeds {
		unreadCount, unreadCountMsg, err := getUnreadPostInfo(s, user, userFeeds[i].FeedName, userFeeds[i].FeedID)
		if err != nil {
			return err
		}

		feedOptionLabel := userFeeds[i].FeedName
		if unreadCount > 0 {
			haveUnreadPosts = true
			feedOptionLabel += "\n\t- " + unreadCountMsg
		}

		feedOptions = append(feedOptions, feedOptionLabel)
	}

	// If there are no new posts at all, exit
	if showOnlyNewPostsFlag && !haveUnreadPosts {
		fmt.Println("- No new posts for any feeds")
		return nil
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}

	feed := userFeeds[choice]
	fmt.Println(feed.FeedName)

	posts := []PrintablePost{}

	// Show only unread posts...
	if showOnlyNewPostsFlag {
		queryResult, err := s.db.GetUnreadPostsForFeed(context.Background(), database.GetUnreadPostsForFeedParams{
			UserID: user.ID,
			FeedID: feed.FeedID,
		})
		if err != nil {
			return fmt.Errorf("error getting unread posts")
		}

		if len(queryResult) == 0 {
			return nil
		}

		for _, r := range queryResult {
			posts = append(posts, PrintablePost{
				PostID:         r.PostID,
				FeedID:         r.FeedID,
				Title:          r.Title,
				PublishedAt:    r.PublishedAt,
				FuzzyTimestamp: fuzzytimestamp.FuzzyTimestamp(r.PublishedAt),
				Description:    r.Description,
				Url:            r.Url,
			})
		}
		// ...or whatever number of posts the user wants
	} else {
		queryResult, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(numReadPosts),
		})
		if err != nil {
			return err
		}

		for _, r := range queryResult {
			posts = append(posts, PrintablePost{
				PostID:         r.ID,
				FeedID:         r.FeedID,
				Title:          r.Title,
				PublishedAt:    r.PublishedAt,
				FuzzyTimestamp: fuzzytimestamp.FuzzyTimestamp(r.PublishedAt),
				Description:    r.Description,
				Url:            r.Url,
			})
		}
	}

	// Read all posts
	if sequentialReadFlag {
		return readPost(s, user, posts)
	}

	// Make option picker of post titles, adding unread markers as needed
	postOptions := make([]string, len(posts))
	for i, post := range posts {
		unreadpost := ""
		if postIsUnread(post.PostID, post.FeedID, s, user) {
			unreadpost = fmt.Sprintf(" (%v) ", UNREADPOSTMARKER)
		}
		postOptions[i] = fmt.Sprintf("%v\t|%v %v", fuzzytimestamp.FuzzyTimestamp(post.PublishedAt.Local()), unreadpost, post.Title)
	}

	// Choose a post
	choice, err = listOptionsReadChoice(postOptions, "Choose a post:")
	if err != nil {
		return err
	}
	post := posts[choice]

	// Mark as read
	if err = markPostAsRead(s, user, post.FeedID, post.PostID); err != nil {
		return err
	}

	// Read single post
	return readPost(s, user, []PrintablePost{post})
}

// Navigate through each post
func readPost(s *state, user database.User, posts []PrintablePost) error {
	// Posts are returned newest to oldest
	// But we want to read from oldest to newest
	// So, reverse the list
	slices.Reverse(posts)

	// Get user's choice
	fmt.Println()
	var navChoice string
	navQuit := "q"
	navNext := "f"
	navPrev := "b"

	// Start with first post
	currPostIndex := 0

	// Navigate through posts or exit
	for navChoice != navQuit {
		// Mark as read
		if err := markPostAsRead(s, user, posts[currPostIndex].FeedID, posts[currPostIndex].PostID); err != nil {
			return err
		}

		// Print post
		post := posts[currPostIndex]
		fmt.Println("---")
		fmt.Println(formatPost(post))

		// Display 1-based page numbers at bottom of post
		fmt.Printf("Post %v of %v\n\n", currPostIndex+1, len(posts))
		fmt.Println("---")

		// Accept command for navigating between posts
		fmt.Printf("%v - forward, %v - back, %v - quit\n\n", navNext, navPrev, navQuit)

		// Read user navigate command
		_, err := fmt.Scan(&navChoice)
		if err != nil {
			return err
		}
		navChoice = strings.ToLower(navChoice)

		// Navigate through posts
		switch navChoice {
		case navNext:
			if currPostIndex == len(posts)-1 {
				fmt.Println("- Reached end of unread posts")
				continue
			}
			currPostIndex++
		case navPrev:
			if currPostIndex <= 0 {
				fmt.Println("- Reached beginning of unread posts")
				continue
			}
			currPostIndex--
		}
	}

	return nil
}

func markPostAsRead(s *state, user database.User, feedID, postID uuid.UUID) error {
	_, err := s.db.GetPostFromUserReadHistory(context.Background(), database.GetPostFromUserReadHistoryParams{
		UserID: user.ID,
		FeedID: feedID,
		PostID: postID,
	})

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error getting post %v from history, %v", postID, err)
	}

	// Save post to history if not present
	if err == sql.ErrNoRows {
		_, err = s.db.CreatePostInUserReadHistory(context.Background(), database.CreatePostInUserReadHistoryParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			UserID:       user.ID,
			FeedID:       feedID,
			PostID:       postID,
			HasViewed:    false,
			IsBookmarked: false,
		})

		if err != nil {
			return fmt.Errorf("error saving post %v to history, %v", postID, err)
		}
	}

	// Mark as read
	err = s.db.MarkPostAsViewed(context.Background(), database.MarkPostAsViewedParams{
		UserID: user.ID,
		FeedID: feedID,
		PostID: postID,
	})
	if err != nil {
		return fmt.Errorf("error marking post %v as viewed %v", postID, err)
	}

	return nil
}

func formatPost(post PrintablePost) string {
	postText := fmt.Sprintf("%v\n", post.Title)
	postText += fmt.Sprintf("%v\n\n", post.PublishedAt.In(time.Local).Format("03:04 PM EST, Monday, 02 Jan 2006"))
	postText += fmt.Sprintf("%v\n\n", post.Description)
	postText += fmt.Sprintf("Full post (opens browser):\n$ open \"%v\"\n", post.Url)
	return postText
}

func postIsUnread(postId, feedId uuid.UUID, s *state, user database.User) bool {
	post, err := s.db.GetPostFromUserReadHistory(context.Background(), database.GetPostFromUserReadHistoryParams{
		UserID: user.ID,
		FeedID: feedId,
		PostID: postId,
	})

	// Nothing returned or any other error is considered unread
	if err != nil {
		return true
	}

	// Negate column
	return !post.HasViewed
}
