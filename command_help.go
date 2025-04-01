package main

import (
	"fmt"
	"slices"
)

func helpCommandInfo() commandInfo {
	return commandInfo{
		description: "Display this help message, 'help COMMAND' for command help",
		usage:       "help",
		examples: []string{
			"help",
			"help agg",
			"help addFeed",
		},
	}
}

// Prints help info for the program commands
// Command-specific help:
//
//	$ gator help COMMAND
//
// All commands:
//
//	$ gator help
func (c *commands) handlerInfo(s *state, cmd command) error {
	programName := "gator"

	programOverview := "Run gator in the background (check `gator help agg`). Then, login or create a user for yourself.\n"
	programOverview += "View and follow available RSS feeds, or add new feeds.\n"
	programOverview += "See updates across all your followed feeds and see posts' descriptions\n"

	// Individual command help
	if len(cmd.args) > 0 {
		// Args: commandName
		commandName := cmd.args[0]

		fmt.Println()
		fmt.Println("Usage:")
		fmt.Printf("\t%v %v\n", programName, c.cmds[commandName].info.usage)
		fmt.Println()
		fmt.Println("Description:")
		fmt.Printf("\t%v\n", c.cmds[commandName].info.description)

		if len(c.cmds[commandName].info.examples) > 0 {
			fmt.Println()
			fmt.Println("Examples:")
			for _, example := range c.cmds[commandName].info.examples {
				fmt.Printf("\t%v %v\n", programName, example)
			}
		}
		// All command help
	} else {
		fmt.Printf("\n%v - a simple CLI RSS reader\n\n", programName)

		fmt.Println(programOverview)

		fmt.Println("Commands:")

		// Sort commands alphabetically before print
		sortedCommandNames := []string{}
		for key := range c.cmds {
			sortedCommandNames = append(sortedCommandNames, key)
		}
		slices.Sort(sortedCommandNames)

		for _, cmd := range sortedCommandNames {
			fmt.Printf("\t%v\t\t%v\n", cmd, c.cmds[cmd].info.description)
		}
	}

	return nil
}
