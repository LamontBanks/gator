package customtimestamp

import (
	"fmt"
	"time"
)

func relativeTime(currDate time.Time, publishDate time.Time) string {
	timeSince := currDate.Sub(publishDate)

	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		return currDate.String()
	}

	// days
	if timeSince >= oneDay {
		numDays := int64(timeSince / (24 * time.Hour))
		msgDays := appendMultiplesSuffixS("day", numDays)
		return fmt.Sprintf("%v %v ago", numDays, msgDays)
	}

	// hours
	if timeSince >= time.Hour && timeSince < oneDay {
		numHours := int64(timeSince / time.Hour)
		msgHours := appendMultiplesSuffixS("hour", numHours)
		return fmt.Sprintf("%v %v ago", numHours, msgHours)
	}

	// minutes
	if timeSince >= time.Minute && timeSince < time.Hour {
		numMinutes := int64(timeSince / time.Minute)
		msgMinutes := appendMultiplesSuffixS("minute", numMinutes)
		return fmt.Sprintf("%v %v ago", numMinutes, msgMinutes)
	}

	// seconds
	numSecond := int64(timeSince / time.Second)
	msgSeconds := appendMultiplesSuffixS("second", numSecond)
	return fmt.Sprintf("%v %v ago", numSecond, msgSeconds)
}

// Return `word` with appended "s" if count > 1
// Otherwise return `word` unchanged
func appendMultiplesSuffixS(word string, count int64) string {
	if count > 1 {
		word += "s"
	}

	return word
}
