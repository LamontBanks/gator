package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/config"
	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"

	// Undercore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

// Application state to be passed to the commands:
// Config, database connection, etc.
type state struct {
	config *config.Config
	db     *database.Queries
}

// CLI command
type command struct {
	name string
	args []string
}

// Maps commands -> handler functions
type commands struct {
	cmds map[string]func(*state, command) error
}

// --- Main

func main() {
	// Register the CLI commands
	appCommands := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	appCommands.register("login", handlerLogin)
	appCommands.register("register", handlerRegister)

	// Initialize the application state
	// - Config
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// - Database connection (from config)
	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	dbQueries := database.New(db) // Use the SQLC wrapper database instead of the SQL db directly
	appState := state{
		config: &cfg,
		db:     dbQueries,
	}

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

// CLI Command Handlers
// Save the username to the config file
// Usage:
//
//	$ go run . login <username>
//	$ go run . login alice
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username required")
	}

	// Save username to the config file
	username := cmd.args[0]
	s.config.CurrentUserName = username

	if err := s.config.SetConfig(); err != nil {
		return err
	}

	fmt.Printf("Logged in as %v\n", username)

	return nil
}

// Register a user in on the server. Updates the config with the user.
// Usage:
//
//	$ go run . register <username>
//	$ go run . register alice
func handlerRegister(s *state, cmd command) error {
	// Get name from args
	if len(cmd.args) < 1 {
		return fmt.Errorf("name required")
	}
	name := cmd.args[0]

	// Insert user
	queryResult, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created user %v: %v\n", name, queryResult)

	// Update the config as well
	return handlerLogin(s, cmd)
}

// --- Command functions

// Adds a new CLI command
// Command name is normalized to lowercase.
// Returns an errors if the command with the same name already exists
// Requires:
// - Name of the command
// - Function that accepts the application state and command details containing program args
func (c *commands) register(name string, f func(*state, command) error) error {
	name = strings.ToLower(name)

	_, exists := c.cmds[name]
	if exists {
		return fmt.Errorf("command already exists: %v", name)
	}

	c.cmds[name] = f

	return nil
}

// Runs the function mapped to the named command
func (c *commands) run(s *state, cmd command) error {
	// Search the mapping for the assoicated handler function
	handlerFunc, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %v", cmd.name)
	}

	return handlerFunc(s, cmd)
}
