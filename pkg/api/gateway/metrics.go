package gateway

import (
	gocontext "context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap/zapcore"
)

type Metrics struct {
	// Records API request durations in microseconds by various dimensions
	// e.g. grpc_method, grpc_error_code, client name/version, app name/version.
	apiRequestHistogram *prometheus.HistogramVec
}

var requestDurationLabels = []string{
	"grpc_method", "grpc_error_code", "app_client", "app_client_version", "api_app", "api_app_version",
}

func NewMetrics() *Metrics {
	return &Metrics{
		apiRequestHistogram: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "xmtpd",
				Subsystem: "api",
				Name:      "request_duration_us",
				Help:      "duration of api requests (microseconds)",
				Buckets:   prometheus.ExponentialBuckets(10, 10, 10),
			},
			requestDurationLabels,
		),
	}
}

func (m *Metrics) recordRequest(ctx gocontext.Context, duration time.Duration, fields []zapcore.Field) {
	if m == nil || m.apiRequestHistogram == nil {
		return
	}
	met, err := m.apiRequestHistogram.GetMetricWithLabelValues(requestDurationLabelValuesFromFields(fields)...)
	if err != nil {
		return
	}
	met.Observe(float64(duration.Microseconds()))
}

func requestDurationLabelValuesFromFields(fields []zapcore.Field) (values []string) {
	for _, label := range requestDurationLabels {
		var val string
		for _, field := range fields {
			if field.Key == label {
				val = field.String
				break
			}
		}
		values = append(values, val)
	}
	return values
}
