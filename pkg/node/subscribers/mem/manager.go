package memsubs

import (
	"context"
	"sync"

	crdttypes "github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type Manager struct {
	log        *zap.Logger
	bufferSize int

	topicSubs     map[string]map[chan *crdttypes.Event]struct{}
	topicSubsLock sync.RWMutex

	ctx       context.Context
	ctxCancel func()
}

func New(log *zap.Logger, bufferSize int) *Manager {
	m := &Manager{
		log:        log,
		bufferSize: bufferSize,

		topicSubs: map[string]map[chan *crdttypes.Event]struct{}{},
	}
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
	return m
}

func (m *Manager) Close() error {
	if m.ctxCancel != nil {
		m.ctxCancel()
	}
	m.topicSubsLock.Lock()
	defer m.topicSubsLock.Unlock()
	for _, subs := range m.topicSubs {
		for ch := range subs {
			close(ch)
		}
	}
	return nil
}

func (m *Manager) OnNewEvent(topicId string, ev *crdttypes.Event) {
	m.topicSubsLock.RLock()
	defer m.topicSubsLock.RUnlock()

	// TODO: should this check that the topicId matches ev.Envelope.ContentTopic?

	if _, ok := m.topicSubs[topicId]; !ok {
		return
	}
	for ch := range m.topicSubs[topicId] {
		select {
		case <-m.ctx.Done():
		case ch <- ev:
			// TODO: what to do if channel is full here
		}

	}
}

func (m *Manager) Subscribe(ctx context.Context, topicId string) chan *crdttypes.Event {
	m.topicSubsLock.Lock()
	defer m.topicSubsLock.Unlock()

	ch := make(chan *crdttypes.Event, m.bufferSize)
	if _, ok := m.topicSubs[topicId]; !ok {
		m.topicSubs[topicId] = map[chan *crdttypes.Event]struct{}{}
	}
	m.topicSubs[topicId][ch] = struct{}{}
	return ch
}

func (m *Manager) Unsubscribe(topicId string, ch chan *crdttypes.Event) {
	m.topicSubsLock.Lock()
	defer m.topicSubsLock.Unlock()

	if _, ok := m.topicSubs[topicId]; !ok {
		return
	}
	delete(m.topicSubs[topicId], ch)
	close(ch)
	if len(m.topicSubs[topicId]) == 0 {
		delete(m.topicSubs, topicId)
	}
}
