package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var blockchainWaitForTransaction = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_blockchain_wait_for_transaction_seconds",
		Help: "Time spent waiting for transaction receipt",
		Buckets: []float64{
			0.05,
			0.1,
			0.15,
			0.2,
			0.25,
			0.3,
			0.35,
			0.4,
			0.45,
			0.5,
			0.55,
			0.6,
			0.65,
			0.7,
			0.75,
			0.8,
			0.85,
			0.9,
			1.0,
			2.5,
			5,
			10,
		},
	},
	[]string{"status"},
)

func EmitBlockchainWaitForTransaction(status string, duration float64) {
	blockchainWaitForTransaction.With(prometheus.Labels{"status": status}).
		Observe(duration)
}

func MeasureWaitForTransaction[Return any](
	fn func() (Return, error),
) (Return, error) {
	var (
		start    = time.Now()
		status   = "success"
		ret, err = fn()
	)

	if err != nil {
		status = "failure"
	}

	defer func() {
		EmitBlockchainWaitForTransaction(status, time.Since(start).Seconds())
	}()

	return ret, err
}

var blockchainPublishPayload = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_blockchain_publish_payload_seconds",
		Help: "Time to publish a payload to the blockchain",
		Buckets: []float64{
			0.05,
			0.1,
			0.15,
			0.2,
			0.25,
			0.3,
			0.35,
			0.4,
			0.45,
			0.5,
			0.55,
			0.6,
			0.65,
			0.7,
			0.75,
			0.8,
			0.85,
			0.9,
			1.0,
			2.5,
			5,
			10,
		},
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

var blockchainBroadcastTransaction = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "xmtp_blockchain_broadcast_transaction_seconds",
		Help: "Time to publish a payload to the blockchain",
		Buckets: []float64{
			0.05,
			0.1,
			0.15,
			0.2,
			0.25,
			0.3,
			0.35,
			0.4,
			0.45,
			0.5,
			0.55,
			0.6,
			0.65,
			0.7,
			0.75,
			0.8,
			0.85,
			0.9,
			1.0,
			2.5,
			5,
			10,
		},
	},
	[]string{"payload_type", "status"},
)

func EmitBlockchainBroadcastTransaction(
	payloadType string,
	status string,
	duration time.Duration,
) {
	blockchainBroadcastTransaction.With(prometheus.Labels{
		"payload_type": payloadType,
		"status":       status,
	}).
		Observe(duration.Seconds())
}

func MeasureBroadcastTransaction[Return any](
	payloadType string,
	fn func() (Return, error),
) (Return, error) {
	var (
		start    = time.Now()
		status   = "success"
		ret, err = fn()
	)

	if err != nil {
		status = "failure"
	}

	defer func() {
		EmitBlockchainBroadcastTransaction(payloadType, status, time.Since(start))
	}()

	return ret, err
}
