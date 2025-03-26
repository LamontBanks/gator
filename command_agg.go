package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	// Args: RSS feed url
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <RSS Feed URL>", cmd.name)
	}
	feedUrl := cmd.args[0]

	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)

	return nil
}
