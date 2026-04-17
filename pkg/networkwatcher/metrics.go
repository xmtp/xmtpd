package networkwatcher

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cursorGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_cursor",
			Help: "Last sequence ID reported by publisher for envelopes from originator.",
		},
		[]string{"publisher_node_id", "originator_node_id"},
	)

	cursorDivergence = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_cursor_divergence",
			Help: "max-min cursor across live publishers for an originator.",
		},
		[]string{"originator_node_id"},
	)

	cursorMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_cursor_max",
			Help: "Network-wide max cursor value for an originator.",
		},
		[]string{"originator_node_id"},
	)

	nodeUp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_node_up",
			Help: "1 if the subscribe stream to this publisher is currently connected, 0 otherwise.",
		},
		[]string{"publisher_node_id"},
	)

	nodeStreamErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_network_watcher_node_stream_errors_total",
			Help: "Total stream errors per publisher node, labeled by reason.",
		},
		[]string{"publisher_node_id", "reason"},
	)

	nodeLastUpdateSeconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_node_last_update_seconds",
			Help: "Unix timestamp of the last cursor push from this publisher.",
		},
		[]string{"publisher_node_id"},
	)

	knownNodes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_known_nodes",
			Help: "Current node count from the on-chain registry.",
		},
	)

	registryErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "xmtpd_network_watcher_registry_errors_total",
			Help: "Total registry refresh/read errors.",
		},
	)

	buildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "xmtpd_network_watcher_build_info",
			Help: "Build information, value is always 1.",
		},
		[]string{"version"},
	)
)

// RegisterMetrics registers all network-watcher collectors on reg.
// Panics on duplicate registration (via MustRegister).
func RegisterMetrics(reg prometheus.Registerer) {
	collectors := []prometheus.Collector{
		cursorGauge,
		cursorDivergence,
		cursorMax,
		nodeUp,
		nodeStreamErrors,
		nodeLastUpdateSeconds,
		knownNodes,
		registryErrors,
		buildInfo,
	}
	for _, c := range collectors {
		reg.MustRegister(c)
	}
}

// nodeIDLabel renders a uint32 node ID as a metric label value.
//
//nolint:unused // used by the aggregator added in a later task
func nodeIDLabel(id uint32) string {
	return strconv.FormatUint(uint64(id), 10)
}
