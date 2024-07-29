// Package tracing enables [Datadog APM tracing](https://docs.datadoghq.com/tracing/) capabilities,
// focusing specifically on [Error Tracking](https://docs.datadoghq.com/tracing/error_tracking/)
package tracing

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
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
func Wrap(ctx context.Context, log *zap.Logger, operation string, action func(context.Context, *zap.Logger, Span) error) error {
	span, ctx := tracer.StartSpanFromContext(ctx, operation)
	defer span.Finish()
	log = Link(span, log.With(zap.String("span", operation)))
	err := action(ctx, log, span)
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

func SpanTag(span Span, key string, value interface{}) {
	span.SetTag(key, value)
}

// GoPanicWrap extends PanicWrap by running the body in a goroutine and
// synchronizing the goroutine exit with the WaitGroup.
// The body must respect cancellation of the Context.
func GoPanicWrap(ctx context.Context, wg *sync.WaitGroup, name string, body func(context.Context)) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		PanicWrap(ctx, name, body)
	}()
}
