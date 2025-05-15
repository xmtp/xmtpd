package indexer

import (
	"context"
	"database/sql"
	"sync"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"go.uber.org/zap"
)

type Indexer struct {
	ctx      context.Context
	log      *zap.Logger
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	appChain *app_chain.AppChain
}

func NewIndexer(
	ctx context.Context,
	log *zap.Logger,
	db *sql.DB,
	cfg config.ContractsOptions,
) (*Indexer, error) {
	ctx, cancel := context.WithCancel(ctx)

	indexerLogger := log.Named("indexer")

	appChain, err := app_chain.NewAppChain(
		ctx,
		indexerLogger,
		cfg.AppChain,
		db,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Indexer{
		ctx:      ctx,
		log:      indexerLogger,
		cancel:   cancel,
		appChain: appChain,
	}, nil
}

func (i *Indexer) Close() {
	i.log.Debug("Closing")

	if i.appChain != nil {
		i.appChain.Stop()
	}

	i.cancel()
	i.wg.Wait()

	i.log.Debug("Closed")
}

func (i *Indexer) StartIndexer(
	db *sql.DB,
	validationService mlsvalidate.MLSValidationService,
) error {
	i.appChain.Start(db, validationService)

	return nil
}
