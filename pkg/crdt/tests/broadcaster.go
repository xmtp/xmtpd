package tests

import (
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/zap"
)

// In-memory broadcaster that uses channels to broadcast Events between Nodes.
type chanBroadcaster struct {
	log         *zap.Logger
	subscribers map[*crdt.Node]bool
}

func NewChanBroadcaster(log *zap.Logger) crdt.NodeBroadcaster {
	return &chanBroadcaster{
		log:         log.Named("chanbc"),
		subscribers: make(map[*crdt.Node]bool),
	}
}

func (b *chanBroadcaster) NewTopic(name string, n *crdt.Node) crdt.TopicBroadcaster {
	return &topicChanBroadcaster{
		node:            n,
		log:             n.LogNamed(name),
		chanBroadcaster: b,
	}
}

func (b *chanBroadcaster) Broadcast(ev *crdt.Event, from *crdt.Node) {
	for sub := range b.subscribers {
		if sub == from {
			continue
		}
		t, err := sub.GetOrCreateTopic(ev.ContentTopic)
		if err != nil {
			b.log.Error("receiving event", zap.Error(err), zap.Cid("event", ev.Cid), zap.String("topic", ev.ContentTopic))
		}
		t.ReceiveEvent(ev)
	}
}

func (b *chanBroadcaster) AddNode(n *crdt.Node) {
	b.subscribers[n] = true
}

func (b *chanBroadcaster) RemoveNode(n *crdt.Node) {
	delete(b.subscribers, n)
}

type topicChanBroadcaster struct {
	*chanBroadcaster
	node *crdt.Node
	log  *zap.Logger
}

func (tb *topicChanBroadcaster) Broadcast(ev *crdt.Event) {
	tb.log.Debug("broadcasting", zap.Cid("event", ev.Cid))
	tb.chanBroadcaster.Broadcast(ev, tb.node)
}
