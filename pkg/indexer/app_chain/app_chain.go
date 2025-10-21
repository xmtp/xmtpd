// Package appchain implements the AppChain Indexer.
// It's responsible for indexing the GroupMessageBroadcaster and IdentityUpdateBroadcaster contracts.
package appchain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain/contracts"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	rpcstreamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	// The app chain can't have a lag from the highest block.
	lagFromHighestBlock = 0

	// XMTP app chain is based on Arbitrum Orbit, which manages a maximum of 1.5M gas/block.
	// Given the average size of identity updates and group messages, 50 is a reasonable number.
	defaultLogsPerBlock = 50
)

var ErrInitializingAppChain = errors.New("initializing app chain")

// An AppChain has a GroupMessageBroadcaster and IdentityUpdateBroadcaster contract.
type AppChain struct {
	ctx                       context.Context
	cancel                    context.CancelFunc
	wg                        sync.WaitGroup
	rpcClient                 *ethclient.Client
	wsClient                  *ethclient.Client
	logger                    *zap.Logger
	streamer                  c.ILogStreamer
	groupMessageBroadcaster   *contracts.GroupMessageBroadcaster
	identityUpdateBroadcaster *contracts.IdentityUpdateBroadcaster
	chainID                   int64
}

func NewAppChain(
	ctxwc context.Context,
	logger *zap.Logger,
	cfg config.AppChainOptions,
	db *sql.DB,
	validationService mlsvalidate.MLSValidationService,
) (*AppChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	chainLogger := logger.Named(utils.AppChainIndexerLoggerName).
		With(utils.ChainIDField(cfg.ChainID))

	rpcClient, err := blockchain.NewRPCClient(
		ctxwc,
		cfg.RPCURL,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, fmt.Errorf("%v: %w", ErrInitializingAppChain, err)
	}

	querier := queries.New(db)

	groupMessageBroadcaster, err := contracts.NewGroupMessageBroadcaster(
		ctxwc,
		rpcClient,
		querier,
		chainLogger,
		common.HexToAddress(cfg.GroupMessageBroadcasterAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, fmt.Errorf("%v: %w", ErrInitializingAppChain, err)
	}

	groupMessageLatestBlockNumber, groupMessageLatestBlockHash := groupMessageBroadcaster.GetLatestBlock()

	identityUpdateBroadcaster, err := contracts.NewIdentityUpdateBroadcaster(
		ctxwc,
		rpcClient,
		db,
		chainLogger,
		validationService,
		common.HexToAddress(cfg.IdentityUpdateBroadcasterAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, fmt.Errorf("%v: %w", ErrInitializingAppChain, err)
	}

	identityUpdateLatestBlockNumber, identityUpdateLatestBlockHash := identityUpdateBroadcaster.GetLatestBlock()

	wsClient, err := blockchain.NewWebsocketClient(
		ctxwc,
		cfg.WssURL,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, fmt.Errorf("%v: %w", ErrInitializingAppChain, err)
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		ctxwc,
		rpcClient,
		wsClient,
		chainLogger,
		rpcstreamer.WithLagFromHighestBlock(lagFromHighestBlock),
		rpcstreamer.WithContractConfig(
			&rpcstreamer.ContractConfig{
				ID:                   contracts.GroupMessageBroadcasterName(cfg.ChainID),
				FromBlockNumber:      groupMessageLatestBlockNumber,
				FromBlockHash:        groupMessageLatestBlockHash,
				Address:              groupMessageBroadcaster.Address(),
				Topics:               groupMessageBroadcaster.Topics(),
				MaxDisconnectTime:    cfg.MaxChainDisconnectTime,
				ExpectedLogsPerBlock: defaultLogsPerBlock,
			},
		),
		rpcstreamer.WithContractConfig(
			&rpcstreamer.ContractConfig{
				ID:                   contracts.IdentityUpdateBroadcasterName(cfg.ChainID),
				FromBlockNumber:      identityUpdateLatestBlockNumber,
				FromBlockHash:        identityUpdateLatestBlockHash,
				Address:              identityUpdateBroadcaster.Address(),
				Topics:               identityUpdateBroadcaster.Topics(),
				MaxDisconnectTime:    cfg.MaxChainDisconnectTime,
				ExpectedLogsPerBlock: defaultLogsPerBlock,
			},
		),
		rpcstreamer.WithBackfillBlockPageSize(cfg.BackfillBlockPageSize),
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		wsClient.Close()
		return nil, fmt.Errorf("%v: %w", ErrInitializingAppChain, err)
	}

	return &AppChain{
		ctx:                       ctxwc,
		cancel:                    cancel,
		rpcClient:                 rpcClient,
		wsClient:                  wsClient,
		logger:                    chainLogger,
		streamer:                  streamer,
		chainID:                   cfg.ChainID,
		groupMessageBroadcaster:   groupMessageBroadcaster,
		identityUpdateBroadcaster: identityUpdateBroadcaster,
	}, nil
}

func (a *AppChain) Start() error {
	err := a.streamer.Start()
	if err != nil {
		return err
	}

	tracing.GoPanicWrap(
		a.ctx,
		&a.wg,
		"indexer-group-message-broadcaster",
		func(ctx context.Context) {
			c.IndexLogs(
				ctx,
				a.GroupMessageBroadcasterEventChannel(),
				a.groupMessageBroadcaster,
			)
		})

	tracing.GoPanicWrap(
		a.ctx,
		&a.wg,
		"indexer-identity-update-broadcaster",
		func(ctx context.Context) {
			c.IndexLogs(
				ctx,
				a.IdentityUpdateBroadcasterEventChannel(),
				a.identityUpdateBroadcaster,
			)
		})

	return nil
}

func (a *AppChain) Stop() {
	a.logger.Debug("stopping")

	if a.streamer != nil {
		a.streamer.Stop()
	}

	if a.rpcClient != nil {
		a.rpcClient.Close()
	}

	if a.wsClient != nil {
		a.wsClient.Close()
	}

	a.cancel()
	a.wg.Wait()

	a.logger.Debug("stopped")
}

func (a *AppChain) GroupMessageBroadcasterEventChannel() <-chan types.Log {
	return a.streamer.GetEventChannel(contracts.GroupMessageBroadcasterName(a.chainID))
}

func (a *AppChain) IdentityUpdateBroadcasterEventChannel() <-chan types.Log {
	return a.streamer.GetEventChannel(contracts.IdentityUpdateBroadcasterName(a.chainID))
}
