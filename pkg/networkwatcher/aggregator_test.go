package networkwatcher

import (
	"sync"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

// metricValue returns the current value of a labelled counter/gauge.
func metricValue(t *testing.T, c interface{ Write(*dto.Metric) error }) float64 {
	t.Helper()
	var m dto.Metric
	require.NoError(t, c.Write(&m))
	if g := m.GetGauge(); g != nil {
		return g.GetValue()
	}
	if g := m.GetCounter(); g != nil {
		return g.GetValue()
	}
	return 0
}

// resetAggregatorMetrics wipes state between tests. Metrics are process-
// globals so reuse across tests must be explicit.
func resetAggregatorMetrics() {
	cursorGauge.Reset()
	cursorDivergence.Reset()
	cursorMax.Reset()
	nodeUp.Reset()
	nodeLastUpdateSeconds.Reset()
}

func TestAggregator_ApplyCursor_WritesRawGauges(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	a.SetNodeUp(1, true)

	a.Apply(1, map[uint32]uint64{100: 500, 200: 600})

	require.InDelta(t, 500.0, metricValue(t, cursorGauge.WithLabelValues("1", "100")), 0)
	require.InDelta(t, 600.0, metricValue(t, cursorGauge.WithLabelValues("1", "200")), 0)
}

func TestAggregator_ConcurrentApplyIsRaceFree(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	for i := uint32(1); i <= 4; i++ {
		a.SetNodeUp(i, true)
	}

	var wg sync.WaitGroup
	for pub := uint32(1); pub <= 4; pub++ {
		wg.Add(1)
		go func(pub uint32) {
			defer wg.Done()
			for i := range uint64(100) {
				a.Apply(pub, map[uint32]uint64{100: i, 200: i * 2})
			}
		}(pub)
	}
	wg.Wait()

	_ = time.Now
}
