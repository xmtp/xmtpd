package memsubs

import (
	"context"
	"sync"

	crdttypes "github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type MemoryManager struct {
	log *zap.Logger

	topicSubs     map[string]map[chan *crdttypes.Event]struct{}
	topicSubsLock sync.RWMutex

	ctx       context.Context
	ctxCancel func()
}

func New(log *zap.Logger) *MemoryManager {
	m := &MemoryManager{
		log: log,

		topicSubs: map[string]map[chan *crdttypes.Event]struct{}{},
	}
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
	return m
}

func (m *MemoryManager) Close() error {
	if m.ctxCancel != nil {
		m.ctxCancel()
	}
	return nil
}

func (m *MemoryManager) OnNewEvent(topicId string, ev *crdttypes.Event) {
	m.topicSubsLock.RLock()
	defer m.topicSubsLock.RUnlock()

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

func (m *MemoryManager) Subscribe(ctx context.Context, topicId string, buffer int) chan *crdttypes.Event {
	m.topicSubsLock.Lock()
	defer m.topicSubsLock.Unlock()

	ch := make(chan *crdttypes.Event, buffer)
	if _, ok := m.topicSubs[topicId]; !ok {
		m.topicSubs[topicId] = map[chan *crdttypes.Event]struct{}{}
	}
	m.topicSubs[topicId][ch] = struct{}{}
	return ch
}

func (m *MemoryManager) Unsubscribe(topicId string, ch chan *crdttypes.Event) {
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
