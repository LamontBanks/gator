package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/LamontBanks/blog-aggregator/internal/config"
	"github.com/LamontBanks/blog-aggregator/internal/database"

	// Leading underscore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

// Application state to be passed to the commands:
// Config, database connection, etc.
type state struct {
	config *config.Config
	db     *database.Queries
}

// --- Main

func main() {
	// Initialize info for the application state
	// Config
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Database connection
	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dbQueries := database.New(db) // Use the SQLC database wrapper instead of the SQL db directly

	// Set state
	appState := state{
		config: &cfg,
		db:     dbQueries,
	}

	// Register the CLI commands
	appCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)
	appCommands.register("reset", handlerReset)
	appCommands.register("users", handlerGetUsers)
	appCommands.register("agg", handlerAggregator)
	appCommands.register("addFeed", handlerAddFeed)

	// Read the CLI args to take action
	// os.Args includes the program name, then the command, and (possibly) args
	if len(os.Args) < 2 {
		log.Fatal("not enough args provided - need <command> <args>")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Run command
	cmdErr := appCommands.run(&appState, command{
		name: cmdName,
		args: cmdArgs,
	})
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

// DEV/TESTING ONLY
// Deletes all users
func handlerReset(s *state, cmd command) error {
	return s.db.Reset(context.Background())
}
