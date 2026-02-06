package tracing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

// TestIntegration_SpanHierarchy verifies that spans are created with correct
// parent-child relationships, simulating a real request flow.
func TestIntegration_SpanHierarchy(t *testing.T) {
	cleanup := tracing.SetEnabledForTesting(true)
	defer cleanup()
	mt := mocktracer.Start()
	defer mt.Stop()

	// Simulate a publish request flow
	ctx := context.Background()

	// 1. Parent span (API request)
	parentSpan, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanNodePublishPayerEnvelopes)
	tracing.SpanTag(parentSpan, "num_envelopes", 1)

	// 2. Child span (staging transaction)
	txSpan, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanNodeStageTransaction)
	tracing.SpanTag(txSpan, "staged_id", 12345)

	// 3. Grandchild span (DB query) - should be child of transaction
	dbSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanDBQuery)
	tracing.SpanTag(dbSpan, "db.role", "writer")
	tracing.SpanTag(dbSpan, "db.statement", "INSERT INTO staged_originator_envelopes...")
	dbSpan.Finish()

	txSpan.Finish()
	parentSpan.Finish()

	// Verify spans
	spans := mt.FinishedSpans()
	require.Len(t, spans, 3, "expected 3 spans")

	// Find spans by operation name
	var parent, tx, db mocktracer.Span
	for _, s := range spans {
		switch s.OperationName() {
		case tracing.SpanNodePublishPayerEnvelopes:
			parent = s
		case tracing.SpanNodeStageTransaction:
			tx = s
		case tracing.SpanDBQuery:
			db = s
		}
	}

	require.NotNil(t, parent, "parent span not found")
	require.NotNil(t, tx, "transaction span not found")
	require.NotNil(t, db, "db span not found")

	// Verify hierarchy: db -> tx -> parent
	assert.Equal(t, tx.SpanID(), db.ParentID(), "db span should be child of tx span")
	assert.Equal(t, parent.SpanID(), tx.ParentID(), "tx span should be child of parent span")
	assert.Equal(t, uint64(0), parent.ParentID(), "parent span should have no parent")

	// Verify all spans share the same trace ID
	assert.Equal(t, parent.TraceID(), tx.TraceID(), "tx should share trace ID with parent")
	assert.Equal(t, parent.TraceID(), db.TraceID(), "db should share trace ID with parent")

	// Verify tags (mocktracer stores numbers as float64)
	assert.Equal(t, float64(1), parent.Tag("num_envelopes"))
	assert.Equal(t, float64(12345), tx.Tag("staged_id"))
	assert.Equal(t, "writer", db.Tag("db.role"))
}

// TestIntegration_AsyncContextPropagation verifies that TraceContextStore
// correctly propagates context across async boundaries.
func TestIntegration_AsyncContextPropagation(t *testing.T) {
	cleanup := tracing.SetEnabledForTesting(true)
	defer cleanup()
	mt := mocktracer.Start()
	defer mt.Stop()

	store := tracing.NewTraceContextStore()
	stagedID := int64(67890)

	// Simulate staging side (API request)
	ctx := context.Background()
	apiSpan, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanNodePublishPayerEnvelopes)
	tracing.SpanTag(apiSpan, "staged_id", stagedID)

	// Store context for async propagation
	store.Store(stagedID, apiSpan)
	apiSpan.Finish()

	// Simulate worker side (async processing)
	parentCtx := store.Retrieve(stagedID)
	require.NotNil(t, parentCtx, "should retrieve stored context")

	// Create child span linked to parent
	workerSpan := tracing.StartSpanWithParent(tracing.SpanPublishWorkerProcess, parentCtx)
	tracing.SpanTag(workerSpan, "trace_linked", true)
	workerSpan.Finish()

	// Verify spans
	spans := mt.FinishedSpans()
	require.Len(t, spans, 2, "expected 2 spans")

	var api, worker mocktracer.Span
	for _, s := range spans {
		switch s.OperationName() {
		case tracing.SpanNodePublishPayerEnvelopes:
			api = s
		case tracing.SpanPublishWorkerProcess:
			worker = s
		}
	}

	require.NotNil(t, api, "api span not found")
	require.NotNil(t, worker, "worker span not found")

	// Verify async linking worked
	assert.Equal(t, api.SpanID(), worker.ParentID(), "worker should be child of api span")
	assert.Equal(t, api.TraceID(), worker.TraceID(), "worker should share trace ID")
	assert.Equal(t, "true", worker.Tag("trace_linked")) // mocktracer stores bools as strings
}

// TestIntegration_ErrorTagging verifies that errors are properly tagged on spans.
func TestIntegration_ErrorTagging(t *testing.T) {
	cleanup := tracing.SetEnabledForTesting(true)
	defer cleanup()
	mt := mocktracer.Start()
	defer mt.Stop()

	ctx := context.Background()
	span, _ := tracing.StartSpanFromContext(ctx, "test.error_operation")

	// Simulate an error
	testErr := assert.AnError
	span.Finish(tracing.WithError(testErr))

	spans := mt.FinishedSpans()
	require.Len(t, spans, 1)

	// Verify error is tagged (the error value is stored, not just true)
	assert.NotNil(t, spans[0].Tag("error"))
}

// TestIntegration_DBSubscriptionTriggerTags verifies the trigger tag pattern
// that would catch the read-replica bug.
func TestIntegration_DBSubscriptionTriggerTags(t *testing.T) {
	cleanup := tracing.SetEnabledForTesting(true)
	defer cleanup()
	mt := mocktracer.Start()
	defer mt.Stop()

	// Simulate notification-triggered poll
	ctx := context.Background()
	span1, _ := tracing.StartSpanFromContext(ctx, tracing.SpanDBSubscriptionPoll)
	tracing.SpanTag(span1, tracing.TagTrigger, tracing.TriggerNotification)
	tracing.SpanTag(span1, "num_results", 5)
	span1.Finish()

	// Simulate timer fallback poll (the bug indicator!)
	span2, _ := tracing.StartSpanFromContext(ctx, tracing.SpanDBSubscriptionPoll)
	tracing.SpanTag(span2, tracing.TagTrigger, tracing.TriggerTimerFallback)
	tracing.SpanTag(span2, "num_results", 0)
	tracing.SpanTag(span2, tracing.TagNotificationMiss, true)
	span2.Finish()

	spans := mt.FinishedSpans()
	require.Len(t, spans, 2)

	// Find the timer fallback span
	var timerFallbackSpan mocktracer.Span
	for _, s := range spans {
		if s.Tag(tracing.TagTrigger) == tracing.TriggerTimerFallback {
			timerFallbackSpan = s
			break
		}
	}

	require.NotNil(t, timerFallbackSpan, "timer fallback span not found")
	assert.Equal(t, tracing.TriggerTimerFallback, timerFallbackSpan.Tag(tracing.TagTrigger))
	assert.Equal(t, "true", timerFallbackSpan.Tag(tracing.TagNotificationMiss)) // mocktracer stores bools as strings
	assert.Equal(t, float64(0), timerFallbackSpan.Tag("num_results"))           // mocktracer stores ints as float64
}

// TestIntegration_CrossNodeReplication verifies the sync worker span flow.
func TestIntegration_CrossNodeReplication(t *testing.T) {
	cleanup := tracing.SetEnabledForTesting(true)
	defer cleanup()
	mt := mocktracer.Start()
	defer mt.Stop()

	ctx := context.Background()

	// Simulate sync flow
	connectSpan := tracing.StartSpan(tracing.SpanSyncConnectToNode)
	tracing.SpanTag(connectSpan, tracing.TagTargetNode, 200)
	connectSpan.Finish()

	setupSpan, ctx := tracing.StartSpanFromContext(ctx, tracing.SpanSyncSetupStream)
	tracing.SpanTag(setupSpan, tracing.TagTargetNode, 200)

	subscribeSpan, _ := tracing.StartSpanFromContext(ctx, tracing.SpanSyncSubscribe)
	subscribeSpan.Finish()

	setupSpan.Finish()

	// Simulate receiving a batch
	batchSpan := tracing.StartSpan(tracing.SpanSyncReceiveBatch)
	tracing.SpanTag(batchSpan, tracing.TagSourceNode, 200)
	tracing.SpanTag(batchSpan, tracing.TagNumEnvelopes, 10)
	batchSpan.Finish()

	spans := mt.FinishedSpans()
	require.Len(t, spans, 4)

	// Verify we can find all the sync spans
	opNames := make(map[string]bool)
	for _, s := range spans {
		opNames[s.OperationName()] = true
	}

	assert.True(t, opNames[tracing.SpanSyncConnectToNode])
	assert.True(t, opNames[tracing.SpanSyncSetupStream])
	assert.True(t, opNames[tracing.SpanSyncSubscribe])
	assert.True(t, opNames[tracing.SpanSyncReceiveBatch])
}
