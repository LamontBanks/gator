/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Logs in the user",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return login(appState, args[0])
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func login(s *state, username string) error {
	// User must be registered to log in
	if !userRegistered(s, username) {
		return fmt.Errorf("%v not registered", username)
	}

	// "Log in" the user by writing their name to the config file
	s.config.CurrentUserName = username
	if err := s.config.SetConfig(); err != nil {
		return err
	}

	fmt.Printf("Logged in as %v\n", username)

	return nil
}
