package customtimestamp

import "time"

func relativeTime(currDate time.Time, publishDate time.Time) string {
	return currDate.Sub(publishDate).String()
}
