package app_chain

import (
	"context"
	"database/sql"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain/contracts"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const (
	lagFromHighestBlock = 0
)

// An AppChain has a GroupMessageBroadcaster and IdentityUpdateBroadcaster contract.
type AppChain struct {
	ctx                       context.Context
	cancel                    context.CancelFunc
	wg                        sync.WaitGroup
	client                    *ethclient.Client
	log                       *zap.Logger
	streamer                  *blockchain.RpcLogStreamer
	groupMessageBroadcaster   *contracts.GroupMessageBroadcaster
	identityUpdateBroadcaster *contracts.IdentityUpdateBroadcaster
	chainID                   int
}

func NewAppChain(
	ctxwc context.Context,
	log *zap.Logger,
	cfg config.AppChainOptions,
	db *sql.DB,
	validationService mlsvalidate.MLSValidationService,
) (*AppChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	chainLogger := log.Named("app-chain").
		With(zap.Int("chainID", cfg.ChainID))

	client, err := blockchain.NewClient(ctxwc, cfg.RpcURL)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	querier := queries.New(db)

	groupMessageBroadcaster, err := contracts.NewGroupMessageBroadcaster(
		ctxwc,
		client,
		querier,
		chainLogger,
		common.HexToAddress(cfg.GroupMessageBroadcasterAddress),
		cfg.ChainID,
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	groupMessageLatestBlockNumber, _ := groupMessageBroadcaster.GetLatestBlock()

	identityUpdateBroadcaster, err := contracts.NewIdentityUpdateBroadcaster(
		ctxwc,
		client,
		db,
		chainLogger,
		validationService,
		common.HexToAddress(cfg.IdentityUpdateBroadcasterAddress),
		cfg.ChainID,
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	identityUpdateLatestBlockNumber, _ := identityUpdateBroadcaster.GetLatestBlock()

	streamer, err := blockchain.NewRpcLogStreamer(
		ctxwc,
		client,
		log,
		cfg.ChainID,
		blockchain.WithLagFromHighestBlock(lagFromHighestBlock),
		blockchain.WithContractConfig(
			blockchain.ContractConfig{
				ID:                contracts.GroupMessageBroadcasterName(cfg.ChainID),
				FromBlock:         groupMessageLatestBlockNumber,
				ContractAddress:   groupMessageBroadcaster.Address(),
				Topics:            groupMessageBroadcaster.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		blockchain.WithContractConfig(
			blockchain.ContractConfig{
				ID:                contracts.IdentityUpdateBroadcasterName(cfg.ChainID),
				FromBlock:         identityUpdateLatestBlockNumber,
				ContractAddress:   identityUpdateBroadcaster.Address(),
				Topics:            identityUpdateBroadcaster.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		blockchain.WithBackfillBlockSize(cfg.BackfillBlockSize),
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	return &AppChain{
		ctx:                       ctxwc,
		cancel:                    cancel,
		client:                    client,
		log:                       chainLogger,
		streamer:                  streamer,
		chainID:                   cfg.ChainID,
		groupMessageBroadcaster:   groupMessageBroadcaster,
		identityUpdateBroadcaster: identityUpdateBroadcaster,
	}, nil
}

func (a *AppChain) Start() {
	a.streamer.Start()

	tracing.GoPanicWrap(
		a.ctx,
		&a.wg,
		"indexer-group-message-broadcaster",
		func(ctx context.Context) {
			c.IndexLogs(
				ctx,
				a.streamer.Client(),
				a.GroupMessageBroadcasterEventChannel(),
				a.GroupMessageBroadcasterReorgChannel(),
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
				a.streamer.Client(),
				a.IdentityUpdateBroadcasterEventChannel(),
				a.IdentityUpdateBroadcasterReorgChannel(),
				a.identityUpdateBroadcaster,
			)
		})
}

func (a *AppChain) Stop() {
	a.log.Debug("Stopping app chain")

	if a.streamer != nil {
		a.streamer.Stop()
	}

	if a.client != nil {
		a.client.Close()
	}

	a.cancel()
	a.wg.Wait()

	a.log.Debug("App chain stopped")
}

func (a *AppChain) GroupMessageBroadcasterEventChannel() <-chan types.Log {
	return a.streamer.GetEventChannel(contracts.GroupMessageBroadcasterName(a.chainID))
}

func (a *AppChain) GroupMessageBroadcasterReorgChannel() chan uint64 {
	return a.streamer.GetReorgChannel(contracts.GroupMessageBroadcasterName(a.chainID))
}

func (a *AppChain) IdentityUpdateBroadcasterEventChannel() <-chan types.Log {
	return a.streamer.GetEventChannel(contracts.IdentityUpdateBroadcasterName(a.chainID))
}

func (a *AppChain) IdentityUpdateBroadcasterReorgChannel() chan uint64 {
	return a.streamer.GetReorgChannel(contracts.IdentityUpdateBroadcasterName(a.chainID))
}
