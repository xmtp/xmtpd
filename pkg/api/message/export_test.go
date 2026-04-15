package message

import (
	"github.com/xmtp/xmtpd/pkg/db"
)

// DispatchedMet reports whether the subscribe worker has *dispatched* every
// sequence ID in vc to its listeners. This is stronger than "polled past" —
// a true return guarantees the start() loop has already handed those rows
// off, so any listener registered after this call will not retroactively
// receive them. Tests use this predicate with require.Eventually to pre-seed
// envelopes before opening a stream without racing the dispatch loop.
func (s *Service) DispatchedMet(vc db.VectorClock) bool {
	return s.subscribeWorker.dispatchedMet(vc)
}

// GlobalListenerCount returns the number of global listeners currently
// registered with the subscribe worker. Tests poll this (via require.Eventually)
// to wait for proof that a server-side listener is registered before
// triggering the code under test — no time.Sleep needed.
func (s *Service) GlobalListenerCount() int {
	return s.subscribeWorker.countGlobalListeners()
}

func (s *subscribeWorker) countGlobalListeners() int {
	count := 0
	s.globalListeners.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}
