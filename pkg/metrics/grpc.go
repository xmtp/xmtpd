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

// GRPCType represents the type of RPC call.
type GRPCType string

const (
	GRPCTypeUnary        GRPCType = "unary"
	GRPCTypeServerStream GRPCType = "server_stream"
	GRPCTypeClientStream GRPCType = "client_stream"
	GRPCTypeBidiStream   GRPCType = "bidi_stream"
)

// StreamTypeToGRPCType converts connect.StreamType to our GRPCType.
func StreamTypeToGRPCType(streamType connect.StreamType) GRPCType {
	switch streamType {
	case connect.StreamTypeUnary:
		return GRPCTypeUnary
	case connect.StreamTypeServer:
		return GRPCTypeServerStream
	case connect.StreamTypeClient:
		return GRPCTypeClientStream
	case connect.StreamTypeBidi:
		return GRPCTypeBidiStream
	default:
		return GRPCTypeUnary
	}
}

// EmitGRPCServerStarted increments the started counter for an RPC.
func EmitGRPCServerStarted(grpcType GRPCType, service, method string) {
	grpcServerStartedTotal.With(prometheus.Labels{
		"grpc_type":    string(grpcType),
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerHandled increments the handled counter for a completed RPC.
func EmitGRPCServerHandled(grpcType GRPCType, service, method string, code connect.Code) {
	grpcServerHandledTotal.With(prometheus.Labels{
		"grpc_type":    string(grpcType),
		"grpc_service": service,
		"grpc_method":  method,
		"grpc_code":    code.String(),
	}).Inc()
}

// EmitGRPCServerMsgReceived increments the received message counter.
func EmitGRPCServerMsgReceived(grpcType GRPCType, service, method string) {
	grpcServerMsgReceivedTotal.With(prometheus.Labels{
		"grpc_type":    string(grpcType),
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerMsgSent increments the sent message counter.
func EmitGRPCServerMsgSent(grpcType GRPCType, service, method string) {
	grpcServerMsgSentTotal.With(prometheus.Labels{
		"grpc_type":    string(grpcType),
		"grpc_service": service,
		"grpc_method":  method,
	}).Inc()
}

// EmitGRPCServerHandlingTime records the handling duration for an RPC.
func EmitGRPCServerHandlingTime(grpcType GRPCType, service, method string, duration time.Duration) {
	grpcServerHandlingSeconds.With(prometheus.Labels{
		"grpc_type":    string(grpcType),
		"grpc_service": service,
		"grpc_method":  method,
	}).Observe(duration.Seconds())
}

// GRPCServerMetrics returns all gRPC server metrics for registration.
func GRPCServerMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		grpcServerStartedTotal,
		grpcServerHandledTotal,
		grpcServerMsgReceivedTotal,
		grpcServerMsgSentTotal,
		grpcServerHandlingSeconds,
	}
}
