package prune_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/prune"
)

const (
	DEFAULT_EXPIRED_CNT = 10
	DEFAULT_VALID_CNT   = 5
)

func setupTestData(t *testing.T, db *sql.DB, expired int, valid int) {
	// Insert expired envelopes
	for i := 0; i < expired; i++ {
		testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV2Params{{
			OriginatorNodeID:     100,
			OriginatorSequenceID: int64(i),
			Topic:                []byte("topic"),
			OriginatorEnvelope:   []byte("payload"),
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(-1 * time.Hour).Unix(),
		}})
	}

	// Insert non-expired envelopes
	for i := 0; i < valid; i++ {
		testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeV2Params{{
			OriginatorNodeID:     100,
			OriginatorSequenceID: int64(i + expired),
			Topic:                []byte("topic"),
			OriginatorEnvelope:   []byte("payload"),
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(1 * time.Hour).Unix(),
		}})
	}
}

func makeTestExecutor(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	config *config.PruneConfig,
) *prune.Executor {
	config.BatchSize = 1000

	return prune.NewPruneExecutor(
		ctx,
		testutils.NewLog(t),
		db,
		config,
	)
}

func TestExecutor_PrunesExpired(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, db, DEFAULT_EXPIRED_CNT, DEFAULT_VALID_CNT)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 5,
	})
	err := exec.Run()
	assert.NoError(t, err)

	q := queries.New(db)
	cnt, err := q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt, "All expired envelopes should be deleted")

	// Ensure non-expired remain
	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta_v2`)
	err = row.Scan(&total)
	assert.NoError(t, err)
	assert.EqualValues(t, DEFAULT_VALID_CNT, total, "Only non-expired envelopes should remain")
}

func TestExecutor_DryRun_NoPrune(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, db, DEFAULT_EXPIRED_CNT, DEFAULT_VALID_CNT)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    true,
		MaxCycles: 5,
	})
	err := exec.Run()
	assert.NoError(t, err)

	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta_v2`)
	err = row.Scan(&total)
	assert.NoError(t, err)

	assert.EqualValues(
		t,
		DEFAULT_VALID_CNT+DEFAULT_EXPIRED_CNT,
		total,
		"All envelopes should still be present",
	)
}

func openAndHoldLock(t *testing.T, ctx context.Context, db *sql.DB) *sql.Tx {
	// Begin a transaction and lock sequence_id = 1
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `
		SELECT * FROM gateway_envelopes_meta_v2
		WHERE originator_sequence_id = 1 
		FOR UPDATE
	`)
	require.NoError(t, err)

	return tx
}

func TestExecutor_PrunesExpired_WithConcurrentLock(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	setupTestData(t, db, DEFAULT_EXPIRED_CNT, 0)

	tx := openAndHoldLock(t, ctx, db)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 5,
	})
	err := exec.Run()
	assert.NoError(t, err)

	remainingIDs := getRemainingSequenceIds(t, ctx, db)

	assert.Contains(t, remainingIDs, int64(1), "Locked row should still exist after pruning")
	assert.Len(t, remainingIDs, 1, "Only locked row should remain during lock")

	// Commit the lock transaction
	require.NoError(t, tx.Commit())

	err = exec.Run()
	assert.NoError(t, err)

	// Confirm DB is empty now
	cnt, err := q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt, "All expired envelopes should now be deleted")
}

func getRemainingSequenceIds(t *testing.T, ctx context.Context, db *sql.DB) []int64 {
	var remainingIDs []int64
	rows, err := db.QueryContext(ctx, `
		SELECT originator_sequence_id FROM gateway_envelopes_meta_v2
	`)
	require.NoError(t, err)
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var seqID int64
		require.NoError(t, rows.Scan(&seqID))
		remainingIDs = append(remainingIDs, seqID)
	}
	return remainingIDs
}

func TestExecutor_PrunesExpired_LargePayload(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	const KEEP_THIS_MANY = 10

	setupTestData(t, db, 1000+KEEP_THIS_MANY, 0)

	// only allow for 1 cycle, which deletes at most 1000 envelopes
	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})
	err := exec.Run()
	assert.NoError(t, err)

	q := queries.New(db)
	cnt, err := q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, KEEP_THIS_MANY, cnt)

	// 2nd cycle should finish off
	err = exec.Run()
	assert.NoError(t, err)
	cnt, err = q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt, "All expired envelopes should be deleted")
}

func TestExecutor_PruneCountWorks(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, db, DEFAULT_EXPIRED_CNT, DEFAULT_VALID_CNT)

	logger := testutils.NewCapturingLogger(zapcore.DebugLevel)

	exec := prune.NewPruneExecutor(
		ctx,
		logger.Logger,
		db,
		&config.PruneConfig{
			BatchSize:      1000,
			MaxCycles:      5,
			CountDeletable: true,
		},
	)
	err := exec.Run()
	assert.NoError(t, err)

	if !logger.Contains("count of envelopes eligible for pruning") {
		t.Errorf("expected log message not found, got: %s", logger.Logs())
	}
}
