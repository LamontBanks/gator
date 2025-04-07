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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := appState.db.GetUsers(context.Background())
		if err != nil {
			return err
		}

		if len(users) == 0 {
			fmt.Println("No users registered")
			return nil
		}

		// List users, with indicator for current user
		var usersList string
		for _, user := range users {
			u := fmt.Sprintf("* %v", user)

			if user == appState.config.CurrentUserName {
				u += " (current)"
			}
			u += "\n"

			usersList += u
		}

		fmt.Print(usersList)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func userRegistered(s *state, username string) bool {
	_, err := s.db.GetUser(context.Background(), username)

	return err == nil
}
