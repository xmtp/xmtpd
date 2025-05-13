package indexerpoc

import (
	"context"
	"fmt"
	"sync"

	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

// Indexer coordinates multiple indexing tasks.
type Indexer struct {
	ctx      context.Context
	log      *zap.Logger
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	networks map[int64]*Network
	tasks    []*task

	// Storage is temporary, to be substituted with BlockTracker.
	storage     Storage
	batchSize   uint64
	concurrency int
}

func NewIndexer(
	ctx context.Context,
	log *zap.Logger,
	storage Storage,
	contractsCfg config.ContractsOptions,
	indexerCfg config.IndexerOptions,
) (*Indexer, error) {
	if indexerCfg.BatchSize <= 0 {
		return nil, fmt.Errorf("batch size must be greater than 0")
	}

	if indexerCfg.Concurrency <= 0 {
		return nil, fmt.Errorf("concurrency must be greater than 0")
	}

	ctx, cancel := context.WithCancel(ctx)

	manager := &Indexer{
		ctx:         ctx,
		cancel:      cancel,
		log:         log.Named("indexer"),
		networks:    make(map[int64]*Network, 0),
		storage:     storage,
		batchSize:   indexerCfg.BatchSize,
		concurrency: indexerCfg.Concurrency,
	}

	if err := manager.configureNetworks(contractsCfg); err != nil {
		return nil, fmt.Errorf("failed to add networks: %w", err)
	}

	return manager, nil
}

func (i *Indexer) configureNetworks(cfg config.ContractsOptions) error {
	appChain, err := NewNetwork(
		i.ctx,
		NetworkConfig{
			Name:         "app-chain",
			ChainID:      int64(cfg.AppChain.ChainID),
			RpcURL:       cfg.AppChain.RpcURL,
			PollInterval: cfg.AppChain.PollInterval,
			// TODO: Make use of MaxDisconnectionTime.
		},
		i.log)
	if err != nil {
		return fmt.Errorf("failed to create app chain network: %w", err)
	}

	i.networks[int64(cfg.AppChain.ChainID)] = appChain

	// TODO: Add the settlement chain.

	return nil
}

// AddContracts adds multiple contracts to be indexed.
func (i *Indexer) AddContracts(contracts []Contract) error {
	if len(contracts) == 0 {
		return fmt.Errorf("no contracts to add")
	}

	for _, c := range contracts {
		if c.GetChainID() == 0 {
			return fmt.Errorf("contract %s must have a valid ChainID", c.GetName())
		}

		if c.GetStartBlock() == 0 {
			return fmt.Errorf("contract %s must have a StartBlock specified", c.GetName())
		}

		network, exists := i.networks[c.GetChainID()]
		if !exists {
			i.log.Error("Network not configured for chain",
				zap.Int64("chainID", c.GetChainID()),
			)
			return fmt.Errorf("network with chainID %d not configured", c.GetChainID())
		}

		task, err := getOrCreateTask(
			i.ctx,
			network,
			c,
			i.storage,
			i.batchSize,
			WithConcurrency(i.concurrency),
		)
		if err != nil {
			return fmt.Errorf("creating task for contract %s on chain %d: %w",
				c.GetName(), c.GetChainID(), err)
		}

		i.tasks = append(i.tasks, task)
		i.log.Info("Added contract",
			zap.String("name", c.GetName()),
			zap.String("address", c.GetAddress()),
			zap.Int64("chainID", c.GetChainID()),
			zap.Uint64("startBlock", c.GetStartBlock()),
			zap.Int("concurrency", i.concurrency),
		)
	}

	return nil
}

// Run starts all indexing tasks.
func (i *Indexer) Run() {
	for _, network := range i.networks {
		i.wg.Add(1)

		tracing.GoPanicWrap(
			i.ctx,
			&i.wg,
			fmt.Sprintf("indexer-network-%v", network.GetName()),
			func(ctx context.Context) {
				network.start(ctx)
			})

		i.log.Info("Started indexing network",
			zap.String("name", network.GetName()),
			zap.Int64("chainID", network.GetChainID()),
		)
	}

	for _, task := range i.tasks {
		i.wg.Add(1)

		tracing.GoPanicWrap(
			i.ctx,
			&i.wg,
			fmt.Sprintf("indexer-task-%v", task.contract.GetName()),
			func(ctx context.Context) {
				task.run()
			})

		i.log.Info("Started indexing task",
			zap.String("contract", task.contract.GetName()),
			zap.String("address", task.contract.GetAddress()),
			zap.Int64("chainID", task.contract.GetChainID()),
		)
	}

	i.wg.Wait()
}

func (i *Indexer) Close() {
	i.log.Debug("closing indexer")
	i.cancel()
	i.wg.Wait()
	i.log.Debug("closed indexer")
}
