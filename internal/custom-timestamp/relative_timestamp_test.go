package customtimestamp

import (
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	tests := []struct {
		name            string
		publishedAtDate time.Time
		timeElapsed     string
		expected        string
	}{
		// Days
		{
			name:            "1 day",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "24h",
			expected:        "1 day ago",
		},
		{
			name:            "More than 1, less than 2 days",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "36h",
			expected:        "1 day ago",
		},
		{
			name:            "2 days",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "48h",
			expected:        "2 days ago",
		},
		// Hours
		{
			name:            "1 hour",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "1h",
			expected:        "1 hour ago",
		},
		{
			name:            "More than 1 hour, less than 1 day",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "11h",
			expected:        "11 hours ago",
		},
		// Minutes
		{
			name:            "1 minute",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "1m",
			expected:        "1 minute ago",
		},
		{
			name:            "More than 1 minute, less than 1 hour",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "59m",
			expected:        "59 minutes ago",
		},
		// Seconds
		{
			name:            "1 second",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "1s",
			expected:        "1 second ago",
		},
		{
			name:            "More than 1 second, less than 1 minute",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "59s",
			expected:        "59 seconds ago",
		},
	}

	for _, test := range tests {
		// Setup
		timeElapsed, err := time.ParseDuration(test.timeElapsed)
		if err != nil {
			t.Error(err)
		}
		mockCurrTime := test.publishedAtDate.Add(timeElapsed)

		// Test
		actual := relativeTime(mockCurrTime, test.publishedAtDate)

		// Check
		if actual != test.expected {
			t.Errorf("\nCase:\t\t\t%v\npublishedAtDate:\t%v\ntimeElapsed:\t\t%v\nExpected:\t\t%v\nActual:\t\t\t%v", test.name, test.publishedAtDate, test.timeElapsed, test.expected, actual)
		}
	}

}
