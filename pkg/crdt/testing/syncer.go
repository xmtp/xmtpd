package crdttest

import (
	"context"
	"sort"
	"testing"

	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/crdt"
	"github.com/xmtp/xmtpd/pkg/crdt/types"
	test "github.com/xmtp/xmtpd/pkg/testing"
	"github.com/xmtp/xmtpd/pkg/zap"
)

type ITestSyncer interface {
	crdt.Syncer

	AddPeer(peer interface{})
	AddStoreEvents(context.Context, []*types.Event) error
}

type TestSyncerMaker func(t *testing.T) *TestSyncer

type TestSyncer struct {
	ITestSyncer

	ctx context.Context
	log *zap.Logger
}

func NewTestSyncer(ctx context.Context, log *zap.Logger, bc ITestSyncer) *TestSyncer {
	return &TestSyncer{
		ITestSyncer: bc,

		ctx: ctx,
		log: log,
	}
}

func RunSyncerTests(t *testing.T, syncerMaker TestSyncerMaker) {
	t.Helper()

	t.Run("new close", func(t *testing.T) {
		t.Parallel()

		s := syncerMaker(t)
		err := s.Close()
		require.NoError(t, err)
	})

	t.Run("fetch from local", func(t *testing.T) {
		t.Parallel()

		s1 := syncerMaker(t)
		defer s1.Close()

		events, cids := s1.addManyRandom(t, 5)
		s1.requireFetchEqual(t, cids, events)
		s1.requireFetchEqual(t, []multihash.Multihash{events[0].Cid}, events[:1])
	})

	t.Run("fetch from peer", func(t *testing.T) {
		t.Parallel()

		s1 := syncerMaker(t)
		defer s1.Close()

		s2 := syncerMaker(t)
		defer s2.Close()
		s1.addPeer(t, s2)

		events, cids := s2.addManyRandom(t, 5)
		s1.requireFetchEqual(t, cids, events)
		s1.requireFetchEqual(t, []multihash.Multihash{events[0].Cid}, events[:1])
	})
}

func (b *TestSyncer) addPeer(t *testing.T, peer *TestSyncer) {
	t.Helper()
	b.AddPeer(peer.ITestSyncer)
}

func (s *TestSyncer) addManyRandom(t *testing.T, count int) ([]*types.Event, []multihash.Multihash) {
	events := make([]*types.Event, count)
	for i := 0; i < count; i++ {
		ev, err := types.NewEvent([]byte("event-"+test.RandomStringLower(13)), nil)
		require.NoError(t, err)
		events[i] = ev
	}
	err := s.AddStoreEvents(s.ctx, events)
	require.NoError(t, err)
	cids := make([]multihash.Multihash, len(events))
	for i, ev := range events {
		cids[i] = ev.Cid
	}
	return events, cids
}

func (s *TestSyncer) requireFetchEqual(t *testing.T, cids []multihash.Multihash, expected []*types.Event) {
	events, err := s.Fetch(s.ctx, cids)
	require.NoError(t, err)

	sort.Slice(events, func(i, j int) bool {
		return events[i].Cid.String() < events[j].Cid.String()
	})

	expected = expected[:]
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].Cid.String() < expected[j].Cid.String()
	})

	require.Equal(t, expected, events)
}
