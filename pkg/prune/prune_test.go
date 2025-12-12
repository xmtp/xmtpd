package prune_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/payerreport"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/prune"
)

const (
	DefaultOriginatorID = 100
	DefaultExpiredCnt   = 10
	DefaultValidCnt     = 5
	DefaultSubmittedCnt = 10000
)

func setupTestData(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	expired int,
	valid int,
	submitted int,
) {
	// Insert expired envelopes
	for i := 0; i < expired; i++ {
		testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{{
			OriginatorNodeID:     DefaultOriginatorID,
			OriginatorSequenceID: int64(i + 1),
			Topic:                []byte("topic"),
			OriginatorEnvelope:   []byte("payload"),
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(-1 * time.Hour).Unix(),
		}})
	}

	// Insert non-expired envelopes
	for i := 0; i < valid; i++ {
		testutils.InsertGatewayEnvelopes(t, db, []queries.InsertGatewayEnvelopeParams{{
			OriginatorNodeID:     DefaultOriginatorID,
			OriginatorSequenceID: int64(i + expired + 1),
			Topic:                []byte("topic"),
			OriginatorEnvelope:   []byte("payload"),
			GatewayTime:          time.Now(),
			Expiry:               time.Now().Add(1 * time.Hour).Unix(),
		}})
	}

	createPrunableReport(t, ctx, db, submitted)
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

func createPrunableReport(t *testing.T, ctx context.Context, db *sql.DB, endSequence int) {
	if endSequence == 0 {
		return
	}

	q := queries.New(db)

	reportID := testutils.RandomReportID()

	_, err := q.InsertOrIgnorePayerReport(ctx, queries.InsertOrIgnorePayerReportParams{
		ID:                  reportID,
		OriginatorNodeID:    DefaultOriginatorID,
		StartSequenceID:     0,
		EndSequenceID:       int64(endSequence),
		EndMinuteSinceEpoch: 0,
		PayersMerkleRoot:    make([]byte, 0),
		ActiveNodeIds:       []int32{DefaultOriginatorID},
	})
	require.NoError(t, err)
	err = q.SetReportSubmitted(ctx, queries.SetReportSubmittedParams{
		ReportID:             reportID,
		NewStatus:            payerreport.SubmissionSubmitted,
		PrevStatus:           []int16{int16(payerreport.SubmissionPending)},
		SubmittedReportIndex: 0,
	})
	require.NoError(t, err)
}

func TestExecutor_PrunesExpired(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]
	q := queries.New(db)

	setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, DefaultSubmittedCnt)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 5,
	})

	err := exec.Run()
	assert.NoError(t, err)

	cnt, err := q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, cnt, "All expired envelopes should be deleted")

	// Ensure non-expired remain
	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta`)
	err = row.Scan(&total)
	assert.NoError(t, err)
	assert.EqualValues(t, DefaultValidCnt, total, "Only non-expired envelopes should remain")
}

func TestExecutor_DryRun_NoPrune(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, DefaultSubmittedCnt)

	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    true,
		MaxCycles: 5,
	})
	err := exec.Run()
	assert.NoError(t, err)

	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta`)
	err = row.Scan(&total)
	assert.NoError(t, err)

	assert.EqualValues(
		t,
		DefaultValidCnt+DefaultExpiredCnt,
		total,
		"All envelopes should still be present",
	)
}

func openAndHoldLock(t *testing.T, ctx context.Context, db *sql.DB) *sql.Tx {
	// Begin a transaction and lock sequence_id = 1
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	require.NoError(t, err)
	_, err = tx.ExecContext(ctx, `
		SELECT * FROM gateway_envelopes_meta
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

	setupTestData(t, ctx, db, DefaultExpiredCnt, 0, DefaultSubmittedCnt)

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
		SELECT originator_sequence_id FROM gateway_envelopes_meta
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
	q := queries.New(db)

	const KeepThisMany = 10

	setupTestData(t, ctx, db, 1000+KeepThisMany, 0, DefaultSubmittedCnt)

	// only allow for 1 cycle, which deletes at most 1000 envelopes
	exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
		DryRun:    false,
		MaxCycles: 1,
	})
	err := exec.Run()
	assert.NoError(t, err)

	cnt, err := q.CountExpiredEnvelopes(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, KeepThisMany, cnt)

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

	setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, DefaultSubmittedCnt)

	logger := testutils.NewCapturingLogger(zapcore.DebugLevel)

	exec := prune.NewPruneExecutor(
		ctx,
		logger.Logger,
		db,
		&config.PruneConfig{
			BatchSize: 1000,
			MaxCycles: 5,
		},
	)
	err := exec.Run()
	assert.NoError(t, err)

	if !logger.Contains("count of envelopes eligible for pruning") {
		t.Errorf("expected log message not found, got: %s", logger.Logs())
	}
}

func TestExecutor_CantPruneWithoutReport(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, 0)

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

	// Ensure all remain
	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta`)
	err = row.Scan(&total)
	assert.NoError(t, err)
	assert.EqualValues(
		t,
		DefaultValidCnt+DefaultExpiredCnt,
		total,
		"All envelopes should remain",
	)
}

func TestExecutor_MultipleOverlappingReportsOK(t *testing.T) {
	ctx := context.Background()
	dbs := testutils.NewDBs(t, ctx, 1)
	db := dbs[0]

	setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, 0)

	createPrunableReport(t, ctx, db, 10)
	createPrunableReport(t, ctx, db, 50)
	createPrunableReport(t, ctx, db, 100)

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

	var total int64
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta`)
	err = row.Scan(&total)
	assert.NoError(t, err)
	assert.EqualValues(t, DefaultValidCnt, total, "Valid envelopes should remain")
}

func TestExecutor_ReportStatusVariants(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name   string
		status int16
		pruned int
	}{
		{
			name:   "pending",
			status: int16(payerreport.SubmissionPending),
			pruned: 0,
		},
		{
			name:   "rejected",
			status: int16(payerreport.SubmissionRejected),
			pruned: 0,
		},
		{
			name:   "submitted",
			status: int16(payerreport.SubmissionSubmitted),
			pruned: DefaultExpiredCnt,
		},
		{
			name:   "settled",
			status: int16(payerreport.SubmissionSettled),
			pruned: DefaultExpiredCnt,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("report-status-state-%s", tt.name), func(t *testing.T) {
			dbs := testutils.NewDBs(t, ctx, 1)
			db := dbs[0]

			setupTestData(t, ctx, db, DefaultExpiredCnt, DefaultValidCnt, 0)

			q := queries.New(db)
			reportID := testutils.RandomReportID()
			_, err := q.InsertOrIgnorePayerReport(ctx, queries.InsertOrIgnorePayerReportParams{
				ID:                  reportID,
				OriginatorNodeID:    DefaultOriginatorID,
				StartSequenceID:     0,
				EndSequenceID:       DefaultSubmittedCnt,
				EndMinuteSinceEpoch: 0,
				PayersMerkleRoot:    make([]byte, 0),
				ActiveNodeIds:       []int32{DefaultOriginatorID},
			})
			require.NoError(t, err)

			t.Logf("Setting report to state %d", tt.status)

			err = q.SetReportSubmissionStatus(ctx, queries.SetReportSubmissionStatusParams{
				ReportID:   reportID,
				NewStatus:  tt.status,
				PrevStatus: []int16{int16(payerreport.SubmissionPending)},
			})
			require.NoError(t, err)

			cnt, err := q.CountExpiredEnvelopes(ctx)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.pruned, cnt)

			exec := makeTestExecutor(t, ctx, db, &config.PruneConfig{
				DryRun:    false,
				MaxCycles: 5,
			})

			err = exec.Run()
			assert.NoError(t, err)

			var total int64
			row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gateway_envelopes_meta`)
			err = row.Scan(&total)
			assert.NoError(t, err)
			assert.EqualValues(
				t,
				DefaultValidCnt+DefaultExpiredCnt-tt.pruned,
				total,
				"All envelopes should remain",
			)
		})
	}
}
