package ratelimiter

import "github.com/prometheus/client_golang/prometheus"

var (
	StreamDecisionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_stream_limit_decisions_total",
			Help: "Stream-limit decisions broken down by service and outcome",
		},
		[]string{"service", "outcome"},
	)
	StreamActiveStreams = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtpd_stream_limit_active_streams",
			Help: "Number of active streams tracked by this process (local count, not Redis)",
		},
	)
)
