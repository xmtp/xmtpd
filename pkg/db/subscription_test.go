package db_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
)

var topicA = topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).Bytes()

func setup(t *testing.T) (*db.Handler, *zap.Logger) {
	ctx := context.Background()
	store, _ := testutils.NewDB(t, ctx)
	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	return store, log
}

func insertInitialRows(t *testing.T, store *db.Handler) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope1"),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope2"),
		},
	})
}

func envelopesQuery(
	store *db.Handler,
) db.PollableDBQuery[queries.GatewayEnvelopesView, db.VectorClockRecord] {
	return func(ctx context.Context, lastSeen db.VectorClockRecord, numRows int32) ([]queries.GatewayEnvelopesView, db.VectorClockRecord, error) {
		envs, err := store.Query().
			SelectGatewayEnvelopesByOriginators(ctx, *db.SetVectorClockByOriginators(&queries.SelectGatewayEnvelopesByOriginatorsParams{
				OriginatorNodeIds: []int32{100},
				RowLimit:          numRows,
			}, lastSeen))
		if err != nil {
			return nil, lastSeen, err
		}
		for _, env := range envs {
			lastSeen[uint32(env.OriginatorNodeID)] = uint64(env.OriginatorSequenceID)
		}
		return db.TransformRowsByOriginator(envs), lastSeen, nil
	}
}

func insertAdditionalRows(t *testing.T, store *db.Handler, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 2,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope3"),
		},
		{
			OriginatorNodeID:     200,
			OriginatorSequenceID: 2,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope4"),
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope5"),
		},
	}, notifyChan...)
}

func validateUpdates(
	t *testing.T,
	updates <-chan []queries.GatewayEnvelopesView,
	ctxCancel func(),
) {
	envs := <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int32(100), envs[0].OriginatorNodeID)
	require.Equal(t, int64(2), envs[0].OriginatorSequenceID)
	require.Equal(t, []byte("envelope3"), envs[0].OriginatorEnvelope)

	envs = <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int32(100), envs[0].OriginatorNodeID)
	require.Equal(t, int64(3), envs[0].OriginatorSequenceID)
	require.Equal(t, []byte("envelope5"), envs[0].OriginatorEnvelope)

	ctxCancel()
	_, more := <-updates
	require.False(t, more)
}

// flakyEnvelopesQuery returns a query that fails every other time
// to simulate a transient database error
func flakyEnvelopesQuery(
	store *db.Handler,
) db.PollableDBQuery[queries.GatewayEnvelopesView, db.VectorClockRecord] {
	numQueries := 0
	query := envelopesQuery(store)
	return func(ctx context.Context, lastSeen db.VectorClockRecord, numRows int32) ([]queries.GatewayEnvelopesView, db.VectorClockRecord, error) {
		numQueries++
		if numQueries%2 == 1 {
			return nil, lastSeen, errors.New("flaky query")
		}

		return query(ctx, lastSeen, numRows)
	}
}

func createGatewayEnvelopes(t *testing.T, n int) []queries.InsertGatewayEnvelopeParams {
	t.Helper()

	envelopes := make([]queries.InsertGatewayEnvelopeParams, n)
	for i := range n {
		e := queries.InsertGatewayEnvelopeParams{
			OriginatorNodeID: 100,
			Topic: topic.NewTopic(topic.TopicKindGroupMessagesV1, []byte("topicA")).
				Bytes(),
			OriginatorSequenceID: int64(1 + i),
			OriginatorEnvelope:   fmt.Appendf([]byte{}, "envelope%v", i+1),
		}

		envelopes[i] = e
	}

	return envelopes
}

func TestChunkSizeSubscription(t *testing.T) {
	var (
		store, logger = setup(t)
		count         = 100
		chunkSize     = 7
		envelopes     = createGatewayEnvelopes(t, count)

		interval = 10 * time.Millisecond
	)

	// Insert envelopes in the DB
	testutils.InsertGatewayEnvelopes(t, store, envelopes)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	sub := db.NewDBSubscription(
		ctx,
		logger,
		envelopesQuery(store),
		db.VectorClockRecord{},
		db.PollingOptions{Interval: interval, NumRows: int32(chunkSize)},
	)

	updates, err := sub.Start()
	require.NoError(t, err)

	var (
		retrieved  []queries.GatewayEnvelopesView
		chunkCount = 0
	)
loop:
	for {
		select {
		case chunk, ok := <-updates:
			if !ok {
				break loop
			}

			require.LessOrEqual(t, len(chunk), chunkSize)

			chunkCount++
			retrieved = append(retrieved, chunk...)

			// Not the foolproof way to signal completion.
			if len(retrieved) == len(envelopes) {
				cancel()
			}

		case <-ctx.Done():
			break loop
		}
	}

	require.Len(t, retrieved, count)

	for i := range len(envelopes) {
		e := envelopes[i]
		record := retrieved[i]

		require.Equal(t, e.OriginatorNodeID, record.OriginatorNodeID)
		require.Equal(t, e.OriginatorSequenceID, record.OriginatorSequenceID)
		require.Equal(t, e.Topic, record.Topic)
		require.Equal(t, e.OriginatorEnvelope, record.OriginatorEnvelope)
	}
}

func TestIntervalSubscription(t *testing.T) {
	store, log := setup(t)

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := db.NewDBSubscription(
		ctx,
		log,
		envelopesQuery(store),
		db.VectorClockRecord{100: 1},
		db.PollingOptions{
			Interval: 100 * time.Millisecond,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, updates, ctxCancel)
}

func TestNotifiedSubscription(t *testing.T) {
	deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	store, log := setup(t)

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(deadline)
	notifyChan := make(chan bool)
	subscription := db.NewDBSubscription(
		ctx,
		log,
		envelopesQuery(store),
		db.VectorClockRecord{100: 1},
		db.PollingOptions{
			Notifier: notifyChan,
			Interval: 30 * time.Second,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, store, notifyChan)
	validateUpdates(t, updates, ctxCancel)
}

func TestTemporaryDBError(t *testing.T) {
	store, log := setup(t)

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := db.NewDBSubscription(
		ctx,
		log,
		flakyEnvelopesQuery(store),
		db.VectorClockRecord{100: 1},
		db.PollingOptions{
			Interval: 100 * time.Millisecond,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, store)
	validateUpdates(t, updates, ctxCancel)
}

func TestSubscriptionDeliversContiguousSequencesPerOriginator(t *testing.T) {
	store, log := setup(t)

	// Insert baseline seq=1 (already seen), and then seq=2 + seq=3
	// but with gateway_time out of order:
	// seq=3 has earlier gateway_time than seq=2.
	now := time.Now()
	seq3Earlier := now.Add(-10 * time.Second)
	seq2Later := now.Add(-5 * time.Second)

	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope1"),
			GatewayTime:          now,
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 2,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope2"),
			GatewayTime:          seq2Later, // later
		},
		{
			OriginatorNodeID:     100,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope3"),
			GatewayTime:          seq3Earlier, // earlier (should NOT come first)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start polling at seq=1, 1 row per poll so we can observe ordering across batches.
	sub := db.NewDBSubscription(
		ctx,
		log,
		envelopesQuery(store),
		db.VectorClockRecord{100: 1},
		db.PollingOptions{
			Interval: 10 * time.Millisecond,
			NumRows:  1,
		},
	)

	updates, err := sub.Start()
	require.NoError(t, err)

	// Expect seq=2 then seq=3 (strictly increasing, no gaps).
	// This is the correct behavior and must hold for vector-clock paging to be safe.

	first := <-updates
	require.Len(t, first, 1)
	require.Equal(t, int32(100), first[0].OriginatorNodeID)
	require.Equal(t, int64(2), first[0].OriginatorSequenceID,
		"subscription must deliver the next sequence ID (2) first to avoid cursor skipping")
	require.Equal(t, []byte("envelope2"), first[0].OriginatorEnvelope)

	second := <-updates
	require.Len(t, second, 1)
	require.Equal(t, int32(100), second[0].OriginatorNodeID)
	require.Equal(t, int64(3), second[0].OriginatorSequenceID,
		"subscription must deliver sequence IDs contiguously (3 after 2)")
	require.Equal(t, []byte("envelope3"), second[0].OriginatorEnvelope)
}
