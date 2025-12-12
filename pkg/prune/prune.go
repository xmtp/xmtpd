// Package prune implements the DB prune executor.
package prune

import (
	"context"
	"database/sql"
	"time"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/xmtp/xmtpd/pkg/db/queries"

	"go.uber.org/zap"
)

type Executor struct {
	ctx      context.Context
	logger   *zap.Logger
	writerDB *sql.DB
	config   *config.PruneConfig
}

func NewPruneExecutor(
	ctx context.Context,
	logger *zap.Logger,
	writerDB *sql.DB,
	config *config.PruneConfig,
) *Executor {
	if config.BatchSize <= 0 {
		logger.Panic("batch size must be greater than zero")
	}

	if config.MaxCycles <= 0 {
		logger.Panic("max cycles must be greater than zero")
	}

	return &Executor{
		ctx:      ctx,
		logger:   logger,
		writerDB: writerDB,
		config:   config,
	}
}

func (e *Executor) Run() error {
	var (
		querier        = queries.New(e.writerDB)
		start          = time.Now()
		envelopesCount int64
		migratedCount  int64
		err            error
	)

	envelopesCount, err = querier.CountExpiredEnvelopes(e.ctx)
	if err != nil {
		return err
	}

	migratedCount, err = querier.CountExpiredMigratedEnvelopes(e.ctx)
	if err != nil {
		return err
	}

	total := envelopesCount + migratedCount

	e.logger.Info("count of envelopes eligible for pruning", utils.CountField(total))

	if total == 0 {
		e.logger.Info("no envelopes found for pruning")
		return nil
	}

	if e.config.DryRun {
		e.logger.Info("dry run mode enabled, nothing to do")
		return nil
	}

	var (
		cyclesCompleted    = 0
		totalDeletionCount = 0
	)

	for {
		if cyclesCompleted >= e.config.MaxCycles {
			e.logger.Warn(
				"reached maximum pruning cycles",
				zap.Int("max_cycles", e.config.MaxCycles),
			)
			break
		}

		var deletedThisCycle int

		if envelopesCount > 0 {
			rows, err := querier.DeleteExpiredEnvelopesBatch(e.ctx, e.config.BatchSize)
			if err != nil {
				return err
			}

			deletedThisCycle += len(rows)
			envelopesCount -= int64(len(rows))
		}

		if migratedCount > 0 {
			rows, err := querier.DeleteExpiredMigratedEnvelopesBatch(
				e.ctx,
				e.config.BatchSize,
			)
			if err != nil {
				return err
			}

			deletedThisCycle += len(rows)
			migratedCount -= int64(len(rows))
		}

		totalDeletionCount += deletedThisCycle

		e.logger.Info("pruned expired envelopes batch", utils.CountField(int64(deletedThisCycle)))

		cyclesCompleted++

		if deletedThisCycle < int(e.config.BatchSize) {
			break
		}
	}

	if totalDeletionCount == 0 {
		e.logger.Info("no expired envelopes found")
	}

	e.logger.Info(
		"done",
		utils.CountField(int64(totalDeletionCount)),
		utils.DurationMsField(time.Since(start)),
	)

	return nil
}
