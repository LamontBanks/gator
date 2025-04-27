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
		{
			name:            "Basic",
			publishedAtDate: time.Date(2025, time.June, 15, 0, 0, 0, 0, time.UTC),
			timeElapsed:     "24h",
			expected:        "1 day ago",
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
			t.Errorf("\n%v\nActual: %v", test, actual)
		}
	}

}
