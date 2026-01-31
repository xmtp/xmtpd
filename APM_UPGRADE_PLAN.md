# APM Instrumentation A+ Upgrade Plan

## Current Grade: B+ (7.5/10)
## Target Grade: A+ (9.5/10)

---

## Tasks

### 1. Sampling Configuration (Production Readiness)
**Priority: Critical**
**Files:** `pkg/tracing/tracing.go`

- [ ] Add `APM_SAMPLE_RATE` environment variable (0.0 - 1.0)
- [ ] Add `APM_ENABLED` environment variable (true/false)
- [ ] Default to 100% in dev/test, 10% in production
- [ ] Document in code comments

### 2. Consistent Span Naming (Code Quality)
**Priority: High**
**Files:** `pkg/tracing/spans.go` (new)

- [ ] Create constants for all span operation names
- [ ] Use consistent `service.operation` pattern
- [ ] Refactor existing code to use constants
- [ ] Pattern: `xmtpd.{component}.{operation}`

### 3. Subscribe Worker Tracing (Coverage)
**Priority: High**
**Files:** `pkg/api/message/subscribe_worker.go`

- [ ] Add span for batch dispatch
- [ ] Add span for listener management
- [ ] Tag with subscriber count, envelope count
- [ ] Track channel full / dropped messages

### 4. Unit Tests (Testing)
**Priority: High**
**Files:** `pkg/tracing/tracing_test.go` (new)

- [ ] Test TraceContextStore Store/Retrieve
- [ ] Test TTL expiration
- [ ] Test cleanup runs correctly
- [ ] Test Size() method

### 5. Integration Test (Testing)
**Priority: Medium**
**Files:** `pkg/tracing/integration_test.go` (new)

- [ ] Verify spans are created with correct parent/child relationships
- [ ] Verify DB spans are children of request spans
- [ ] Use mock tracer for testing

### 6. Documentation (Operational)
**Priority: Medium**
**Files:** `pkg/tracing/README.md` (new)

- [ ] Overview of instrumentation
- [ ] Configuration options
- [ ] Example Datadog queries for common debugging
- [ ] Troubleshooting guide

---

## Execution Order

1. Sampling configuration (unblocks production)
2. Span naming constants (improves maintainability)
3. Subscribe worker tracing (completes coverage)
4. Unit tests (validates core logic)
5. Documentation (enables operators)

---

## Success Criteria

- [ ] All tasks completed
- [ ] Build passes
- [ ] Tests pass
- [ ] No new linter warnings
