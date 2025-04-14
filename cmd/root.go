// Default 'gator' command will print RSS feeds the user is following
// It can also be used to see all feed available to follw.
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/LamontBanks/gator/internal/config"
	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"

	// Leading underscore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

var (
	// Used by gator
	appState  *state
	db        *sql.DB
	resetFlag bool // TODO: DEV ONLY

	// Command flags/parameters
	showAllFeeds    bool
	numPostsPerFeed int

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gator",
		Short: "Gator is a terminal-based RSS reader",
		Long: `Gator is a terminal-based RSS reader.
It is best ran as a terminal background process (ex: gator ... &).
Then, interact with the tool to read and manage RSS feeds.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if resetFlag {
				return reset()
			} else {
				if showAllFeeds {
					return printAllFeeds(appState)
				} else {
					return userAuthCall(printFollowedFeeds)(appState)
				}
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initAppState)
	cobra.OnFinalize(closeDB) // Closing database using Cobra instead of the usual "defer ..."

	rootCmd.Flags().BoolVarP(&showAllFeeds, "all", "a", false, "Show all feeds added to gator")
	rootCmd.Flags().IntVarP(&numPostsPerFeed, "numPosts", "n", 2, "maximum number of posts per feed")

	rootCmd.Flags().BoolVar(&resetFlag, "reset", false, "Deletes all users, effectively clearing the database (DEV ONLY)")
}

// Initialize info for the application state
func initAppState() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Database
	connStr := cfg.DbUrl
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Using the SQLC database wrapper instead of the native Go SQL db directly
	dbQueries := database.New(db)

	// Set state
	appState = &state{
		config: &cfg,
		db:     dbQueries,
	}
}

func closeDB() {
	db.Close()
}

// Return the user's choice from a 2D slice of labels-values
// Ex:
//
//	labelsValues := [][]string{
//		{"Label 1", {"Value 1"},
//		{"Label 2", {"Value 2"},
//		...
//	}
//
// Returns:
//
//	int - the index of the choice
//	error - if choice out of range
func listOptionsReadChoice(labelsValues [][]string, message string) (int, error) {
	fmt.Println(message)

	// List options, start index with "1"; easier to select than "0" for choosing the first option (the most common case)
	for i, lblVal := range labelsValues {
		fmt.Printf("%v: %v\n", i+1, lblVal[0])
	}

	// Get user's choice
	fmt.Println()
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		return 0, err
	}

	// Normalize to 0-based indexing
	choice -= 1
	if choice < 0 || choice >= len(labelsValues) {
		return 0, errors.New("invalid choice")
	}

	// Return
	return choice, nil
}

// A closure for wrapping functions requiring a user
func userAuthCall(f func(s *state, user database.User) error) func(*state) error {
	return func(s *state) error {
		username := s.config.CurrentUserName

		if username == "" {
			return fmt.Errorf("logged in user required for this command")
		}

		u, err := s.db.GetUser(context.Background(), username)
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %v not registered", username)
		}
		if err != nil {
			return err
		}

		return f(s, u)
	}

}

func printFollowedFeeds(s *state, user database.User) error {
	if numPostsPerFeed < 0 {
		return fmt.Errorf("number of posts must be >= 0")
	}

	// Get feeds followed by user
	feeds, err := s.db.GetFeedsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("- Not following any feeds")
		return nil
	}

	// Pull posts for each feed
	fmt.Println("Your Feeds:")
	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.FeedID,
			Limit:  int32(numPostsPerFeed),
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.FeedUrl)
		if len(posts) > 0 {
			for _, post := range posts {
				fmt.Printf("\t- %v\n", post.Title)
			}
		} else {
			fmt.Println("\t- No recent posts")
		}
		fmt.Println()
	}

	return nil
}

func printAllFeeds(s *state) error {
	if numPostsPerFeed < 0 {
		return fmt.Errorf("number of posts must be >= 0")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Print posts for each feed
	fmt.Println("All RSS Feeds:")
	if len(feeds) == 0 {
		fmt.Println("- No feeds added")
		return nil
	}

	for _, feed := range feeds {
		posts, err := s.db.GetPostsFromFeed(context.Background(), database.GetPostsFromFeedParams{
			FeedID: feed.ID,
			Limit:  int32(numPostsPerFeed),
		})
		if err != nil {
			return err
		}

		// Print feeds, posts
		fmt.Printf("%v | %v\n", feed.FeedName, feed.Url)
		if len(posts) > 0 {
			for _, post := range posts {
				fmt.Printf("\t- %v\n", post.Title)
			}
		} else {
			fmt.Println("\t- No recent posts")
		}
		fmt.Println()
	}

	return nil
}

// DEV ONLY - Delete all users
func reset() error {
	return appState.db.Reset(context.Background())
}
