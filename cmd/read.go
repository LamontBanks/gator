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
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var numReadPosts int
var newPosts bool

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use: "read",

	Short: "Read posts in your feeds",
	Long: `Read posts in your feeds.
A interactive menu will help navigate through your followed feeds, then to the posts within a feed.

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
		if newPosts {
			return userAuthCall(readNewPosts)(appState)
		} else {
			numReadPosts = 3

			if len(args) == 1 {
				i, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}

				numReadPosts = i
			}
			return userAuthCall(readPosts)(appState)
		}
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&newPosts, "new", "n", false, "Read new posts from oldest to latest")
}

// Display option picker for user to select a feed, then select the post to read it's RSS description.
// Some RSS feeds contain the full post test in the description, others have only a snippet.
// This is dependent on the creators of the feed itself, not a limitation of this program.
func readPosts(s *state, user database.User) error {
	userFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		return err
	}

	// Copy feedNames, feedUrl into a label-value 2D slice, pass to the option picker, select the feed
	feedOptions := make([][]string, len(userFeeds))
	for i := range userFeeds {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = userFeeds[i].FeedName
		feedOptions[i][1] = userFeeds[i].FeedUrl
	}

	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}
	feed := userFeeds[choice]
	fmt.Println(feed.FeedName)

	posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
		FeedID: feed.FeedID,
		Limit:  int32(numReadPosts),
	})
	if err != nil {
		return err
	}

	// Copy postTitle, postId into a label-value 2D slice, pass to the option picker, select the post
	postOptions := make([][]string, len(posts))
	for i := range posts {
		postOptions[i] = make([]string, 2)
		postOptions[i][0] = posts[i].Title + "\n\t" + posts[i].PublishedAt.In(time.Local).Format("03:04 PM, Mon, 02 Jan 06")
		postOptions[i][1] = posts[i].ID.String()
	}

	choice, err = listOptionsReadChoice(postOptions, "Choose a post:")
	if err != nil {
		return err
	}

	// Mark as read for user
	if err = markPostAsRead(s, user, posts[choice].FeedID, posts[choice].ID); err != nil {
		return err
	}

	// Display the post
	postText := fmt.Sprintf("%v\n", posts[choice].Title)
	postText += fmt.Sprintf("%v\n\n", posts[choice].PublishedAt.In(time.Local).Format("03:04 PM, Monday, 02 Jan"))
	postText += fmt.Sprintf("%v\n\n", posts[choice].Description)
	postText += fmt.Sprintf("%v\n", posts[choice].Url)
	fmt.Println(postText)

	return nil
}

// Sequential display only the newest posts
func readNewPosts(s *state, user database.User) error {
	userFeeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		return err
	}

	// Copy feedNames, feedUrl into a label-value 2D slice, pass to the option picker, select the feed
	feedOptions := make([][]string, len(userFeeds))
	for i := range userFeeds {
		feedOptions[i] = make([]string, 2)
		feedOptions[i][0] = userFeeds[i].FeedName
		feedOptions[i][1] = userFeeds[i].FeedUrl
	}
	choice, err := listOptionsReadChoice(feedOptions, "Choose a feed:")
	if err != nil {
		return err
	}
	feed := userFeeds[choice]

	unreadPosts, err := s.db.GetUnreadPostsForFeed(context.Background(), database.GetUnreadPostsForFeedParams{
		UserID: user.ID,
		FeedID: feed.FeedID,
	})
	if err != nil {
		return fmt.Errorf("error getting unread posts, %v", err)
	}
	if len(unreadPosts) == 0 {
		fmt.Printf("- No new posts in %v\n", feed.FeedName)
		return nil
	}

	// Posts are returned newest to oldest
	// But we want to read from oldest to newest
	// So, reverse the list
	slices.Reverse(unreadPosts)

	// Get user's choice
	fmt.Println()
	var navChoice string
	navQuit := "q"
	navNext := "n"
	navPrev := "p"

	// Start with first post
	currPostIndex := 0

	// Navigate through posts or exit
	for navChoice != navQuit {
		// Print post, marked as read
		post := unreadPosts[currPostIndex]
		fmt.Println("---")
		postText := fmt.Sprintf("%v\n", post.Title)
		postText += fmt.Sprintf("%v\n\n", post.PublishedAt.In(time.Local).Format("03:04 PM EST, Monday, 02 Jan"))
		postText += fmt.Sprintf("%v\n\n", post.Description)
		postText += fmt.Sprintf("%v\n", post.Url)
		fmt.Println(postText)

		if err := markPostAsRead(s, user, unreadPosts[currPostIndex].FeedID, unreadPosts[currPostIndex].PostID); err != nil {
			return err
		}

		// Display 1-based page numbers
		fmt.Printf("Post %v of %v\n\n", currPostIndex+1, len(unreadPosts))

		// Commands
		fmt.Printf("'%v' - next, '%v' - back, '%v' - quit\n\n", navNext, navPrev, navQuit)

		// Read user command
		_, err = fmt.Scan(&navChoice)
		if err != nil {
			return err
		}
		navChoice = strings.ToLower(navChoice)

		// Navigate through posts
		switch navChoice {
		case navNext:
			if currPostIndex == len(unreadPosts)-1 {
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
	_, err := s.db.GetPostFromUserReadHisory(context.Background(), database.GetPostFromUserReadHisoryParams{
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
