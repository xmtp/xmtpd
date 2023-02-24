package crdttest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type ITestBroadcaster interface {
	crdt.Broadcaster

	AddPeer(t *testing.T, peer interface{})
}

type TestBroadcasterMaker func(t *testing.T) *TestBroadcaster

type TestBroadcaster struct {
	ITestBroadcaster

	ctx context.Context
	log *zap.Logger
}

func NewTestBroadcaster(ctx context.Context, log *zap.Logger, bc ITestBroadcaster) *TestBroadcaster {
	return &TestBroadcaster{
		ITestBroadcaster: bc,

		ctx: ctx,
		log: log,
	}
}

func RunBroadcasterTests(t *testing.T, broadcasterMaker TestBroadcasterMaker) {
	t.Helper()

	t.Run("broadcast", func(t *testing.T) {
		t.Parallel()

		bc1 := broadcasterMaker(t)
		defer bc1.Close()

		bc2 := broadcasterMaker(t)
		defer bc2.Close()
		bc1.addPeer(t, bc2)

		events := bc1.broadcastRandom(t, 1)

		ev1 := bc1.next(t)
		require.Equal(t, events, []*types.Event{ev1})

		ev2 := bc2.next(t)
		require.Equal(t, events, []*types.Event{ev2})
	})
}

func (b *TestBroadcaster) broadcastRandom(t *testing.T, count int) []*types.Event {
	events := make([]*types.Event, count)
	for i := 0; i < count; i++ {
		ev, err := types.NewEvent(newRandomEnvelopeWithRandomTopic(t), nil)
		require.NoError(t, err)

		err = b.Broadcast(b.ctx, ev)
		require.NoError(t, err)

		events[i] = ev
	}
	return events
}

func (b *TestBroadcaster) next(t *testing.T) *types.Event {
	t.Helper()
	ctx, cancel := context.WithTimeout(b.ctx, time.Second)
	defer cancel()
	ev, err := b.Next(ctx)
	require.NoError(t, err)
	require.NotNil(t, ev)
	return ev
}

func (b *TestBroadcaster) addPeer(t *testing.T, peer *TestBroadcaster) {
	t.Helper()
	b.AddPeer(t, peer.ITestBroadcaster)
}
