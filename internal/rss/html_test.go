package rss

import (
	"context"
	"reflect"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		url      string
		expected []string
	}{
		{
			name: "Smoke",
			html: `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`,
			expected: []string{
				"foo",
				"/bar/baz",
			},
		},
	}

	for _, test := range tests {
		actual, err := download_post_html_body(context.Background(), "https://pivot-to-ai.com/2025/05/01/seattle-worldcon-science-fiction-convention-vets-panelists-with-chatgpt/")
		// actual, err := download_post_html_body(context.Background(), "https://phys.org/news/2025-04-laser-communication-mars.html")
		// actual, err := download_post_html_body(context.Background(), "https://massivelyop.com/2025/04/30/were-finally-getting-our-first-glimpse-of-12-year-old-camelot-unchained-since-last-year-in-todays-dev-stream/")
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("%v\nExpected: %v\nActual: %v\n", test.name, test.expected, actual)
		}
	}
}
