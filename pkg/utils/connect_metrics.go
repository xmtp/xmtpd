package utils

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
)

// ConnectClientMetrics holds Prometheus metrics for Connect client calls.
type ConnectClientMetrics struct {
	clientStartedCounter   *prometheus.CounterVec
	clientHandledCounter   *prometheus.CounterVec
	clientHandledHistogram *prometheus.HistogramVec
}

// NewConnectClientMetrics creates a new ConnectClientMetrics with default buckets.
func NewConnectClientMetrics() *ConnectClientMetrics {
	return &ConnectClientMetrics{
		clientStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "connect_client_started_total",
				Help: "Total number of RPCs started on the client.",
			},
			[]string{"connect_type", "connect_service", "connect_method"},
		),
		clientHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "connect_client_handled_total",
				Help: "Total number of RPCs completed on the client, regardless of success or failure.",
			},
			[]string{"connect_type", "connect_service", "connect_method", "connect_code"},
		),
		clientHandledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "connect_client_handling_seconds",
				Help:    "Histogram of response latency (seconds) of the Connect call.",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
			},
			[]string{"connect_type", "connect_service", "connect_method", "connect_code"},
		),
	}
}

// Describe implements prometheus.Collector.
func (m *ConnectClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.clientStartedCounter.Describe(ch)
	m.clientHandledCounter.Describe(ch)
	m.clientHandledHistogram.Describe(ch)
}

// Collect implements prometheus.Collector.
func (m *ConnectClientMetrics) Collect(ch chan<- prometheus.Metric) {
	m.clientStartedCounter.Collect(ch)
	m.clientHandledCounter.Collect(ch)
	m.clientHandledHistogram.Collect(ch)
}

// Interceptor returns a connect.Interceptor that implements both WrapUnary and WrapStreamingClient.
func (m *ConnectClientMetrics) Interceptor() connect.Interceptor {
	return &connectMetricsInterceptor{metrics: m}
}

type connectMetricsInterceptor struct {
	metrics *ConnectClientMetrics
}

// WrapUnary implements connect.Interceptor.
func (i *connectMetricsInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		startTime := time.Now()
		service, method := splitProcedure(req.Spec().Procedure)
		streamType := "unary"

		i.metrics.clientStartedCounter.WithLabelValues(streamType, service, method).Inc()

		// Call the next handler
		resp, err := next(ctx, req)

		// Record the metrics
		duration := time.Since(startTime).Seconds()
		code := getConnectCode(err)

		i.metrics.clientHandledCounter.WithLabelValues(streamType, service, method, code).Inc()
		i.metrics.clientHandledHistogram.WithLabelValues(streamType, service, method, code).
			Observe(duration)

		return resp, err
	}
}

// WrapStreamingClient implements connect.Interceptor.
func (i *connectMetricsInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		service, method := splitProcedure(spec.Procedure)
		streamType := getStreamTypeName(spec.StreamType)

		i.metrics.clientStartedCounter.WithLabelValues(streamType, service, method).Inc()

		startTime := time.Now()
		conn := next(ctx, spec)

		return &metricsStreamingClientConn{
			StreamingClientConn: conn,
			metrics:             i.metrics,
			startTime:           startTime,
			service:             service,
			method:              method,
			streamType:          streamType,
		}
	}
}

// WrapStreamingHandler implements connect.Interceptor (unused for client).
func (i *connectMetricsInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return next
}

type metricsStreamingClientConn struct {
	connect.StreamingClientConn
	metrics    *ConnectClientMetrics
	startTime  time.Time
	service    string
	method     string
	streamType string
}

func (m *metricsStreamingClientConn) CloseResponse() error {
	err := m.StreamingClientConn.CloseResponse()

	// Record metrics when the stream closes
	duration := time.Since(m.startTime).Seconds()
	code := getConnectCode(err)

	m.metrics.clientHandledCounter.WithLabelValues(m.streamType, m.service, m.method, code).Inc()
	m.metrics.clientHandledHistogram.WithLabelValues(m.streamType, m.service, m.method, code).
		Observe(duration)

	return err
}

// splitProcedure splits a Connect procedure like "/package.Service/Method" into service and method.
func splitProcedure(procedure string) (service, method string) {
	// Procedure format: /package.Service/Method
	// We want to extract "Service" and "Method"
	if len(procedure) == 0 {
		return "unknown", "unknown"
	}

	// Remove leading slash
	if procedure[0] == '/' {
		procedure = procedure[1:]
	}

	// Find the last slash
	for i := len(procedure) - 1; i >= 0; i-- {
		if procedure[i] == '/' {
			servicePath := procedure[:i]
			method = procedure[i+1:]

			// Extract just the service name from the full path
			// e.g., "xmtpv4.message_api.ReplicationApi" -> "ReplicationApi"
			for j := len(servicePath) - 1; j >= 0; j-- {
				if servicePath[j] == '.' {
					service = servicePath[j+1:]
					return service, method
				}
			}
			service = servicePath
			return service, method
		}
	}

	return procedure, "unknown"
}

// getConnectCode extracts the Connect error code or returns "OK" for success.
func getConnectCode(err error) string {
	if err == nil {
		return "OK"
	}

	var connectErr *connect.Error
	if errors.As(err, &connectErr) {
		return connectErr.Code().String()
	}

	return "Unknown"
}

// getStreamTypeName converts Connect stream type to a string label.
func getStreamTypeName(streamType connect.StreamType) string {
	switch streamType {
	case connect.StreamTypeUnary:
		return "unary"
	case connect.StreamTypeClient:
		return "client_stream"
	case connect.StreamTypeServer:
		return "server_stream"
	case connect.StreamTypeBidi:
		return "bidi_stream"
	default:
		return "unknown"
	}
}
