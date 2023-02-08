package crdt_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	crdt "github.com/xmtp/xmtpd/pkg/crdt"
	chanbroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/chan"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	chansyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/chan"
	crdttest "github.com/xmtp/xmtpd/pkg/crdt/testing"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

var (
	visualize bool
)

func init() {
	flag.BoolVar(&visualize, "visualize", false, "output graphviz depiction of replicas")
}

func TestReplica_NewClose(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	replica, err := crdt.NewReplica(ctx, log, store, bc, syncer, nil)
	require.NoError(t, err)
	require.NotNil(t, replica)

	err = replica.Close()
	require.NoError(t, err)
}

func TestReplica_BroadcastStore_SingleReplica(t *testing.T) {
	replica := crdttest.NewTestReplica(t)
	defer replica.Close()

	events := replica.BroadcastRandom(t, 1)
	replica.RequireEventuallyCapturedEvents(t, events)
	replica.RequireEventuallyStoredEvents(t, events)
}

func TestReplica_BroadcastStore_TwoReplicas(t *testing.T) {
	replica1 := crdttest.NewTestReplica(t)
	defer replica1.Close()

	replica2 := crdttest.NewTestReplica(t)
	defer replica2.Close()

	// Add replica2 as peer of replica1, broadcast events via replica1, and
	// expect that both replicas eventually capture and store the events.
	replica1.AddPeer(t, replica2)

	events1 := replica1.BroadcastRandom(t, 1)

	replica1.RequireEventuallyCapturedEvents(t, events1)
	replica1.RequireEventuallyStoredEvents(t, events1)

	replica2.RequireEventuallyCapturedEvents(t, events1)
	replica2.RequireEventuallyStoredEvents(t, events1)

	// Broadcaster events via replica2, but with no peers yet, so expect that
	// replica1 captures and stores just it's originally broadcasted events,
	// and not the newly broadcasted events via replica2.
	events2 := replica2.BroadcastRandom(t, 1)

	replica1.RequireEventuallyCapturedEvents(t, events1)
	replica1.RequireEventuallyStoredEvents(t, events1)

	events := append(events1, events2...)
	replica2.RequireEventuallyCapturedEvents(t, events)
	replica2.RequireEventuallyStoredEvents(t, events)

	// Add replica1 as peer of replica2, and expect that both replicas
	// eventually capture and store all events.
	// replica2.addPeer(t, replica1)

	// replica1.requireEventuallyCapturedEvents(t, events)
	// replica1.requireEventuallyStoredEvents(t, events)

	// replica2.requireEventuallyCapturedEvents(t, events)
	// replica2.requireEventuallyStoredEvents(t, events)
}

func TestReplica_BroadcastStore_ReplicaSet(t *testing.T) {
	rs := crdttest.NewTestReplicaSet(t, 3)
	events := rs.BroadcastRandom(t, 5)
	rs.RequireEventuallyCapturedEvents(t, events)
	rs.RequireEventuallyStoredEvents(t, events)

	if visualize {
		for i, replica := range rs.Replicas() {
			replica.Visualize(os.Stdout, fmt.Sprintf("replica%d", i+1))
		}
	}
}
