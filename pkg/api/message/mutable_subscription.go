package message

import (
	"context"

	"github.com/xmtp/xmtpd/pkg/envelopes"
)

// mutableSubscription is an in-place-mutable topic subscription on the subscribeWorker, used by
// the XIP-83 bidirectional Subscribe RPC. Unlike subscribeWorker.listen (a fixed query), its
// topic set is grown and shrunk over the life of the stream via addTopics / removeTopics. It is
// never global: an empty topic set delivers nothing (a fresh stream that has subscribed to
// nothing yet), in contrast to newListener, which treats an empty query as "all envelopes".
//
// All mutation methods are intended to be called from the single owning goroutine (the Subscribe
// handler's writer loop). The worker's reap path (closeListener) may also touch the underlying
// listener's topic set concurrently; both sides take listener.topicsMu, so that is safe.
type mutableSubscription struct {
	worker *subscribeWorker
	l      *listener
	// ch is the receive end of the listener channel. The worker pushes envelope batches here;
	// it CLOSES the channel when it reaps the listener (ctx done or the consumer fell behind),
	// which the writer loop observes as "torn down, reconnect from cursors".
	ch <-chan []*envelopes.OriginatorEnvelope
}

// newMutableSubscription registers an empty, non-global topic subscription on the worker and
// returns a handle the caller mutates over the stream's lifetime. Call close when done.
func (s *subscribeWorker) newMutableSubscription(ctx context.Context) *mutableSubscription {
	ch := make(chan []*envelopes.OriginatorEnvelope, subscriptionBufferSize)
	l := &listener{
		ctx:         ctx,
		ch:          ch,
		topics:      make(map[string]struct{}),
		originators: make(map[uint32]struct{}),
		isGlobal:    false, // empty topic set => delivers nothing, NOT all-envelopes
	}
	return &mutableSubscription{worker: s, l: l, ch: ch}
}

// addTopics begins delivering the given topic keys to this subscription. Idempotent: a key
// already subscribed is a no-op. A no-op once the worker has reaped the listener.
func (m *mutableSubscription) addTopics(keys map[string]struct{}) {
	if len(keys) == 0 {
		return
	}
	m.l.topicsMu.Lock()
	defer m.l.topicsMu.Unlock()
	if m.l.closed {
		return
	}
	m.worker.topicListeners.addListener(keys, m.l)
	for k := range keys {
		m.l.topics[k] = struct{}{}
	}
}

// removeTopics stops delivering the given topic keys. Idempotent.
func (m *mutableSubscription) removeTopics(keys map[string]struct{}) {
	if len(keys) == 0 {
		return
	}
	m.l.topicsMu.Lock()
	defer m.l.topicsMu.Unlock()
	m.worker.topicListeners.removeListener(keys, m.l)
	for k := range keys {
		delete(m.l.topics, k)
	}
}

// close unregisters the subscription from every topic it still holds. Safe to call even after
// the worker has already reaped the listener (removeListener is idempotent).
func (m *mutableSubscription) close() {
	m.l.topicsMu.Lock()
	defer m.l.topicsMu.Unlock()
	if len(m.l.topics) > 0 {
		m.worker.topicListeners.removeListener(m.l.topics, m.l)
		m.l.topics = make(map[string]struct{})
	}
}
