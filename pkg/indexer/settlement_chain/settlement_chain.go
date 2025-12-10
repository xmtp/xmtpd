// Package settlementchain implements the SettlementChain Indexer.
// It's responsible for indexing the PayerRegistry and PayerReportManager contracts.
package settlementchain

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/db"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	rpcstreamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/indexer/settlement_chain/contracts"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const (
	// TODO: Modify after introducing changes to rpc log streamer.
	lagFromHighestBlock = 0
)

type SettlementChain struct {
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
	rpcClient          *ethclient.Client
	wsClient           *ethclient.Client
	logger             *zap.Logger
	streamer           c.ILogStreamer
	payerRegistry      *contracts.PayerRegistry
	payerReportManager *contracts.PayerReportManager
	chainID            int64
}

func NewSettlementChain(
	ctxwc context.Context,
	logger *zap.Logger,
	cfg config.SettlementChainOptions,
	db *db.Handler,
) (*SettlementChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	chainLogger := logger.Named(utils.SettlementChainIndexerLoggerName).
		With(utils.SettlementChainChainIDField(cfg.ChainID))

	rpcClient, err := blockchain.NewRPCClient(
		ctxwc,
		cfg.RPCURL,
	)
	if err != nil {
		cancel()
		return nil, err
	}

	payerRegistry, err := contracts.NewPayerRegistry(
		ctxwc,
		rpcClient,
		db,
		chainLogger,
		common.HexToAddress(cfg.PayerRegistryAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, err
	}

	payerRegistryLatestBlockNumber, payerRegistryLatestBlockHash := payerRegistry.GetLatestBlock()

	payerReportManager, err := contracts.NewPayerReportManager(
		ctxwc,
		rpcClient,
		db,
		chainLogger,
		common.HexToAddress(cfg.PayerReportManagerAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, err
	}

	payerReportManagerLatestBlockNumber, payerReportManagerLatestBlockHash := payerReportManager.GetLatestBlock()

	wsClient, err := blockchain.NewWebsocketClient(
		ctxwc,
		cfg.WssURL,
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		return nil, err
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		ctxwc,
		rpcClient,
		wsClient,
		chainLogger,
		rpcstreamer.WithLagFromHighestBlock(lagFromHighestBlock),
		rpcstreamer.WithContractConfig(
			&rpcstreamer.ContractConfig{
				ID:                contracts.PayerRegistryName(cfg.ChainID),
				FromBlockNumber:   payerRegistryLatestBlockNumber,
				FromBlockHash:     payerRegistryLatestBlockHash,
				Address:           payerRegistry.Address(),
				Topics:            payerRegistry.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		rpcstreamer.WithContractConfig(
			&rpcstreamer.ContractConfig{
				ID:                contracts.PayerReportManagerName(cfg.ChainID),
				FromBlockNumber:   payerReportManagerLatestBlockNumber,
				FromBlockHash:     payerReportManagerLatestBlockHash,
				Address:           payerReportManager.Address(),
				Topics:            payerReportManager.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		rpcstreamer.WithBackfillBlockPageSize(cfg.BackfillBlockPageSize),
	)
	if err != nil {
		cancel()
		rpcClient.Close()
		wsClient.Close()
		return nil, err
	}

	return &SettlementChain{
		ctx:                ctxwc,
		cancel:             cancel,
		rpcClient:          rpcClient,
		wsClient:           wsClient,
		logger:             chainLogger,
		streamer:           streamer,
		chainID:            cfg.ChainID,
		payerRegistry:      payerRegistry,
		payerReportManager: payerReportManager,
	}, nil
}

func (s *SettlementChain) Start() error {
	err := s.streamer.Start()
	if err != nil {
		return err
	}

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-payer-registry",
		func(ctx context.Context) {
			c.IndexLogs(
				ctx,
				s.PayerRegistryEventChannel(),
				s.payerRegistry,
			)
		})

	tracing.GoPanicWrap(
		s.ctx,
		&s.wg,
		"indexer-payer-report-manager",
		func(ctx context.Context) {
			c.IndexLogs(
				ctx,
				s.PayerReportManagerEventChannel(),
				s.payerReportManager,
			)
		})

	return nil
}

func (s *SettlementChain) Stop() {
	s.logger.Debug("stopping")

	if s.streamer != nil {
		s.streamer.Stop()
	}

	if s.rpcClient != nil {
		s.rpcClient.Close()
	}

	if s.wsClient != nil {
		s.wsClient.Close()
	}

	s.cancel()
	s.wg.Wait()

	s.logger.Debug("stopped")
}

func (s *SettlementChain) PayerRegistryEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(contracts.PayerRegistryName(s.chainID))
}

func (s *SettlementChain) PayerReportManagerEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(contracts.PayerReportManagerName(s.chainID))
}
