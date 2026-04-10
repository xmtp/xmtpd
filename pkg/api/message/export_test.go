package message

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
)

// AwaitCursor blocks until the subscribe worker has *dispatched* every
// sequence ID in vc to its listeners. This is stronger than "polled past" —
// it guarantees the start() loop has already handed those rows off, so any
// listener registered after this call will not retroactively receive them.
// Tests rely on this guarantee to pre-seed envelopes before opening a stream
// without racing the worker's dispatch loop.
func (s *Service) AwaitCursor(ctx context.Context, vc db.VectorClock) error {
	const checkInterval = 5 * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		if s.subscribeWorker.dispatchedMet(vc) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// AwaitGlobalListeners blocks until at least n global listeners are registered
// with the subscribe worker. Tests use this to eliminate the races between
// opening a SubscribeAllEnvelopes stream and inserting envelopes — instead of
// sleeping, they wait for proof that the server-side listener is registered
// before triggering the code under test.
func (s *Service) AwaitGlobalListeners(ctx context.Context, n int) error {
	const checkInterval = 5 * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		if s.subscribeWorker.countGlobalListeners() >= n {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (s *subscribeWorker) countGlobalListeners() int {
	count := 0
	s.globalListeners.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}
