package node

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/zap"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type OpenTelemetryOptions struct {
	TraceCollectorEndpoint   string `long:"trace-collector" env:"OTEL_TRACE_COLLECTOR" default:"localhost:4317"`
	MetricsCollectorEndpoint string `long:"metrics-collector" env:"OTEL_METRICS_COLLECTOR" default:"localhost:4317"`
}

type openTelemetry struct {
	ctx context.Context
	log *zap.Logger

	traceExporter   *otlptrace.Exporter
	traceProvider   *sdktrace.TracerProvider
	metricsExporter sdkmetric.Exporter
	metricsProvider *sdkmetric.MeterProvider
}

func newOpenTelemetry(ctx context.Context, log *zap.Logger, opts *OpenTelemetryOptions) (*openTelemetry, error) {
	ot := &openTelemetry{
		ctx: ctx,
		log: log.Named("otel"),
	}

	extraResources, err := sdkresource.New(
		ctx,
		sdkresource.WithOS(),
		sdkresource.WithProcess(),
		sdkresource.WithContainer(),
		sdkresource.WithHost(),
	)
	if err != nil {
		return nil, err
	}
	resource, err := sdkresource.Merge(
		sdkresource.Default(),
		extraResources,
	)
	if err != nil {
		return nil, err
	}

	ot.traceExporter, err = otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(opts.TraceCollectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	ot.traceProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(ot.traceExporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(ot.traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	ot.metricsExporter, err = otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(opts.MetricsCollectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	ot.metricsProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(ot.metricsExporter)),
		sdkmetric.WithResource(resource),
	)
	global.SetMeterProvider(ot.metricsProvider)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		return nil, err
	}

	return ot, nil
}

func (ot *openTelemetry) Close() error {
	err := ot.traceExporter.Shutdown(ot.ctx)
	if err != nil {
		ot.log.Error("error shutting down trace exporter")
	}

	err = ot.traceProvider.Shutdown(ot.ctx)
	if err != nil {
		ot.log.Error("error shutting down trace provider")
	}

	err = ot.metricsExporter.Shutdown(ot.ctx)
	if err != nil {
		ot.log.Error("error shutting down metrics exporter")
	}

	err = ot.metricsProvider.Shutdown(ot.ctx)
	if err != nil {
		ot.log.Error("error shutting down metrics provider")
	}

	return nil
}
