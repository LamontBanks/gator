/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all users",
	Long:  `List all users`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printAllUsers(appState)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}

func userRegistered(s *state, username string) bool {
	_, err := s.db.GetUser(context.Background(), username)

	return err == nil
}

func printAllUsers(s *state) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users registered")
		return nil
	}

	// Print users, with indicator for current user
	var usersList string
	for _, user := range users {
		u := fmt.Sprintf("* %v", user)

		if user == s.config.CurrentUserName {
			u += " (current)"
		}
		u += "\n"

		usersList += u
	}

	fmt.Print(usersList)

	return nil
}
