package crdt_test

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

func TestReplicaSet_BroadcastStore(t *testing.T) {
	rs := newTestReplicaSet(t, 3)
	events := rs.broadcastRandom(t, 10)
	rs.requireEventuallyCapturedEvents(t, events)
	rs.requireEventuallyStoredEvents(t, events)

	if visualize {
		rs.visualize(os.Stdout)
	}
}

type testReplicaSet struct {
	replicas []*testReplica
}

func newTestReplicaSet(t *testing.T, count int) *testReplicaSet {
	t.Helper()
	replicas := make([]*testReplica, count)
	for i := 0; i < count; i++ {
		replicas[i] = newTestReplica(t)
	}
	for _, a := range replicas {
		for _, b := range replicas {
			if a == b {
				continue
			}
			a.addPeer(t, b)
			b.addPeer(t, a)
		}
	}
	return &testReplicaSet{
		replicas: replicas,
	}
}

func (rs *testReplicaSet) visualize(w io.Writer) {
	for i, replica := range rs.replicas {
		replica.visualize(w, fmt.Sprintf("replica%d", i+1))
	}
}

func (rs *testReplicaSet) broadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	events := make([]*types.Event, count)
	// Emulate concurrent appends across replicas.
	for i := 0; i < count; i++ {
		replica := rs.randomReplica(t)
		evs := replica.broadcastRandom(t, 1)
		events[i] = evs[0]
		if i%count == 0 {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		}
	}
	return events
}

func (rs *testReplicaSet) randomReplica(t *testing.T) *testReplica {
	t.Helper()
	i := rand.Intn(len(rs.replicas))
	return rs.replicas[i]
}

func (rs *testReplicaSet) requireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	for _, replica := range rs.replicas {
		replica.requireEventuallyCapturedEvents(t, expected)
	}
}

func (rs *testReplicaSet) requireEventuallyStoredEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	for _, replica := range rs.replicas {
		replica.requireEventuallyStoredEvents(t, expected)
	}
}
