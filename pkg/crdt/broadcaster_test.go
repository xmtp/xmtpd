package crdt

import (
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/zap"
)

func Test_BasicBroadcast(t *testing.T) {
	net := newNetwork(t, 5, 1)
	defer net.Close()
	net.Publish(t, 0, t0, "hi")
	net.AssertEventuallyConsistent(t, time.Second)
}

// In-memory broadcaster that uses channels to broadcast Events between Nodes.
type chanBroadcaster struct {
	log         *zap.Logger
	subscribers map[*Node]bool
}

func newChanBroadcaster(log *zap.Logger) *chanBroadcaster {
	return &chanBroadcaster{
		log:         log.Named("chanbc"),
		subscribers: make(map[*Node]bool),
	}
}

func (b *chanBroadcaster) NewTopic(name string, n *Node) TopicBroadcaster {
	return &topicChanBroadcaster{
		node:            n,
		log:             n.log.Named(name),
		chanBroadcaster: b,
	}
}

func (b *chanBroadcaster) Broadcast(ev *Event, from *Node) {
	for sub := range b.subscribers {
		if sub == from {
			continue
		}
		t := sub.getOrCreateTopic(ev.ContentTopic)
		t.pendingReceiveEvents <- ev
	}
}

func (b *chanBroadcaster) AddNode(n *Node) {
	b.subscribers[n] = true
}

func (b *chanBroadcaster) RemoveNode(n *Node) {
	delete(b.subscribers, n)
}

type topicChanBroadcaster struct {
	*chanBroadcaster
	node *Node
	log  *zap.Logger
}

func (tb *topicChanBroadcaster) Broadcast(ev *Event) {
	tb.log.Debug("broadcasting", zap.Cid("event", ev.cid))
	tb.chanBroadcaster.Broadcast(ev, tb.node)
}
