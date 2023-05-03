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

	// Total number of topics on the node.
	topicsGauge *prometheus.GaugeVec

	API      *gateway.Metrics
	Replicas *crdt.Metrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		API:      gateway.NewMetrics(),
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
		topicsGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "xmtpd",
				Subsystem: "sync",
				Name:      "topic_count",
				Help:      "total number of topics",
			},
			[]string{"topic_type"},
		),
	}
}

func (m *Metrics) recordFetch(ctx context.Context, topic string, duration time.Duration) {
	if m == nil || m.syncFetchHistogram == nil {
		return
	}
	topicType := utils.CategoryFromTopic(topic)
	met, err := m.syncFetchHistogram.GetMetricWithLabelValues(topicType)
	if err != nil {
		ctx.Logger().Warn("metric observe",
			zap.Error(err),
			zap.String("metric", "fetch_duration_ms"),
			zap.String("topic_type", topicType),
		)
		return
	}
	met.Observe(float64(duration.Microseconds()))
}

func (m *Metrics) recordTopicAdd(ctx context.Context, topic string) {
	m.recordTopicCountChange(ctx, topic, 1)
}

func (m *Metrics) recordTopicRemove(ctx context.Context, topic string) {
	m.recordTopicCountChange(ctx, topic, -1)
}

func (m *Metrics) recordTopicCountChange(ctx context.Context, topic string, change int) {
	if m == nil || m.topicsGauge == nil {
		return
	}
	topicType := utils.CategoryFromTopic(topic)
	met, err := m.topicsGauge.GetMetricWithLabelValues(topicType)
	if err != nil {
		ctx.Logger().Warn("metric observe",
			zap.Error(err),
			zap.String("metric", "topic_count"),
			zap.String("topic_type", topicType),
		)
		return
	}
	met.Add(float64(change))
}
