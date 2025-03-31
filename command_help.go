package main

import "fmt"

func helpCommandInfo() commandInfo {
	return commandInfo{
		description: "Show commands",
		usage:       "help",
		examples:    []string{},
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
	cmdName := cmd.name

	if len(cmd.args) > 0 {
		cmdName = cmd.args[0]
	}

	programName := "gator"

	fmt.Println(c.cmds[cmdName].info.description)
	fmt.Println()
	fmt.Printf("Usage:\t%v %v\n", programName, c.cmds[cmdName].info.usage)

	if len(c.cmds[cmdName].info.examples) > 0 {
		fmt.Println()
		fmt.Println("Examples:")
		for _, example := range c.cmds[cmdName].info.examples {
			fmt.Printf("\t%v %v\n", programName, example)
		}
	}

	return nil
}
