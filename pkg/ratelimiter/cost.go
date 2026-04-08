package ratelimiter

import (
	"math"
	"time"
)

// CostQuery returns the rate-limit token cost of a query against `numTopics` topics.
// Cost is sublinear: ceil(sqrt(max(numTopics, 1))). A 0-topic query is malformed
// but charged the baseline cost of 1 rather than rejected separately.
func CostQuery(numTopics int) uint64 {
	if numTopics < 1 {
		numTopics = 1
	}
	return uint64(math.Ceil(math.Sqrt(float64(numTopics))))
}

// CostSubscribeDrain returns the retrospective drain cost for a subscription
// that was held open for `elapsed` time. The interval count is rounded up:
// ceil(elapsed / intervalMinutes) intervals, each costing `drainAmount` tokens.
//
// A stream that closes at exactly elapsed=0 pays no drain cost. Any non-zero
// elapsed time within the first interval rounds up to one full interval —
// rounding up (rather than down) means partial intervals are billed and a
// client cannot avoid drain by closing just before each interval boundary.
// Open-and-immediately-close abuse is bounded separately by the admission
// cost paid at stream open and the subscribe-opens-per-minute sub-limit.
func CostSubscribeDrain(elapsed time.Duration, intervalMinutes, drainAmount int) uint64 {
	if elapsed <= 0 || intervalMinutes <= 0 || drainAmount <= 0 {
		return 0
	}
	intervals := uint64(math.Ceil(elapsed.Minutes() / float64(intervalMinutes)))
	return intervals * uint64(drainAmount)
}
