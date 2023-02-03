package tests

import (
	"math/rand"
	"sync"

	mh "github.com/multiformats/go-multihash"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// In-memory syncer that implements fetching by
// reaching directly into a random Node's store.
type randomSyncer struct {
	sync.RWMutex
	nodes []*crdt.Node
}

func NewRandomSyncer() *randomSyncer {
	return &randomSyncer{}
}

func (s *randomSyncer) AddNode(n *crdt.Node) {
	s.Lock()
	defer s.Unlock()
	s.nodes = append(s.nodes, n)
}

func (s *randomSyncer) GetRandomNode() *crdt.Node {
	s.RLock()
	defer s.RUnlock()
	return s.nodes[rand.Intn(len(s.nodes))]
}

func (s *randomSyncer) NewTopic(name string, n *crdt.Node) crdt.TopicSyncer {
	return &randomTopicSyncer{
		randomSyncer: s,
		node:         n,
		topic:        name,
		log:          n.LogNamed(name),
	}
}

type randomTopicSyncer struct {
	*randomSyncer
	node  *crdt.Node
	topic string
	log   *zap.Logger
}

func (s *randomTopicSyncer) Fetch(cids []mh.Multihash) (results []*crdt.Event, err error) {
	node := s.GetRandomNode()
	s.log.Debug("fetching", zap.Cids("cids", cids...))
	for _, cid := range cids {
		ev, err := node.Get(s.topic, cid)
		if err != nil {
			return nil, err
		}
		results = append(results, ev)
	}
	return results, nil
}
