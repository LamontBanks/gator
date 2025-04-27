package relativetimestamp

import (
	"testing"
	"time"
)

func TestSmoke(t *testing.T) {
	tests := []struct {
		name        string
		timeElapsed string
		expected    string
	}{
		// Days
		{
			name:        "1 day",
			timeElapsed: "24h",
			expected:    "1 day ago",
		},
		{
			name:        "More than 1, less than 2 days",
			timeElapsed: "36h",
			expected:    "1 day ago",
		},
		{
			name:        "2 days",
			timeElapsed: "48h",
			expected:    "2 days ago",
		},
		// Hours
		{
			name:        "1 hour",
			timeElapsed: "1h",
			expected:    "1 hour ago",
		},
		{
			name:        "More than 1 hour, less than 1 day",
			timeElapsed: "11h",
			expected:    "11 hours ago",
		},
		// Minutes
		{
			name:        "1 minute",
			timeElapsed: "1m",
			expected:    "1 minute ago",
		},
		{
			name:        "More than 1 minute, less than 1 hour",
			timeElapsed: "59m",
			expected:    "59 minutes ago",
		},
		// Seconds
		{
			name:        "1 second",
			timeElapsed: "1s",
			expected:    "1 second ago",
		},
		{
			name:        "More than 1 second, less than 1 minute",
			timeElapsed: "59s",
			expected:    "59 seconds ago",
		},
	}

	for _, test := range tests {
		// Setup
		timeElapsed, err := time.ParseDuration(test.timeElapsed)
		if err != nil {
			t.Error(err)
		}

		mockPublishedAtDate := time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC)
		mockCurrTime := mockPublishedAtDate.Add(timeElapsed)

		// Test
		actual := relativeTime(mockCurrTime, mockPublishedAtDate)

		// Check
		if actual != test.expected {
			t.Errorf("\nCase:\t\t\t%v\npublishedAtDate:\t%v\ntimeElapsed:\t\t%v\nExpected:\t\t%v\nActual:\t\t\t%v", test.name, mockPublishedAtDate, test.timeElapsed, test.expected, actual)
		}
	}

}
