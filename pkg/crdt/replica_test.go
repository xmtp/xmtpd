package crdt_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crdt "github.com/xmtp/xmtpd/pkg/crdt"
	chanbroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/chan"
	memstore "github.com/xmtp/xmtpd/pkg/crdt/stores/mem"
	chansyncer "github.com/xmtp/xmtpd/pkg/crdt/syncers/chan"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestCRDT_NewClose(t *testing.T) {
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

func TestCRDT_BroadcastStore(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	var events []*types.Event
	var eventsLock sync.RWMutex
	crdt, err := crdt.NewReplica(ctx, log, store, bc, syncer, func(ev *types.Event) {
		eventsLock.Lock()
		defer eventsLock.Unlock()
		events = append(events, ev)
	})
	require.NoError(t, err)
	defer crdt.Close()

	ev, err := types.NewEvent([]byte("payload"), nil)
	require.NoError(t, err)
	err = bc.Broadcast(ev)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		eventsLock.RLock()
		defer eventsLock.RUnlock()
		return len(events) == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, []*types.Event{ev}, events)

	assert.Eventually(t, func() bool {
		return len(store.Events()) == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, map[string]*types.Event{
		ev.Cid.String(): ev,
	}, store.Events())
}
