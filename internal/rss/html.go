package rss

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func download_post_html_body(ctx context.Context, url string) ([]string, error) {
	// Get page
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("User-Agent", "Go-Demo-Aggregator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	// Read body
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}

	paragraphs := []string{}

	for n := range doc.Descendants() {
		// https://pkg.go.dev/golang.org/x/net@v0.39.0/html#Node
		if n.Type == html.ElementNode && n.DataAtom == atom.Article {
			for c := range n.ChildNodes() {
				for d := range c.ChildNodes() {
					if d.Type == html.ElementNode && (d.DataAtom == atom.P) {
						paragraphs = append(paragraphs, d.FirstChild.Data)
					}
				}

			}

			break
		}
	}

	return paragraphs, fmt.Errorf("fail to extract post body")
}
