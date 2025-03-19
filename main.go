package main

import (
	"fmt"
	"log"

	"github.com/LamontBanks/blog-aggregator/internal/config"
)

func main() {
	// Test reading and writing the config
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.CurrentUserName = "Lamont"
	cfg.SetConfig()
	fmt.Println(cfg)
}
