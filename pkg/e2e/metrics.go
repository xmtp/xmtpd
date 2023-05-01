package e2e

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Metrics struct {
	runDuration *prometheus.HistogramVec
}

func newMetrics() *Metrics {
	return &Metrics{
		runDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "e2e",
				Name:      "run_duration_us",
				Help:      "duration of test case run (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"test", "status"},
		),
	}
}

func (m *Metrics) recordRun(ctx context.Context, test, status string, duration time.Duration) {
	if m == nil || m.runDuration == nil {
		return
	}
	met, err := m.runDuration.GetMetricWithLabelValues(test, status)
	if err != nil {
		ctx.Logger().Warn("error observing metric",
			zap.Error(err),
			zap.String("metric", "run_duration_us"),
			zap.String("test", test),
			zap.String("status", status),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}
