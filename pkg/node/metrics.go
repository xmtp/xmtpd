package node

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Metrics struct {
	// Records syncer fetch durations in microseconds by topic_type.
	syncFetchHistogram *prometheus.HistogramVec

	Api      *gateway.Metrics
	Replicas *crdt.Metrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		Api:      gateway.NewMetrics(),
		Replicas: crdt.NewMetrics(),
		syncFetchHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "sync",
				Name:      "fetch_duration_us",
				Help:      "duration of fetch requests from replica syncers (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			[]string{"topic_type"},
		),
	}
}

func (m *Metrics) recordFetch(ctx context.Context, topic string, duration time.Duration) {
	if m == nil || m.syncFetchHistogram == nil {
		return
	}
	topic_type := utils.CategoryFromTopic(topic)
	met, err := m.syncFetchHistogram.GetMetricWithLabelValues(topic_type)
	if err != nil {
		ctx.Logger().Warn("metric observe",
			zap.Error(err),
			zap.String("metric", "fetch_duration_ms"),
			zap.String("topic_type", topic_type),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}
