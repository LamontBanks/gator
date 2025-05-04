package fuzzytimestamp

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
		// hour+minute if bewtween 1-6 hours
		{
			name:        "Original hour:minute if between 1-5 hours (no rounding)",
			timeElapsed: "1h",
			expected:    "1:23 AM",
		},
		{
			name:        "Original hour:minute if between 1-5 hours (no rounding)",
			timeElapsed: "5h",
			expected:    "1:23 AM",
		},
		{
			name:        "Hours if between 6-23 hours, 6h",
			timeElapsed: "6h",
			expected:    "6h",
		},
		{
			name:        "Hours if between 6-23 hours, 23h",
			timeElapsed: "23h",
			expected:    "23h",
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

		mockPublishedAtDate := time.Date(2025, time.June, 15, 1, 23, 45, 0, time.UTC)
		mockCurrTime := mockPublishedAtDate.Add(timeElapsed)

		// Test
		actual := fuzzyTime(mockCurrTime, mockPublishedAtDate)

		// Check
		if actual != test.expected {
			t.Errorf("\nCase:\t\t\t%v\npublishedAtDate:\t%v\ntimeElapsed:\t\t%v\nExpected:\t\t%v\nActual:\t\t\t%v", test.name, mockPublishedAtDate, test.timeElapsed, test.expected, actual)
		}
	}

}
