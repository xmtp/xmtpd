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

func insertGatewayEnvelopes(
	t *testing.T,
	db *sql.DB,
	rows []queries.InsertGatewayEnvelopeParams,
	notifyChan ...chan bool,
) {
	println("insertGatewayEnvelopes")
	q := queries.New(db)
	for _, row := range rows {
		inserted, err := q.InsertGatewayEnvelope(context.Background(), row)
		require.Equal(t, int64(1), inserted)
		require.NoError(t, err)

		if len(notifyChan) > 0 {
			select {
			case notifyChan[0] <- true:
			default:
			}
		}
	}
}

func setup(t *testing.T) (*sql.DB, *zap.Logger, func()) {
	ctx := context.Background()
	db, _, dbCleanup := testutils.NewDB(t, ctx)
	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	return db, log, dbCleanup
}

func insertInitialRows(t *testing.T, db *sql.DB) {
	insertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{
		{
			// Auto-generated ID: 1
			OriginatorID:         1,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope1"),
		},
		{
			// Auto-generated ID: 2
			OriginatorID:         2,
			OriginatorSequenceID: 1,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope2"),
		},
	})
}

func envelopesQuery(db *sql.DB) PollableDBQuery[queries.GatewayEnvelope] {
	return func(ctx context.Context, lastSeenID int64, numRows int32) ([]queries.GatewayEnvelope, int64, error) {
		println("envelopesQuery", lastSeenID, numRows)
		envs, err := queries.New(db).
			SelectGatewayEnvelopes(ctx, queries.SelectGatewayEnvelopesParams{
				OriginatorNodeID:  NullInt32(1),
				GatewaySequenceID: NullInt64(lastSeenID),
				RowLimit:          NullInt32(numRows),
			})
		if err != nil {
			return nil, 0, err
		}
		if len(envs) > 0 {
			lastSeenID = envs[len(envs)-1].ID
		}
		println("Envs length", len(envs))
		return envs, lastSeenID, nil
	}
}

func insertAdditionalRows(t *testing.T, db *sql.DB, notifyChan ...chan bool) {
	insertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{
		{
			// Auto-generated ID: 3
			OriginatorID:         1,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope3"),
		},
		{
			// Auto-generated ID: 4
			OriginatorID:         2,
			OriginatorSequenceID: 2,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope4"),
		},
		{
			// Auto-generated ID: 5
			OriginatorID:         1,
			OriginatorSequenceID: 3,
			Topic:                []byte("topicA"),
			OriginatorEnvelope:   []byte("envelope5"),
		},
	}, notifyChan...)
}

func validateUpdates(t *testing.T, updates <-chan []queries.GatewayEnvelope, ctxCancel func()) {
	envs := <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int64(3), envs[0].ID)
	require.Equal(t, []byte("envelope3"), envs[0].OriginatorEnvelope)

	envs = <-updates
	require.Equal(t, 1, len(envs))
	require.Equal(t, int64(5), envs[0].ID)
	require.Equal(t, []byte("envelope5"), envs[0].OriginatorEnvelope)

	ctxCancel()
	_, more := <-updates
	require.False(t, more)
}

// flakyEnvelopesQuery returns a query that fails every other time
// to simulate a transient database error
func flakyEnvelopesQuery(db *sql.DB) PollableDBQuery[queries.GatewayEnvelope] {
	numQueries := 0
	query := envelopesQuery(db)
	return func(ctx context.Context, lastSeenID int64, numRows int32) ([]queries.GatewayEnvelope, int64, error) {
		numQueries++
		if numQueries%2 == 1 {
			return nil, 0, fmt.Errorf("flaky query")
		}

		return query(ctx, lastSeenID, numRows)
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
		1, // lastSeenID
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
		1, // lastSeenID
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
		1, // lastSeenID
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
