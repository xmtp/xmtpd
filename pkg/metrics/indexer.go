package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var indexerNumLogsFound = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtpd_indexer_log_streamer_logs",
		Help: "Number of logs found by the log streamer",
	},
	[]string{"contract_address"},
)

var indexerCurrentBlock = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "xmtpd_indexer_log_streamer_current_block",
		Help: "Current block being processed by the log streamer",
	},
	[]string{"contract_address"},
)

var indexerGetLogsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "xmtpd_indexer_log_streamer_get_logs_duration",
		Help:    "Duration of the get logs call",
		Buckets: []float64{1, 10, 100, 500, 1000, 5000, 10000, 50000, 100000},
	},
	[]string{"contract_address"},
)

var indexerGetLogsRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "xmtpd_indexer_log_streamer_get_logs_requests",
		Help: "Number of get logs requests",
	},
	[]string{"contract_address", "success"},
)

func EmitIndexerNumLogsFound(contractAddress string, numLogs int) {
	indexerNumLogsFound.With(prometheus.Labels{"contract_address": contractAddress}).
		Add(float64(numLogs))
}

func EmitIndexerCurrentBlock(contractAddress string, block int) {
	indexerCurrentBlock.With(prometheus.Labels{"contract_address": contractAddress}).
		Set(float64(block))
}

func EmitIndexerGetLogsDuration(contractAddress string, duration time.Duration) {
	indexerGetLogsDuration.With(prometheus.Labels{"contract_address": contractAddress}).
		Observe(float64(duration.Milliseconds()))
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
