package merklecrdt_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crdt "github.com/xmtp/xmtpd/pkg/merklecrdt"
	chanbroadcaster "github.com/xmtp/xmtpd/pkg/merklecrdt/broadcasters/chan"
	memstore "github.com/xmtp/xmtpd/pkg/merklecrdt/stores/mem"
	chansyncer "github.com/xmtp/xmtpd/pkg/merklecrdt/syncers/chan"
	"github.com/xmtp/xmtpd/pkg/merklecrdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

func TestMerkleCRDT_NewClose(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	crdt, err := crdt.New(ctx, log, store, bc, syncer)
	require.NoError(t, err)
	require.NotNil(t, crdt)

	err = crdt.Close()
	require.NoError(t, err)
}

func TestMerkleCRDT_BroadcastStore(t *testing.T) {
	log := test.NewLogger(t)
	ctx := context.Background()

	store := memstore.New(log)
	bc := chanbroadcaster.New(log)
	syncer := chansyncer.New(log)

	crdt, err := crdt.New(ctx, log, store, bc, syncer)
	require.NoError(t, err)
	defer crdt.Close()

	ev, err := types.NewEvent([]byte("payload"), nil)
	require.NoError(t, err)
	err = bc.Broadcast(ev)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return len(store.Events()) == 1
	}, time.Second, 10*time.Millisecond)
	require.Equal(t, map[string]*types.Event{
		ev.Cid.String(): ev,
	}, store.Events())
}
