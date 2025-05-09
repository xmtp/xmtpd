package prune

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"time"

	"go.uber.org/zap"
)

type Executor struct {
	ctx      context.Context
	log      *zap.Logger
	writerDB *sql.DB
}

func NewPruneExecutor(
	ctx context.Context,
	log *zap.Logger,
	writerDB *sql.DB,
) *Executor {
	return &Executor{
		ctx:      ctx,
		log:      log,
		writerDB: writerDB,
	}
}

func (e *Executor) Run() error {
	querier := queries.New(e.writerDB)
	totalDeletionCount := 0
	start := time.Now()

	cnt, err := querier.CountExpiredEnvelopes(e.ctx)
	if err != nil {
		return err
	}
	e.log.Info("Count of envelopes eligible for pruning", zap.Int64("count", cnt))

	//TODO(mkysel) limit the number of cycles
	for {
		rows, err := querier.DeleteExpiredEnvelopesBatch(e.ctx)
		if err != nil {
			return err
		}

		if len(rows) == 0 {
			break
		}

		totalDeletionCount = totalDeletionCount + len(rows)

		e.log.Info("Pruned expired envelopes batch", zap.Int("count", len(rows)))

		for row := range rows {
			e.log.Debug(fmt.Sprintf("Pruning expired envelopes batch row: %d", row))
		}
	}

	if totalDeletionCount == 0 {
		e.log.Info("No expired envelopes found")
	}

	e.log.Info("Done", zap.Int("pruned count", totalDeletionCount), zap.Duration("elapsed", time.Since(start)))

	return nil
}
