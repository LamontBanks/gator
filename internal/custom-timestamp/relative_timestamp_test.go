package customtimestamp

import (
	"fmt"
	"testing"
	"time"
)

func TestRelativeHour(t *testing.T) {
	tests := []struct {
		name         string
		hourRange    []int
		expectedTime string
	}{
		{
			name:         "Overnight",
			hourRange:    []int{0, 3},
			expectedTime: "overnight",
		},
		{
			name:         "Early Morning",
			hourRange:    []int{4, 7},
			expectedTime: "early morning",
		},
		{
			name:         "Morning",
			hourRange:    []int{8, 11},
			expectedTime: "morning",
		},
		{
			name:         "Afternoon",
			hourRange:    []int{12, 15},
			expectedTime: "afternoon",
		},
		{
			name:         "Evening",
			hourRange:    []int{16, 19},
			expectedTime: "evening",
		},
		{
			name:         "Night",
			hourRange:    []int{20, 23},
			expectedTime: "night",
		},
	}

	for _, test := range tests {

		// Loop through the hourRange (inclusive), checking the timestamp
		for hour := test.hourRange[0]; hour <= test.hourRange[1]; hour++ {
			actualTime := relativeTimeOfDay(time.Date(2025, time.April, 10, hour, 59, 59, 0, time.UTC))

			if actualTime != test.expectedTime {
				t.Error(printTestError(test.name, hour, test.expectedTime, actualTime))
			}
		}
	}
}

func TestRelativeTimestamp(t *testing.T) {
	tests := []struct {
		name         string
		mockTimeThen time.Time
		mockTimeNow  time.Time
		expected     string
	}{
		{
			name:         "Years",
			mockTimeThen: time.Date(2021, time.April, 10, 0, 0, 0, 0, time.UTC),
			mockTimeNow:  time.Date(2024, time.April, 10, 0, 0, 0, 0, time.UTC),
			expected:     "4 years ago",
		},
	}

	for _, test := range tests {
		actual := relativeTimestamp(test.mockTimeNow, test.mockTimeThen)
		if actual != test.expected {
			t.Error(printTestError(test.name, test.mockTimeThen, actual, test.expected))
		}
	}
}

func printTestError(testName string, input, actual, expected any) string {
	return fmt.Sprintf("\n%v\nInput:\t\t[%v]\nActual:\t\t[%v]\nExpected:\t[%v]", testName, input, actual, expected)
}
