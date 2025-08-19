// Package tracing provides OpenTelemetry tracing capabilities for XMTPD
package tracing

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.28.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	defaultServiceName    = "xmtpd"
	defaultInstrumentName = "github.com/xmtp/xmtpd"
)

// OTelConfig holds configuration for OpenTelemetry tracing
type OTelConfig struct {
	ServiceName        string
	ServiceVersion     string
	Environment        string
	OTLPEndpoint       string
	OTLPHeaders        map[string]string
	Enabled            bool
	SamplingRatio      float64
	SamplerType        string   // "ratio", "always_on", "always_off", "parent_based" (default)
	Propagators        []string // Propagators to use: "tracecontext", "baggage", "b3", "jaeger"
	BatchTimeout       time.Duration
	MaxExportBatchSize int
	MaxQueueSize       int
	ShutdownTimeout    time.Duration // Timeout for shutdown operations
	Insecure           bool          // Whether to use insecure HTTP transport
	UseSimpleProcessor bool          // Whether to use SimpleSpanProcessor instead of BatchSpanProcessor
	UseStdout          bool          // Whether to export traces to stdout (for development)
	IncludeHostInfo    bool          // Whether to include host information in resource
	IncludeProcessInfo bool          // Whether to include process information in resource
	StrictMode         bool          // Whether to fail fast on exporter creation errors
}

// Validate validates the OpenTelemetry configuration and returns any errors
func (c *OTelConfig) Validate() []error {
	var errors []error

	if c.ServiceName == "" {
		errors = append(errors, fmt.Errorf("service name cannot be empty"))
	}

	if c.SamplingRatio < 0.0 || c.SamplingRatio > 1.0 {
		errors = append(errors, fmt.Errorf(
			"sampling ratio must be between 0.0 and 1.0, got %f",
			c.SamplingRatio,
		))
	}

	if c.BatchTimeout <= 0 {
		errors = append(errors, fmt.Errorf(
			"batch timeout must be positive, got %v",
			c.BatchTimeout,
		))
	}

	if c.MaxExportBatchSize <= 0 {
		errors = append(errors, fmt.Errorf(
			"max export batch size must be positive, got %d",
			c.MaxExportBatchSize,
		))
	}

	if c.MaxQueueSize <= 0 {
		errors = append(errors, fmt.Errorf(
			"max queue size must be positive, got %d",
			c.MaxQueueSize,
		))
	}

	if c.ShutdownTimeout <= 0 {
		errors = append(errors, fmt.Errorf(
			"shutdown timeout must be positive, got %v",
			c.ShutdownTimeout,
		))
	}

	validSamplers := map[string]bool{
		"always_on":    true,
		"always_off":   true,
		"ratio":        true,
		"parent_based": true,
		"":             true, // empty defaults to parent_based
	}
	if !validSamplers[strings.ToLower(c.SamplerType)] {
		errors = append(errors, fmt.Errorf(
			"invalid sampler type %q, must be one of: always_on, always_off, ratio, parent_based",
			c.SamplerType,
		))
	}

	return errors
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBoolOrDefault returns the environment variable as a boolean or a default value
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvFloatOrDefault returns the environment variable as a float64 or a default value
func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvDurationOrDefault returns the environment variable as a time.Duration or a default value
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvIntOrDefault returns the environment variable as an int or a default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// parseOTelHeaders parses OTEL_EXPORTER_OTLP_HEADERS environment variable
func parseOTelHeaders() map[string]string {
	headers := make(map[string]string)
	headersEnv := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
	if headersEnv == "" {
		return headers
	}

	// Parse headers in format: key1=value1,key2=value2
	pairs := strings.Split(headersEnv, ",")
	for _, pair := range pairs {
		trimmedPair := strings.TrimSpace(pair)
		if trimmedPair == "" {
			continue
		}

		if kv := strings.SplitN(trimmedPair, "=", 2); len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key != "" {
				headers[key] = value
			}
		} else {
			// Log warning for malformed header pairs
			fmt.Fprintf(os.Stderr, "Warning: malformed OTLP header pair ignored: %q\n", trimmedPair)
		}
	}
	return headers
}

// parsePropagators parses OTEL_PROPAGATORS environment variable
func parsePropagators() []string {
	propagatorsEnv := os.Getenv("OTEL_PROPAGATORS")
	if propagatorsEnv == "" {
		// Default propagators
		return []string{"tracecontext", "baggage"}
	}

	// Parse comma-separated list
	var propagators []string
	for _, p := range strings.Split(propagatorsEnv, ",") {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			propagators = append(propagators, strings.ToLower(trimmed))
		}
	}
	return propagators
}

// parseResourceAttributes parses OTEL_RESOURCE_ATTRIBUTES environment variable
// Format: key1=value1,key2=value2,key3=value3
func parseResourceAttributes() []attribute.KeyValue {
	var attrs []attribute.KeyValue
	resourceEnv := os.Getenv("OTEL_RESOURCE_ATTRIBUTES")
	if resourceEnv == "" {
		return attrs
	}

	// Parse comma-separated key=value pairs
	pairs := strings.Split(resourceEnv, ",")
	for _, pair := range pairs {
		// Split on first = to handle values that might contain =
		if kv := strings.SplitN(strings.TrimSpace(pair), "=", 2); len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if key != "" && value != "" {
				attrs = append(attrs, attribute.String(key, value))
			}
		}
	}
	return attrs
}

// createPropagators creates propagators based on configuration
func createPropagators(propagatorNames []string) propagation.TextMapPropagator {
	var propagators []propagation.TextMapPropagator
	originalCount := len(propagatorNames)

	for _, name := range propagatorNames {
		switch strings.ToLower(name) {
		case "tracecontext":
			propagators = append(propagators, propagation.TraceContext{})
		case "baggage":
			propagators = append(propagators, propagation.Baggage{})
		case "b3":
			propagators = append(propagators, b3.New())
		case "jaeger":
			propagators = append(propagators, jaeger.Jaeger{})
		case "b3multi":
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
		default:
			// Log warning for unknown propagator
			fmt.Fprintf(os.Stderr, "Warning: unknown propagator %q ignored\n", name)
		}
	}

	if len(propagators) == 0 && originalCount > 0 {
		// Log warning when falling back due to all invalid propagators
		fmt.Fprintf(
			os.Stderr,
			"Warning: all propagators were invalid, falling back to default [tracecontext, baggage]\n",
		)
	}

	if len(propagators) == 0 {
		// Fallback to default if no valid propagators
		propagators = []propagation.TextMapPropagator{
			propagation.TraceContext{},
			propagation.Baggage{},
		}
	}

	return propagation.NewCompositeTextMapPropagator(propagators...)
}

// normalizeOTLPEndpoint normalizes an OTLP endpoint to be compatible with otlptracehttp.WithEndpoint
// which expects host[:port] format, not full URLs
func normalizeOTLPEndpoint(endpoint string) string {
	// If it's already just host:port, return as-is
	if !strings.Contains(endpoint, "://") {
		return endpoint
	}

	// Parse the URL to extract host and port
	if strings.HasPrefix(endpoint, "http://") {
		endpoint = strings.TrimPrefix(endpoint, "http://")
	} else if strings.HasPrefix(endpoint, "https://") {
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	// Remove path components if present
	if idx := strings.Index(endpoint, "/"); idx != -1 {
		endpoint = endpoint[:idx]
	}

	return endpoint
}

// createSampler creates the appropriate sampler based on configuration
func createSampler(config *OTelConfig) sdktrace.Sampler {
	// Clamp ratio for safety
	clampedRatio := math.Max(0.0, math.Min(1.0, config.SamplingRatio))

	switch strings.ToLower(config.SamplerType) {
	case "always_on":
		return sdktrace.AlwaysSample()
	case "always_off":
		return sdktrace.NeverSample()
	case "ratio":
		return sdktrace.TraceIDRatioBased(clampedRatio)
	case "parent_based", "":
		// Default: parent-based with ratio-based root sampling
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(clampedRatio))
	default:
		// Fallback to parent-based if unknown sampler type
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(clampedRatio))
	}
}

// NewOTelConfig creates an OpenTelemetry configuration by reading environment variables
// with sensible defaults when environment variables are not set.
//
// Sampling behavior:
//   - Default sampler is "parent_based" with 10% ratio for production efficiency
//   - Parent-based means: follow parent span's sampling decision if present,
//     otherwise apply ratio-based sampling (10% by default)
//   - For development, use EnableDevTracer() which samples 100% of traces
//   - Override with OTEL_TRACES_SAMPLER and OTEL_TRACES_SAMPLER_ARG environment variables
//
// Common sampler configurations:
//   - OTEL_TRACES_SAMPLER=always_on: Sample all traces (development)
//   - OTEL_TRACES_SAMPLER=always_off: Sample no traces
//   - OTEL_TRACES_SAMPLER=parentbased_always_off: Never sample root traces
//   - OTEL_TRACES_SAMPLER=ratio OTEL_TRACES_SAMPLER_ARG=0.01: Sample 1% of traces
func NewOTelConfig() *OTelConfig {
	// Determine if tracing is enabled (OTEL SDK auto-detection or explicit enable)
	enabled := getEnvBoolOrDefault("XMTPD_TRACING_ENABLE", false) ||
		os.Getenv("OTEL_SERVICE_NAME") != "" ||
		os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != ""

	// Use OTEL_EXPORTER_OTLP_TRACES_ENDPOINT or fallback to OTEL_EXPORTER_OTLP_ENDPOINT
	endpoint := getEnvOrDefault("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "")
	if endpoint == "" {
		// Use base endpoint, will be normalized for otlptracehttp.WithEndpoint
		endpoint = getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
	}

	// Security: Only allow insecure mode when explicitly enabled
	// Default to TLS for production safety
	insecure := getEnvBoolOrDefault("OTEL_EXPORTER_OTLP_INSECURE", false)

	// Development options
	useSimpleProcessor := getEnvBoolOrDefault("XMTPD_TRACING_SIMPLE_PROCESSOR", false)
	useStdout := getEnvBoolOrDefault("XMTPD_TRACING_STDOUT", false)

	return &OTelConfig{
		ServiceName:    getEnvOrDefault("OTEL_SERVICE_NAME", defaultServiceName),
		ServiceVersion: getEnvOrDefault("OTEL_SERVICE_VERSION", "unknown"),
		Environment:    getEnvOrDefault("OTEL_DEPLOYMENT_ENVIRONMENT", "development"),
		OTLPEndpoint:   normalizeOTLPEndpoint(endpoint),
		Enabled:        enabled,
		SamplingRatio: getEnvFloatOrDefault(
			"OTEL_TRACES_SAMPLER_ARG",
			0.1,
		), // Default to 10% sampling for production efficiency
		SamplerType: getEnvOrDefault("OTEL_TRACES_SAMPLER", "parent_based"),
		Propagators: parsePropagators(),
		// Use standard OTEL batch span processor environment variables
		BatchTimeout: getEnvDurationOrDefault(
			"OTEL_BSP_SCHEDULE_DELAY",
			5*time.Second,
		),
		MaxExportBatchSize: getEnvIntOrDefault("OTEL_BSP_MAX_EXPORT_BATCH_SIZE", 512),
		MaxQueueSize:       getEnvIntOrDefault("OTEL_BSP_MAX_QUEUE_SIZE", 2048),
		ShutdownTimeout: getEnvDurationOrDefault(
			"XMTPD_TRACING_SHUTDOWN_TIMEOUT",
			5*time.Second,
		),
		OTLPHeaders:        parseOTelHeaders(),
		Insecure:           insecure,
		UseSimpleProcessor: useSimpleProcessor,
		UseStdout:          useStdout,
		IncludeHostInfo:    getEnvBoolOrDefault("XMTPD_TRACING_INCLUDE_HOST", true),
		IncludeProcessInfo: getEnvBoolOrDefault("XMTPD_TRACING_INCLUDE_PROCESS", true),
		StrictMode:         getEnvBoolOrDefault("XMTPD_TRACING_STRICT_MODE", false),
	}
}

// DefaultOTelConfig returns a static default configuration without reading environment variables
// Use NewOTelConfig() instead for environment-aware configuration
func DefaultOTelConfig() *OTelConfig {
	return &OTelConfig{
		ServiceName:        defaultServiceName,
		ServiceVersion:     "unknown",
		Environment:        "development",
		OTLPEndpoint:       "localhost:4318",
		Enabled:            false, // Disabled by default when not reading env vars
		SamplingRatio:      1.0,
		SamplerType:        "parent_based",
		Propagators:        []string{"tracecontext", "baggage"},
		BatchTimeout:       5 * time.Second,
		MaxExportBatchSize: 512,
		MaxQueueSize:       2048,
		ShutdownTimeout:    5 * time.Second,
		OTLPHeaders:        make(map[string]string),
		Insecure:           true, // Default to insecure for backward compatibility
		UseSimpleProcessor: false,
		UseStdout:          false,
		IncludeHostInfo:    true,
		IncludeProcessInfo: true,
		StrictMode:         true, // Production should fail fast on errors
	}
}

// EnableDevTracer creates a configuration optimized for local development
// with stdout exporter and immediate span processing
func EnableDevTracer() *OTelConfig {
	return &OTelConfig{
		ServiceName:        defaultServiceName,
		ServiceVersion:     "dev",
		Environment:        "development",
		OTLPEndpoint:       "localhost:4318", // Not used when UseStdout=true
		Enabled:            true,
		SamplingRatio:      1.0, // Sample everything in dev
		SamplerType:        "always_on",
		Propagators:        []string{"tracecontext", "baggage"},
		BatchTimeout:       100 * time.Millisecond,
		MaxExportBatchSize: 10,
		MaxQueueSize:       100,
		ShutdownTimeout:    2 * time.Second,
		OTLPHeaders:        make(map[string]string),
		Insecure:           true,
		UseSimpleProcessor: true,  // Immediate processing for dev
		UseStdout:          true,  // Print spans to console
		IncludeHostInfo:    false, // Less noise in dev
		IncludeProcessInfo: false, // Less noise in dev
		StrictMode:         false, // Allow errors in development
	}
}

// InitializeOTel initializes OpenTelemetry tracing with the provided configuration.
//
// This function sets up the global tracer provider, configures exporters (OTLP or stdout),
// sets up propagators for distributed tracing, and configures resource information.
//
// Parameters:
//   - ctx: Context for initialization (used for exporter creation)
//   - config: OpenTelemetry configuration (use NewOTelConfig() for environment-aware defaults)
//   - logger: Logger for tracing setup messages and errors
//
// Returns:
//   - A cleanup function that should be called before application shutdown
//   - An error if initialization fails (only in strict mode)
//
// Example usage:
//
//	config := NewOTelConfig()
//	cleanup, err := InitializeOTel(ctx, config, logger)
//	if err != nil {
//	    log.Fatal("Failed to initialize tracing:", err)
//	}
//	defer cleanup()
//
// In non-strict mode, errors are logged but don't prevent startup (graceful degradation).
// In strict mode (production), errors cause initialization to fail fast.
func InitializeOTel(ctx context.Context, config *OTelConfig, logger *zap.Logger) (func(), error) {
	if !config.Enabled {
		logger.Info("OpenTelemetry tracing is disabled")
		// Return a no-op shutdown function
		return func() {}, nil
	}

	logger.Info("Initializing OpenTelemetry tracing",
		zap.String("service", config.ServiceName),
		zap.String("version", config.ServiceVersion),
		zap.String("environment", config.Environment),
		zap.String("endpoint", config.OTLPEndpoint),
	)

	// Create resource with service information
	resourceAttrs := []attribute.KeyValue{
		semconv.ServiceName(config.ServiceName),
		semconv.ServiceVersion(config.ServiceVersion),
		semconv.DeploymentEnvironmentName(config.Environment),
	}

	// Add OTEL_RESOURCE_ATTRIBUTES from environment
	envResourceAttrs := parseResourceAttributes()
	resourceAttrs = append(resourceAttrs, envResourceAttrs...)

	// Add host information if requested
	if config.IncludeHostInfo {
		if hostname, err := os.Hostname(); err == nil {
			resourceAttrs = append(resourceAttrs, semconv.HostName(hostname))
		}
	}

	// Add process information if requested
	if config.IncludeProcessInfo {
		resourceAttrs = append(resourceAttrs, semconv.ProcessPID(os.Getpid()))
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(resourceAttrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var exporter sdktrace.SpanExporter

	// Choose exporter based on configuration
	if config.UseStdout {
		// Use stdout exporter for development
		stdoutExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			if config.StrictMode {
				return func() {}, fmt.Errorf("failed to create stdout exporter: %w", err)
			}
			logger.Error("Failed to create stdout exporter, falling back to no-op", zap.Error(err))
			return func() {}, nil
		}
		exporter = stdoutExporter
		logger.Info("Using stdout trace exporter for development")
	} else {
		// Create OTLP HTTP exporter - TLS by default unless explicitly set to insecure
		exporterOptions := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(config.OTLPEndpoint),
			otlptracehttp.WithHeaders(config.OTLPHeaders),
		}

		// Only add insecure option if explicitly requested
		if config.Insecure {
			exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
			logger.Warn("Using insecure OTLP transport", zap.String("endpoint", config.OTLPEndpoint))
		}

		otlpExporter, err := otlptracehttp.New(ctx, exporterOptions...)
		if err != nil {
			if config.StrictMode {
				return func() {}, fmt.Errorf("failed to create OTLP exporter: %w", err)
			}
			logger.Error("Failed to create OTLP exporter, falling back to no-op",
				zap.Error(err),
				zap.String("endpoint", config.OTLPEndpoint))
			return func() {}, nil
		}
		exporter = otlpExporter
	}

	// Create sampler based on configuration
	sampler := createSampler(config)

	// Choose span processor based on configuration
	var processor sdktrace.SpanProcessor
	if config.UseSimpleProcessor {
		// Simple processor for immediate processing (good for dev/debugging)
		processor = sdktrace.NewSimpleSpanProcessor(exporter)
		logger.Info("Using simple span processor for immediate processing")
	} else {
		// Batch processor for production efficiency
		processor = sdktrace.NewBatchSpanProcessor(exporter,
			sdktrace.WithBatchTimeout(config.BatchTimeout),
			sdktrace.WithMaxExportBatchSize(config.MaxExportBatchSize),
			sdktrace.WithMaxQueueSize(config.MaxQueueSize),
		)
		logger.Info("Using batch span processor for efficient processing")
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(processor),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global propagator based on configuration
	propagator := createPropagators(config.Propagators)
	otel.SetTextMapPropagator(propagator)

	logger.Info("OpenTelemetry tracing initialized successfully")

	// Return shutdown function with configurable timeout
	return func() {
		logger.Info("Shutting down OpenTelemetry tracing")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("Error shutting down tracer provider", zap.Error(err))
		} else {
			logger.Info("OpenTelemetry tracer provider shut down cleanly")
		}
	}, nil
}

// StartOTelSpan starts a new span with the given name and options.
//
// This is a convenience function that gets the global tracer and starts a span.
// The span will automatically be linked to the parent span in the context.
//
// Parameters:
//   - ctx: Parent context (may contain a parent span)
//   - spanName: Name of the operation being traced
//   - opts: Additional span options (attributes, span kind, etc.)
//
// Returns:
//   - A new context containing the created span
//   - The created span (remember to call span.End() when the operation completes)
//
// Example usage:
//
//	ctx, span := StartOTelSpan(ctx, "database.query",
//	    trace.WithAttributes(
//	        attribute.String("query", sql),
//	        attribute.String("table", "users"),
//	    ),
//	)
//	defer span.End()
//
//	// Perform database operation
//	result, err := db.Query(ctx, sql)
//	if err != nil {
//	    span.RecordError(err)
//	    span.SetStatus(codes.Error, err.Error())
//	}
//
// The span name should describe the operation being performed, following
// OpenTelemetry naming conventions (e.g., "service.method", "component.operation").
func StartOTelSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	tracer := otel.Tracer(defaultInstrumentName)
	return tracer.Start(ctx, spanName, opts...)
}

// AddEvent adds an event to the span in the current context
//
// Example usage:
//
//	AddEvent(ctx, "cache.hit", attribute.String("key", cacheKey))
func AddEvent(ctx context.Context, eventName string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(eventName, trace.WithAttributes(attrs...))
	}
}

// SetAttributes sets attributes on the span in the current context
//
// Example usage:
//
//	SetAttributes(ctx,
//	  attribute.String("user.id", userID),
//	  attribute.Int("batch.size", batchSize),
//	)
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// RecordError records an error on the span in the current context with structured attributes
//
// Example usage:
//
//	if err := someOperation(); err != nil {
//	  RecordError(ctx, err, trace.WithAttributes(attribute.String("operation", "someOperation")))
//	}
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() && err != nil {
		// Add structured error attributes
		attrs := []attribute.KeyValue{
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.message", err.Error()),
		}

		// Combine with any additional options
		eventOpts := append(opts, trace.WithAttributes(attrs...))
		span.RecordError(err, eventOpts...)

		// Also set span status to error
		span.SetStatus(codes.Error, err.Error())
	}
}

// SetOTelStatus sets the status of the span in the current context
//
// Example usage:
//
//	SetOTelStatus(ctx, codes.Error, "validation failed")
func SetOTelStatus(ctx context.Context, code codes.Code, description string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetStatus(code, description)
	}
}

// WrapWithSpan wraps a function with a span, automatically handling errors and panics
//
// This function provides comprehensive error handling:
//   - Automatically records errors and sets error status
//   - Panic recovery with error recording
//   - Guaranteed span cleanup even if fn never returns normally
//   - Context leak protection
//
// Example usage:
//
//	err := WrapWithSpan(ctx, "process-data", func(ctx context.Context) error {
//	  return processData(ctx, data)
//	}, trace.WithAttributes(attribute.String("data.type", "user")))
func WrapWithSpan(ctx context.Context, spanName string, fn func(context.Context) error, opts ...trace.SpanStartOption) (err error) {
	ctx, span := StartOTelSpan(ctx, spanName, opts...)

	// Panic recovery and guaranteed cleanup
	defer func() {
		if r := recover(); r != nil {
			// Record panic as error
			panicErr := fmt.Errorf("panic: %v", r)
			span.RecordError(panicErr)
			span.SetStatus(codes.Error, "panic occurred")

			// Set return error
			err = panicErr
		}

		// Always end span, even on panic
		span.End()

		// Re-panic to maintain original panic behavior
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	err = fn(ctx)
	if err != nil {
		RecordError(ctx, err)
		SetOTelStatus(ctx, codes.Error, err.Error())
	}

	return err
}

// AddTraceToLogger adds trace context (trace ID, span ID) to the logger for correlation
//
// Example usage:
//
//	logger := AddTraceToLogger(ctx, baseLogger)
//	logger.Info("This log will include trace and span IDs")
func AddTraceToLogger(ctx context.Context, logger *zap.Logger) *zap.Logger {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return logger
	}

	sc := span.SpanContext()
	if !sc.IsValid() {
		return logger
	}

	return logger.With(
		zap.String("trace_id", sc.TraceID().String()),
		zap.String("span_id", sc.SpanID().String()),
		zap.String("trace_flags", sc.TraceFlags().String()),
	)
}

// AddErrorEvent records an error as an event without changing span status
//
// Example usage:
//
//	// For non-fatal errors that shouldn't mark the span as failed
//	AddErrorEvent(ctx, err, "warning", "rate limit exceeded but retrying")
func AddErrorEvent(ctx context.Context, err error, level string, description string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() && err != nil {
		attrs := []attribute.KeyValue{
			attribute.String("error.type", fmt.Sprintf("%T", err)),
			attribute.String("error.message", err.Error()),
			attribute.String("level", level),
		}
		if description != "" {
			attrs = append(attrs, attribute.String("description", description))
		}

		span.AddEvent("error", trace.WithAttributes(attrs...))
	}
}

// HTTP Tracing Utilities

// OTelHTTPConfig holds configuration for HTTP OpenTelemetry instrumentation
type OTelHTTPConfig struct {
	// SpanNameFormatter allows customization of span naming
	SpanNameFormatter func(operation string, r *http.Request) string
	// RouteExtractor extracts route template from request (e.g., "/users/{id}" instead of "/users/123")
	RouteExtractor func(*http.Request) string
	// Additional attributes to add to all spans
	Attributes []attribute.KeyValue
}

// DefaultRouteExtractor provides a basic route extraction that just uses the path
func DefaultRouteExtractor(r *http.Request) string {
	return r.URL.Path
}

// WrapHTTPHandler wraps an HTTP handler with OpenTelemetry instrumentation
func WrapHTTPHandler(handler http.Handler, config *OTelHTTPConfig) http.Handler {
	if config == nil {
		config = &OTelHTTPConfig{}
	}

	// Default span name formatter using route template if available
	spanNameFormatter := config.SpanNameFormatter
	if spanNameFormatter == nil {
		routeExtractor := config.RouteExtractor
		if routeExtractor == nil {
			routeExtractor = DefaultRouteExtractor
		}

		spanNameFormatter = func(operation string, r *http.Request) string {
			route := routeExtractor(r)
			return fmt.Sprintf("%s %s", r.Method, route)
		}
	}

	// Build base attributes
	// Note: HTTP semantic conventions are automatically added by otelhttp
	// Service name belongs at the resource level, not span attributes
	attrs := config.Attributes

	// Use the otelhttp package to wrap the handler with custom attributes
	wrappedHandler := otelhttp.NewHandler(handler, "xmtpd-http-server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attrs...),
		),
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
	)

	return wrappedHandler
}

// Testing Utilities

// NewInMemoryTracer creates an in-memory tracer provider for testing
// Returns the tracer provider and a function to retrieve recorded spans
//
// Example usage:
//
//	tp, getSpans := NewInMemoryTracer()
//	defer tp.Shutdown(context.Background())
//
//	// Use tracing in the code
//	ctx, span := StartOTelSpan(ctx, "test-operation")
//	span.End()
//
//	// Assert spans in tests
//	spans := getSpans()
//	assert.Len(t, spans, 1)
//	assert.Equal(t, "test-operation", spans[0].Name)
func NewInMemoryTracer() (*sdktrace.TracerProvider, func() []sdktrace.ReadOnlySpan) {
	exporter := &inMemoryExporter{}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter)),
	)

	// Set as global for convenience in tests
	otel.SetTracerProvider(tp)

	return tp, exporter.getSpans
}

// inMemoryExporter is a simple in-memory span exporter for testing
type inMemoryExporter struct {
	spans []sdktrace.ReadOnlySpan
}

func (e *inMemoryExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	e.spans = append(e.spans, spans...)
	return nil
}

func (e *inMemoryExporter) Shutdown(ctx context.Context) error {
	e.spans = nil
	return nil
}

func (e *inMemoryExporter) getSpans() []sdktrace.ReadOnlySpan {
	return e.spans
}
