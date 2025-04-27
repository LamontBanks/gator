package relativetimestamp

import (
	"testing"
	"time"
)

func TestRelativeTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		timeElapsed string
		expected    string
	}{
		// Days
		{
			name:        "1 day",
			timeElapsed: "24h",
			expected:    "1d",
		},
		{
			name:        "1 day if less than halfway to next day",
			timeElapsed: "35h",
			expected:    "1d",
		},
		{
			name:        "Round up when halfway to next day",
			timeElapsed: "36h",
			expected:    "2d",
		},
		{
			name:        "2 days",
			timeElapsed: "48h",
			expected:    "2d",
		},
		// Hours
		{
			name:        "1 hour",
			timeElapsed: "1h",
			expected:    "1h",
		},
		{
			name:        "1 hour if less than halfway to next hour",
			timeElapsed: "89m",
			expected:    "1h",
		},
		{
			name:        "Round up when halfway to next hour",
			timeElapsed: "90m",
			expected:    "2h",
		},
		// Minutes
		{
			name:        "1 minute",
			timeElapsed: "1m",
			expected:    "1m",
		},
		{
			name:        "59 minutes",
			timeElapsed: "59m",
			expected:    "59m",
		},
		{
			name:        "1 minute if less than halfway to next minute",
			timeElapsed: "89s",
			expected:    "1m",
		},
		{
			name:        "Round up when halfway to next minute",
			timeElapsed: "90s",
			expected:    "2m",
		},
		// Seconds (no rounding)
		{
			name:        "1 second",
			timeElapsed: "1s",
			expected:    "1s",
		},
		{
			name:        "59 seconds",
			timeElapsed: "59s",
			expected:    "59s",
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
