package rss

import "testing"

func TestExtractPost(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name: "Smoke Test",
			url:  "https://www.guildwars2.com/en/news/expand-your-homestead-with-our-sea-raider-homestead-pack/",
			// url:  "https://www.nasa.gov/image-detail/pia26352/",
			// url: "https://pivot-to-ai.com/2025/05/05/uk-backtracks-on-ai-opt-out-scheme-for-creative-works/",
			expected: `
SPHEREx Starts Scanning Entire Sky

NASA's SPHEREx mission is observing the entire sky in 102 infrared colors, or wavelengths of light not visible to the human eye. This image shows a section of sky in one wavelength (3.29 microns), revealing a cloud of dust made of a molecule similar to soot or smoke.

Image Credit: NASA/JPL-Caltech
`,
		},
	}

	for _, test := range tests {
		actual, err := ExtractFullPost(test.url)
		if err != nil {
			t.Error(err)
		}

		if actual != test.expected {
			t.Errorf("Name: %v\nUrl:\t%v\nActual:\t%v\nExpected:\t%v", test.name, test.url, actual, test.expected)
		}
	}
}
