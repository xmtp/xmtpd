package chainwatcher

import (
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// gaugeValue reads the current value of a Prometheus Gauge.
func gaugeValue(g interface{ Write(*dto.Metric) error }) float64 {
	var m dto.Metric
	if err := g.Write(&m); err != nil {
		return 0
	}
	return m.GetGauge().GetValue()
}

func TestUpdateActiveOriginatorCount_ExcludesExpiredNodes(t *testing.T) {
	w := &Watcher{
		activeOriginatorWindow: 150 * time.Minute,
		activeOriginators: map[uint32]time.Time{
			// Current nodes — recent block timestamps (within window)
			100: time.Now().Add(-10 * time.Minute),
			200: time.Now().Add(-30 * time.Minute),
			300: time.Now().Add(-60 * time.Minute),
			// Old decommissioned nodes — old block timestamps (outside window)
			10: time.Now().Add(-24 * time.Hour),
			11: time.Now().Add(-48 * time.Hour),
			13: time.Now().Add(-72 * time.Hour),
		},
	}

	w.updateActiveOriginatorCount()

	// With block timestamps, only the 3 recent nodes should be counted.
	// The old bug (using time.Now()) would have shown 6.
	got := gaugeValue(activeOriginatorNodes)
	assert.InDelta(t, 3, got, 0, "expected 3 active originators, got %v", got)
}

func TestUpdateActiveOriginatorCount_AllExpired(t *testing.T) {
	w := &Watcher{
		activeOriginatorWindow: 150 * time.Minute,
		activeOriginators: map[uint32]time.Time{
			10: time.Now().Add(-24 * time.Hour),
			11: time.Now().Add(-48 * time.Hour),
		},
	}

	w.updateActiveOriginatorCount()

	got := gaugeValue(activeOriginatorNodes)
	assert.InDelta(t, 0, got, 0, "expected 0 active originators when all expired")
}

func TestUpdateActiveOriginatorCount_AllActive(t *testing.T) {
	w := &Watcher{
		activeOriginatorWindow: 150 * time.Minute,
		activeOriginators: map[uint32]time.Time{
			100: time.Now().Add(-1 * time.Minute),
			200: time.Now().Add(-5 * time.Minute),
		},
	}

	w.updateActiveOriginatorCount()

	got := gaugeValue(activeOriginatorNodes)
	assert.InDelta(t, 2, got, 0, "expected 2 active originators")
}

func TestCleanupStaleEntries_RemovesExpiredOriginators(t *testing.T) {
	w := &Watcher{
		activeOriginatorWindow: 150 * time.Minute,
		activeOriginators: map[uint32]time.Time{
			100: time.Now().Add(-10 * time.Minute), // fresh
			10:  time.Now().Add(-24 * time.Hour),   // stale
		},
		submissionTimeByKey: make(map[string]time.Time),
		blockTimestampCache: map[uint64]time.Time{
			1000: time.Now().Add(-10 * time.Minute), // recent block
			50:   time.Now().Add(-24 * time.Hour),   // old block
		},
	}

	w.cleanupStaleEntries()

	// Stale originator removed
	require.Len(t, w.activeOriginators, 1)
	_, exists := w.activeOriginators[100]
	assert.True(t, exists, "fresh originator should remain")
	_, exists = w.activeOriginators[10]
	assert.False(t, exists, "stale originator should be removed")

	// Old block timestamp cache entry removed
	require.Len(t, w.blockTimestampCache, 1)
	_, exists = w.blockTimestampCache[1000]
	assert.True(t, exists, "recent cache entry should remain")
	_, exists = w.blockTimestampCache[50]
	assert.False(t, exists, "old cache entry should be removed")
}
