package main

import (
	"context"
	"fmt"
	"time"

	"github.com/LamontBanks/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

// Register a user on the server, then pdates the config with the user.
// Usage:
//
//	$ go run . register <username>
//	$ go run . register alice
func handlerRegister(s *state, cmd command) error {
	// Args: username
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <username>", cmd.name)
	}
	username := cmd.args[0]

	// Insert user
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

	// Update the config as well
	return handlerLogin(s, cmd)
}

// Lists all users
// Indicates which usersis logged in
func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currUser := s.config.CurrentUserName
	var usersList string

	for _, user := range users {
		u := fmt.Sprintf("* %v", user)

		if user == currUser {
			u += " (current)"
		}
		u += "\n"

		usersList += u
	}

	fmt.Print(usersList)

	return nil
}
