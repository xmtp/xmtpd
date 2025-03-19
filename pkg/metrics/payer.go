package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var payerNodePublishDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_payer_node_publish_duration_seconds",
		Help: "Duration of the node publish call",
	},
	[]string{"originator_id"},
)

func EmitPayerNodePublishDuration(originatorId uint32, duration float64) {
	payerNodePublishDuration.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Observe(duration)
}

var payerCursorBlockTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_payer_read_own_commit_in_time_seconds",
		Help: "Read your own commit duration in seconds",
	},
	[]string{"originator_id"},
)

func EmitPayerBlockUntilDesiredCursorReached(originatorId uint32, duration float64) {
	payerCursorBlockTime.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Observe(duration)
}

var payerCurrentNonce = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "xmtp_payer_lru_nonce",
		Help: "Least recently used blockchain nonce of the payer (not guaranteed to be the highest nonce).",
	},
)

func EmitPayerCurrentNonce(nonce float64) {
	// Set is thread-safe
	payerCurrentNonce.Set(nonce)
}

var payerBanlistRetry = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_payer_failed_attempts_to_publish_to_node_via_banlist",
		Help:    "Number of failed attempts to publish to a node via banlist",
		Buckets: []float64{0, 1, 2, 3, 4, 5},
	},
	[]string{"originator_id"},
)

func EmitPayerBanlistRetries(originatorId uint32, retries int) {
	payerBanlistRetry.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Observe(float64(retries))
}

var payerMessagesOriginated = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_payer_messages_originated",
		Help: "Number of messages originated by the payer.",
	},
	[]string{"originator_id"},
)

func EmitPayerMessageOriginated(originatorId uint32, count int) {
	payerMessagesOriginated.With(prometheus.Labels{"originator_id": strconv.Itoa(int(originatorId))}).
		Add(float64(count))
}
