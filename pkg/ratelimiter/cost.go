package ratelimiter

import "math"

// CostQuery returns the rate-limit token cost of a query against `numTopics` topics.
// Cost is sublinear: ceil(sqrt(max(numTopics, 1))). A 0-topic query is malformed
// but charged the baseline cost of 1 rather than rejected separately.
func CostQuery(numTopics int) uint64 {
	if numTopics < 1 {
		numTopics = 1
	}
	return uint64(math.Ceil(math.Sqrt(float64(numTopics))))
}
