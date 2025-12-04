package metrics

import (
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
)

// gRPC server metrics following the standard grpc-ecosystem prometheus naming.
// These are compatible with existing gRPC Prometheus dashboards.

var grpcServerStartedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_server_started_total",
		Help: "Total number of RPCs started on the server.",
	},
	[]string{"grpc_type", "grpc_service", "grpc_method"},
)

var grpcServerHandledTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_server_handled_total",
		Help: "Total number of RPCs completed on the server, regardless of success or failure.",
	},
	[]string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"},
)

var grpcServerMsgReceivedTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_server_msg_received_total",
		Help: "Total number of RPC stream messages received on the server.",
	},
	[]string{"grpc_type", "grpc_service", "grpc_method"},
)

var grpcServerMsgSentTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "grpc_server_msg_sent_total",
		Help: "Total number of gRPC stream messages sent by the server.",
	},
	[]string{"grpc_type", "grpc_service", "grpc_method"},
)

var grpcServerHandlingSeconds = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "grpc_server_handling_seconds",
		Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"grpc_type", "grpc_service", "grpc_method"},
)

// EmitGRPCServerStarted increments the started counter for an RPC.
func EmitGRPCServerStarted(grpcType string, service, method string) {
	grpcServerStartedTotal.With(prometheus.Labels{
		"grpc_type":    grpcType,
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerHandled increments the handled counter for a completed RPC.
func EmitGRPCServerHandled(grpcType string, service, method string, code connect.Code) {
	grpcServerHandledTotal.With(prometheus.Labels{
		"grpc_type":    grpcType,
		"grpc_service": service,
		"grpc_method":  method,
		"grpc_code":    code.String(),
	}).Inc()
}

// EmitGRPCServerMsgReceived increments the received message counter.
func EmitGRPCServerMsgReceived(grpcType string, service, method string) {
	grpcServerMsgReceivedTotal.With(prometheus.Labels{
		"grpc_type":    grpcType,
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerMsgSent increments the sent message counter.
func EmitGRPCServerMsgSent(grpcType string, service, method string) {
	grpcServerMsgSentTotal.With(prometheus.Labels{
		"grpc_type":    grpcType,
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerHandlingTime records the handling duration for an RPC.
func EmitGRPCServerHandlingTime(grpcType string, service, method string, duration time.Duration) {
	grpcServerHandlingSeconds.With(prometheus.Labels{
		"grpc_type":    grpcType,
		"grpc_service": service,
		"grpc_method":  method,
	}).Observe(duration.Seconds())
}
