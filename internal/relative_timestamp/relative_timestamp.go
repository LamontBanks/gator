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

	// weeks
	if timeSince >= oneWeek {
		return fmt.Sprintf("%vw", int64(timeSince.Round(oneWeek)/oneWeek))
	}

	// day of the week
	if timeSince >= oneDay*3 && timeSince < oneWeek {
		return strings.ToLower(publishDate.Format("Mon"))
	}

	// days
	if timeSince >= oneDay && timeSince < oneDay*3 {
		return fmt.Sprintf("%vd", int64(timeSince.Round(24*time.Hour)/(24*time.Hour)))
	}

	// hours
	if timeSince >= time.Hour && timeSince < oneDay {
		return fmt.Sprintf("%vh", int64(timeSince.Round(time.Hour)/time.Hour))
	}

	// minutes
	if timeSince >= time.Minute && timeSince < time.Hour {
		return fmt.Sprintf("%vm", int64(timeSince.Round(time.Minute)/time.Minute))
	}

	// seconds
	// No rounding
	return fmt.Sprintf("%vs", int64(timeSince/time.Second))
}
