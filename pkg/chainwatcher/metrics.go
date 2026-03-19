// Package chainwatcher implements a standalone chain watcher that monitors
// payer report events on the settlement chain and emits Prometheus metrics.
package chainwatcher

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	reportSubmittedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtp_chain_report_submitted_total",
			Help: "Total number of payer reports submitted on-chain, by originator node.",
		},
		[]string{"originator_node_id"},
	)

	reportSettledTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtp_chain_report_settled_total",
			Help: "Total number of payer reports fully settled on-chain, by originator node.",
		},
		[]string{"originator_node_id"},
	)

	timeSinceLastSubmissionSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtp_chain_time_since_last_submission_seconds",
			Help: "Seconds since the last PayerReportSubmitted event was observed.",
		},
	)

	submissionToSettlementSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "xmtp_chain_submission_to_settlement_seconds",
			Help: "Duration between report submission and full settlement.",
			Buckets: []float64{
				60, 120, 300, 600, 900, 1800, 3600, 5400, 7200,
			},
		},
		[]string{"originator_node_id"},
	)

	envelopeRangeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtp_chain_envelope_range_total",
			Help: "Total envelopes covered by payer reports (endSeq - startSeq).",
		},
		[]string{"originator_node_id"},
	)

	envelopeRangeGap = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtp_chain_envelope_range_gap",
			Help: "Gap between consecutive report ranges for the same originator. >0 indicates potential data loss.",
		},
		[]string{"originator_node_id"},
	)

	attestingNodeCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtp_chain_attesting_node_count",
			Help: "Number of nodes that signed the report.",
		},
		[]string{"originator_node_id"},
	)

	activeOriginatorNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtp_chain_active_originator_nodes",
			Help: "Number of distinct originator node IDs seen in the sliding window.",
		},
	)

	feesSettledPicodollars = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtp_chain_fees_settled_picodollars",
			Help: "Total fees settled on-chain in picodollars.",
		},
		[]string{"originator_node_id"},
	)

	usageSettledTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "xmtp_chain_usage_settled_total",
			Help: "Total number of UsageSettled events observed.",
		},
	)

	eventsProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtp_chain_events_processed_total",
			Help: "Total chain events processed by type.",
		},
		[]string{"event_type"},
	)
)

func RegisterMetrics(reg prometheus.Registerer) {
	collectors := []prometheus.Collector{
		reportSubmittedTotal,
		reportSettledTotal,
		timeSinceLastSubmissionSeconds,
		submissionToSettlementSeconds,
		envelopeRangeTotal,
		envelopeRangeGap,
		attestingNodeCount,
		activeOriginatorNodes,
		feesSettledPicodollars,
		usageSettledTotal,
		eventsProcessedTotal,
	}
	for _, c := range collectors {
		reg.MustRegister(c)
	}
}

func nodeIDLabel(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}
