package networkwatcher

import (
	"sync"
	"time"
)

// Aggregator holds the latest-known cursor state from every publisher and
// writes Prometheus metrics on each update. It is safe for concurrent use.
//
// A publisher's last-known cursor is retained even after its subscribe
// stream drops, and continues to contribute to divergence, max, and lag.
// That's the point of the liveness signal: a dead publisher's cursor
// falls behind live peers and the growing lag is the observable symptom.
type Aggregator struct {
	mu sync.RWMutex
	// state[publisherID][originatorID] = sequenceID
	state map[uint32]map[uint32]uint64
	// now is injectable for tests.
	now func() time.Time
}

// NewAggregator returns an empty Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		state: make(map[uint32]map[uint32]uint64),
		now:   time.Now,
	}
}

// Apply records a cursor snapshot reported by publisherID and updates
// derived metrics (raw gauges, divergence, max, lag, last-update timestamp).
func (a *Aggregator) Apply(publisherID uint32, cursor map[uint32]uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	pubLabel := nodeIDLabel(publisherID)

	existing, ok := a.state[publisherID]
	if !ok {
		existing = make(map[uint32]uint64, len(cursor))
		a.state[publisherID] = existing
	}

	touched := make(map[uint32]struct{}, len(cursor))
	for originatorID, seq := range cursor {
		existing[originatorID] = seq
		cursorGauge.
			WithLabelValues(pubLabel, nodeIDLabel(originatorID)).
			Set(float64(seq))
		touched[originatorID] = struct{}{}
	}

	nodeLastUpdateSeconds.WithLabelValues(pubLabel).
		Set(float64(a.now().Unix()))

	for originatorID := range touched {
		a.recomputeDerivedLocked(originatorID)
	}
}

// SetNodeUp records whether publisherID's subscribe stream is currently
// connected. The node_up gauge is the only liveness signal; divergence,
// max, and lag continue to reflect the publisher's last-known cursor so
// a dead node's growing gap remains visible.
func (a *Aggregator) SetNodeUp(publisherID uint32, up bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	val := 0.0
	if up {
		val = 1
	}
	nodeUp.WithLabelValues(nodeIDLabel(publisherID)).Set(val)
}

// recomputeDerivedLocked updates divergence, max, and per-publisher lag
// gauges for originatorID. Callers must hold a.mu.
func (a *Aggregator) recomputeDerivedLocked(originatorID uint32) {
	var (
		minSeq, maxSeq uint64
		seen           bool
	)
	for _, pubState := range a.state {
		seq, ok := pubState[originatorID]
		if !ok {
			continue
		}
		if !seen {
			minSeq, maxSeq, seen = seq, seq, true
			continue
		}
		if seq < minSeq {
			minSeq = seq
		}
		if seq > maxSeq {
			maxSeq = seq
		}
	}
	originatorLabel := nodeIDLabel(originatorID)
	if !seen {
		cursorDivergence.WithLabelValues(originatorLabel).Set(0)
		cursorMax.WithLabelValues(originatorLabel).Set(0)
		return
	}
	cursorDivergence.WithLabelValues(originatorLabel).Set(float64(maxSeq - minSeq))
	cursorMax.WithLabelValues(originatorLabel).Set(float64(maxSeq))

	for publisherID, pubState := range a.state {
		seq, ok := pubState[originatorID]
		if !ok {
			continue
		}
		cursorLag.
			WithLabelValues(nodeIDLabel(publisherID), originatorLabel).
			Set(float64(maxSeq - seq))
	}
}
