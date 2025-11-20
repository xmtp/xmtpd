package workers

import (
	"time"
)

var (
	// Spread nodes across N minutes (0 .. 59).
	distributionSpreadMinutes uint32 = 10

	// Run every M minutes.
	repeatIntervalMinutes uint32 = 10
)

// Distributes the minute of the hour that each node will run operations on.
// This is to avoid having many nodes do identical work at the same time, and potentially leading to
// duplicate reports or submissions.
//
// TODO: Make distributionInterval and distributionRange configurable via blockchain parameters.
func findNextRunTime(myNodeID uint32, workerID uint32) time.Time {
	// Use a hash function to distribute node IDs evenly across minutes
	// regardless of their numerical distribution
	// Uses a Knuth multiplicative hash function
	minuteOffset := int((((myNodeID + workerID) * 2654435761) >> 16) % distributionSpreadMinutes)

	now := time.Now().UTC()

	// Next run is at the „minuteOffset" of the current hour unless that
	// time has already passed.
	nextRun := time.Date(
		now.Year(), now.Month(), now.Day(), now.Hour(), minuteOffset, 0, 0, time.UTC,
	)

	// Advance next run.
	// The for loop cover corner cases where the next run is in the past.
	for !nextRun.After(now) {
		nextRun = nextRun.Add(time.Duration(repeatIntervalMinutes) * time.Minute)
	}

	return nextRun
}
