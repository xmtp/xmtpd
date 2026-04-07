package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCostQuery(t *testing.T) {
	tests := []struct {
		topics int
		want   uint64
	}{
		{0, 1}, // clamp to 1
		{1, 1},
		{4, 2},
		{100, 10},
		{1000, 32},
	}
	for _, tt := range tests {
		require.Equal(t, tt.want, CostQuery(tt.topics), "topics=%d", tt.topics)
	}
}

func TestCostSubscribeDrain(t *testing.T) {
	intervalMinutes := 5
	drainAmount := 1
	tests := []struct {
		elapsed time.Duration
		want    uint64
	}{
		{0, 0},                 // no time held → no drain
		{1 * time.Minute, 1},   // partial first interval → 1
		{5 * time.Minute, 1},   // exactly one interval → 1
		{6 * time.Minute, 2},   // into second interval → 2
		{60 * time.Minute, 12}, // 1 hour → 12 intervals
	}
	for _, tt := range tests {
		got := CostSubscribeDrain(tt.elapsed, intervalMinutes, drainAmount)
		require.Equal(t, tt.want, got, "elapsed=%s", tt.elapsed)
	}
}
