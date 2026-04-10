// Package tracing enables [Datadog APM tracing](https://docs.datadoghq.com/tracing/) capabilities,
// focusing specifically on [Error Tracking](https://docs.datadoghq.com/tracing/error_tracking/)
package tracing

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// EnvAPMSampleRate is the environment variable for configuring APM sample rate.
const EnvAPMSampleRate = "APM_SAMPLE_RATE"

// Re-export tracer types and options that don't need gating.
var (
	ChildOf         = tracer.ChildOf
	WithError       = tracer.WithError
	ContextWithSpan = tracer.ContextWithSpan
)

// Span is the tracer span type (ddtrace.Span interface).
type Span = ddtrace.Span

// noopSpan implements ddtrace.Span with zero-cost no-ops.
// A single instance is shared across all callers when tracing is disabled.
type noopSpan struct{}

var noopSpanInstance Span = &noopSpan{}

// noopSpanContext satisfies ddtrace.SpanContext for the no-op span.
// A single instance is shared to avoid heap allocations.
type noopSpanContext struct{}

var noopSpanCtxInstance ddtrace.SpanContext = &noopSpanContext{}

func (*noopSpan) SetTag(string, any)             {}
func (*noopSpan) SetOperationName(string)        {}
func (*noopSpan) BaggageItem(string) string      { return "" }
func (*noopSpan) SetBaggageItem(string, string)  {}
func (*noopSpan) Finish(...ddtrace.FinishOption) {}

func (*noopSpan) Context() ddtrace.SpanContext {
	return noopSpanCtxInstance
}

func (*noopSpanContext) SpanID() uint64  { return 0 }
func (*noopSpanContext) TraceID() uint64 { return 0 }

func (*noopSpanContext) ForeachBaggageItem(
	func(k, v string) bool,
) {
}

// StartSpanFromContext creates a span as a child of the context's active span.
// Returns a no-op span and the unchanged context when tracing is disabled.
func StartSpanFromContext(
	ctx context.Context,
	operationName string,
	opts ...ddtrace.StartSpanOption,
) (Span, context.Context) {
	if !apmEnabled.Load() {
		return noopSpanInstance, ctx
	}
	return tracer.StartSpanFromContext(ctx, operationName, opts...)
}

// StartSpan creates a new root span.
// Returns a no-op span when tracing is disabled.
func StartSpan(operationName string, opts ...ddtrace.StartSpanOption) Span {
	if !apmEnabled.Load() {
		return noopSpanInstance
	}
	return tracer.StartSpan(operationName, opts...)
}

type logger struct{ *zap.Logger }

func (l logger) Log(msg string) {
	l.Error(msg)
}

// apmEnabled tracks whether tracing was started (for conditional span creation).
// Controlled exclusively by XMTPD_TRACING_ENABLE at the application level.
// When false, all span creation functions return no-ops with zero overhead.
//
// Stored as atomic.Bool so tests that flip it on/off via SetEnabledForTesting
// don't race with concurrent readers in tracing hot paths.
var apmEnabled atomic.Bool

// Start boots the datadog tracer, run this once early in the startup sequence.
// Tracing is gated by XMTPD_TRACING_ENABLE at the application config level;
// callers should only invoke Start() when the feature flag is on.
//
// Configuration via environment variables:
//   - ENV: Environment name (default: "test")
//   - DD_AGENT_HOST: Datadog agent host (standard DD env var, default: "localhost")
//   - DD_TRACE_AGENT_PORT: Datadog agent port (standard DD env var, default: "8126")
//   - APM_SAMPLE_RATE: Sample rate 0.0-1.0 (default: 1.0 dev/test, 0.1 prod)
func Start(version string, l *zap.Logger) {
	apmEnabled.Store(true)

	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}

	sampleRate := getSampleRate(env, l)

	opts := []tracer.StartOption{
		tracer.WithEnv(env),
		tracer.WithService("xmtpd"),
		tracer.WithServiceVersion(version),
		tracer.WithLogger(logger{l}),
		tracer.WithRuntimeMetrics(),
	}

	// Add sampler if not 100%
	if sampleRate < 1.0 {
		opts = append(opts, tracer.WithSampler(tracer.NewRateSampler(sampleRate)))
		l.Info("APM tracing enabled with sampling",
			zap.Float64("sample_rate", sampleRate),
			zap.String("env", env),
		)
	} else {
		l.Info("APM tracing enabled (100% sampling)",
			zap.String("env", env),
		)
	}

	tracer.Start(opts...)
}

// getSampleRate returns the configured sample rate.
// Defaults: 1.0 (100%) for dev/test, 0.1 (10%) for production/staging.
func getSampleRate(env string, l *zap.Logger) float64 {
	// Check for explicit configuration
	if rateStr := os.Getenv(EnvAPMSampleRate); rateStr != "" {
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil || rate < 0 || rate > 1 {
			l.Warn("Invalid APM_SAMPLE_RATE, using default",
				zap.String("value", rateStr),
			)
		} else {
			return rate
		}
	}

	// Environment-based defaults
	switch env {
	case "production", "prod", "staging":
		return 0.1 // 10% sampling in production
	default:
		return 1.0 // 100% sampling in dev/test
	}
}

// IsEnabled returns whether APM tracing is currently enabled.
// Use this to conditionally skip expensive span creation.
func IsEnabled() bool {
	return apmEnabled.Load()
}

// SetEnabledForTesting overrides the apmEnabled flag for use in tests.
// Returns a cleanup function that restores the previous state.
// This must only be called from test code.
func SetEnabledForTesting(enabled bool) func() {
	prev := apmEnabled.Swap(enabled)
	return func() { apmEnabled.Store(prev) }
}

// Stop shuts down the datadog tracer, defer this right after Start().
func Stop() {
	tracer.Stop()
}

// Wrap executes action in the context of a span.
// Tags the span with the error if action returns one.
// When tracing is disabled, the action runs without span overhead.
func Wrap(
	ctx context.Context,
	logger *zap.Logger,
	operation string,
	action func(context.Context, *zap.Logger, Span) error,
) error {
	span, ctx := StartSpanFromContext(ctx, operation)
	logger = Link(span, logger.With(zap.String("span", operation)))

	var err error
	defer func() {
		r := recover()
		if r != nil {
			// A panic occurred - finish span with error and re-panic
			if panicErr, ok := r.(error); ok {
				span.Finish(WithError(panicErr))
			} else {
				span.Finish(WithError(fmt.Errorf("panic: %v", r)))
			}
			panic(r)
		}
		// Normal path - finish span with error if action returned one
		if err != nil {
			span.Finish(WithError(err))
		} else {
			span.Finish()
		}
	}()

	err = action(ctx, logger, span)
	return err
}

// PanicWrap executes the body guarding for panics.
// If panic happens it emits a span with the error attached.
// This should trigger DD APM's Error Tracking to record the error.
func PanicWrap(ctx context.Context, name string, body func(context.Context)) {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			StartSpan("panic: " + name).Finish(
				WithError(err),
			)
		}
		if r != nil {
			// Repanic so that we don't suppress normal panic behavior.
			panic(r)
		}
	}()
	body(ctx)
}

// Link connects a logger to a particular trace and span.
// DD APM should provide some additional functionality based on that.
// Returns the logger unchanged when tracing is disabled.
func Link(span Span, l *zap.Logger) *zap.Logger {
	if !apmEnabled.Load() {
		return l
	}
	return l.With(
		zap.Uint64("dd.trace_id", span.Context().TraceID()),
		zap.Uint64("dd.span_id", span.Context().SpanID()))
}

func SpanType(span Span, typ string) {
	if !apmEnabled.Load() {
		return
	}
	span.SetTag(ext.SpanType, typ)
}

func SpanResource(span Span, resource string) {
	if !apmEnabled.Load() {
		return
	}
	span.SetTag(ext.ResourceName, resource)
}

// SpanTag sets a tag on a span with production safety limits.
// String values longer than MaxTagValueLength (in runes) are truncated.
// Uses rune-based truncation to safely handle multi-byte UTF-8 characters.
func SpanTag(span Span, key string, value any) {
	if !apmEnabled.Load() {
		return
	}
	// Truncate long strings to prevent excessive payload sizes
	if s, ok := value.(string); ok {
		runes := []rune(s)
		if len(runes) > MaxTagValueLength {
			value = string(runes[:MaxTagValueLength]) + "...[truncated]"
		}
	}
	span.SetTag(key, value)
}

// GoPanicWrap extends PanicWrap by running the body in a goroutine and
// synchronizing the goroutine exit with the WaitGroup.
// The body must respect cancellation of the Context.
func GoPanicWrap(
	ctx context.Context,
	wg *sync.WaitGroup,
	name string,
	body func(context.Context),
	labels ...string,
) {
	wg.Add(1)

	expandedLabels := append(labels, "name", name)

	go pprof.Do(ctx, pprof.Labels(expandedLabels...), func(ctx context.Context) {
		defer wg.Done()
		PanicWrap(ctx, name, body)
	})
}

// StartSpanWithParent creates a new span, optionally linked to a parent context.
// If parentCtx is nil, creates a new root span. This is useful for async
// workflows where the parent context may or may not be available.
// Returns a no-op span when tracing is disabled.
func StartSpanWithParent(operationName string, parentCtx ddtrace.SpanContext) Span {
	if !apmEnabled.Load() {
		return noopSpanInstance
	}
	if parentCtx != nil {
		return tracer.StartSpan(operationName, tracer.ChildOf(parentCtx))
	}
	return tracer.StartSpan(operationName)
}
