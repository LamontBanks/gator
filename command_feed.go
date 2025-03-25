package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	// Args: url
	// if len(cmd.args) < 1 {
	// 	return fmt.Errorf("feed url required")
	// }
	// feedUrl := cmd.args[0]

	feedUrl := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)

	return nil
}
