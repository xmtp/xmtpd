package tracing

import (
	"sync"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

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

// Span limits for production safety - prevent runaway memory/payload sizes.
const (
	// DefaultTraceContextTTL is the default time-to-live for stored span contexts.
	// 5 minutes is generous - publish_worker typically processes within seconds.
	DefaultTraceContextTTL = 5 * time.Minute

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
	if !apmEnabled.Load() || span == nil {
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
