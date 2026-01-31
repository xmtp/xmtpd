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

// reimport relevant bits of the tracer API
var (
	StartSpanFromContext = tracer.StartSpanFromContext
	StartSpan            = tracer.StartSpan
	ChildOf              = tracer.ChildOf
	WithError            = tracer.WithError
	ContextWithSpan      = tracer.ContextWithSpan
)

type Span = tracer.Span

type logger struct{ *zap.Logger }

func (l logger) Log(msg string) {
	l.Error(msg)
}

// Start boots the datadog tracer, run this once early in the startup sequence.
func Start(version string, l *zap.Logger) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}
	tracer.Start(
		tracer.WithEnv(env),
		tracer.WithService("xmtpd"),
		tracer.WithServiceVersion(version),
		tracer.WithLogger(logger{l}),
		tracer.WithRuntimeMetrics(),
	)
}

// Stop shuts down the datadog tracer, defer this right after Start().
func Stop() {
	tracer.Stop()
}

// Wrap executes action in the context of a span.
// Tags the span with the error if action returns one.
func Wrap(
	ctx context.Context,
	logger *zap.Logger,
	operation string,
	action func(context.Context, *zap.Logger, Span) error,
) error {
	span, ctx := tracer.StartSpanFromContext(ctx, operation)
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
func Link(span tracer.Span, l *zap.Logger) *zap.Logger {
	return l.With(
		zap.Uint64("dd.trace_id", span.Context().TraceID()),
		zap.Uint64("dd.span_id", span.Context().SpanID()))
}

func SpanType(span Span, typ string) {
	span.SetTag(ext.SpanType, typ)
}

func SpanResource(span Span, resource string) {
	span.SetTag(ext.ResourceName, resource)
}

func SpanTag(span Span, key string, value any) {
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
func (s *TraceContextStore) Store(stagedID int64, span Span) {
	if span == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	// Lazy cleanup: run every minute to prevent unbounded growth
	if time.Since(s.lastCleanup) > time.Minute {
		s.cleanupExpiredLocked()
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
func StartSpanWithParent(operationName string, parentCtx ddtrace.SpanContext) Span {
	if parentCtx != nil {
		return tracer.StartSpan(operationName, tracer.ChildOf(parentCtx))
	}
	return tracer.StartSpan(operationName)
}
