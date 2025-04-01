package main

import (
	"fmt"
)

// CLI command
type command struct {
	name string
	args []string
}

type commandDetails struct {
	handlerFunc func(*state, command) error
	info        commandInfo
}

type commandInfo struct {
	description string
	usage       string
	examples    []string
}

// Mapping of commands -> handler functions
type commands struct {
	cmds map[string]commandDetails
}

// Maps a command name to its details, notably its handler function
// Returns an errors if the command with the same name already exists
func (c *commands) register(name string, info commandInfo, f func(*state, command) error) error {
	_, exists := c.cmds[name]
	if exists {
		return fmt.Errorf("command already exists: %v", name)
	}
	c.cmds[name] = commandDetails{
		handlerFunc: f,
		info:        info,
	}

	return nil
}

// Runs the function mapped to the named command
func (c *commands) run(s *state, cmd command) error {
	cmdDetails, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %v", cmd.name)
	}

	return cmdDetails.handlerFunc(s, cmd)
}
