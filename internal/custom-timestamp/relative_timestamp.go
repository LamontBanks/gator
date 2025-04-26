package customtimestamp

import (
	"fmt"
	"strconv"
	"time"
)

func RelativeTimestamp(then time.Time) string {
	return relativeTimestamp(time.Now(), then)
}

// Internal function for unit test, allows "now" to be set
func relativeTimestamp(now time.Time, then time.Time) string {
	timeDiff := now.Sub(then)

	// Ignoring errors
	year, _ := time.ParseDuration("8760h")
	month, _ := time.ParseDuration("720h")
	// day, _ := time.ParseDuration("24h")
	// hour, _ := time.ParseDuration("1h")
	// minute, _ := time.ParseDuration("1m")

	// > year
	if timeDiff >= year {
		roundedTimeDiff, err := convertDurationToInt(timeDiff / year)
		if err != nil {
			return ""
		}

		period := appendMultiplesSuffix("year", roundedTimeDiff)
		return fmt.Sprintf("%v %v ago", roundedTimeDiff, period)
	}

	// 1 - 11 months
	if timeDiff < year && timeDiff >= month {
		roundedTimeDiff, err := convertDurationToInt(timeDiff / month)
		if err != nil {
			return ""
		}

		period := appendMultiplesSuffix("month", roundedTimeDiff)
		return fmt.Sprintf("%v %v ago", roundedTimeDiff, period)
	}

	return timeDiff.String()
}

// Append "s" to `word` if num is > "0"
func appendMultiplesSuffix(word string, num int) string {
	if num > 1 {
		word += "s"
	}

	return word
}

func convertDurationToInt(t time.Duration) (int, error) {
	i, err := strconv.Atoi(fmt.Sprintf("%d", t))
	if err != nil {
		return 0, err
	}

	return i, nil
}

func relativeTimeOfDay(t time.Time) string {
	hour := t.Hour()

	// 12 AM - 3:59 AM
	if hour >= 0 && hour < 4 {
		return "overnight"
	}
	// 4 AM - 7:59 AM
	if hour >= 4 && hour < 8 {
		return "early morning"
	}
	// 8 AM - 11:59 AM
	if hour >= 8 && hour < 12 {
		return "morning"
	}
	// 12 PM - 3:59 PM
	if hour >= 12 && hour < 16 {
		return "afternoon"
	}
	// 4 PM - 7:59 PM
	if hour >= 16 && hour < 20 {
		return "evening"
	}
	// 8 PM - 11:59 PM
	if hour >= 20 && hour <= 23 {
		return "night"
	}

	return "unknown"
}
