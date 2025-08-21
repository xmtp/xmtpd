package sync

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap/zaptest"
)

func TestTraceNodeConnection_Success(t *testing.T) {
	executed := false
	ctx, err := traceNodeConnection(context.Background(), 123, "http://test.com", nil,
		func(ctx context.Context) error {
			executed = true
			return nil
		})

	assert.NoError(t, err)
	assert.True(t, executed)
	assert.NotNil(t, ctx)
}

func TestTraceNodeConnection_Error(t *testing.T) {
	testErr := errors.New("connection failed")
	ctx, err := traceNodeConnection(context.Background(), 123, "http://test.com", nil,
		func(ctx context.Context) error {
			return testErr
		})

	assert.Equal(t, testErr, err)
	assert.NotNil(t, ctx)
}

func TestTraceSyncWorkerLifecycle(t *testing.T) {
	logger := zaptest.NewLogger(t)

	ctx, span, cleanup := traceSyncWorkerLifecycle(context.Background(), logger)

	assert.NotNil(t, ctx)
	assert.NotNil(t, span)
	assert.NotNil(t, cleanup)

	// Should be safe to call cleanup multiple times
	cleanup(ctx)
	cleanup(ctx)
}

// Add missing test functions that were referenced in previous implementations
func TestTraceNodeRegistryOperation_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	executed := false

	ctx, err := traceNodeRegistryOperation(
		context.Background(),
		"test_operation",
		logger,
		func(ctx context.Context) error {
			executed = true
			return nil
		},
	)

	assert.NoError(t, err)
	assert.True(t, executed)
	assert.NotNil(t, ctx)
}

func TestTraceNodeRegistryOperation_Error(t *testing.T) {
	logger := zaptest.NewLogger(t)
	testErr := errors.New("registry operation failed")

	ctx, err := traceNodeRegistryOperation(
		context.Background(),
		"test_operation",
		logger,
		func(ctx context.Context) error {
			return testErr
		},
	)

	assert.Equal(t, testErr, err)
	assert.NotNil(t, ctx)
}

func TestTraceSyncOperation_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	executed := false

	ctx, err := traceSyncOperation(
		context.Background(),
		"test.operation",
		logger,
		nil,
		func(ctx context.Context) error {
			executed = true
			return nil
		},
	)

	assert.NoError(t, err)
	assert.True(t, executed)
	assert.NotNil(t, ctx)
}

func TestTraceSyncOperation_Error(t *testing.T) {
	logger := zaptest.NewLogger(t)
	testErr := errors.New("sync operation failed")

	ctx, err := traceSyncOperation(
		context.Background(),
		"test.operation",
		logger,
		nil,
		func(ctx context.Context) error {
			return testErr
		},
	)

	assert.Equal(t, testErr, err)
	assert.NotNil(t, ctx)
}

func TestTraceSyncOperation_Panic(t *testing.T) {
	logger := zaptest.NewLogger(t)

	assert.Panics(t, func() {
		_, _ = traceSyncOperation(
			context.Background(),
			"test.operation",
			logger,
			nil,
			func(ctx context.Context) error {
				panic("test panic")
			},
		)
	})
}

func TestAddSpanEvent(t *testing.T) {
	ctx := context.Background()
	// Should not panic when called with context without span
	addSpanEvent(ctx, "test.event")
}

func TestSyncSpanOptions(t *testing.T) {
	options := syncSpanOptions("test")

	assert.Len(t, options, 2)
	assert.Contains(t, options, attribute.String(AttrSyncOperation, "test"))
	assert.Contains(t, options, attribute.String(AttrSyncComponent, "sync-service"))
}
