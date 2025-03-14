package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

// this should be a counter, but it does not support set and we can not rely on memory state
var syncOriginatorSequenceId = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_sync_originator_sequence_id",
		Help: "Last synced sequence id of the originator",
	},
	[]string{"originator_id"},
)

func EmitSyncLastSeenOriginatorSequenceId(originatorId uint32, lastSequence uint64) {
	syncOriginatorSequenceId.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Set(float64(lastSequence))
}

var syncOriginatorErrorMessages = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_sync_messages_received_error_count",
		Help: "Count of failed/errored messages received from the originator",
	},
	[]string{"originator_id"},
)

func EmitSyncOriginatorErrorMessages(originatorId uint32, count int) {
	syncOriginatorErrorMessages.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Add(float64(count))
}

var syncOriginatorMessagesReceived = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_sync_messages_received_count",
		Help: "Count of messages received from the originator",
	},
	[]string{"originator_id"},
)

func EmitSyncOriginatorReceivedMessagesCount(originatorId uint32, count int) {
	syncOriginatorMessagesReceived.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Add(float64(count))
}

var syncOutgoingSyncConnections = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "xmtp_sync_outgoing_sync_connections",
		Help: "Gauge of open outgoing sync connections",
	},
)

var syncFailedOutgoingSyncConnections = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "xmtp_sync_failed_outgoing_sync_connections",
		Help: "Gauge of current failed outgoing sync connections",
	},
)

var syncFailedOutgoingSyncConnectionCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_sync_failed_outgoing_sync_connections_counter",
		Help: "Counter of total number of failed outgoing sync connection attempts",
	},
	[]string{"originator_id"},
)

type SyncConnectionsStatusCounter struct {
	hasFailed    bool
	hasSucceeded bool
	originatorId uint32
}

func NewSyncConnectionsStatusCounter(originatorId uint32) *SyncConnectionsStatusCounter {
	return &SyncConnectionsStatusCounter{
		hasFailed:    false,
		hasSucceeded: false,
		originatorId: originatorId,
	}
}

func (fc *SyncConnectionsStatusCounter) MarkFailure() {
	if !fc.hasFailed {
		fc.hasFailed = true
		syncFailedOutgoingSyncConnections.Inc()
	}
	if fc.hasSucceeded {
		fc.hasSucceeded = false
		syncOutgoingSyncConnections.Dec()
	}
	syncFailedOutgoingSyncConnectionCounter.With(prometheus.Labels{"originator_id": strconv.Itoa(int(fc.originatorId))}).
		Inc()
}

func (fc *SyncConnectionsStatusCounter) MarkSuccess() {
	if fc.hasFailed {
		fc.hasFailed = false
		syncFailedOutgoingSyncConnections.Dec()
	}
	if !fc.hasSucceeded {
		fc.hasSucceeded = true
		syncOutgoingSyncConnections.Inc()
	}
}

func (fc *SyncConnectionsStatusCounter) Close() {
	if fc.hasFailed {
		fc.hasFailed = false
		syncFailedOutgoingSyncConnections.Dec()
	}
	if fc.hasSucceeded {
		fc.hasSucceeded = false
		syncOutgoingSyncConnections.Dec()
	}
}
