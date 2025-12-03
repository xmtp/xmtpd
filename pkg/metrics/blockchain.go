package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var blockchainWaitForTransaction = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_blockchain_wait_for_transaction_seconds",
		Help:    "Time spent waiting for transaction receipt",
		Buckets: []float64{0.01, 0.05, 0.075, 0.1, 0.25, 0.5},
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

var blockchainGasPriceGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "xmtp_blockchain_oracle_gas_price",
	Help: "Current gas price in wei",
}, []string{"chain_id"})

func EmitBlockchainGasPrice(chainID int64, gasPrice uint64) {
	blockchainGasPriceGauge.WithLabelValues(strconv.FormatInt(chainID, 10)).
		Set(float64(gasPrice))
}

var blockchainGasPriceUpdatesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "xmtp_blockchain_oracle_gas_price_updates_total",
	Help: "Total number of gas price updates",
}, []string{"chain_id"})

func EmitBlockchainGasPriceUpdatesTotal(chainID int64) {
	blockchainGasPriceUpdatesTotal.WithLabelValues(strconv.FormatInt(chainID, 10)).
		Inc()
}

var blockchainGasPriceDefaultFallbackTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "xmtp_blockchain_oracle_gas_price_default_fallback_total",
	Help: "Total times default gas price was used due to staleness",
}, []string{"chain_id"})

func EmitBlockchainGasPriceDefaultFallbackTotal(chainID int64) {
	blockchainGasPriceDefaultFallbackTotal.WithLabelValues(strconv.FormatInt(chainID, 10)).
		Inc()
}

var blockchainGasPriceLastUpdateTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "xmtp_blockchain_oracle_gas_price_last_update_timestamp_unix",
	Help: "Unix timestamp of last gas price update",
}, []string{"chain_id"})

func EmitBlockchainGasPriceLastUpdateTimestamp(chainID int64, timestamp int64) {
	blockchainGasPriceLastUpdateTimestamp.WithLabelValues(strconv.FormatInt(chainID, 10)).
		Set(float64(timestamp))
}
