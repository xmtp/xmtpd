package prune

import (
	"context"
	"database/sql"
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
	return nil
}
