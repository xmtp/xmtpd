package indexerpoc

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

// Indexer coordinates multiple indexing tasks.
type Indexer struct {
	ctx     context.Context
	log     *zap.Logger
	wg      sync.WaitGroup
	sources map[int64]Source
	tasks   []*task

	// Storage is temporary, to be substituted with BlockTracker.
	storage     Storage
	batchSize   uint64
	concurrency int
}

// NewIndexer creates a new indexer manager.
func NewIndexer(
	ctx context.Context,
	log *zap.Logger,
	storage Storage,
	pollInterval time.Duration,
	batchSize uint64,
	concurrency int,
	networks []*Network,
) (*Indexer, error) {
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
	}

	manager := &Indexer{
		ctx:         ctx,
		log:         log.Named("indexer"),
		sources:     make(map[int64]Source),
		storage:     storage,
		batchSize:   batchSize,
		concurrency: concurrency,
	}

	if err := manager.AddNetworks(networks); err != nil {
		return nil, fmt.Errorf("failed to add networks: %w", err)
	}

	return manager, nil
}

func (i *Indexer) AddNetworks(networks []*Network) error {
	if len(networks) == 0 {
		return fmt.Errorf("no networks to add")
	}

	for _, network := range networks {
		if _, exists := i.sources[network.ChainID]; exists {
			i.log.Info("Network already added, skipping",
				zap.String("name", network.Name),
				zap.Int64("chainID", network.ChainID),
			)
			continue
		}

		source, err := NewGethSource(network, i.log)
		if err != nil {
			return fmt.Errorf("creating source for network %s: %w", network.Name, err)
		}

		i.sources[network.ChainID] = source
		i.log.Info("Added network",
			zap.String("name", network.Name),
			zap.Int64("chainID", network.ChainID),
			zap.String("rpcURL", network.RpcURL),
		)
	}

	return nil
}

// AddContracts adds multiple contracts to be indexed
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

		source, exists := i.sources[c.GetChainID()]
		if !exists {
			i.log.Error("Network not configured for chain",
				zap.Int64("chainID", c.GetChainID()),
			)
			return fmt.Errorf("network with chainID %d not configured", c.GetChainID())
		}

		task, err := getOrCreateTask(
			i.ctx,
			source,
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
	for _, task := range i.tasks {
		i.wg.Add(1)
		task := task

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
