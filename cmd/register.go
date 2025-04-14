/*
 */
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/gator/internal/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a user",
	Long: `Username must be unique and are case-sensitive.
	There is currenetly no option to delete users.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return register(appState, args[0])
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func register(s *state, username string) error {
	// Check if user is already registered
	if userRegistered(s, username) {
		return fmt.Errorf("%v is already registered", username)
	}

	// If not registered, add them
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		return fmt.Errorf("error saving the user %v: %v", username, err)
	}

	fmt.Printf("Registered %v\n", user.Name)

	login(s, username)

	return nil
}
