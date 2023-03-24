package node

import (
	"time"

	"github.com/xmtp/xmtpd/pkg/api/gateway"
	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Metrics struct {
	// Records syncer fetch durations in microseconds by topic_type.
	syncFetchHistogram instrument.Int64Histogram

	Api      *gateway.Metrics
	Replicas *crdt.Metrics
}

func NewMetrics(meter metric.Meter) (*Metrics, error) {
	var m Metrics
	var err error
	m.Api, err = gateway.NewMetrics(meter)
	if err != nil {
		return nil, err
	}
	m.Replicas, err = crdt.NewMetrics(meter)
	if err != nil {
		return nil, err
	}
	m.syncFetchHistogram, err = meter.Int64Histogram(
		"xmtpd.sync.fetch_duration_us",
		instrument.WithDescription(`duration of fetch requests from replica syncers (microseconds)`),
		instrument.WithUnit("microsecond"),
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Metrics) recordFetch(ctx context.Context, topic string, duration time.Duration) {
	if m == nil || m.syncFetchHistogram == nil {
		return
	}
	m.syncFetchHistogram.Record(ctx, duration.Microseconds(),
		attribute.String("topic_type", utils.CategoryFromTopic(topic)),
	)
}
