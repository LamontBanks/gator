package relativetimestamp

import (
	"fmt"
	"strings"
	"time"
)

func RelativeTimestamp(publishDate time.Time) string {
	return relativeTime(time.Now(), publishDate)
}

func relativeTime(currDate time.Time, publishDate time.Time) string {
	timeSince := currDate.Sub(publishDate)

	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		return currDate.String()
	}

	oneWeek, err := time.ParseDuration("168h")
	if err != nil {
		return currDate.String()
	}

	// weeks (> 7 days)
	if timeSince >= oneWeek {
		return fmt.Sprintf("%vw", int64(timeSince.Round(oneWeek)/oneWeek))
	}

	// day of the week (between 4-6 days)
	if timeSince > oneDay*3 && timeSince < oneWeek {
		return strings.ToLower(publishDate.Format("Mon"))
	}

	// days (between 1 - 3 days)
	if timeSince >= oneDay && timeSince <= oneDay*3 {
		return fmt.Sprintf("%vd", int64(timeSince.Round(24*time.Hour)/(24*time.Hour)))
	}

	// number of hours (between 6 - 24 hours)
	if timeSince >= time.Hour*6 && timeSince < oneDay {
		return fmt.Sprintf("%vh", int64(timeSince.Round(time.Hour)/time.Hour))
	}

	// hour:minute (between 1 - 6 hours)
	if timeSince >= time.Hour && timeSince < time.Hour*6 {
		return publishDate.Format("3:04 PM")
	}

	// minutes
	if timeSince >= time.Minute && timeSince < time.Hour {
		return fmt.Sprintf("%vm", int64(timeSince.Round(time.Minute)/time.Minute))
	}

	// seconds, no rounding
	return fmt.Sprintf("%vs", int64(timeSince/time.Second))
}

func relativeTimeOfDay(t time.Time) string {
	hour := t.Hour()

	// 12 AM - 6 AM
	if hour >= 0 && hour < 6 {
		return "overnight"
	}
	// 6 AM - 12 PM
	if hour >= 6 && hour < 12 {
		return "this morning"
	}

	// 12 PM - 4:59 PM
	if hour >= 12 && hour < 17 {
		return "this afternoon"
	}
	// 4 PM - 7:59 PM
	if hour >= 17 && hour < 20 {
		return "this evening"
	}
	// 8 PM - 11:59 PM
	if hour >= 20 && hour <= 23 {
		return "last night"
	}

	return "unknown"
}
