package vectorclock_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/db/vectorclock"
	"github.com/xmtp/xmtpd/pkg/testutils"
	envelopeTestUtils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
	"github.com/xmtp/xmtpd/pkg/topic"
)

func TestVectorClock(t *testing.T) {
	t.Run("setting and getting values", func(t *testing.T) {
		var (
			ts = newTestScaffold(t)
			vc = ts.vc

			seed = newRandomMap()
		)

		for k, v := range seed {
			vc.Save(k, v)
		}

		for k, v := range seed {
			require.Equal(t, v, vc.Get(k))
		}
	})
	t.Run("returned map is complete", func(t *testing.T) {
		var (
			ts = newTestScaffold(t)
			vc = ts.vc

			seed = newRandomMap()
		)

		for k, v := range seed {
			vc.Save(k, v)
		}

		read := vc.Values()
		require.Equal(t, seed, read)
	})
}

func TestVectorClock_DB(t *testing.T) {
	t.Run("read", func(t *testing.T) {
		var (
			ctx = t.Context()
			ts  = newTestScaffold(t)
			vc  = ts.vc

			seed = newSequenceIDSet()
		)

		// Create envelopes in DB.
		createEnvelopes(t, ts.db, seed)

		read, err := vc.ReadFromDB(ctx)
		require.NoError(t, err)

		require.Equal(t, len(read), len(seed))

		// Validate values are correct.
		for id, seqIDs := range seed {
			max := slices.Max(seqIDs)

			val, ok := read[id]
			require.True(t, ok)
			require.Equal(t, max, val)
		}
	})
}

func TestVectorClock_Sync(t *testing.T) {
	t.Run("force sync", func(t *testing.T) {
		var (
			ctx = t.Context()
			ts  = newTestScaffold(t)
			vc  = ts.vc

			seed = newSequenceIDSet()
		)

		// Create envelopes in DB.
		createEnvelopes(t, ts.db, seed)

		// Force sync should update our vector clock.
		err := vc.ForceSync(ctx)
		require.NoError(t, err)

		// Validate values are correct.
		read := vc.Values()
		for id, seqIDs := range seed {
			max := slices.Max(seqIDs)

			val, ok := read[id]
			require.True(t, ok)
			require.Equal(t, max, val)
		}
	})
	t.Run("force sync update", func(t *testing.T) {
		var (
			ctx = t.Context()
			ts  = newTestScaffold(t)
			vc  = ts.vc

			seed = map[uint32][]uint64{
				100: {1001, 1002, 1003, 1004, 1005},
				200: {2001, 2002, 2003, 2004, 2005},
				300: {3001, 3002, 3003, 3004, 3005},
			}
		)

		// Create envelopes in DB.
		createEnvelopes(t, ts.db, seed)

		// Force sync should update our vector clock.
		err := vc.ForceSync(ctx)
		require.NoError(t, err)

		// Verify the expected sequence ID.
		require.Equal(t, uint64(1005), vc.Get(100))

		// Create a new envelope that updates the sequence ID.
		createEnvelopes(
			t,
			ts.db,
			map[uint32][]uint64{
				100: {1010},
			})

		err = vc.ForceSync(ctx)
		require.NoError(t, err)

		// Verify vector clock is updated.
		require.Equal(t, uint64(1010), vc.Get(100))
	})
	t.Run("sync loop", func(t *testing.T) {
		var (
			ctx         = t.Context()
			syncTimeout = 500 * time.Millisecond
			ts          = newTestScaffold(
				t,
				vectorclock.WithResolveStrategy(
					vectorclock.ResolveReconcile,
				), // Default option but lets be explicit
				vectorclock.WithSyncTimeout(syncTimeout),
			)
			vc = ts.vc

			seed = map[uint32][]uint64{
				100: {1001, 1002, 1003, 1004, 1005},
				200: {2001, 2002, 2003, 2004, 2005},
				300: {3001, 3002, 3003, 3004, 3005},
			}
		)

		err := vc.Start(ctx)
		require.NoError(t, err)

		// Create envelopes in DB.
		createEnvelopes(t, ts.db, seed)

		// Wait a little longer than the sync timeout.
		time.Sleep(syncTimeout + 100*time.Millisecond)

		require.Equal(t, uint64(1005), vc.Get(100))
		require.Equal(t, uint64(2005), vc.Get(200))
		require.Equal(t, uint64(3005), vc.Get(300))

		// Update DB out of band and make sure that sync loop updates our values.
		createEnvelopes(
			t,
			ts.db,
			map[uint32][]uint64{
				100: {1100},
				200: {2100},
				300: {3100},
			})

		// Again, wait a little longer than the sync timeout.
		time.Sleep(syncTimeout + 100*time.Millisecond)

		require.Equal(t, uint64(1100), vc.Get(100))
		require.Equal(t, uint64(2100), vc.Get(200))
		require.Equal(t, uint64(3100), vc.Get(300))
	})
}

type testScaffold struct {
	db *sql.DB
	vc *vectorclock.VectorClock
}

func newTestScaffold(t *testing.T, opts ...vectorclock.ConfigOption) testScaffold {
	t.Helper()

	var (
		db, _ = testutils.NewRawDB(t, t.Context())
		log   = testutils.NewLog(t)
	)

	vectorClockReadFunc := func(ctx context.Context) (map[uint32]uint64, error) {
		el, err := queries.New(db).SelectVectorClock(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not read vector clock: %w", err)
		}

		out := make(map[uint32]uint64)
		for _, e := range el {
			out[uint32(e.OriginatorNodeID)] = uint64(e.OriginatorSequenceID)
		}

		return out, nil
	}

	vc := vectorclock.New(log, vectorClockReadFunc, opts...)

	s := testScaffold{
		db: db,
		vc: vc,
	}

	return s
}

func newRandomMap() map[uint32]uint64 {
	var (
		minKeys = 5
		keys    = minKeys + rand.Intn(10)
	)

	out := make(map[uint32]uint64)
	for range keys {
		id := rand.Uint32()
		seqID := rand.Uint64()

		out[uint32(id)] = seqID
	}

	return out
}

func newSequenceIDSet() map[uint32][]uint64 {
	var (
		minKeys   = 5
		minSeqIDs = 10

		keys = minKeys + rand.Intn(10)
	)

	out := make(map[uint32][]uint64)
	for i := range keys {
		id := (uint32(i+1) * 100)          // nodeIDs are multiples of 100 - make a realistic looking set.
		count := minSeqIDs + rand.Intn(10) // 10-20 envelopes

		used := make(map[uint64]struct{})
		s := make([]uint64, count)
		for i := range count {
			val := uint64(rand.Int63())

			// Pretty rare scenario that we generate a random value that is 0 (invalid sequenceID)
			// or we generate the same value twice (will fail due to DB constraint) so
			// have a fallback to increase the value until we get a fitting one.
			for {
				_, ok := used[val]
				if val > 0 && !ok {
					s[i] = val
					used[val] = struct{}{}
					break
				}

				// Try to find the next fitting value.
				val++
			}

			s[i] = val
		}

		out[id] = s
	}

	return out
}

func createEnvelopes(t *testing.T, sqlDB *sql.DB, values map[uint32][]uint64) {
	t.Helper()

	// NOTE: We need to use lower level DB functions to avoid import cycle between testutils <=> db <=> vectorclock packages.

	var (
		payerAddress = testutils.RandomString(42)
		querier      = queries.New(sqlDB)
	)

	payerID, err := querier.FindOrCreatePayer(t.Context(), payerAddress)
	require.NoError(t, err)

	// TODO: Does a query per-envelope instead of batches so it's not super fast, switch to batches.
	for id, seqIDs := range values {
		for _, seqID := range seqIDs {

			env := createEnvelopeParams(t, id, seqID, payerID)
			insertEnvelope(t, querier, env)
		}
	}
}

func insertEnvelope(
	t *testing.T,
	querier *queries.Queries,
	params queries.InsertGatewayEnvelopeParams,
) {
	t.Helper()

	inserted, err := querier.InsertGatewayEnvelope(t.Context(), params)
	if err == nil {
		require.Equal(t, int64(1), inserted.InsertedMetaRows)
		return
	}

	require.Contains(t, err.Error(), "no partition of relation")

	err = querier.EnsureGatewayParts(t.Context(), queries.EnsureGatewayPartsParams{
		OriginatorNodeID:     params.OriginatorNodeID,
		OriginatorSequenceID: params.OriginatorSequenceID,
		BandWidth:            1_000_000,
	})
	require.NoError(t, err)

	inserted, err = querier.InsertGatewayEnvelope(t.Context(), params)
	require.Equal(t, int64(1), inserted.InsertedMetaRows)
}

func createEnvelopeParams(
	t *testing.T,
	nodeID uint32,
	seqID uint64,
	payerID int32,
) queries.InsertGatewayEnvelopeParams {
	t.Helper()

	topic := topic.NewTopic(
		topic.TopicKindGroupMessagesV1,
		[]byte(fmt.Sprintf("generic-topic-%v", rand.Int())),
	)

	oe := testutils.Marshal(t,
		envelopeTestUtils.CreateOriginatorEnvelopeWithTopic(
			t,
			uint32(nodeID),
			uint64(seqID),
			topic.Bytes(),
		),
	)

	env := queries.InsertGatewayEnvelopeParams{
		OriginatorNodeID:     int32(nodeID),
		OriginatorSequenceID: int64(seqID),
		Topic:                topic.Bytes(),
		PayerID:              sql.NullInt32{Int32: payerID, Valid: true},
		OriginatorEnvelope:   oe,
	}

	return env
}
