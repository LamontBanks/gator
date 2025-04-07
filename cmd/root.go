/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"os"

	"github.com/LamontBanks/gator/internal/config"
	"github.com/LamontBanks/gator/internal/database"
	"github.com/spf13/cobra"

	// Leading underscore means the package will be used, but not directly
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

var (
	// Used by gator
	appState *state
	db       *sql.DB

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gator",
		Short: "A CLI RSS Feed reader",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cobra.OnInitialize(initAppState)
	cobra.OnFinalize(closeDB)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gator.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initAppState() {
	// Initialize info for the application state
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Database
	connStr := cfg.DbUrl
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	// Database is closed by Cobra using 'closeDB()'

	// Using the SQLC `database` wrapper instead of the native Go SQL db directly
	dbQueries := database.New(db)

	// Set state
	appState = &state{
		config: &cfg,
		db:     dbQueries,
	}
}

func closeDB() {
	db.Close()
}
