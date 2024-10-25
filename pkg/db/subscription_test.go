package db_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/topic"
	"go.uber.org/zap"
)

var topicA = topic.NewTopic(topic.TOPIC_KIND_GROUP_MESSAGES_V1, []byte("topicA")).Bytes()

func setup(t *testing.T) (*sql.DB, *zap.Logger, func()) {
	ctx := context.Background()
	store, _, storeCleanup := testutils.NewDB(t, ctx)
	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	return store, log, storeCleanup
}

func insertInitialRows(t *testing.T, store *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope1"),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope2"),
		},
	})
}

func envelopesQuery(store *sql.DB) db.PollableDBQuery[queries.GatewayEnvelope, db.VectorClock] {
	return func(ctx context.Context, lastSeen db.VectorClock, numRows int32) ([]queries.GatewayEnvelope, db.VectorClock, error) {
		envs, err := queries.New(store).
			SelectGatewayEnvelopes(ctx, *db.SetVectorClock(&queries.SelectGatewayEnvelopesParams{
				OriginatorNodeIds: []int32{1},
				RowLimit:          numRows,
			}, lastSeen))
		if err != nil {
			return nil, lastSeen, err
		}
		for _, env := range envs {
			lastSeen[uint32(env.OriginatorNodeID)] = uint64(env.OriginatorSequenceID)
		}
		return envs, lastSeen, nil
	}
}

func insertAdditionalRows(t *testing.T, store *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, store, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope3"),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope4"),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                topicA,
			OriginatorEnvelope:   []byte("envelope5"),
		},
	}, notifyChan...)
}

func validateUpdates(t *testing.T, updates <-chan []queries.GatewayEnvelope, ctxCancel func()) {
	envs := <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int32(1), envs[0].OriginatorNodeID)
	require.Equal(t, int64(2), envs[0].OriginatorSequenceID)
	require.Equal(t, []byte("envelope3"), envs[0].OriginatorEnvelope)

	envs = <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int32(1), envs[0].OriginatorNodeID)
	require.Equal(t, int64(3), envs[0].OriginatorSequenceID)
	require.Equal(t, []byte("envelope5"), envs[0].OriginatorEnvelope)

	ctxCancel()
	_, more := <-updates
	require.False(t, more)
}

// flakyEnvelopesQuery returns a query that fails every other time
// to simulate a transient database error
func flakyEnvelopesQuery(
	store *sql.DB,
) db.PollableDBQuery[queries.GatewayEnvelope, db.VectorClock] {
	numQueries := 0
	query := envelopesQuery(store)
	return func(ctx context.Context, lastSeen db.VectorClock, numRows int32) ([]queries.GatewayEnvelope, db.VectorClock, error) {
		numQueries++
		if numQueries%2 == 1 {
			return nil, lastSeen, fmt.Errorf("flaky query")
		}

		return query(ctx, lastSeen, numRows)
	}
}

func TestIntervalSubscription(t *testing.T) {
	store, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := db.NewDBSubscription(
		ctx,
		log,
		envelopesQuery(store),
		db.VectorClock{1: 1},
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
	store, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	notifyChan := make(chan bool)
	subscription := db.NewDBSubscription(
		ctx,
		log,
		envelopesQuery(store),
		db.VectorClock{1: 1},
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
	store, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, store)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := db.NewDBSubscription(
		ctx,
		log,
		flakyEnvelopesQuery(store),
		db.VectorClock{1: 1},
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
