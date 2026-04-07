package ratelimiter

import "github.com/prometheus/client_golang/prometheus"

var (
	DecisionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_decisions_total",
			Help: "Rate-limit decisions broken down by service, method, tier, and outcome",
		},
		[]string{"service", "method", "tier", "outcome"},
	)
	BreakerStateGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "xmtpd_rate_limit_circuit_breaker_state",
			Help: "Circuit breaker state: 0=closed, 1=half_open, 2=open",
		},
	)
	BreakerTripsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_circuit_breaker_trips_total",
			Help: "Number of times the circuit breaker has tripped open",
		},
	)
	StreamTerminationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "xmtpd_rate_limit_stream_terminations_total",
			Help: "Subscribe stream terminations broken down by reason",
		},
		[]string{"reason"},
	)
)

// Register registers all rate-limit metrics with the provided registry.
// Safe to call once per process. Re-registration errors are ignored.
func Register(reg prometheus.Registerer) {
	for _, c := range []prometheus.Collector{
		DecisionsTotal, BreakerStateGauge, BreakerTripsTotal, StreamTerminationsTotal,
	} {
		_ = reg.Register(c)
	}
}
