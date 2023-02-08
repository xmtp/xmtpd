package crdt_test

import (
	"os"
	"testing"

	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
)

func TestReplicaSet_BroadcastStore(t *testing.T) {
	rs := crdttest.NewTestReplicaSet(t, 3)
	events := rs.BroadcastRandom(t, 5)
	rs.RequireEventuallyCapturedEvents(t, events)
	rs.RequireEventuallyStoredEvents(t, events)

	if visualize {
		rs.Visualize(os.Stdout)
	}
}
