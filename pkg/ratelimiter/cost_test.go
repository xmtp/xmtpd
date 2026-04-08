package ratelimiter

import (
	"testing"

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
