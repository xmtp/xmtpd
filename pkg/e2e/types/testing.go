// Package types defines shared types for the E2E test framework.
package types

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// TestFailedError is a sentinel error used by the TestingT adapter when
// FailNow() is called. The runner recovers this panic and reports
// the test as failed without crashing the process.
type TestFailedError struct {
	Messages []string
}

func (e *TestFailedError) Error() string {
	if len(e.Messages) == 0 {
		return "test failed"
	}
	return "test failed: " + e.Messages[len(e.Messages)-1]
}

// TestingT implements testify's TestingT interface so that require/assert
// work in E2E tests that run outside of `go test`.
//
// Usage in tests:
//
//	func (t *MyTest) Run(ctx context.Context, env *Environment) error {
//	    require := require.New(env.T())
//	    require.NoError(err)
//	}
//
// The environment's T() is set automatically by the runner before each test.
// Errorf logs the failure message and marks the test as failed.
// FailNow panics with TestFailedError, which the runner recovers.
type TestingT struct {
	logger   *zap.Logger
	mu       sync.Mutex
	failed   bool
	messages []string
}

// NewTestingT creates a new TestingT adapter backed by the given logger.
func NewTestingT(logger *zap.Logger) *TestingT {
	return &TestingT{logger: logger}
}

// Errorf logs a test failure message and marks the test as failed.
// This satisfies the testify TestingT interface.
func (t *TestingT) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	t.mu.Lock()
	t.failed = true
	t.messages = append(t.messages, msg)
	t.mu.Unlock()
	t.logger.Error("assertion failed", zap.String("message", msg))
}

// FailNow marks the test as failed and aborts execution by panicking
// with a TestFailedError. The runner catches this panic and reports
// the test as a failure.
// This satisfies the testify TestingT interface.
func (t *TestingT) FailNow() {
	t.mu.Lock()
	t.failed = true
	msgs := make([]string, len(t.messages))
	copy(msgs, t.messages)
	t.mu.Unlock()
	panic(&TestFailedError{Messages: msgs})
}

// Failed reports whether the test has been marked as failed.
func (t *TestingT) Failed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.failed
}

// Messages returns all failure messages recorded so far.
func (t *TestingT) Messages() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	msgs := make([]string, len(t.messages))
	copy(msgs, t.messages)
	return msgs
}
