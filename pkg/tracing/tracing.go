// Package tracing enables [Datadog APM tracing](https://docs.datadoghq.com/tracing/) capabilities,
// focusing specifically on [Error Tracking](https://docs.datadoghq.com/tracing/error_tracking/)
package tracing

import (
	"context"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

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

func (*noopSpan) SetTag(string, interface{})     {}
func (*noopSpan) SetOperationName(string)        {}
func (*noopSpan) BaggageItem(string) string      { return "" }
func (*noopSpan) SetBaggageItem(string, string)  {}
func (*noopSpan) Finish(...ddtrace.FinishOption) {}

func (*noopSpan) Context() ddtrace.SpanContext {
	return &noopSpanContext{}
}

// noopSpanContext satisfies ddtrace.SpanContext for the no-op span.
type noopSpanContext struct{}

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
	if !apmEnabled {
		return noopSpanInstance, ctx
	}
	return tracer.StartSpanFromContext(ctx, operationName, opts...)
}

// StartSpan creates a new root span.
// Returns a no-op span when tracing is disabled.
func StartSpan(operationName string, opts ...ddtrace.StartSpanOption) Span {
	if !apmEnabled {
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
var apmEnabled bool

// Start boots the datadog tracer, run this once early in the startup sequence.
// Tracing is gated by XMTPD_TRACING_ENABLE at the application config level;
// callers should only invoke Start() when the feature flag is on.
// When enabled, all traces are collected deterministically (100%, no sampling).
//
// Configuration via environment variables:
//   - ENV: Environment name (default: "test")
//   - DD_AGENT_HOST: Datadog agent host (standard DD env var, default: "localhost")
//   - DD_TRACE_AGENT_PORT: Datadog agent port (standard DD env var, default: "8126")
func Start(version string, l *zap.Logger) {
	apmEnabled = true

	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}

	l.Info("APM tracing enabled (deterministic, no sampling)",
		zap.String("env", env),
	)

	tracer.Start(
		tracer.WithEnv(env),
		tracer.WithService("xmtpd"),
		tracer.WithServiceVersion(version),
		tracer.WithLogger(logger{l}),
		tracer.WithRuntimeMetrics(),
	)
}

// IsEnabled returns whether APM tracing is currently enabled.
// Use this to conditionally skip expensive span creation.
func IsEnabled() bool {
	return apmEnabled
}

// SetEnabledForTesting overrides the apmEnabled flag for use in tests.
// Returns a cleanup function that restores the previous state.
// This must only be called from test code.
func SetEnabledForTesting(enabled bool) func() {
	prev := apmEnabled
	apmEnabled = enabled
	return func() { apmEnabled = prev }
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
	defer span.Finish()
	logger = Link(span, logger.With(zap.String("span", operation)))
	err := action(ctx, logger, span)
	if err != nil {
		span.Finish(WithError(err))
	}
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
	if !apmEnabled {
		return l
	}
	return l.With(
		zap.Uint64("dd.trace_id", span.Context().TraceID()),
		zap.Uint64("dd.span_id", span.Context().SpanID()))
}

func SpanType(span Span, typ string) {
	if !apmEnabled {
		return
	}
	span.SetTag(ext.SpanType, typ)
}

func SpanResource(span Span, resource string) {
	if !apmEnabled {
		return
	}
	span.SetTag(ext.ResourceName, resource)
}

// SpanTag sets a tag on a span with production safety limits.
// String values longer than MaxTagValueLength (in runes) are truncated.
// Uses rune-based truncation to safely handle multi-byte UTF-8 characters.
// No-ops when tracing is disabled.
func SpanTag(span Span, key string, value any) {
	if !apmEnabled {
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

// traceContextEntry holds a span context with its creation time for TTL cleanup.
type traceContextEntry struct {
	ctx       ddtrace.SpanContext
	createdAt time.Time
}

// TraceContextStore provides async context propagation by mapping
// staged envelope IDs to their originating span contexts. This allows
// the publish_worker to create child spans linked to the original
// staging request, enabling end-to-end distributed tracing across
// async boundaries.
//
// Includes TTL-based cleanup to prevent memory leaks from orphaned entries.
type TraceContextStore struct {
	mu           sync.RWMutex
	contexts     map[int64]traceContextEntry
	ttl          time.Duration
	lastCleanup  time.Time
	cleanupCount int // Track cleanups for testing/monitoring
}

// DefaultTraceContextTTL is the default time-to-live for stored span contexts.
// 5 minutes is generous - publish_worker typically processes within seconds.
const DefaultTraceContextTTL = 5 * time.Minute

// Span limits for production safety - prevent runaway memory/payload sizes.
const (
	// MaxTagValueLength is the maximum length for string tag values.
	// Longer strings are truncated to prevent excessive trace payload sizes.
	// 1KB is generous for most use cases while preventing abuse.
	MaxTagValueLength = 1024

	// MaxStoreSize is the maximum number of entries in TraceContextStore.
	// Prevents unbounded memory growth if publish_worker falls behind.
	// 10K entries at ~100 bytes each = ~1MB max memory.
	MaxStoreSize = 10000
)

// NewTraceContextStore creates a new store for async trace context propagation.
func NewTraceContextStore() *TraceContextStore {
	return &TraceContextStore{
		contexts:    make(map[int64]traceContextEntry),
		ttl:         DefaultTraceContextTTL,
		lastCleanup: time.Now(),
	}
}

// Store saves the span context for a staged envelope ID.
// Call this after staging an envelope to enable trace linking.
// Performs lazy cleanup of expired entries to prevent memory leaks.
// Drops new entries if store is at capacity (production safety).
// No-ops when tracing is disabled.
func (s *TraceContextStore) Store(stagedID int64, span Span) {
	if !apmEnabled || span == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// Lazy cleanup: run every minute to prevent unbounded growth
	if time.Since(s.lastCleanup) > time.Minute {
		s.cleanupExpiredLocked()
	}

	// Production safety: refuse new entries if at capacity
	// This indicates publish_worker is falling behind and needs investigation
	if len(s.contexts) >= MaxStoreSize {
		return
	}

	s.contexts[stagedID] = traceContextEntry{
		ctx:       span.Context(),
		createdAt: time.Now(),
	}
}

// Retrieve gets and removes the span context for a staged envelope ID.
// Returns nil if no context was stored for this ID or if it expired.
func (s *TraceContextStore) Retrieve(stagedID int64) ddtrace.SpanContext {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.contexts[stagedID]
	if !ok {
		return nil
	}

	// Always delete the entry (retrieved or expired)
	delete(s.contexts, stagedID)

	// Check if expired
	if time.Since(entry.createdAt) > s.ttl {
		return nil
	}

	return entry.ctx
}

// cleanupExpiredLocked removes entries older than TTL.
// Must be called with lock held.
func (s *TraceContextStore) cleanupExpiredLocked() {
	now := time.Now()
	for id, entry := range s.contexts {
		if now.Sub(entry.createdAt) > s.ttl {
			delete(s.contexts, id)
		}
	}
	s.lastCleanup = now
	s.cleanupCount++
}

// Size returns the current number of stored contexts (for monitoring).
func (s *TraceContextStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.contexts)
}

// StartSpanWithParent creates a new span, optionally linked to a parent context.
// If parentCtx is nil, creates a new root span. This is useful for async
// workflows where the parent context may or may not be available.
// Returns a no-op span when tracing is disabled.
func StartSpanWithParent(operationName string, parentCtx ddtrace.SpanContext) Span {
	if !apmEnabled {
		return noopSpanInstance
	}
	if parentCtx != nil {
		return tracer.StartSpan(operationName, tracer.ChildOf(parentCtx))
	}
	return tracer.StartSpan(operationName)
}
