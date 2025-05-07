package rss

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ExtractFullPost(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var onlySpaces = regexp.MustCompile(`^[\s]{1,}$`)

	postText := ""
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Article {
			for articleNodes := range n.Descendants() {
				if articleNodes.Type == html.TextNode {

					if onlySpaces.MatchString(articleNodes.Data) {
						continue
					}

					postText += articleNodes.Data + " "
				}
			}
			break
		}
	}

	return postText, nil
}
