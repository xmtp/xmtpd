package networkwatcher

import (
	"sync"
	"time"
)

// Aggregator holds the latest-known cursor state from every publisher and
// writes Prometheus metrics on each update. It is safe for concurrent use.
type Aggregator struct {
	mu sync.RWMutex
	// state[publisherID][originatorID] = sequenceID
	state map[uint32]map[uint32]uint64
	// live[publisherID] = true if the subscribe stream is currently connected.
	live map[uint32]bool
	// now is injectable for tests.
	now func() time.Time
}

// NewAggregator returns an empty Aggregator.
func NewAggregator() *Aggregator {
	return &Aggregator{
		state: make(map[uint32]map[uint32]uint64),
		live:  make(map[uint32]bool),
		now:   time.Now,
	}
}

// Apply records a cursor snapshot reported by publisherID and updates
// derived metrics (raw gauges, divergence, max, last-update timestamp).
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
// connected. Writes the node_up gauge and recomputes divergence/max for
// every originator known to this publisher.
func (a *Aggregator) SetNodeUp(publisherID uint32, up bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.live[publisherID] = up

	val := 0.0
	if up {
		val = 1
	}
	nodeUp.WithLabelValues(nodeIDLabel(publisherID)).Set(val)

	if pubState, ok := a.state[publisherID]; ok {
		for originatorID := range pubState {
			a.recomputeDerivedLocked(originatorID)
		}
	}
}

// recomputeDerivedLocked updates divergence + max gauges for originatorID.
// Callers must hold a.mu.
func (a *Aggregator) recomputeDerivedLocked(originatorID uint32) {
	var (
		minSeq, maxSeq uint64
		seen           bool
	)
	for publisherID, pubState := range a.state {
		if !a.live[publisherID] {
			continue
		}
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
}
