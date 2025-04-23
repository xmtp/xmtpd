package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var indexerNumLogsFound = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_indexer_log_streamer_logs",
		Help: "Number of logs found by the log streamer",
	},
	[]string{"contract_address"},
)

var indexerCurrentBlock = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_indexer_log_streamer_current_block",
		Help: "Current block being processed by the log streamer",
	},
	[]string{"contract_address"},
)

var indexerMaxBlock = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_indexer_log_streamer_max_block",
		Help: "Max block on the chain to be processed by the log streamer",
	},
	[]string{"contract_address"},
)

var indexerCurrentBlockLag = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtp_indexer_log_streamer_block_lag",
		Help: "Lag between current block and max block",
	},
	[]string{"contract_address"},
)

var indexerCountRetryableStorageErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_indexer_retryable_storage_error_count",
		Help: "Number of retryable storage errors",
	},
	[]string{"contract_address"},
)

var indexerGetLogsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtp_indexer_log_streamer_get_logs_duration",
		Help:    "Duration of the get logs call",
		Buckets: []float64{1, 10, 100, 500, 1000, 5000, 10000, 50000, 100000},
	},
	[]string{"contract_address"},
)

var indexerGetLogsRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtp_indexer_log_streamer_get_logs_requests",
		Help: "Number of get logs requests",
	},
	[]string{"contract_address", "success"},
)

var indexerLogProcessingTime = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "xmtp_indexer_log_processing_time_seconds",
		Help: "Time to process a blockchain log",
	},
)

func EmitIndexerNumLogsFound(contractAddress string, numLogs int) {
	indexerNumLogsFound.With(prometheus.Labels{"contract_address": contractAddress}).
		Add(float64(numLogs))
}

func EmitIndexerCurrentBlock(contractAddress string, block uint64) {
	indexerCurrentBlock.With(prometheus.Labels{"contract_address": contractAddress}).
		Set(float64(block))
}

func EmitIndexerMaxBlock(contractAddress string, block uint64) {
	indexerMaxBlock.With(prometheus.Labels{"contract_address": contractAddress}).
		Set(float64(block))
}

func EmitIndexerGetLogsDuration(contractAddress string, duration time.Duration) {
	indexerGetLogsDuration.With(prometheus.Labels{"contract_address": contractAddress}).
		Observe(float64(duration.Milliseconds()))
}

func EmitIndexerCurrentBlockLag(contractAddress string, lag uint64) {
	indexerCurrentBlockLag.With(prometheus.Labels{"contract_address": contractAddress}).
		Set(float64(lag))
}

func EmitIndexerRetryableStorageError(contractAddress string) {
	indexerCountRetryableStorageErrors.With(prometheus.Labels{"contract_address": contractAddress}).
		Inc()
}

func MeasureGetLogs[Return any](contractAddress string, fn func() (Return, error)) (Return, error) {
	start := time.Now()
	ret, err := fn()
	if err == nil {
		EmitIndexerGetLogsDuration(contractAddress, time.Since(start))
	}
	indexerGetLogsRequests.With(prometheus.Labels{"contract_address": contractAddress, "success": strconv.FormatBool(err == nil)}).
		Inc()
	return ret, err
}

func EmitIndexerLogProcessingTime(duration time.Duration) {
	indexerLogProcessingTime.Observe(duration.Seconds())
}
