package message

import (
	"context"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
)

// AwaitCursor blocks until the subscribe worker has polled past all sequence IDs in vc.
// Only compiled during testing (export_test.go pattern).
func (s *Service) AwaitCursor(ctx context.Context, vc db.VectorClock) error {
	const checkInterval = 5 * time.Millisecond
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		if s.subscribeWorker.subscriptions.cursorMet(vc) {
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

func (s *subscriptionHandler) cursorMet(vc db.VectorClock) bool {
	s.Lock()
	defer s.Unlock()
	for nodeID, minSeq := range vc {
		poller, ok := s.subs[nodeID]
		if !ok || uint64(poller.sub.LastSeen()) < minSeq {
			return false
		}
	}
	return true
}
