package tests

import (
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/crdt"
)

func Test_BasicSyncing(t *testing.T) {
	// 3 nodes, one topic "t0"
	net := NewNetwork(t, 3)
	defer net.Close()
	net.Publish(t, 0, t0, "hi")
	net.Publish(t, 1, t0, "hi back")
	// wait for things to settle
	net.AssertEventuallyConsistent(t, time.Second)
	// suspend broadcasts to n1/t0 and publish few things
	net.WithSuspendedTopic(t, 1, t0, func(n *crdt.Node) {
		net.Publish(t, 2, t0, "oh hello")
		net.Publish(t, 2, t0, "how goes")
		net.Publish(t, 1, t0, "how are you")
	})
	// wait for things to settle but ignore n1
	// because it needs a new broadcast to trigger syncing.
	net.AssertEventuallyConsistent(t, time.Second, 1)
	net.Publish(t, 0, t0, "not bad")
	net.AssertEventuallyConsistent(t, time.Second)
}
