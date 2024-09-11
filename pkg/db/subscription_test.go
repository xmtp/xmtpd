package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"go.uber.org/zap"
)

func setup(t *testing.T) (*sql.DB, *zap.Logger, func()) {
	ctx := context.Background()
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	return db, log, dbCleanup
}

func insertInitialRows(t *testing.T, db *sql.DB) {
	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope1"),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope2"),
		},
	})
}

func envelopesQuery(db *sql.DB) PollableDBQuery[queries.GatewayEnvelope, VectorClock] {
	return func(ctx context.Context, lastSeen VectorClock, numRows int32) ([]queries.GatewayEnvelope, VectorClock, error) {
		envs, err := queries.New(db).
			SelectGatewayEnvelopes(ctx, *SetVectorClock(&queries.SelectGatewayEnvelopesParams{
				OriginatorNodeID: NullInt32(1),
				RowLimit:         NullInt32(numRows),
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

func insertAdditionalRows(t *testing.T, db *sql.DB, notifyChan ...chan bool) {
	testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope3"),
		},
		{
			OriginatorNodeID:     2,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope4"),
		},
		{
			OriginatorNodeID:     1,
			OriginatorSequenceID: 3,
			Topic:                []byte("topicA"),
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
func flakyEnvelopesQuery(db *sql.DB) PollableDBQuery[queries.GatewayEnvelope, VectorClock] {
	numQueries := 0
	query := envelopesQuery(db)
	return func(ctx context.Context, lastSeen VectorClock, numRows int32) ([]queries.GatewayEnvelope, VectorClock, error) {
		numQueries++
		if numQueries%2 == 1 {
			return nil, lastSeen, fmt.Errorf("flaky query")
		}

		return query(ctx, lastSeen, numRows)
	}
}

func TestIntervalSubscription(t *testing.T) {
	db, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, db)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := NewDBSubscription(
		ctx,
		log,
		envelopesQuery(db),
		VectorClock{1: 1},
		PollingOptions{
			Interval: 100 * time.Millisecond,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, updates, ctxCancel)
}

func TestNotifiedSubscription(t *testing.T) {
	db, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, db)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	notifyChan := make(chan bool)
	subscription := NewDBSubscription(
		ctx,
		log,
		envelopesQuery(db),
		VectorClock{1: 1},
		PollingOptions{
			Notifier: notifyChan,
			Interval: 30 * time.Second,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, db, notifyChan)
	validateUpdates(t, updates, ctxCancel)
}

func TestTemporaryDBError(t *testing.T) {
	db, log, cleanup := setup(t)
	defer cleanup()

	insertInitialRows(t, db)

	// Create a subscription that polls every 100ms
	ctx, ctxCancel := context.WithCancel(context.Background())
	subscription := NewDBSubscription(
		ctx,
		log,
		flakyEnvelopesQuery(db),
		VectorClock{1: 1},
		PollingOptions{
			Interval: 100 * time.Millisecond,
			NumRows:  1,
		},
	)
	updates, err := subscription.Start()
	require.NoError(t, err)

	insertAdditionalRows(t, db)
	validateUpdates(t, updates, ctxCancel)
}
