package crdt

import (
	"sync"
	"time"

	"github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/utils"
	"github.com/xmtp/xmtpd/pkg/zap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Metrics struct {
	// Counts events received for Replicas.
	// This can be either from a broadcast (action=broadcast),
	// or from a syncer fetch response (action=sync)
	crdtReceivedEventCounter instrument.Int64Counter

	// Samples free space in internal Replica channels (1-len(ch)/cap(ch)) by channel type.
	// The recorded value is the percentage of free space in the channel represented as an int 0-100.
	// The samples are taken whenever channels are being sent to.
	// The bucket counts reflect the number of samples taken, not the number of channels.
	crdtChannelFreeSpaceHistogram instrument.Int64Histogram

	// used to throttle warnings
	lastWarning     time.Time
	lastWarningLock sync.Mutex
}

func NewMetrics(meter metric.Meter) (m *Metrics, err error) {
	m = &Metrics{}
	m.crdtReceivedEventCounter, err = meter.Int64Counter(
		"xmtpd.crdt.received_events",
		instrument.WithDescription("received event counter"),
	)
	if err != nil {
		return nil, err
	}
	m.crdtChannelFreeSpaceHistogram, err = meter.Int64Histogram(
		"xmtpd.crdt.channel_free_space",
		instrument.WithDescription("samples free space in replica channels (cap - len) by channel type"),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Metrics) recordReceivedEvent(ctx context.Context, ev *types.Event, sync bool) {
	if m == nil || m.crdtReceivedEventCounter == nil {
		return
	}
	action := "sync"
	if !sync {
		action = "broadcast"
	}
	m.crdtReceivedEventCounter.Add(ctx, 1,
		attribute.String("action", action),
		attribute.String("topic_type", utils.CategoryFromTopic(ev.ContentTopic)),
	)
}

func (m *Metrics) recordFreeSpaceInLinks(ctx context.Context, ch chan multihash.Multihash) {
	if m == nil || m.crdtChannelFreeSpaceHistogram == nil {
		return
	}
	percentFull := len(ch) * 100 / cap(ch)
	if percentFull > 80 {
		m.warnChannelFillingUp(ctx, "links", percentFull)

	}
	m.crdtChannelFreeSpaceHistogram.Record(ctx,
		int64(100-percentFull),
		attribute.String("channel", "links"),
	)
}

func (m *Metrics) recordFreeSpaceInEvents(ctx context.Context, ch chan *types.Event, sync bool) {
	if m == nil || m.crdtChannelFreeSpaceHistogram == nil {
		return
	}
	chType := "sync_events"
	if !sync {
		chType = "cast_events"
	}
	percentFull := len(ch) * 100 / cap(ch)
	if percentFull > 80 {
		m.warnChannelFillingUp(ctx, chType, percentFull)
	}
	m.crdtChannelFreeSpaceHistogram.Record(ctx,
		int64(100-percentFull),
		attribute.String("channel", chType),
	)
}

func (m *Metrics) warnChannelFillingUp(ctx context.Context, chType string, percentFull int) {
	m.lastWarningLock.Lock()
	defer m.lastWarningLock.Unlock()
	if time.Since(m.lastWarning) > time.Second {
		m.lastWarning = time.Now()
		ctx.Logger().Warn("channel filling up", zap.String("channel", chType), zap.Percent("fullness", percentFull))
	}
}
