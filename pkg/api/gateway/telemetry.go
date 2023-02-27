package gateway

import (
	"context"
	"strings"

	"github.com/xmtp/xmtpd/pkg/zap"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	clientVersionMetadataKey = "x-client-version"
	appVersionMetadataKey    = "x-app-version"
)

type TelemetryInterceptor struct {
	log *zap.Logger

	metrics        metric.Meter
	requestCounter instrument.Int64Counter
}

func NewTelemetryInterceptor(log *zap.Logger) (*TelemetryInterceptor, error) {
	metrics := global.MeterProvider().Meter("xmtpd")
	requestCounter, err := metrics.Int64Counter(
		"xmtpd.api_requests",
		instrument.WithDescription("API requests"),
	)
	if err != nil {
		return nil, err
	}

	return &TelemetryInterceptor{
		log: log,

		metrics:        metrics,
		requestCounter: requestCounter,
	}, nil
}

func (ti *TelemetryInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		res, err := handler(ctx, req)
		ti.record(ctx, info.FullMethod, err)
		return res, err
	}
}

func (ti *TelemetryInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		res := handler(srv, stream)
		ti.record(stream.Context(), info.FullMethod, nil)
		return res
	}
}

func (ti *TelemetryInterceptor) record(ctx context.Context, fullMethod string, err error) {
	serviceName, methodName := splitMethodName(fullMethod)
	if serviceName == "grpc.health.v1.Health" {
		return
	}

	md, _ := metadata.FromIncomingContext(ctx)
	clientName, _, clientVersion := parseVersionHeaderValue(md.Get(clientVersionMetadataKey))
	appName, _, appVersion := parseVersionHeaderValue(md.Get(appVersionMetadataKey))
	fields := zap.Fields(
		zap.String("grpc_service", serviceName),
		zap.String("grpc_method", methodName),
		zap.String("api_client", clientName),
		zap.String("api_client_version", clientVersion),
		zap.String("api_app", appName),
		zap.String("api_app_version", appVersion),
	)

	if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
		// There are potentially multiple comma separated IPs bundled in that first value
		ips := strings.Split(ips[0], ",")
		fields = append(fields, zap.String("client_ip", strings.TrimSpace(ips[0])))
	}

	logFn := ti.log.Debug

	if err != nil {
		logFn = ti.log.Info
		fields = append(fields, zap.Error(err))
		grpcErr, _ := status.FromError(err)
		if grpcErr != nil {
			errCode := grpcErr.Code().String()
			fields = append(fields,
				zap.String("grpc_error_code", errCode),
				zap.String("grpc_error_message", grpcErr.Message()),
			)
		}
	}

	logFn("api request", fields...)

	attrs := make([]attribute.KeyValue, 0, len(fields))
	exclude := map[string]bool{
		"grpc_error_code": true,
		"client_ip":       true,
	}
	for _, field := range fields {
		if exclude[field.Key] {
			continue
		}
		attrs = append(attrs, attribute.String(field.Key, field.String))
	}
	ti.requestCounter.Add(ctx, 1, attrs...)
}

func splitMethodName(fullMethodName string) (serviceName string, methodName string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}

func parseVersionHeaderValue(vals []string) (name string, version string, full string) {
	if len(vals) == 0 {
		return
	}
	full = vals[0]
	parts := strings.Split(full, "/")
	if len(parts) > 0 {
		name = parts[0]
		if len(parts) > 1 {
			version = parts[1]
		}
	}
	return
}
