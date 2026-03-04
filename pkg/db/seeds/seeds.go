// Package seeds provides functions to seed the database with test data.
package seeds

import (
	"context"
	cryptorand "crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/db"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

// Config controls the size and shape of the seeded dataset.
type Config struct {
	NumEnvelopes    uint64
	NumOriginators  uint64
	NumTopics       uint64
	NumPayers       uint64
	BlobSize        uint64
	LogInterval     uint64
	NumUsageMinutes int32
}

// DefaultConfig returns a config matching the bench suite defaults.
func DefaultConfig() Config {
	return Config{
		NumEnvelopes:    100_000,
		NumOriginators:  3,
		NumTopics:       100,
		NumPayers:       5,
		BlobSize:        500,
		LogInterval:     10_000,
		NumUsageMinutes: 40,
	}
}

// SeedResult carries metadata about the seeded dataset for use in benchmarks or follow-up queries.
type SeedResult struct {
	Topics        [][]byte
	OriginatorIDs []int32
	PayerIDs      []int32
}

// SeedEnvelopes inserts gateway envelopes, payers, and the necessary partitions.
// Envelopes are distributed round-robin across originators and topics.
func SeedEnvelopes(
	ctx context.Context,
	dbConn *sql.DB,
	cfg Config,
	logger *zap.Logger,
) (SeedResult, error) {
	q := queries.New(dbConn)

	topics := make([][]byte, cfg.NumTopics)
	for i := range cfg.NumTopics {
		topics[i] = randomBytes(32)
	}

	payerIDs := make([]int32, cfg.NumPayers)
	for i := range cfg.NumPayers {
		addr := utils.HexEncode(randomBytes(20))
		id, err := q.FindOrCreatePayer(ctx, addr)
		if err != nil {
			return SeedResult{}, fmt.Errorf("create payer %d: %w", i, err)
		}
		payerIDs[i] = id
	}

	originators := make([]int32, cfg.NumOriginators)
	for i := range cfg.NumOriginators {
		originators[i] = int32((i + 1) * 100)
	}

	// Pre-create partitions for all originators up to the expected sequence range.
	perOriginator := cfg.NumEnvelopes / cfg.NumOriginators
	for seqID := int64(0); seqID < int64(perOriginator)+db.GatewayEnvelopeBandWidth; seqID += db.GatewayEnvelopeBandWidth {
		for _, origID := range originators {
			_ = q.EnsureGatewayParts(ctx, queries.EnsureGatewayPartsParams{
				OriginatorNodeID:     origID,
				OriginatorSequenceID: seqID,
				BandWidth:            db.GatewayEnvelopeBandWidth,
			})
		}
	}

	var (
		blob   = randomBytes(int(cfg.BlobSize))
		seqIDs = make([]int64, cfg.NumOriginators)
	)

	for i := range cfg.NumEnvelopes {
		var (
			origIdx = i % cfg.NumOriginators
			origID  = originators[origIdx]
		)

		seqIDs[origIdx]++

		_, err := db.InsertGatewayEnvelopeWithChecksStandalone(
			ctx,
			q,
			queries.InsertGatewayEnvelopeParams{
				OriginatorNodeID:     origID,
				OriginatorSequenceID: seqIDs[origIdx],
				PayerID:              sql.NullInt32{Int32: payerIDs[i%cfg.NumPayers], Valid: true},
				Topic:                topics[i%cfg.NumTopics],
				Expiry:               time.Now().Add(24 * time.Hour).Unix(),
				OriginatorEnvelope:   blob,
			},
		)
		if err != nil {
			return SeedResult{}, fmt.Errorf("insert envelope %d: %w", i, err)
		}

		if cfg.LogInterval > 0 && (i+1)%cfg.LogInterval == 0 {
			if logger != nil {
				logger.Info("seeding envelopes",
					zap.Int("seeded", int(i+1)),
					zap.Int("total", int(cfg.NumEnvelopes)),
				)
			}
		}
	}

	if logger != nil {
		logger.Info("envelopes seeded", zap.Int("total", int(cfg.NumEnvelopes)))
	}

	return SeedResult{
		Topics:        topics,
		OriginatorIDs: originators,
		PayerIDs:      payerIDs,
	}, nil
}

// SeedUsage inserts unsettled usage rows for each payer × originator × minute combination
// from the provided SeedResult.
func SeedUsage(
	ctx context.Context,
	dbConn *sql.DB,
	result SeedResult,
	cfg Config,
	logger *zap.Logger,
) error {
	q := queries.New(dbConn)

	total := len(result.PayerIDs) * len(result.OriginatorIDs) * int(cfg.NumUsageMinutes)

	for _, payerID := range result.PayerIDs {
		for _, origID := range result.OriginatorIDs {
			for minute := range cfg.NumUsageMinutes {
				err := q.IncrementUnsettledUsage(
					ctx,
					queries.IncrementUnsettledUsageParams{
						PayerID:           payerID,
						OriginatorID:      origID,
						MinutesSinceEpoch: minute,
						SpendPicodollars:  1_000_000,
						SequenceID:        int64(minute),
						MessageCount:      1,
					},
				)
				if err != nil {
					return fmt.Errorf("insert usage (payer=%d originator=%d minute=%d): %w",
						payerID, origID, minute, err)
				}
			}
		}
	}

	if logger != nil {
		logger.Info("usage seeded",
			zap.Int("payers", len(result.PayerIDs)),
			zap.Int("originators", len(result.OriginatorIDs)),
			zap.Int("minutes", int(cfg.NumUsageMinutes)),
			zap.Int("total_rows", total),
		)
	}

	return nil
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = cryptorand.Read(b)
	return b
}
