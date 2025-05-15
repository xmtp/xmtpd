package workers

import "time"

// Distributes the minute of the hour that each node will run operations on.
// This is to avoid having many nodes do identical work at the same time, and potentially leading to
// duplicate reports or submissions.
func findNextRunTime(myNodeID uint32, workerID uint32) time.Time {
	// Use a hash function to distribute node IDs evenly across minutes
	// regardless of their numerical distribution
	// Uses a Knuth multiplicative hash function
	minuteOffset := int((((myNodeID + workerID) * 2654435761) >> 16) % 60)

	now := time.Now().UTC()
	// Next run is at the „minuteOffset" of the current hour unless that
	// time has already passed.
	nextRun := time.Date(
		now.Year(), now.Month(), now.Day(), now.Hour(), minuteOffset, 0, 0, time.UTC,
	)
	if !nextRun.After(now) {
		nextRun = nextRun.Add(time.Hour)
	}

	return nextRun
}
