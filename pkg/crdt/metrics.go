package crdt

import (
	"sync"
	"time"

	"github.com/multiformats/go-multihash"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Metrics struct {
	// Counts events received for Replicas.
	// This can be either from a broadcast (action=broadcast),
	// or from a syncer fetch response (action=sync)
	receivedEventCounter *prometheus.CounterVec

	// Samples free space in internal Replica channels (1-len(ch)/cap(ch)) by channel type.
	// The recorded value is the percentage of free space in the channel represented as an int 0-100.
	// The samples are taken whenever channels are being sent to.
	// The bucket counts reflect the number of samples taken, not the number of channels.
	channelFreeSpaceHistogram *prometheus.HistogramVec

	// used to throttle warnings
	lastWarning     time.Time
	lastWarningLock sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		receivedEventCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "xmtpd",
				Subsystem: "crdt",
				Name:      "received_events",
				Help:      "received event counter",
			},
			[]string{"topic_type", "action"},
		),
		channelFreeSpaceHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "crdt",
				Name:      "channel_free_space",
				Help:      "samples free space in replica channels (cap - len) by channel type",
				Buckets:   prometheus.LinearBuckets(10, 10, 10),
			},
			[]string{"channel"},
		),
	}
}

func (m *Metrics) recordReceivedEvent(ctx context.Context, ev *types.Event, sync bool) {
	if m == nil || m.receivedEventCounter == nil {
		return
	}
	action := "sync"
	if !sync {
		action = "broadcast"
	}
	topic_type := utils.CategoryFromTopic(ev.ContentTopic)
	met, err := m.receivedEventCounter.GetMetricWithLabelValues(topic_type, action)
	if err != nil {
		ctx.Logger().Warn("metric observe",
			zap.Error(err),
			zap.String("metric", "receive_events"),
			zap.String("topic_type", topic_type),
			zap.String("action", action),
		)
		return
	}
	met.Add(1)
}

func (m *Metrics) recordFreeSpaceInLinks(ctx context.Context, ch chan multihash.Multihash) {
	if m == nil || m.channelFreeSpaceHistogram == nil {
		return
	}
	percentFull := len(ch) * 100 / cap(ch)
	m.recordFreeSpace(ctx, "links", percentFull)
}

func (m *Metrics) recordFreeSpaceInEvents(ctx context.Context, ch chan *types.Event, sync bool) {
	if m == nil || m.channelFreeSpaceHistogram == nil {
		return
	}
	chType := "sync_events"
	if !sync {
		chType = "cast_events"
	}
	percentFull := len(ch) * 100 / cap(ch)
	m.recordFreeSpace(ctx, chType, percentFull)
}

func (m *Metrics) recordFreeSpace(ctx context.Context, chType string, percentFull int) {
	if percentFull > 80 {
		m.warnChannelFillingUp(ctx, chType, percentFull)
	}
	met, err := m.channelFreeSpaceHistogram.GetMetricWithLabelValues(chType)
	if err != nil {
		ctx.Logger().Warn("metric observe",
			zap.Error(err),
			zap.String("metric", "channel_free_space"),
			zap.String("channel", chType),
		)
		return
	}
	met.Observe(float64(100 - percentFull))
}

func (m *Metrics) warnChannelFillingUp(ctx context.Context, chType string, percentFull int) {
	m.lastWarningLock.Lock()
	defer m.lastWarningLock.Unlock()
	if time.Since(m.lastWarning) > time.Second {
		m.lastWarning = time.Now()
		ctx.Logger().Warn("channel filling up", zap.String("channel", chType), zap.Percent("fullness", percentFull))
	}
}
