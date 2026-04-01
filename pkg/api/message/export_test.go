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
