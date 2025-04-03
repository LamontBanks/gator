package main

import (
	"fmt"
	"slices"
)

func helpCommandInfo() commandInfo {
	return commandInfo{
		description: "Display this help message",
		usage:       "help <command>",
		examples: []string{
			"help",
			"help agg",
			"help addFeed",
		},
	}
}

// Prints info for the program commands
func (c *commands) handlerInfo(s *state, cmd command) error {
	programName := "gator"

	// Individual command help
	if len(cmd.args) > 0 {
		commandName := cmd.args[0]
		fmt.Printf("usage: %v %v\n", programName, c.cmds[commandName].info.usage)
		fmt.Println()
		fmt.Printf("%v\n", c.cmds[commandName].info.description)

		// Print examples, if there are any
		if len(c.cmds[commandName].info.examples) > 0 {
			fmt.Println()
			fmt.Println("Examples:")
			for _, example := range c.cmds[commandName].info.examples {
				fmt.Printf("\t%v %v\n", programName, example)
			}
		}
		// All command help
	} else {
		fmt.Printf("%v is a tool for viewing RSS feeds in the console.\n\n", programName)

		fmt.Printf("Usage:\n\n")
		fmt.Printf("\t%v <command> [arguments]\n\n", programName)

		// Sort commands alphabetically before print
		sortedCommandNames := []string{}
		for key := range c.cmds {
			sortedCommandNames = append(sortedCommandNames, key)
		}
		slices.Sort(sortedCommandNames)

		for _, cmd := range sortedCommandNames {
			fmt.Printf("\t%v\t\t%v\n", cmd, c.cmds[cmd].info.description)
		}

		fmt.Printf("Use \"%v help <command>\" for more information about a command.\n", programName)
	}

	return nil
}
