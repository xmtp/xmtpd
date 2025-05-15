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
	"github.com/xmtp/xmtpd/pkg/indexer/app_chain/storer"
	rh "github.com/xmtp/xmtpd/pkg/indexer/reorg_handler"
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
	chainID                   int
	reorgHandler              rh.ChainReorgHandler
	groupMessageBroadcaster   *GroupMessageBroadcaster
	identityUpdateBroadcaster *IdentityUpdateBroadcaster
}

func NewAppChain(
	ctxwc context.Context,
	log *zap.Logger,
	cfg config.AppChainOptions,
	db *sql.DB,
) (*AppChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	chainLogger := log.Named("app-chain").
		With(zap.Int("chainID", cfg.ChainID))

	client, err := blockchain.NewClient(ctxwc, cfg.RpcURL)
	if err != nil {
		cancel()
		return nil, err
	}

	querier := queries.New(db)

	groupMessageBroadcaster, err := NewGroupMessageBroadcaster(
		ctxwc,
		client,
		querier,
		cfg.GroupMessageBroadcasterAddress,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	groupMessageLatestBlockNumber, _ := groupMessageBroadcaster.blockTracker.GetLatestBlock()

	identityUpdateBroadcaster, err := NewIdentityUpdateBroadcaster(
		ctxwc,
		client,
		querier,
		cfg.IdentityUpdateBroadcasterAddress,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	identityUpdateLatestBlockNumber, _ := identityUpdateBroadcaster.blockTracker.GetLatestBlock()

	streamer := blockchain.NewRpcLogStreamer(
		ctxwc,
		client,
		log,
		blockchain.WithLagFromHighestBlock(lagFromHighestBlock),
		blockchain.WithContractConfig(
			groupMessageBroadcasterName(cfg.ChainID),
			groupMessageLatestBlockNumber,
			common.HexToAddress(cfg.GroupMessageBroadcasterAddress),
			groupMessageBroadcaster.Topics(),
			cfg.MaxChainDisconnectTime,
		),
		blockchain.WithContractConfig(
			identityUpdateBroadcasterName(cfg.ChainID),
			identityUpdateLatestBlockNumber,
			identityUpdateBroadcaster.Address(),
			identityUpdateBroadcaster.Topics(),
			cfg.MaxChainDisconnectTime,
		),
	)

	reorgHandler := rh.NewChainReorgHandler(ctxwc, streamer.Client(), querier)

	return &AppChain{
		ctx:                       ctxwc,
		cancel:                    cancel,
		client:                    client,
		log:                       chainLogger,
		streamer:                  streamer,
		chainID:                   cfg.ChainID,
		reorgHandler:              reorgHandler,
		groupMessageBroadcaster:   groupMessageBroadcaster,
		identityUpdateBroadcaster: identityUpdateBroadcaster,
	}, nil
}

func (s *AppChain) Start(db *sql.DB, validationService mlsvalidate.MLSValidationService) {
	s.streamer.Start()
	s.indexGroupMessageBroadcasterLogs(s.groupMessageBroadcaster, db)
	s.indexIdentityUpdateBroadcasterLogs(s.identityUpdateBroadcaster, validationService, db)
}

func (s *AppChain) Stop() {
	s.streamer.Stop()
	s.cancel()
}

func (s *AppChain) indexGroupMessageBroadcasterLogs(
	broadcaster *GroupMessageBroadcaster,
	db *sql.DB,
) error {
	contractAddress := broadcaster.Address()

	querier := queries.New(db)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-group-message-broadcaster",
		func(ctx context.Context) {
			logger := s.log.Named("group-message-broadcaster").
				With(zap.String("contractAddress", contractAddress.Hex()))

			indexLogs(
				ctx,
				s.streamer.Client(),
				s.GroupMessageBroadcasterEventChannel(),
				s.GroupMessageBroadcasterReorgChannel(),
				logger,
				storer.NewGroupMessageStorer(querier, logger, broadcaster.contract),
				broadcaster.blockTracker,
				s.reorgHandler,
				contractAddress.Hex(),
			)
		})

	return nil
}

func (s *AppChain) indexIdentityUpdateBroadcasterLogs(
	broadcaster *IdentityUpdateBroadcaster,
	validationService mlsvalidate.MLSValidationService,
	db *sql.DB,
) error {
	contractAddress := broadcaster.Address()

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-identity-update-broadcaster",
		func(ctx context.Context) {
			logger := s.log.Named("identity-update-broadcaster").
				With(zap.String("contractAddress", contractAddress.Hex()))

			indexLogs(
				ctx,
				s.streamer.Client(),
				s.IdentityUpdateBroadcasterEventChannel(),
				s.IdentityUpdateBroadcasterReorgChannel(),
				logger,
				storer.NewIdentityUpdateStorer(
					db,
					logger,
					broadcaster.contract,
					validationService,
				),
				broadcaster.blockTracker,
				s.reorgHandler,
				contractAddress.Hex(),
			)
		})

	return nil
}

func (s *AppChain) GroupMessageBroadcasterEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(groupMessageBroadcasterName(s.chainID))
}

func (s *AppChain) GroupMessageBroadcasterReorgChannel() chan uint64 {
	return s.streamer.GetReorgChannel(groupMessageBroadcasterName(s.chainID))
}

func (s *AppChain) IdentityUpdateBroadcasterEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(identityUpdateBroadcasterName(s.chainID))
}

func (s *AppChain) IdentityUpdateBroadcasterReorgChannel() chan uint64 {
	return s.streamer.GetReorgChannel(identityUpdateBroadcasterName(s.chainID))
}
