package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	"github.com/xmtp/xmtpd/pkg/indexer/storer"
	"github.com/xmtp/xmtpd/pkg/mlsvalidate"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	lagFromHighestBlock = 0
)

// An AppChain has a GroupMessageBroadcaster and IdentityUpdateBroadcaster contract.
type AppChain struct {
	ctx                    context.Context
	client                 *ethclient.Client
	cancel                 context.CancelFunc
	wg                     sync.WaitGroup
	log                    *zap.Logger
	streamer               *blockchain.RpcLogStreamer
	chainID                int
	reorgHandler           ChainReorgHandler
	messagesTracker        *BlockTracker
	identityUpdatesTracker *BlockTracker
}

func NewAppChain(
	ctxwc context.Context,
	log *zap.Logger,
	cfg config.ContractsOptions,
	db *sql.DB,
) (*AppChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	client, err := blockchain.NewClient(ctxwc, cfg.AppChain.RpcURL)
	if err != nil {
		cancel()
		return nil, err
	}

	querier := queries.New(db)

	// TODO(borja): Move this to NewGroupMessageBroadcasterContract struct.
	messagesTopic, err := groupMessageBroadcasterTopic()
	if err != nil {
		cancel()
		return nil, err
	}

	// TODO(borja): Move this to NewGroupMessageBroadcasterContract struct.
	messagesTracker, err := NewBlockTracker(
		ctxwc,
		cfg.AppChain.GroupMessageBroadcasterAddress,
		querier,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	messagesLatestBlockNumber, _ := messagesTracker.GetLatestBlock()

	// TODO(borja): Move this to NewIdentityUpdateBroadcasterContract struct.
	identityUpdatesTopic, err := identityUpdateBroadcasterTopic()
	if err != nil {
		cancel()
		return nil, err
	}

	// TODO(borja): Move this to NewIdentityUpdateBroadcasterContract struct.
	identityUpdatesTracker, err := NewBlockTracker(
		ctxwc,
		cfg.AppChain.IdentityUpdateBroadcasterAddress,
		querier,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	identityLatestBlockNumber, _ := identityUpdatesTracker.GetLatestBlock()

	streamer := blockchain.NewRpcLogStreamer(
		ctxwc,
		client,
		log,
		blockchain.WithLagFromHighestBlock(lagFromHighestBlock),
		// TODO(borja): NewGroupMessageBroadcasterContractConfig function.
		blockchain.WithContractConfig(
			groupMessageBroadcasterName(cfg.AppChain.ChainID),
			messagesLatestBlockNumber,
			common.HexToAddress(cfg.AppChain.GroupMessageBroadcasterAddress),
			[]common.Hash{messagesTopic},
			cfg.AppChain.MaxChainDisconnectTime,
		),
		// TODO(borja): NewIdentityUpdateBroadcasterContractConfig function.
		blockchain.WithContractConfig(
			identityUpdateBroadcasterName(cfg.AppChain.ChainID),
			identityLatestBlockNumber,
			common.HexToAddress(cfg.AppChain.IdentityUpdateBroadcasterAddress),
			[]common.Hash{identityUpdatesTopic},
			cfg.AppChain.MaxChainDisconnectTime,
		),
	)

	reorgHandler := NewChainReorgHandler(ctxwc, streamer.Client(), querier)

	return &AppChain{
		ctx:                    ctxwc,
		cancel:                 cancel,
		client:                 client,
		log:                    log.Named(appChainStreamerName(cfg.AppChain.ChainID)),
		streamer:               streamer,
		chainID:                cfg.AppChain.ChainID,
		reorgHandler:           reorgHandler,
		messagesTracker:        messagesTracker,
		identityUpdatesTracker: identityUpdatesTracker,
	}, nil
}

func (s *AppChain) Start(db *sql.DB, validationService mlsvalidate.MLSValidationService) {
	s.streamer.Start()
	s.indexGroupMessageBroadcasterLogs(s.client, s.messagesTracker, db)
	s.indexIdentityUpdateBroadcasterLogs(s.client, s.identityUpdatesTracker, validationService, db)
}

func (s *AppChain) Stop() {
	s.streamer.Stop()
	s.cancel()
}

func (s *AppChain) indexGroupMessageBroadcasterLogs(
	client *ethclient.Client,
	blockTracker *BlockTracker,
	db *sql.DB,
) error {
	contractAddress := s.GetGroupMessageBroadcasterContractAddress()

	messagesContract, err := groupMessageBroadcasterContract(
		contractAddress.Hex(),
		client,
	)
	if err != nil {
		return err
	}

	querier := queries.New(db)

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-messages",
		func(ctx context.Context) {
			indexingLogger := s.log.Named("group-message-broadcaster").
				With(zap.String("contractAddress", contractAddress.Hex()))

			indexLogs(
				ctx,
				s.streamer.Client(),
				s.GetGroupMessageBroadcasterEventChannel(),
				s.GetGroupMessageBroadcasterReorgChannel(),
				indexingLogger,
				storer.NewGroupMessageStorer(querier, indexingLogger, messagesContract),
				blockTracker,
				s.reorgHandler,
				contractAddress.Hex(),
			)
		})

	return nil
}

func (s *AppChain) indexIdentityUpdateBroadcasterLogs(
	client *ethclient.Client,
	blockTracker *BlockTracker,
	validationService mlsvalidate.MLSValidationService,
	db *sql.DB,
) error {
	contractAddress := s.GetIdentityUpdateBroadcasterContractAddress()

	identityUpdatesContract, err := identityUpdateBroadcasterContract(
		contractAddress.Hex(),
		client,
	)
	if err != nil {
		return err
	}

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-identities",
		func(ctx context.Context) {
			indexingLogger := s.log.Named("identity-update-broadcaster").
				With(zap.String("contractAddress", contractAddress.Hex()))
			indexLogs(
				ctx,
				s.streamer.Client(),
				s.GetIdentityUpdateBroadcasterEventChannel(),
				s.GetIdentityUpdateBroadcasterReorgChannel(),
				indexingLogger,
				storer.NewIdentityUpdateStorer(
					db,
					indexingLogger,
					identityUpdatesContract,
					validationService,
				),
				blockTracker,
				s.reorgHandler,
				contractAddress.Hex(),
			)
		})

	return nil
}

// GroupMessageBroadcaster functions
func (s *AppChain) GetGroupMessageBroadcasterEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(groupMessageBroadcasterName(s.chainID))
}

func (s *AppChain) GetGroupMessageBroadcasterReorgChannel() chan uint64 {
	return s.streamer.GetReorgChannel(groupMessageBroadcasterName(s.chainID))
}

func (s *AppChain) GetGroupMessageBroadcasterContractAddress() common.Address {
	return s.streamer.GetContractAddress(groupMessageBroadcasterName(s.chainID))
}

func groupMessageBroadcasterContract(
	address string,
	client *ethclient.Client,
) (*gm.GroupMessageBroadcaster, error) {
	return gm.NewGroupMessageBroadcaster(
		common.HexToAddress(address),
		client,
	)
}

func groupMessageBroadcasterAddress(chainID int) string {
	return fmt.Sprintf("groupMessageBroadcaster-%v", chainID)
}

func groupMessageBroadcasterName(chainID int) string {
	return fmt.Sprintf("groupMessageBroadcaster-%v", chainID)
}

func groupMessageBroadcasterTopic() (common.Hash, error) {
	abi, err := gm.GroupMessageBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "MessageSent")
}

// IdentityUpdateBroadcaster functions
func (s *AppChain) GetIdentityUpdateBroadcasterEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(identityUpdateBroadcasterName(s.chainID))
}

func (s *AppChain) GetIdentityUpdateBroadcasterReorgChannel() chan uint64 {
	return s.streamer.GetReorgChannel(identityUpdateBroadcasterName(s.chainID))
}

func (s *AppChain) GetIdentityUpdateBroadcasterContractAddress() common.Address {
	return s.streamer.GetContractAddress(identityUpdateBroadcasterName(s.chainID))
}

func identityUpdateBroadcasterContract(
	address string,
	client *ethclient.Client,
) (*iu.IdentityUpdateBroadcaster, error) {
	return iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(address),
		client,
	)
}

func identityUpdateBroadcasterName(chainID int) string {
	return fmt.Sprintf("identityUpdateBroadcaster-%v", chainID)
}

func identityUpdateBroadcasterTopic() (common.Hash, error) {
	abi, err := iu.IdentityUpdateBroadcasterMetaData.GetAbi()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.GetEventTopic(abi, "IdentityUpdateCreated")
}

// AppChainStreamer utility functions
func appChainStreamerName(chainID int) string {
	return fmt.Sprintf("app-chain-%v", chainID)
}
