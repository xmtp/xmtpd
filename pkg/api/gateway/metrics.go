package gateway

import (
	gocontext "context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Metrics struct {
	// Records API request durations in microseconds by various dimensions
	// e.g. grpc_method, grpc_error_code, client name/version, app name/version.
	apiRequestHistogram instrument.Int64Histogram
}

func NewMetrics(meter metric.Meter) (m *Metrics, err error) {
	m = &Metrics{}
	m.apiRequestHistogram, err = meter.Int64Histogram(
		"xmtpd.api.request_duration_us",
		instrument.WithDescription("duration of API request (microseconds)"),
		instrument.WithUnit("microsecond"),
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Metrics) recordRequest(ctx gocontext.Context, duration time.Duration, attrs ...attribute.KeyValue) {
	if m == nil || m.apiRequestHistogram == nil {
		return
	}
	m.apiRequestHistogram.Record(ctx, duration.Microseconds(), attrs...)
}
