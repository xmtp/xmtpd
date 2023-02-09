package crdttest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
)

type TestBroadcaster interface {
	crdt.Broadcaster

	AddPeer(peer interface{})
}

type TestBroadcasterMaker func(t *testing.T) TestBroadcaster

func RunBroadcasterTests(t *testing.T, broadcasterMaker TestBroadcasterMaker) {
	t.Helper()

	t.Run("broadcast", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		bc1 := broadcasterMaker(t)
		require.NotNil(t, bc1)
		defer bc1.Close()

		bc2 := broadcasterMaker(t)
		require.NotNil(t, bc2)
		defer bc2.Close()
		bc1.AddPeer(bc2)

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
	})
}
