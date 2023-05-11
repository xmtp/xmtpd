package e2e

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Metrics struct {
	runDuration                  *prometheus.HistogramVec
	subscribeDuration            *prometheus.HistogramVec
	publishDuration              *prometheus.HistogramVec
	subscribeConvergenceDuration *prometheus.HistogramVec
	queryConvergenceDuration     *prometheus.HistogramVec
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
		subscribeDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "e2e",
				Name:      "subscribe_duration_us",
				Help:      "duration of test case subscribe (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"test", "node", "status"},
		),
		publishDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "e2e",
				Name:      "publish_duration_us",
				Help:      "duration of test case publish (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"test", "node", "status"},
		),
		subscribeConvergenceDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "e2e",
				Name:      "subscribe_convergence_duration_us",
				Help:      "duration of test case subscribe convergence (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"test", "node", "status"},
		),
		queryConvergenceDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "e2e",
				Name:      "query_convergence_duration_us",
				Help:      "duration of test case query convergence (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"test", "node", "status"},
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

func (m *Metrics) recordSubscribe(ctx context.Context, test, node, status string, duration time.Duration) {
	if m == nil || m.subscribeDuration == nil {
		return
	}
	met, err := m.subscribeDuration.GetMetricWithLabelValues(test, node, status)
	if err != nil {
		ctx.Logger().Warn("error observing metric",
			zap.Error(err),
			zap.String("metric", "subscribe_duration_us"),
			zap.String("test", test),
			zap.String("node", node),
			zap.String("status", status),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}

func (m *Metrics) recordPublish(ctx context.Context, test, node, status string, duration time.Duration) {
	if m == nil || m.publishDuration == nil {
		return
	}
	met, err := m.publishDuration.GetMetricWithLabelValues(test, node, status)
	if err != nil {
		ctx.Logger().Warn("error observing metric",
			zap.Error(err),
			zap.String("metric", "publish_duration_us"),
			zap.String("test", test),
			zap.String("node", node),
			zap.String("status", status),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}

func (m *Metrics) recordSubscribeConvergence(ctx context.Context, test, node, status string, duration time.Duration) {
	if m == nil || m.subscribeConvergenceDuration == nil {
		return
	}
	met, err := m.subscribeConvergenceDuration.GetMetricWithLabelValues(test, node, status)
	if err != nil {
		ctx.Logger().Warn("error observing metric",
			zap.Error(err),
			zap.String("metric", "subscribe_convergence_duration_us"),
			zap.String("test", test),
			zap.String("node", node),
			zap.String("status", status),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}

func (m *Metrics) recordQueryConvergence(ctx context.Context, test, node, status string, duration time.Duration) {
	if m == nil || m.queryConvergenceDuration == nil {
		return
	}
	met, err := m.queryConvergenceDuration.GetMetricWithLabelValues(test, node, status)
	if err != nil {
		ctx.Logger().Warn("error observing metric",
			zap.Error(err),
			zap.String("metric", "query_convergence_duration_us"),
			zap.String("test", test),
			zap.String("node", node),
			zap.String("status", status),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}
