package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func registerHelp() commandHelp {
	return commandHelp{
		description: "Registers a new users, logs them in",
		usage:       "gator register <unique username>",
		examples: []string{
			"gator register bob",
		},
	}
}

func usersHelp() commandHelp {
	return commandHelp{
		description: "List all users",
		usage:       "gator users",
		examples: []string{
			"gator users",
		},
	}
}

// Saves a user to the database, then updates the config with the user
func handlerRegister(s *state, cmd command) error {
	// Args: username
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <unique username>", cmd.name)
	}
	username := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Registered user %v\n", user.Name)

	// Update the config
	return handlerLogin(s, cmd)
}

// Lists all users, indicates the current user
func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

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
