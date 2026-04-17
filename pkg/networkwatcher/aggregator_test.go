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
	cursorLag.Reset()
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

func TestAggregator_Divergence_MaxMinusMinAcrossLivePublishers(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	a.SetNodeUp(1, true)
	a.SetNodeUp(2, true)
	a.SetNodeUp(3, true)

	a.Apply(1, map[uint32]uint64{100: 500})
	a.Apply(2, map[uint32]uint64{100: 550})
	a.Apply(3, map[uint32]uint64{100: 700})

	require.InDelta(t, 200.0, metricValue(t, cursorDivergence.WithLabelValues("100")), 0)
	require.InDelta(t, 700.0, metricValue(t, cursorMax.WithLabelValues("100")), 0)
}

func TestAggregator_Divergence_IncludesDroppedPublishers(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	a.SetNodeUp(1, true)
	a.SetNodeUp(2, true)

	a.Apply(1, map[uint32]uint64{100: 500})
	a.Apply(2, map[uint32]uint64{100: 700})

	// Node 2 drops. Its last-known cursor stays in the state and
	// continues to contribute to divergence/max/lag — the whole point of
	// the signal is that a dead node's growing gap should be visible.
	a.SetNodeUp(2, false)

	// node_up reflects liveness, but derived metrics still see both.
	require.InDelta(t, 0.0, metricValue(t, nodeUp.WithLabelValues("2")), 0)
	require.InDelta(t, 200.0, metricValue(t, cursorDivergence.WithLabelValues("100")), 0)
	require.InDelta(t, 700.0, metricValue(t, cursorMax.WithLabelValues("100")), 0)
	require.InDelta(t, 700.0, metricValue(t, cursorGauge.WithLabelValues("2", "100")), 0)

	// Advance the live publisher; the down publisher's lag grows.
	a.Apply(1, map[uint32]uint64{100: 900})
	require.InDelta(t, 0.0, metricValue(t, cursorLag.WithLabelValues("1", "100")), 0)
	require.InDelta(t, 200.0, metricValue(t, cursorLag.WithLabelValues("2", "100")), 0)
}

func TestAggregator_Divergence_SinglePublisherIsZero(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	a.SetNodeUp(1, true)

	a.Apply(1, map[uint32]uint64{100: 500})

	require.InDelta(t, 0.0, metricValue(t, cursorDivergence.WithLabelValues("100")), 0)
	require.InDelta(t, 500.0, metricValue(t, cursorMax.WithLabelValues("100")), 0)
}

func TestAggregator_Lag_IdentifiesBehindPublisher(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	a.SetNodeUp(1, true)
	a.SetNodeUp(2, true)
	a.SetNodeUp(3, true)

	a.Apply(1, map[uint32]uint64{100: 500})
	a.Apply(2, map[uint32]uint64{100: 550})
	a.Apply(3, map[uint32]uint64{100: 700})

	// Max is 700 (pub 3). Lag for each publisher = 700 - seq.
	require.InDelta(t, 200.0, metricValue(t, cursorLag.WithLabelValues("1", "100")), 0)
	require.InDelta(t, 150.0, metricValue(t, cursorLag.WithLabelValues("2", "100")), 0)
	require.InDelta(t, 0.0, metricValue(t, cursorLag.WithLabelValues("3", "100")), 0)
}

func TestAggregator_LastUpdateSecondsUsesInjectedNow(t *testing.T) {
	resetAggregatorMetrics()
	a := NewAggregator()
	fixed := time.Unix(1700000000, 0)
	a.now = func() time.Time { return fixed }
	a.SetNodeUp(1, true)

	a.Apply(1, map[uint32]uint64{100: 1})

	require.InDelta(
		t,
		float64(1700000000),
		metricValue(t, nodeLastUpdateSeconds.WithLabelValues("1")),
		0,
	)
}
