package storer

import (
	"context"
	"database/sql"
	"errors"

	"github.com/xmtp/xmtpd/pkg/db/queries"
	"go.uber.org/zap"
)

func GetVersionForAppend(
	ctx context.Context,
	querier *queries.Queries,
	logger *zap.Logger,
	originatorNodeID int32,
	sequenceID int64,
) (sql.NullInt32, error) {
	var version sql.NullInt32
	currentVersion, err := querier.GetEnvelopeVersion(ctx, queries.GetEnvelopeVersionParams{
		OriginatorNodeID:     originatorNodeID,
		OriginatorSequenceID: sequenceID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("Error getting current version", zap.Error(err))
		return version, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return version, err
	}

	if err == nil {
		if err = querier.InvalidateEnvelope(ctx, queries.InvalidateEnvelopeParams{
			OriginatorNodeID:     originatorNodeID,
			OriginatorSequenceID: sequenceID,
		}); err != nil {
			logger.Error("Error invalidating old envelope", zap.Error(err))
			return version, err
		}

		version = sql.NullInt32{
			Int32: currentVersion.Int32 + 1,
			Valid: true,
		}
	}

	return version, nil
}
