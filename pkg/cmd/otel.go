package cmd

import (
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/context"
	"github.com/xmtp/xmtpd/pkg/zap"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type OpenTelemetryOptions struct {
	CollectorAddress string `long:"collector-address" env:"OTEL_COLLECTOR_ADDRESS" default:"localhost"`
	CollectorPort    uint   `long:"collector-port" env:"OTEL_COLLECTOR_PORT" default:"4317"`
}

type openTelemetry struct {
	ctx context.Context
	log *zap.Logger

	traceExporter   *otlptrace.Exporter
	traceProvider   *sdktrace.TracerProvider
	metricsExporter sdkmetric.Exporter
	metricsProvider *sdkmetric.MeterProvider
	meter           metric.Meter
}

func newOpenTelemetry(ctx context.Context, opts *OpenTelemetryOptions) (*openTelemetry, error) {
	ot := &openTelemetry{
		ctx: ctx,
		log: ctx.Logger().Named("otel"),
	}
	collectorEndpoint := fmt.Sprintf("%s:%d", opts.CollectorAddress, opts.CollectorPort)

	extraResources, err := sdkresource.New(
		ctx,
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			attribute.String("service.name", "xmtpd"),
		),
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
		otlptracegrpc.WithEndpoint(collectorEndpoint),
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
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
	)
	if err != nil {
		return nil, err
	}
	ot.metricsProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(ot.metricsExporter)),
		sdkmetric.WithResource(resource),
	)
	global.SetMeterProvider(ot.metricsProvider)

	ot.meter = ot.metricsProvider.Meter("xmtpd")

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		return nil, err
	}

	return ot, nil
}

func (ot *openTelemetry) Close() error {
	err := ot.traceExporter.Shutdown(ot.ctx)
	if err != nil && err != context.Canceled {
		ot.log.Error("error shutting down trace exporter", zap.Error(err))
	}

	err = ot.traceProvider.Shutdown(ot.ctx)
	if err != nil && err != context.Canceled {
		ot.log.Error("error shutting down trace provider", zap.Error(err))
	}

	err = ot.metricsExporter.Shutdown(ot.ctx)
	if err != nil && err != context.Canceled {
		ot.log.Error("error shutting down metrics exporter", zap.Error(err))
	}

	err = ot.metricsProvider.Shutdown(ot.ctx)
	if err != nil && err != context.Canceled {
		ot.log.Error("error shutting down metrics provider", zap.Error(err))
	}

	return nil
}
