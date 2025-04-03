package main

import (
	"errors"
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

// Return the user's choice from a 2D slice of labels-values
// Ex:
//
//	labelsValues := [][]string{
//		{"Label 1", {"Value 1"},
//		{"Label 2", {"Value 2"},
//		...
//	}
func listOptionsReadChoice(labelsValues [][]string, message string) (string, string, error) {
	fmt.Println(message)

	// List options, start index with "1"; easier to select than "0" for choosing the first option (the most common case)
	for i, lblVal := range labelsValues {
		fmt.Printf("%v: %v\n\n", i+1, lblVal[0])
	}

	// Get user's choice
	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		return "", "", err
	}

	// Normalize to 0-based indexing
	choice -= 1
	if choice < 0 || choice >= len(labelsValues) {
		return "", "", errors.New("invalid choice")
	}

	// Return
	return labelsValues[choice][0], labelsValues[choice][1], nil
}
