// Package indexer implements the Indexer.
// It's responsible for coordinating the AppChain and SettlementChain indexers.
// It can be extended to index other chains in the future.
package indexer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/xmtp/xmtpd/pkg/config"
	appchain "github.com/xmtp/xmtpd/pkg/indexer/app_chain"
	settlementchain "github.com/xmtp/xmtpd/pkg/indexer/settlement_chain"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

type IndexerConfig struct {
	ctx               context.Context
	logger            *zap.Logger
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

func WithLogger(logger *zap.Logger) IndexerOption {
	return func(cfg *IndexerConfig) {
		cfg.logger = logger
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
	logger          *zap.Logger
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	appChain        *appchain.AppChain
	settlementChain *settlementchain.SettlementChain
}

func NewIndexer(opts ...IndexerOption) (*Indexer, error) {
	cfg := &IndexerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.ctx == nil {
		return nil, errors.New("indexer: context is required")
	}
	if cfg.logger == nil {
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
	indexerLogger := cfg.logger.Named(utils.IndexerLoggerName)

	appChain, err := appchain.NewAppChain(
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

	settlementChain, err := settlementchain.NewSettlementChain(
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
		logger:          indexerLogger,
		appChain:        appChain,
		settlementChain: settlementChain,
	}, nil
}

func (i *Indexer) Close() {
	i.logger.Debug("closing")

	if i.appChain != nil {
		i.appChain.Stop()
	}

	if i.settlementChain != nil {
		i.settlementChain.Stop()
	}

	i.cancel()
	i.wg.Wait()

	i.logger.Debug("closed")
}

func (i *Indexer) Start() error {
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
