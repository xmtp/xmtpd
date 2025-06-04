package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"errors"
	"sync"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain"
	"github.com/xmtp/xmtpd/pkg/indexer/settlement_chain"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"go.uber.org/zap"
)

type IndexerConfig struct {
	ctx               context.Context
	log               *zap.Logger
	dB                *sql.DB
	contractsConfig   *config.ContractsOptions
	validationService mlsvalidate.MLSValidationService
}

type IndexerOption func(*IndexerConfig)

func WithContext(ctx context.Context) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.ctx = ctx
	}
}

func WithLogger(log *zap.Logger) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.log = log
	}
}

func WithDB(db *sql.DB) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.dB = db
	}
}

func WithContractsOptions(c *config.ContractsOptions) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.contractsConfig = c
	}
}

func WithValidationService(vs mlsvalidate.MLSValidationService) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.validationService = vs
	}
}

type Indexer struct {
	ctx             context.Context
	log             *zap.Logger
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	appChain        *app_chain.AppChain
	settlementChain *settlement_chain.SettlementChain
}

func NewIndexer(opts ...IndexerOption) (*Indexer, error) {
	cfg := &IndexerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ctx == nil {
		return nil, errors.New("indexer: context is required")
	}
	if cfg.log == nil {
		return nil, errors.New("indexer: logger is required")
	}

	if cfg.dB == nil {
		return nil, errors.New("indexer: DB is required")
	}
	if cfg.validationService == nil {
		return nil, errors.New("indexer: ValidationService is required")
	}

	if cfg.contractsConfig == nil {
		return nil, errors.New("indexer: contracts config is required")
	}

	ctx, cancel := context.WithCancel(cfg.ctx)
	indexerLogger := cfg.log.Named("indexer")

	appChain, err := app_chain.NewAppChain(
		ctx,
		indexerLogger,
		cfg.contractsConfig.AppChain,
		cfg.dB,
		cfg.validationService,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	settlementChain, err := settlement_chain.NewSettlementChain(
		ctx,
		indexerLogger,
		cfg.contractsConfig.SettlementChain,
		cfg.dB,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Indexer{
		ctx:             ctx,
		cancel:          cancel,
		log:             indexerLogger,
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

func (i *Indexer) StartIndexer() error {
	err := i.appChain.Start()
	if err != nil {
		return fmt.Errorf("failed to start app chain: %w", err)
	}

	err = i.settlementChain.Start()
	if err != nil {
		return fmt.Errorf("failed to start settlement chain: %w", err)
	}

	return nil
}
