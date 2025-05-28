package indexer

import (
	"context"
	"database/sql"
	"sync"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain"
	"github.com/xmtp/xmtpd/pkg/indexer/settlement_chain"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"go.uber.org/zap"
)

type Indexer struct {
	ctx             context.Context
	log             *zap.Logger
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	appChain        *app_chain.AppChain
	settlementChain *settlement_chain.SettlementChain
}

func NewIndexer(
	ctx context.Context,
	log *zap.Logger,
	db *sql.DB,
	cfg config.ContractsOptions,
	validationService mlsvalidate.MLSValidationService,
) (*Indexer, error) {
	ctx, cancel := context.WithCancel(ctx)

	indexerLogger := log.Named("indexer")

	appChain, err := app_chain.NewAppChain(
		ctx,
		indexerLogger,
		cfg.AppChain,
		db,
		validationService,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	settlementChain, err := settlement_chain.NewSettlementChain(
		ctx,
		indexerLogger,
		cfg.SettlementChain,
		db,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Indexer{
		ctx:             ctx,
		log:             indexerLogger,
		cancel:          cancel,
		appChain:        appChain,
		settlementChain: settlementChain,
	}, nil
}

func (i *Indexer) Close() {
	i.log.Debug("Closing")

	if i.appChain != nil {
		i.appChain.Stop()
	}

	if i.settlementChain != nil {
		i.settlementChain.Stop()
	}

	i.cancel()
	i.wg.Wait()

	i.log.Debug("Closed")
}

func (i *Indexer) StartIndexer() {
	i.appChain.Start()
	i.settlementChain.Start()
}
