package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var gatwayPublishDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_gateway_publish_duration_seconds",
		Help: "Duration of the node publish call",
	},
	[]string{"originator_id"},
)

func EmitGatewayPublishDuration(originatorID uint32, duration float64) {
	gatwayPublishDuration.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorID))}).
		Observe(duration)
}

var gatewayCursorBlockTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_gateway_read_own_commit_in_time_seconds",
		Help: "Read your own commit duration in seconds",
	},
	[]string{"originator_id"},
)

func EmitGatewayBlockUntilDesiredCursorReached(originatorID uint32, duration float64) {
	gatewayCursorBlockTime.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorID))}).
		Observe(duration)
}

var gatewayCurrentNonce = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "xmtp_gateway_lru_nonce",
		Help: "Least recently used blockchain nonce of the gateway (not guaranteed to be the highest nonce).",
	},
)

func EmitGatewayCurrentNonce(nonce float64) {
	// Set is thread-safe
	gatewayCurrentNonce.Set(nonce)
}

var gatewayBanlistRetry = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_gateway_failed_attempts_to_publish_to_node_via_banlist",
		Help:    "Number of failed attempts to publish to a node via banlist",
		Buckets: []float64{0, 1, 2, 3, 4, 5},
	},
	[]string{"originator_id"},
)

func EmitGatewayBanlistRetries(originatorID uint32, retries int) {
	gatewayBanlistRetry.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorID))}).
		Observe(float64(retries))
}

var gatewayMessagesOriginated = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_gateway_messages_originated",
		Help: "Number of messages originated by the gateway.",
	},
	[]string{"originator_id"},
)

func EmitGatewayMessageOriginated(originatorID uint32, count int) {
	gatewayMessagesOriginated.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorID))}).
		Add(float64(count))
}

var gatewayGetNodesAvailableNodes = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "xmtp_gateway_get_nodes_available_nodes",
		Help: "Number of currently available nodes for reader selection",
	},
)

func EmitGatewayGetNodesAvailableNodes(count int) {
	gatewayGetNodesAvailableNodes.Set(float64(count))
}
