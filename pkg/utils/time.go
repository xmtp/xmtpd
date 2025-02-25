package utils

import "time"

func MinutesSinceEpoch(timestamp time.Time) int32 {
	durationSinceEpoch := timestamp.Sub(time.Unix(0, 0))

	return int32(durationSinceEpoch.Minutes())
}

func MinutesSinceEpochNow() int32 {
	return MinutesSinceEpoch(time.Now())
}

func NsToDate(ns int64) time.Time {
	return time.Unix(0, ns)
}
