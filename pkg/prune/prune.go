// Package prune implements the DB prune executor.
package prune

import (
	"context"
	"database/sql"

	"github.com/xmtp/xmtpd/pkg/config"
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
	err := e.PruneRows()
	if err != nil {
		return err
	}

	err = e.DropPrunablePartitions()
	if err != nil {
		return err
	}

	return nil
}
