package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/LamontBanks/gator/internal/config"
	"github.com/LamontBanks/gator/internal/database"

	// Leading underscore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

// Application state to be passed to the commands
// Config, database connection, etc.
type state struct {
	config *config.Config
	db     *database.Queries
}

// -- Main

func main() {
	// Initialize info for the application state
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Database
	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Using the SQLC `database` wrapper instead of the native Go SQL db directly
	dbQueries := database.New(db)

	// Set state
	appState := state{
		config: &cfg,
		db:     dbQueries,
	}

	// Register the CLI commands
	appCommands := commands{
		cmds: make(map[string]commandDetails),
	}
	appCommands.register("agg", aggCommandInfo(), handlerAggregator)
	appCommands.register("login", loginCommandInfo(), handlerLogin)
	appCommands.register("register", registerCommandInfo(), handlerRegister)
	appCommands.register("browse", browseCommandInfo(), middlewareLoggedIn(handlerBrowse))
	appCommands.register("feeds", feedsCommandInfo(), handlerGetFeeds)
	appCommands.register("addFeed", addFeedCommandInfo(), middlewareLoggedIn(handlerAddFeed))
	appCommands.register("browseFeed", browseFeedCommandInfo(), middlewareLoggedIn(handlerBrowseFeed))
	appCommands.register("follow", followCommandInfo(), middlewareLoggedIn(handlerFollow))
	appCommands.register("unfollow", unfollowCommandInfo(), middlewareLoggedIn(handlerUnfollow))
	appCommands.register("following", followingCommandInfo(), middlewareLoggedIn(handlerFollowing))
	appCommands.register("users", usersCommandInfo(), handlerGetUsers)
	appCommands.register("help", helpCommandInfo(), appCommands.handlerInfo)

	appCommands.register("reset", commandInfo{}, handlerReset)

	// Read the CLI args to take action
	// os.Args includes the program name, then the command, and (possible) args
	if len(os.Args) < 2 {
		log.Fatal("not enough args provided - need <command> <args>")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Create and run command
	cmdErr := appCommands.run(&appState, command{
		name: cmdName,
		args: cmdArgs,
	})
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

// Wrapper for CLI commands that require the user to be logged in
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	// Return the needed handler function...
	return func(s *state, cmd command) error {

		// ...but first get the user from the database
		username := s.config.CurrentUserName

		if username == "" {
			return fmt.Errorf("no user logged in")
		}

		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("user z%v not registered", username)
		}

		// pass user into the handler function
		return handler(s, cmd, user)
	}
}

// TODO: DEV/TESTING ONLY
// Deletes all users
func handlerReset(s *state, cmd command) error {
	return s.db.Reset(context.Background())
}
