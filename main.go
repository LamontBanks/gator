package main

import (
	"fmt"
	"log"

	"github.com/LamontBanks/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.CurrentUserName = "Lamont"

	cfg.SetConfig()

	fmt.Println(cfg)
}
