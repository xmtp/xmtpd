package message

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/envelopes"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/topic"
)

type listener struct {
	ctx         context.Context
	ch          chan<- []*envelopes.OriginatorEnvelope
	closed      bool
	topics      map[string]struct{}
	originators map[uint32]struct{}
	isGlobal    bool
}

func newListener(
	ctx context.Context,
	logger *zap.Logger,
	query *message_api.EnvelopesQuery,
	ch chan<- []*envelopes.OriginatorEnvelope,
) *listener {
	l := &listener{
		ctx:         ctx,
		ch:          ch,
		topics:      make(map[string]struct{}),
		originators: make(map[uint32]struct{}),
		isGlobal:    false,
	}
	topics := query.GetTopics()
	originators := query.GetOriginatorNodeIds()

	if len(topics) == 0 && len(originators) == 0 {
		l.isGlobal = true
		return l
	}

	for _, t := range topics {
		validatedTopic, err := topic.ParseTopic(t)
		if err != nil {
			logger.Warn("skipping invalid topic", zap.Binary("topic_bytes", t))
			continue
		}
		logger.Debug("adding topic listener", zap.String("topic", validatedTopic.String()))
		l.topics[validatedTopic.String()] = struct{}{}
	}

	for _, originator := range originators {
		l.originators[originator] = struct{}{}
	}

	return l
}

type listenerSet struct {
	sync.Map // map[*listener]struct{}
}

func (ls *listenerSet) addListener(l *listener) {
	ls.Store(l, struct{}{})
}

func (ls *listenerSet) removeListener(l *listener) {
	ls.Delete(l)
}

func (ls *listenerSet) isEmpty() bool {
	empty := true

	ls.Range(func(_, _ any) bool {
		empty = false
		return false // stop iteration
	})
	return empty
}

// Maps from a key to a set of listeners
type listenersMap[K comparable] struct {
	data sync.Map     // map[K]*listenerSet
	mu   sync.RWMutex // ensures mutations are consistent
}

func (lm *listenersMap[K]) addListener(keys map[K]struct{}, l *listener) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	for key := range keys {
		value, _ := lm.data.LoadOrStore(key, &listenerSet{})
		set := value.(*listenerSet)
		set.addListener(l)
	}
}

func (lm *listenersMap[K]) removeListener(keys map[K]struct{}, l *listener) {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	for key := range keys {
		value, ok := lm.data.Load(key)
		if !ok || value == nil {
			return // Key doesn't exist, nothing to do
		}
		set := value.(*listenerSet)
		set.removeListener(l)
		if set.isEmpty() {
			lm.data.Delete(key)
		}
	}
}

func (lm *listenersMap[K]) getListeners(key K) *listenerSet {
	// No lock needed, because we are not mutating lm.data
	if value, ok := lm.data.Load(key); ok {
		return value.(*listenerSet)
	}
	return nil
}

func (lm *listenersMap[K]) rangeKeys(fn func(key K, listeners *listenerSet) bool) {
	lm.data.Range(func(key, value any) bool {
		return fn(key.(K), value.(*listenerSet))
	})
}
