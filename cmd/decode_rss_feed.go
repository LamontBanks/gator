/*
Reads RSS feed info into custom structs
*/
package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

// https://www.rssboard.org/files/sample-rss-2.xml
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Downloads the feed
func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Go-Demo-Aggregator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return &rssFeed, err
	}

	// Unescape HTML characters in both the Channel and Items
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil
}

// Try parsing time with multiple time format, use what works
func ParseRSSPubDate(pubDate string) (time.Time, error) {
	// https://pkg.go.dev/time@go1.24.1#Layout
	timeLayoutsToTry := []string{
		"Mon, 02 Jan 2006 15:04 MST",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}

	// User current time if none is provided
	if pubDate == "" {
		return time.Now(), nil
	}

	for _, layout := range timeLayoutsToTry {
		convertedDate, err := time.ParseInLocation(layout, pubDate, time.Local)
		if err != nil {
			continue
		} else {
			return convertedDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to convert pubDate %v using these time.Layouts: %v", pubDate, timeLayoutsToTry)
}
