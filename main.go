package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LamontBanks/blog-aggregator/internal/config"
)

// Application state to be passed to the commands:
// Config, database connection, etc.
type state struct {
	config *config.Config
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

	// Initialize the application state
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	appState := state{
		config: &cfg,
	}

	// Read the CLI args to take action
	// os.Args includes the program name, then the command, and (possibly) args
	if len(os.Args) < 2 {
		log.Fatal("not enough args provided - need <command> <args>")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmdErr := appCommands.run(&appState, command{
		name: cmdName,
		args: cmdArgs,
	})

	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

// CLI Command Handlers

func handlerLogin(s *state, cmd command) error {
	// Usage:
	// $ app login <username>
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

// --- Command functions

// Adds a new CLI command to the app
// Names are normalized to lowercase. Errors if the command with the same name already exists ()
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
