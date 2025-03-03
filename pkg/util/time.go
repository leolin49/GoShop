package util

import "time"

const (
	MILLISECOND = 1
	SECOND      = 1000 * MILLISECOND
	MINUTE      = 60 * SECOND
	HOUR        = 60 * MINUTE
	DAY         = 24 * HOUR
)

func TimeToNextHour() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
	return next.Sub(now)
}

func TimeToNextDay() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return next.Sub(now)
}

func TimeToNextMonth() time.Duration {
	now := time.Now()
	next := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	return next.Sub(now)
}
