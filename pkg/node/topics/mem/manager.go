package memtopics

import (
	"context"
	"errors"
	"sync"

	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/zap"
)

var (
	ErrTopicAlreadyExists = errors.New("topic already exists")
)

type Manager struct {
	log               *zap.Logger
	newTopicReplicaFn NewTopicReplicaFunc

	topics     map[string]*crdt.Replica
	topicsLock sync.RWMutex

	ctx       context.Context
	ctxCancel func()
}

type NewTopicReplicaFunc func(topicId string) (*crdt.Replica, error)

func New(log *zap.Logger, newTopicReplicaFn NewTopicReplicaFunc) (*Manager, error) {
	log = log.Named("message/v1")
	m := &Manager{
		log:               log,
		newTopicReplicaFn: newTopicReplicaFn,

		topics: map[string]*crdt.Replica{},
	}
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
	return m, nil
}

func (m *Manager) Close() error {
	if m.ctxCancel != nil {
		m.ctxCancel()
	}
	for _, replica := range m.topics {
		replica.Close()
	}
	return nil
}

func (m *Manager) GetOrCreateTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	topic, err := m.getTopic(ctx, topicId)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		topic, err = m.createTopic(ctx, topicId)
		if err != nil {
			return nil, err
		}
	}
	return topic, nil
}

func (m *Manager) getTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	m.topicsLock.RLock()
	defer m.topicsLock.RUnlock()
	topic, ok := m.topics[topicId]
	if !ok {
		return nil, nil
	}
	return topic, nil
}

func (m *Manager) createTopic(ctx context.Context, topicId string) (*crdt.Replica, error) {
	m.topicsLock.Lock()
	defer m.topicsLock.Unlock()
	if _, ok := m.topics[topicId]; ok {
		return nil, ErrTopicAlreadyExists
	}
	return m.newTopicReplicaFn(topicId)
}
