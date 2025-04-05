package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/LamontBanks/gator/internal/database"
)

func browseCommandInfo() commandInfo {
	return commandInfo{
		description: "Show latest posts for current user's feeds",
		usage:       "browse <max number of posts per feed>",
		examples: []string{
			"browse",
			"browse 5",
		},
	}
}

// Display most recent posts from user's feeds
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Args: <max number of posts per feed, optional, default 3>
	maxNumPosts := 3
	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf(browseCommandInfo().usage)
		}
		maxNumPosts = i
	}

	// Get feeds
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	// Pull posts for each feed
	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(maxNumPosts),
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				printPostTitle(post.Title)
			}
		} else {
			fmt.Println("No posts")
		}

	}

	return nil
}

func browseFeedCommandInfo() commandInfo {
	return commandInfo{
		description: "Read posts from a followed feed",
		usage:       "browseFeed <max number of posts, optional, default 3>",
		examples: []string{
			"browseFeed",
			"browseFeed 10",
		},
	}
}

// Display menus to view specific posts in specific feeds
func handlerBrowseFeed(s *state, cmd command, user database.User) error {
	maxNumPosts := 3

	// Args: Max number of post, optional
	if len(cmd.args) > 0 {
		i, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf(browseFeedCommandInfo().usage)
		}
		maxNumPosts = i
	}

	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("you're not following any feeds")
	}
	if err != nil {
		fmt.Println("error retrieving your followed feeds")
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
	posts, err := s.db.GetRecentPostsFromFeed(context.Background(), database.GetRecentPostsFromFeedParams{
		FeedID: userFeeds[choice].FeedID,
		Limit:  int32(maxNumPosts),
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

	printPost(posts[choice].Title, posts[choice].Description, posts[choice].Url, posts[choice].PublishedAt)

	return nil
}

func printPost(title, desc, link string, published_at time.Time) {
	s := fmt.Sprintf("%v\n", title)
	s += fmt.Sprintf("%v\n\n", published_at.In(time.Local).Format("03:04 PM, Mon, 02 Jan 06"))
	s += fmt.Sprintf("%v\n\n", desc)
	s += fmt.Sprintf("%v\n", link)

	fmt.Println(s)
}

func printPostTitle(title string) {
	fmt.Printf("\t- %v\n", title)
}
