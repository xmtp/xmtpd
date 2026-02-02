package tracing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

func TestTraceContextStore_StoreAndRetrieve(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	store := NewTraceContextStore()

	// Create a span to store
	span := StartSpan("test.operation")
	stagedID := int64(12345)

	// Store the span context
	store.Store(stagedID, span)
	assert.Equal(t, 1, store.Size())

	// Retrieve should return the context and remove it
	ctx := store.Retrieve(stagedID)
	require.NotNil(t, ctx)
	assert.Equal(t, 0, store.Size())

	// Second retrieve should return nil
	ctx2 := store.Retrieve(stagedID)
	assert.Nil(t, ctx2)

	span.Finish()
}

func TestTraceContextStore_RetrieveNonExistent(t *testing.T) {
	store := NewTraceContextStore()

	// Retrieve non-existent ID should return nil
	ctx := store.Retrieve(99999)
	assert.Nil(t, ctx)
}

func TestTraceContextStore_StoreNilSpan(t *testing.T) {
	store := NewTraceContextStore()

	// Storing nil span should be safe and not add entry
	store.Store(12345, nil)
	assert.Equal(t, 0, store.Size())
}

func TestTraceContextStore_TTLExpiration(t *testing.T) {
	store := NewTraceContextStore()
	// Set short TTL for testing
	store.ttl = 50 * time.Millisecond

	mt := mocktracer.Start()
	defer mt.Stop()

	span := StartSpan("test.operation")
	stagedID := int64(12345)

	store.Store(stagedID, span)
	assert.Equal(t, 1, store.Size())

	// Wait for TTL to expire
	time.Sleep(100 * time.Millisecond)

	// Retrieve should return nil for expired entry
	ctx := store.Retrieve(stagedID)
	assert.Nil(t, ctx)

	// Entry should be removed
	assert.Equal(t, 0, store.Size())

	span.Finish()
}

func TestTraceContextStore_CleanupRemovesExpired(t *testing.T) {
	store := NewTraceContextStore()
	store.ttl = 50 * time.Millisecond

	mt := mocktracer.Start()
	defer mt.Stop()

	// Store first span
	span1 := StartSpan("test.operation1")
	store.Store(1, span1)

	// Wait for it to expire
	time.Sleep(100 * time.Millisecond)

	// Force cleanup to run on next store by setting lastCleanup in the past
	store.lastCleanup = time.Now().Add(-2 * time.Minute)

	// Store second span - this should trigger cleanup of expired span1
	span2 := StartSpan("test.operation2")
	store.Store(2, span2)

	// First entry should be cleaned up, only second should remain
	assert.Equal(t, 1, store.Size())

	// Verify first is gone (either cleaned up or expired on retrieve)
	assert.Nil(t, store.Retrieve(1))

	// Second should still be retrievable
	ctx2 := store.Retrieve(2)
	assert.NotNil(t, ctx2)

	span1.Finish()
	span2.Finish()
}

func TestTraceContextStore_ConcurrentAccess(t *testing.T) {
	store := NewTraceContextStore()

	mt := mocktracer.Start()
	defer mt.Stop()

	// Test concurrent store/retrieve operations
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := int64(0); i < 100; i++ {
			span := StartSpan("test.operation")
			store.Store(i, span)
			span.Finish()
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := int64(0); i < 100; i++ {
			store.Retrieve(i)
		}
		done <- true
	}()

	// Wait for both to complete
	<-done
	<-done

	// Should not panic or deadlock
}

func TestGetSampleRate(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected float64
	}{
		{"production defaults to 10%", "production", 0.1},
		{"prod defaults to 10%", "prod", 0.1},
		{"staging defaults to 10%", "staging", 0.1},
		{"dev defaults to 100%", "dev", 1.0},
		{"test defaults to 100%", "test", 1.0},
		{"empty defaults to 100%", "", 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the internal logic without env vars
			rate := getSampleRate(tt.env, nil)
			assert.Equal(t, tt.expected, rate)
		})
	}
}

func TestStartSpanWithParent(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	// Test with nil parent - should create root span
	span1 := StartSpanWithParent("test.root", nil)
	require.NotNil(t, span1)
	span1.Finish()

	// Test with parent context - should create child span
	parentSpan := StartSpan("test.parent")
	childSpan := StartSpanWithParent("test.child", parentSpan.Context())
	require.NotNil(t, childSpan)

	childSpan.Finish()
	parentSpan.Finish()

	// Verify spans were created
	spans := mt.FinishedSpans()
	assert.GreaterOrEqual(t, len(spans), 3)
}

func TestSpanTag_TruncatesLongStrings(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	span := StartSpan("test.operation")

	// Create a string longer than MaxTagValueLength
	longString := make([]byte, MaxTagValueLength+500)
	for i := range longString {
		longString[i] = 'x'
	}

	// Tag with the long string
	SpanTag(span, "long_value", string(longString))
	span.Finish()

	spans := mt.FinishedSpans()
	require.Len(t, spans, 1)

	// Verify the tag was truncated
	tagValue := spans[0].Tag("long_value").(string)
	assert.LessOrEqual(t, len(tagValue), MaxTagValueLength+20) // Allow for "[truncated]" suffix
	assert.Contains(t, tagValue, "...[truncated]")
}

func TestSpanTag_ShortStringsUnchanged(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	span := StartSpan("test.operation")

	shortString := "this is a short string"
	SpanTag(span, "short_value", shortString)
	span.Finish()

	spans := mt.FinishedSpans()
	require.Len(t, spans, 1)

	// Verify the tag was NOT truncated
	tagValue := spans[0].Tag("short_value").(string)
	assert.Equal(t, shortString, tagValue)
}

func TestSpanTag_NonStringsUnchanged(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	span := StartSpan("test.operation")

	// Non-string values should pass through unchanged
	SpanTag(span, "int_value", 12345)
	SpanTag(span, "bool_value", true)
	SpanTag(span, "float_value", 3.14)
	span.Finish()

	spans := mt.FinishedSpans()
	require.Len(t, spans, 1)

	// mocktracer stores numbers as float64
	assert.Equal(t, float64(12345), spans[0].Tag("int_value"))
	assert.Equal(t, "true", spans[0].Tag("bool_value")) // mocktracer converts to string
	assert.Equal(t, 3.14, spans[0].Tag("float_value"))
}

func TestTraceContextStore_MaxSizeLimit(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	store := NewTraceContextStore()

	// Fill store to capacity
	spans := make([]Span, MaxStoreSize+100)
	for i := 0; i < MaxStoreSize; i++ {
		spans[i] = StartSpan("test.operation")
		store.Store(int64(i), spans[i])
	}

	assert.Equal(t, MaxStoreSize, store.Size(), "store should be at capacity")

	// Try to add one more - should be dropped
	extraSpan := StartSpan("test.extra")
	store.Store(int64(MaxStoreSize+1), extraSpan)

	assert.Equal(t, MaxStoreSize, store.Size(), "store should still be at capacity (new entry dropped)")

	// Clean up spans
	for _, s := range spans {
		if s != nil {
			s.Finish()
		}
	}
	extraSpan.Finish()
}
