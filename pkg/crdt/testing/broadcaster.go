package crdttest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	membroadcaster "github.com/xmtp/xmtpd/pkg/crdt/broadcasters/mem"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type TestBroadcaster struct {
	crdt.Broadcaster
}

func (b *TestBroadcaster) addPeer(t *testing.T, peer *TestBroadcaster) {
	switch bc := b.Broadcaster.(type) {
	case *membroadcaster.MemoryBroadcaster:
		switch peerBC := peer.Broadcaster.(type) {
		case *membroadcaster.MemoryBroadcaster:
			bc.AddPeer(peerBC)
		}
	}
}

func NewTestBroadcaster(bc crdt.Broadcaster) *TestBroadcaster {
	return &TestBroadcaster{bc}
}

type TestBroadcasterMaker func(t *testing.T) *TestBroadcaster

func TestBroadcaster_BroadcastNext(t *testing.T, broadcasterMaker TestBroadcasterMaker) {
	t.Helper()

	ctx := context.Background()

	bc1 := broadcasterMaker(t)
	require.NotNil(t, bc1)
	defer bc1.Close()

	bc2 := broadcasterMaker(t)
	require.NotNil(t, bc2)
	defer bc2.Close()
	bc1.addPeer(t, bc2)

	broadcastedEvent, err := types.NewEvent([]byte("event-"+test.RandomStringLower(13)), nil)
	require.NoError(t, err)

	err = bc1.Broadcast(broadcastedEvent)
	require.NoError(t, err)

	bc1Event, err := bc1.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, bc1Event)
	require.Equal(t, broadcastedEvent, bc1Event)

	bc2Event, err := bc2.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, bc2Event)
	require.Equal(t, broadcastedEvent, bc2Event)
}
