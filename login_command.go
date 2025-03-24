package main

import (
	"context"
	"database/sql"
	"fmt"
)

// Log in the user
// User must alrerady be registered
// Usage:
//
//	$ go run . login <username>
//	$ go run . login alice
func handlerLogin(s *state, cmd command) error {
	// Get needed args
	if len(cmd.args) < 1 {
		return fmt.Errorf("username required")
	}
	username := cmd.args[0]

	// Check if the user is registered in the db
	// If nothing is returned, stop
	_, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%v not registered", username)
	}
	if err != nil {
		panic(err)
	}

	// Otherwise, log in the user by writing their name to the config file
	s.config.CurrentUserName = username
	if err := s.config.SetConfig(); err != nil {
		return err
	}
	fmt.Printf("Logged in as %v\n", username)

	return nil
}
