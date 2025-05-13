package prune

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"

	"github.com/xmtp/xmtpd/pkg/db/queries"

	"go.uber.org/zap"
)

type Executor struct {
	ctx      context.Context
	log      *zap.Logger
	writerDB *sql.DB
	config   *config.PruneConfig
}

func NewPruneExecutor(
	ctx context.Context,
	log *zap.Logger,
	writerDB *sql.DB,
	config *config.PruneConfig,
) *Executor {
	return &Executor{
		ctx:      ctx,
		log:      log,
		writerDB: writerDB,
		config:   config,
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

	if e.config.DryRun {
		e.log.Info("Dry run mode enabled. Nothing to do")
		return nil
	}

	cyclesCompleted := 0

	for {
		if cyclesCompleted >= e.config.MaxCycles {
			e.log.Warn("Reached maximum pruning cycles", zap.Int("maxCycles", e.config.MaxCycles))
			break
		}

		rows, err := querier.DeleteExpiredEnvelopesBatch(e.ctx)
		if err != nil {
			return err
		}

		if len(rows) == 0 {
			break
		}

		totalDeletionCount = totalDeletionCount + len(rows)

		e.log.Info("Pruned expired envelopes batch", zap.Int("count", len(rows)))

		for _, row := range rows {
			e.log.Debug(fmt.Sprintf("Pruning expired envelopes batch row: %v", row))
		}
		cyclesCompleted++
	}

	if totalDeletionCount == 0 {
		e.log.Info("No expired envelopes found")
	}

	e.log.Info(
		"Done",
		zap.Int("pruned count", totalDeletionCount),
		zap.Duration("elapsed", time.Since(start)),
	)

	return nil
}
