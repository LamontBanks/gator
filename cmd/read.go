/*
 */
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read posts in your feeds",
	Long: `Read posts in your feeds.
A interactive menu will help navigate through your followed feeds, then to the posts within a feed.

	gator read

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
	RunE: func(cmd *cobra.Command, args []string) error {
		return userAuthCall(readPosts)(appState)
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
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

	// Copy feedNames, feedUrl into a label-value 2D slive, pass to the option picker, select the feed
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
	fmt.Println(userFeeds[choice].FeedName)

	// Get posts for the selected feed
	posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
		FeedID: userFeeds[choice].FeedID,
		Limit:  int32(3),
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
