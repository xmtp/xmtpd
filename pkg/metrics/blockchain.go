package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var blockchainWaitForTransaction = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "xmtp_blockchain_wait_for_transaction_seconds",
		Help:    "Time in seconds to wait before receiving transaction",
		Buckets: []float64{0.01, 0.05, 0.075, 0.1, 0.25, 0.5},
	},
)

func EmitBlockchainWaitForTransaction(duration float64) {
	blockchainWaitForTransaction.Observe(duration)
}

var blockchainPublish = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtpd_blockchain_publish_seconds",
		Help: "Duration of the get logs call",
	},
	[]string{"payload_type"},
)

func EmitBlockchainPublish(payloadType string, duration time.Duration) {
	blockchainPublish.With(prometheus.Labels{"payload_type": payloadType}).
		Observe(duration.Seconds())
}

func MeasurePublishToBlockchainMethod[Return any](payloadType string, fn func() (Return, error)) (Return, error) {
	start := time.Now()
	defer func() {
		EmitBlockchainPublish(payloadType, time.Since(start))
		fmt.Printf("Publishing took %f seconds\n", time.Since(start).Seconds())
	}()
	return fn()
}
