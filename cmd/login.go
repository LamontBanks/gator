/*
Log in a user
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in a register user",
	Long:  `Log in a register user`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return login(appState, args[0])
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

// "Log in" the user by writing their name to the config file
// Users only need to enter their username; there are no passwords
func login(s *state, username string) error {
	if !userRegistered(s, username) {
		return fmt.Errorf("%v not registered", username)
	}

	s.config.CurrentUserName = username
	if err := s.config.SetConfig(); err != nil {
		return err
	}

	fmt.Printf("Logged in as %v\n", username)

	return nil
}
