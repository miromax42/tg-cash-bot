package util

import "time"

func TimeMonthAgo() time.Time {
	return time.Now().Add(-30 * 24 * time.Hour)
}
