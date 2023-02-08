package crdttest

import (
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/xmtp/xmtpd/pkg/crdt/types"
)

type testReplicaSet struct {
	replicas []*testReplica
}

func NewTestReplicaSet(t *testing.T, count int) *testReplicaSet {
	t.Helper()
	replicas := make([]*testReplica, count)
	for i := 0; i < count; i++ {
		replicas[i] = NewTestReplica(t)
	}
	for _, a := range replicas {
		for _, b := range replicas {
			if a == b {
				continue
			}
			a.AddPeer(t, b)
			b.AddPeer(t, a)
		}
	}
	return &testReplicaSet{
		replicas: replicas,
	}
}

func (rs *testReplicaSet) Replicas() []*testReplica {
	return rs.replicas
}

func (rs *testReplicaSet) Visualize(w io.Writer) {
	for i, replica := range rs.Replicas() {
		replica.Visualize(w, fmt.Sprintf("replica%d", i+1))
	}
}

func (rs *testReplicaSet) BroadcastRandom(t *testing.T, count int) []*types.Event {
	t.Helper()
	replica := rs.randomReplica(t, nil)
	return replica.BroadcastRandom(t, count)
}

func (rs *testReplicaSet) randomReplica(t *testing.T, exclude *testReplica) *testReplica {
	t.Helper()
	i := rand.Intn(len(rs.replicas))
	for exclude != nil && rs.replicas[i] == exclude {
		i = rand.Intn(len(rs.replicas))
	}
	return rs.replicas[i]
}

func (rs *testReplicaSet) RequireEventuallyCapturedEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	for _, replica := range rs.replicas {
		replica.RequireEventuallyCapturedEvents(t, expected)
	}
}

func (rs *testReplicaSet) RequireEventuallyStoredEvents(t *testing.T, expected []*types.Event) {
	t.Helper()
	for _, replica := range rs.replicas {
		replica.RequireEventuallyStoredEvents(t, expected)
	}
}
