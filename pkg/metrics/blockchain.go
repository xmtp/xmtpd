package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var blockchainWaitForTransaction = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "xmtp_blockchain_wait_for_transaction_seconds",
		Help:    "Time spent waiting for transaction receipt",
		Buckets: []float64{0.01, 0.05, 0.075, 0.1, 0.25, 0.5},
	},
)

func EmitBlockchainWaitForTransaction(duration float64) {
	blockchainWaitForTransaction.Observe(duration)
}

var blockchainPublishPayload = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtpd_blockchain_publish_payload_seconds",
		Help: "Time to publish a payload to the blockchain",
	},
	[]string{"payload_type"},
)

func EmitBlockchainPublish(payloadType string, duration time.Duration) {
	blockchainPublishPayload.With(prometheus.Labels{"payload_type": payloadType}).
		Observe(duration.Seconds())
}

func MeasurePublishToBlockchainMethod[Return any](
	payloadType string,
	fn func() (Return, error),
) (Return, error) {
	start := time.Now()
	defer func() {
		EmitBlockchainPublish(payloadType, time.Since(start))
	}()
	return fn()
}
