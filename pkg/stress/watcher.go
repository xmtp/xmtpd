package stress

import (
	"context"
	"math/big"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const MAX_RETRIES = 10

type Watcher struct {
	ctx             context.Context
	logger          *zap.Logger
	ethClient       *ethclient.Client
	fromBlock       *big.Int
	watchedContract common.Address
}

func NewWatcher(
	ctx context.Context,
	logger *zap.Logger,
	rpcURL string,
	watchedContract common.Address,
) (*Watcher, error) {
	ethClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	block, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		ethClient.Close()
		return nil, err
	}

	return &Watcher{
		ctx:             ctx,
		logger:          logger.Named("chain-watcher"),
		ethClient:       ethClient,
		fromBlock:       block.Number(),
		watchedContract: watchedContract,
	}, nil
}

func (w *Watcher) Listen() error {
	defer w.ethClient.Close()

	ctxwc, cancel := signal.NotifyContext(w.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	newDone, logCh, err := w.makeSubChannel(ctxwc)
	if err != nil {
		w.logger.Error("failed to subscribe and process new logs.")
		return err
	}

	processingDone := w.processLogs(ctxwc, logCh)

	select {
	case <-ctxwc.Done():
		w.logger.Info("received shutdown signal")
	case <-newDone:
		w.logger.Info("subscription ended")
	case <-processingDone:
		w.logger.Info("log processing ended")
	}

	return nil
}

func (w *Watcher) setupFilterQuery(fromBlock, toBlock *big.Int) ethereum.FilterQuery {
	w.logger.Info("setting up filter query")
	return ethereum.FilterQuery{
		Addresses: []common.Address{w.watchedContract},
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Topics:    [][]common.Hash{},
	}
}

func (w *Watcher) makeSubChannel(
	ctx context.Context,
) (<-chan struct{}, <-chan types.Log, error) {
	query := w.setupFilterQuery(w.fromBlock, nil)
	logCh := make(chan types.Log)
	done := make(chan struct{})

	sub, err := w.ethClient.SubscribeFilterLogs(ctx, query, logCh)
	if err != nil {
		w.logger.Error(
			"unexpected error while creating subscription",
			zap.String("err", err.Error()),
		)
		return nil, nil, err
	}

	go func() {
		defer close(done)
		defer close(logCh)
		defer sub.Unsubscribe()

		w.logger.Info("subscription created")

		for {
			select {
			case err := <-sub.Err():
				if err != nil {
					w.logger.Error("subscription error", zap.String("err", err.Error()))
					sub.Unsubscribe()

					success := false
					for try := range MAX_RETRIES {
						sub, err = w.ethClient.SubscribeFilterLogs(ctx, query, logCh)
						if err == nil {
							w.logger.Info("subscription successfully recreated.")
							success = true

							break
						}

						time.Sleep(time.Second * time.Duration(try))
					}

					if !success {
						w.logger.Error(
							"failed to recreate subscription after retries: shutting down watcher.",
						)
						return
					}
				}

			case <-ctx.Done():
				w.logger.Debug("shutting down subscription")
				return
			}
		}
	}()

	return done, logCh, nil
}

func (w *Watcher) processLogs(
	ctx context.Context,
	newLog <-chan types.Log,
) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		w.logger.Info("starting log processing")

		for {
			select {
			case log := <-newLog:
				w.logger.Info(
					"received log",
					zap.Uint64("block", log.BlockNumber),
					zap.String("address", log.Address.Hex()),
					zap.String("txHash", log.TxHash.Hex()),
				)

			case <-ctx.Done():
				w.logger.Info("context cancelled, stopping log processing")
				return
			}
		}
	}()

	return done
}
