// Package tracing enables Datadog APM tracing (dd-trace-go/v2),
// focusing on Error Tracking.
package tracing

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"go.uber.org/zap"
)

// dd-trace-go expects a minimal logger with Log(string).
type ddLogger struct{ *zap.Logger }

func (l ddLogger) Log(msg string) { l.Error(msg) }

// Start boots the Datadog tracer. Call once early in startup.
func Start(version string, l *zap.Logger) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}

	err := tracer.Start(
		tracer.WithEnv(env),
		tracer.WithService("xmtpd"),
		tracer.WithServiceVersion(version),
		tracer.WithLogger(ddLogger{l}),
		tracer.WithRuntimeMetrics(), // v2 runtime metrics
	)
	if err != nil {
		panic(err)
	}
}

// Stop shuts down the Datadog tracer. Defer right after Start().
func Stop() { tracer.Stop() }

// Wrap executes action in the context of a span and attaches any returned error
// to that span (so Error Tracking can pick it up).
func Wrap(
	ctx context.Context,
	log *zap.Logger,
	operation string,
	action func(context.Context, *zap.Logger, *tracer.Span) error,
) error {
	span, ctx := tracer.StartSpanFromContext(ctx, operation)
	defer span.Finish()

	log = Link(span, log.With(zap.String("span", operation)))

	err := action(ctx, log, span)
	if err != nil {
		// Don't call Finish twice; in v2 just finish once with options.
		span.Finish(tracer.WithError(err))
	}
	return err
}

// PanicWrap guards body for panics. If a panic happens it emits a span with an
// error attached, then re-panics so normal behavior is preserved.
func PanicWrap(ctx context.Context, name string, body func(context.Context)) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		span := tracer.StartSpan("panic." + name)
		switch v := r.(type) {
		case error:
			span.Finish(tracer.WithError(v))
		default:
			span.Finish(tracer.WithError(fmt.Errorf("%v", v)))
		}

		panic(r)
	}()

	body(ctx)
}

// Link connects a zap logger to the current trace/span IDs for log correlation.
func Link(span *tracer.Span, l *zap.Logger) *zap.Logger {
	sc := span.Context()
	return l.With(
		zap.String("dd.trace_id", sc.TraceID()),
		zap.Uint64("dd.span_id", sc.SpanID()),
	)
}

func SpanType(span *tracer.Span, typ string)           { span.SetTag(ext.SpanType, typ) }
func SpanResource(span *tracer.Span, resource string)  { span.SetTag(ext.ResourceName, resource) }
func SpanTag(span *tracer.Span, key string, value any) { span.SetTag(key, value) }

// GoPanicWrap runs body in a goroutine, labels the goroutine with pprof labels,
// and synchronizes exit with the WaitGroup. The body must respect ctx cancelation.
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
