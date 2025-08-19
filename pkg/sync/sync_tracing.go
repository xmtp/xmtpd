// sync_tracing.go provides OpenTelemetry instrumentation for sync operations.
package sync

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Package-level tracer instance to avoid repeated calls
// Using full import path for better trace source identification
var tracer = otel.Tracer("github.com/xmtp/xmtpd/pkg/sync")

// Common attribute keys for consistent naming across spans
const (
	AttrSyncNodeID      = "sync.node.id"
	AttrSyncNodeAddress = "sync.node.address"
	AttrSyncNode        = "sync.node" // Compact node info: "id@address"
	AttrSyncOperation   = "sync.operation.type"
	AttrSyncRegistryOp  = "sync.registry.operation"
	AttrSyncWorkerType  = "sync.worker.type"
	AttrSyncComponent   = "sync.component"
	AttrOperation       = "operation" // Generic operation identifier
)

// syncSpanOptions provides common span attributes for sync operations
func syncSpanOptions(operationType string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String(AttrSyncOperation, operationType),
		attribute.String(AttrSyncComponent, "sync-service"),
	}
}

// recordErrorWithAttributes records an error with additional context attributes
func recordErrorWithAttributes(span trace.Span, err error, attrs ...attribute.KeyValue) {
	span.RecordError(err, trace.WithAttributes(attrs...))
}

// setSpanStatus sets the appropriate span status based on error state
func setSpanStatus(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// logWithTraceContext enriches logger with trace context for correlation
func logWithTraceContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	if logger == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return logger
	}

	spanCtx := span.SpanContext()
	return logger.With(
		zap.String("trace_id", spanCtx.TraceID().String()),
		zap.String("span_id", spanCtx.SpanID().String()),
	)
}

// traceNodeConnection wraps node connection operations with OpenTelemetry tracing
func traceNodeConnection(
	ctx context.Context,
	nodeID uint32,
	address string,
	logger *zap.Logger,
	fn func(context.Context) error,
) (context.Context, error) {
	ctx, span := tracer.Start(ctx, "node.connection")
	defer span.End()

	// Add node connection attributes with consistent naming
	nodeCompact := fmt.Sprintf("%d@%s", nodeID, address)
	span.SetAttributes(
		attribute.Int64(AttrSyncNodeID, int64(nodeID)),
		attribute.String(AttrSyncNodeAddress, address),
		attribute.String(AttrSyncNode, nodeCompact),
		attribute.String(AttrSyncOperation, "connection"),
		attribute.String(AttrSyncComponent, "sync-service"),
	)

	// Get context-enriched logger
	enrichedLogger := logWithTraceContext(ctx, logger)

	// Capture panics and record them as errors
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in node connection: %v", r)
			recordErrorWithAttributes(span, err,
				attribute.Int64(AttrSyncNodeID, int64(nodeID)),
				attribute.String(AttrSyncNodeAddress, address),
				attribute.String("panic_value", fmt.Sprintf("%v", r)),
			)
			span.SetStatus(codes.Error, err.Error())
			if enrichedLogger != nil {
				enrichedLogger.Error("panic during node connection",
					zap.Uint32("node_id", nodeID),
					zap.String("address", address),
					zap.Any("panic", r))
			}
			panic(r) // Re-panic after recording
		}
	}()

	// Execute the function and record any errors
	err := fn(ctx)
	if err != nil {
		recordErrorWithAttributes(span, err,
			attribute.Int64(AttrSyncNodeID, int64(nodeID)),
			attribute.String(AttrSyncNodeAddress, address),
		)
		span.AddEvent("connection.failure", trace.WithAttributes(
			attribute.String(AttrSyncNode, nodeCompact),
			attribute.String("error", err.Error()),
		))
		if enrichedLogger != nil {
			enrichedLogger.Error("node connection failed",
				zap.Uint32("node_id", nodeID),
				zap.String("address", address),
				zap.Error(err))
		}
	} else {
		span.AddEvent("connection.success", trace.WithAttributes(
			attribute.String(AttrSyncNode, nodeCompact),
		))
		if enrichedLogger != nil {
			enrichedLogger.Debug("node connection successful",
				zap.Uint32("node_id", nodeID),
				zap.String("address", address))
		}
	}

	// Set span status based on error state
	setSpanStatus(span, err)
	return ctx, err
}

// traceNodeRegistryOperation wraps node registry operations with OpenTelemetry tracing
func traceNodeRegistryOperation(
	ctx context.Context,
	operation string,
	logger *zap.Logger,
	fn func(context.Context) error,
) (context.Context, error) {
	ctx, span := tracer.Start(ctx, "registry."+operation)
	defer span.End()

	// Add operation attributes with consistent naming
	span.SetAttributes(
		attribute.String(AttrSyncRegistryOp, operation),
		attribute.String(AttrSyncOperation, "registry"),
		attribute.String(AttrSyncComponent, "sync-service"),
	)

	// Get context-enriched logger
	enrichedLogger := logWithTraceContext(ctx, logger)

	// Capture panics and record them as errors
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in registry operation %s: %v", operation, r)
			recordErrorWithAttributes(span, err,
				attribute.String(AttrSyncRegistryOp, operation),
				attribute.String("panic_value", fmt.Sprintf("%v", r)),
			)
			span.SetStatus(codes.Error, err.Error())
			if enrichedLogger != nil {
				enrichedLogger.Error("panic during registry operation",
					zap.String("operation", operation),
					zap.Any("panic", r))
			}
			panic(r) // Re-panic after recording
		}
	}()

	// Execute the function and record any errors
	err := fn(ctx)
	if err != nil {
		recordErrorWithAttributes(span, err,
			attribute.String(AttrSyncRegistryOp, operation),
		)
		span.AddEvent("registry.failure", trace.WithAttributes(
			attribute.String(AttrSyncRegistryOp, operation),
			attribute.String("error", err.Error()),
		))
		if enrichedLogger != nil {
			enrichedLogger.Error("registry operation failed",
				zap.String("operation", operation),
				zap.Error(err))
		}
	} else {
		span.AddEvent("registry.success", trace.WithAttributes(
			attribute.String(AttrSyncRegistryOp, operation),
		))
		if enrichedLogger != nil {
			enrichedLogger.Debug("registry operation successful",
				zap.String("operation", operation))
		}
	}

	// Set span status based on error state
	setSpanStatus(span, err)
	return ctx, err
}

// traceSyncWorkerLifecycle traces sync worker lifecycle events with OpenTelemetry
func traceSyncWorkerLifecycle(ctx context.Context, logger *zap.Logger) (context.Context, trace.Span, func(context.Context)) {
	ctx, span := tracer.Start(ctx, "worker.lifecycle")

	// Add worker attributes with consistent naming
	span.SetAttributes(
		attribute.String(AttrSyncWorkerType, "sync"),
		attribute.String(AttrSyncComponent, "sync-worker"),
		attribute.String(AttrSyncOperation, "lifecycle"),
	)

	// Get context-enriched logger
	enrichedLogger := logWithTraceContext(ctx, logger)

	// Add worker start event for lifecycle consistency
	span.AddEvent("worker.start", trace.WithAttributes(
		attribute.String(AttrSyncWorkerType, "sync"),
	))

	if enrichedLogger != nil {
		enrichedLogger.Info("sync worker started with tracing")
	}

	// Return context, span, and cleanup function that takes context
	return ctx, span, func(cleanupCtx context.Context) {
		// Capture any panics during shutdown
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic during sync worker shutdown: %v", r)
				recordErrorWithAttributes(span, err,
					attribute.String(AttrOperation, "shutdown"),
					attribute.String("panic_value", fmt.Sprintf("%v", r)),
				)
				span.SetStatus(codes.Error, err.Error())
				if enrichedLogger != nil {
					enrichedLogger.Error("panic during sync worker shutdown", zap.Any("panic", r))
				}
				span.End()
				panic(r) // Re-panic after recording
			}
		}()

		// Use OpenTelemetry best practices: codes.Ok with empty message + event
		span.AddEvent("worker.shutdown.success", trace.WithAttributes(
			attribute.String(AttrSyncWorkerType, "sync"),
		))

		// Set successful status and end span
		setSpanStatus(span, nil)
		span.End()

		if enrichedLogger != nil {
			enrichedLogger.Info("sync worker shutdown completed")
		}
	}
}

// traceSyncOperation provides a generic tracing wrapper for sync operations
func traceSyncOperation(
	ctx context.Context,
	spanName string,
	logger *zap.Logger,
	attributes []attribute.KeyValue,
	fn func(context.Context) error,
) (context.Context, error) {
	ctx, span := tracer.Start(ctx, spanName)
	defer span.End()

	// Add provided attributes plus common sync attributes
	allAttrs := append(syncSpanOptions("generic"), attributes...)
	span.SetAttributes(allAttrs...)

	// Get context-enriched logger
	enrichedLogger := logWithTraceContext(ctx, logger)

	// Capture panics and record them as errors
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in %s: %v", spanName, r)
			recordErrorWithAttributes(span, err,
				attribute.String(AttrOperation, spanName),
				attribute.String("panic_value", fmt.Sprintf("%v", r)),
			)
			span.SetStatus(codes.Error, err.Error())
			if enrichedLogger != nil {
				enrichedLogger.Error("panic during sync operation",
					zap.String("operation", spanName),
					zap.Any("panic", r))
			}
			panic(r) // Re-panic after recording
		}
	}()

	// Execute the function and record any errors
	err := fn(ctx)
	if err != nil {
		recordErrorWithAttributes(span, err,
			attribute.String(AttrOperation, spanName),
		)
		span.AddEvent("operation.failure", trace.WithAttributes(
			attribute.String(AttrOperation, spanName),
			attribute.String("error", err.Error()),
		))
		if enrichedLogger != nil {
			enrichedLogger.Error("sync operation failed",
				zap.String("operation", spanName),
				zap.Error(err))
		}
	} else {
		span.AddEvent("operation.success", trace.WithAttributes(
			attribute.String(AttrOperation, spanName),
		))
		if enrichedLogger != nil {
			enrichedLogger.Debug("sync operation successful",
				zap.String("operation", spanName))
		}
	}

	// Set span status based on error state
	setSpanStatus(span, err)
	return ctx, err
}

// addSpanEvent adds a timestamped event to the current span if one exists
func addSpanEvent(ctx context.Context, name string, attributes ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attributes...))
	}
}
