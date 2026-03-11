// Package observe provides a wrapper around the Observer service used for E2E tests.
package observe

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Observer struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) *Observer {
	return &Observer{
		logger: logger,
	}
}

type EnvelopeStats struct {
	TotalEnvelopes   int64
	OriginatorNodeID int32
	LatestSequenceID int64
}

type PayerUsageStats struct {
	PayerAddress        string
	TotalSpendPicoDolls int64
	MessageCount        int64
}

type VectorClockEntry struct {
	OriginatorNodeID     int32
	OriginatorSequenceID int64
}

func (o *Observer) GetEnvelopeCount(ctx context.Context, connStr string) (int64, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	var count int64
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM gateway_envelopes_meta",
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count envelopes: %w", err)
	}

	return count, nil
}

func (o *Observer) GetVectorClock(ctx context.Context, connStr string) ([]VectorClockEntry, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	rows, err := db.QueryContext(
		ctx,
		"SELECT originator_node_id, originator_sequence_id FROM gateway_envelopes_latest ORDER BY originator_node_id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query vector clock: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var entries []VectorClockEntry
	for rows.Next() {
		var e VectorClockEntry
		if err := rows.Scan(&e.OriginatorNodeID, &e.OriginatorSequenceID); err != nil {
			return nil, fmt.Errorf("failed to scan vector clock entry: %w", err)
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

func (o *Observer) GetPayerReportCount(ctx context.Context, connStr string) (int64, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	var count int64
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM payer_reports",
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count payer reports: %w", err)
	}

	return count, nil
}

func (o *Observer) GetUnsettledUsage(
	ctx context.Context,
	connStr string,
) ([]PayerUsageStats, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	rows, err := db.QueryContext(ctx, `
		SELECT p.address,
			COALESCE(SUM(u.spend_picodollars), 0) AS total_spend,
			COALESCE(SUM(u.message_count), 0) AS message_count
		FROM unsettled_usage u
		JOIN payers p ON p.id = u.payer_id
		GROUP BY p.address
		ORDER BY total_spend DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query unsettled usage: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var stats []PayerUsageStats
	for rows.Next() {
		var s PayerUsageStats
		if err := rows.Scan(&s.PayerAddress, &s.TotalSpendPicoDolls, &s.MessageCount); err != nil {
			return nil, fmt.Errorf("failed to scan payer usage: %w", err)
		}
		stats = append(stats, s)
	}

	return stats, rows.Err()
}

type PayerReportStatusCounts struct {
	Total               int64
	AttestationPending  int64
	AttestationApproved int64
	AttestationRejected int64
	SubmissionPending   int64
	SubmissionSubmitted int64
	SubmissionSettled   int64
	SubmissionRejected  int64
}

func (o *Observer) GetPayerReportStatusCounts(
	ctx context.Context,
	connStr string,
) (*PayerReportStatusCounts, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	counts := &PayerReportStatusCounts{}

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payer_reports").Scan(&counts.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to count payer reports: %w", err)
	}

	rows, err := db.QueryContext(ctx, `
		SELECT attestation_status, COUNT(*)
		FROM payer_reports
		GROUP BY attestation_status
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query attestation status: %w", err)
	}
	for rows.Next() {
		var status int16
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("failed to scan attestation status: %w", err)
		}
		switch status {
		case 0:
			counts.AttestationPending = count
		case 1:
			counts.AttestationApproved = count
		case 2:
			counts.AttestationRejected = count
		}
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = db.QueryContext(ctx, `
		SELECT submission_status, COUNT(*)
		FROM payer_reports
		GROUP BY submission_status
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query submission status: %w", err)
	}
	for rows.Next() {
		var status int16
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("failed to scan submission status: %w", err)
		}
		switch status {
		case 0:
			counts.SubmissionPending = count
		case 1:
			counts.SubmissionSubmitted = count
		case 2:
			counts.SubmissionSettled = count
		case 3:
			counts.SubmissionRejected = count
		}
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

// SettledPayerReport holds the key fields from a settled payer report needed
// for claiming from the DistributionManager.
type SettledPayerReport struct {
	OriginatorNodeID     int32
	SubmittedReportIndex int32
}

// GetSettledPayerReports returns all payer reports with submission_status = 2 (settled)
// that have a non-null submitted_report_index.
func (o *Observer) GetSettledPayerReports(
	ctx context.Context,
	connStr string,
) ([]SettledPayerReport, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	rows, err := db.QueryContext(ctx, `
		SELECT originator_node_id, submitted_report_index
		FROM payer_reports
		WHERE submission_status = 2 AND submitted_report_index IS NOT NULL
		ORDER BY originator_node_id, submitted_report_index
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query settled payer reports: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var reports []SettledPayerReport
	for rows.Next() {
		var r SettledPayerReport
		if err := rows.Scan(&r.OriginatorNodeID, &r.SubmittedReportIndex); err != nil {
			return nil, fmt.Errorf("failed to scan settled payer report: %w", err)
		}
		reports = append(reports, r)
	}

	return reports, rows.Err()
}

func (o *Observer) WaitForPayerReports(
	ctx context.Context,
	connStr string,
	checkFn func(*PayerReportStatusCounts) bool,
	description string,
) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		counts, err := o.GetPayerReportStatusCounts(ctx, connStr)
		if err != nil {
			o.logger.Warn("error checking payer report status", zap.Error(err))
		} else if checkFn(counts) {
			o.logger.Info("payer report condition met",
				zap.String("condition", description),
				zap.Int64("total", counts.Total),
				zap.Int64("attestation_approved", counts.AttestationApproved),
				zap.Int64("submission_submitted", counts.SubmissionSubmitted),
				zap.Int64("submission_settled", counts.SubmissionSettled),
			)
			return nil
		}

		select {
		case <-ctx.Done():
			// Include last known counts in the error for debugging
			lastCounts := ""
			if counts != nil {
				lastCounts = fmt.Sprintf(
					" (last: total=%d, att_approved=%d, sub_submitted=%d, sub_settled=%d)",
					counts.Total, counts.AttestationApproved,
					counts.SubmissionSubmitted, counts.SubmissionSettled,
				)
			}
			return fmt.Errorf(
				"timed out waiting for payer reports (%s)%s: %w",
				description, lastCounts, ctx.Err(),
			)
		case <-ticker.C:
		}
	}
}

func (o *Observer) GetStagedEnvelopeCount(ctx context.Context, connStr string) (int64, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	var count int64
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM staged_originator_envelopes",
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count staged envelopes: %w", err)
	}

	return count, nil
}

func (o *Observer) GetNodeInfo(ctx context.Context, connStr string) (nodeID int32, err error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to db: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = db.QueryRowContext(ctx,
		"SELECT node_id FROM node_info LIMIT 1",
	).Scan(&nodeID)
	if err != nil {
		return 0, fmt.Errorf("failed to get node info: %w", err)
	}

	return nodeID, nil
}

func (o *Observer) WaitForEnvelopes(ctx context.Context, connStr string, minCount int64) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		count, err := o.GetEnvelopeCount(ctx, connStr)
		if err != nil {
			o.logger.Warn("error checking envelope count", zap.Error(err))
		} else if count >= minCount {
			o.logger.Info("envelope count threshold met",
				zap.Int64("count", count),
				zap.Int64("min", minCount),
			)
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf(
				"timed out waiting for envelopes (want >= %d): %w",
				minCount,
				ctx.Err(),
			)
		case <-ticker.C:
		}
	}
}
