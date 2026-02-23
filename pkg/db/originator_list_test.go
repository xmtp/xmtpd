package db_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

func TestCachedOriginatorList_Basic(t *testing.T) {
	ctx := t.Context()
	rawDB, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(rawDB)

	// Insert envelopes for 3 originators.
	for _, nodeID := range []int32{100, 200, 300} {
		_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     nodeID,
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(256),
			},
		)
		require.NoError(t, err)
	}

	list := db.NewCachedOriginatorList(querier, 5*time.Minute, zap.NewNop())
	ids, err := list.GetOriginatorNodeIDs(ctx)
	require.NoError(t, err)
	assert.Equal(t, []uint32{100, 200, 300}, ids)
}

func TestCachedOriginatorList_Caching(t *testing.T) {
	ctx := t.Context()
	rawDB, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(rawDB)

	// Insert one originator.
	_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		querier,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(256),
		},
	)
	require.NoError(t, err)

	list := db.NewCachedOriginatorList(querier, 5*time.Minute, zap.NewNop())

	// First call populates cache.
	ids1, err := list.GetOriginatorNodeIDs(ctx)
	require.NoError(t, err)
	require.Len(t, ids1, 1)

	// Insert another originator.
	_, err = db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		querier,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(256),
		},
	)
	require.NoError(t, err)

	// Second call returns cached result (still 1 originator).
	ids2, err := list.GetOriginatorNodeIDs(ctx)
	require.NoError(t, err)
	assert.Len(t, ids2, 1, "should return cached result")
}

func TestCachedOriginatorList_CacheExpiry(t *testing.T) {
	ctx := t.Context()
	rawDB, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(rawDB)

	_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		querier,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(256),
		},
	)
	require.NoError(t, err)

	// Use very short TTL for testing.
	list := db.NewCachedOriginatorList(querier, 10*time.Millisecond, zap.NewNop())

	ids1, err := list.GetOriginatorNodeIDs(ctx)
	require.NoError(t, err)
	require.Len(t, ids1, 1)

	// Insert second originator.
	_, err = db.InsertGatewayEnvelopeWithChecksStandalone(
		ctx,
		querier,
		queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                testutils.RandomBytes(32),
			OriginatorEnvelope:   testutils.RandomBytes(256),
		},
	)
	require.NoError(t, err)

	// Wait for cache to expire and re-fetch.
	require.Eventually(t, func() bool {
		ids, err := list.GetOriginatorNodeIDs(ctx)
		return err == nil && len(ids) == 2
	}, time.Second, 5*time.Millisecond, "should re-fetch after TTL expires")
}

func TestCachedOriginatorList_Concurrent(t *testing.T) {
	ctx := t.Context()
	rawDB, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(rawDB)

	for _, nodeID := range []int32{100, 200} {
		_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			querier,
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     nodeID,
				OriginatorSequenceID: 1,
				Topic:                testutils.RandomBytes(32),
				OriginatorEnvelope:   testutils.RandomBytes(256),
			},
		)
		require.NoError(t, err)
	}

	// Very short TTL so cache expires mid-flight.
	list := db.NewCachedOriginatorList(querier, 1*time.Millisecond, zap.NewNop())

	var wg sync.WaitGroup
	for range 50 {
		wg.Go(func() {
			ids, err := list.GetOriginatorNodeIDs(ctx)
			assert.NoError(t, err) //nolint:testifylint // require not safe in goroutine
			assert.NotEmpty(t, ids)
		})
	}
	wg.Wait()
}

func TestCachedOriginatorList_Empty(t *testing.T) {
	ctx := t.Context()
	rawDB, _ := testutils.NewRawDB(t, ctx)
	querier := queries.New(rawDB)

	list := db.NewCachedOriginatorList(querier, 5*time.Minute, zap.NewNop())
	ids, err := list.GetOriginatorNodeIDs(ctx)
	require.NoError(t, err)
	assert.Empty(t, ids)
}
