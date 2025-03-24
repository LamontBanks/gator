package main

import (
	"fmt"
	"strings"
)

// CLI command
type command struct {
	name string
	args []string
}

// Maps commands -> handler functions
type commands struct {
	cmds map[string]func(*state, command) error
}

// Adds a new CLI command
// Command name is normalized to lowercase.
// Returns an errors if the command with the same name already exists
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
