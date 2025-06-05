package settlement_chain

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
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	streamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/indexer/settlement_chain/contracts"
	"github.com/xmtp/xmtpd/pkg/tracing"
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
	client             *ethclient.Client
	log                *zap.Logger
	streamer           c.ILogStreamer
	payerRegistry      *contracts.PayerRegistry
	payerReportManager *contracts.PayerReportManager
	chainID            int
}

func NewSettlementChain(
	ctxwc context.Context,
	log *zap.Logger,
	cfg config.SettlementChainOptions,
	db *sql.DB,
) (*SettlementChain, error) {
	ctxwc, cancel := context.WithCancel(ctxwc)

	chainLogger := log.Named("settlement-chain").
		With(zap.Int("chainID", cfg.ChainID))

	client, err := blockchain.NewClient(ctxwc, cfg.RpcURL)
	if err != nil {
		cancel()
		return nil, err
	}

	querier := queries.New(db)

	payerRegistry, err := contracts.NewPayerRegistry(
		ctxwc,
		client,
		querier,
		chainLogger,
		common.HexToAddress(cfg.PayerRegistryAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	payerRegistryLatestBlockNumber, payerRegistryLatestBlockHash := payerRegistry.GetLatestBlock()

	payerReportManager, err := contracts.NewPayerReportManager(
		ctxwc,
		client,
		querier,
		chainLogger,
		common.HexToAddress(cfg.PayerReportManagerAddress),
		cfg.ChainID,
		cfg.DeploymentBlock,
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	payerReportManagerLatestBlockNumber, payerReportManagerLatestBlockHash := payerReportManager.GetLatestBlock()

	streamer, err := streamer.NewRpcLogStreamer(
		ctxwc,
		client,
		chainLogger,
		streamer.WithLagFromHighestBlock(lagFromHighestBlock),
		streamer.WithContractConfig(
			streamer.ContractConfig{
				ID:                contracts.PayerRegistryName(cfg.ChainID),
				FromBlockNumber:   payerRegistryLatestBlockNumber,
				FromBlockHash:     payerRegistryLatestBlockHash,
				Address:           payerRegistry.Address(),
				Topics:            payerRegistry.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		streamer.WithContractConfig(
			streamer.ContractConfig{
				ID:                contracts.PayerReportManagerName(cfg.ChainID),
				FromBlockNumber:   payerReportManagerLatestBlockNumber,
				FromBlockHash:     payerReportManagerLatestBlockHash,
				Address:           payerReportManager.Address(),
				Topics:            payerReportManager.Topics(),
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		streamer.WithBackfillBlockSize(cfg.BackfillBlockSize),
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	return &SettlementChain{
		ctx:                ctxwc,
		cancel:             cancel,
		client:             client,
		log:                chainLogger,
		streamer:           streamer,
		chainID:            cfg.ChainID,
		payerRegistry:      payerRegistry,
		payerReportManager: payerReportManager,
	}, nil
}

func (s *SettlementChain) Start() {
	s.streamer.Start()

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
}

func (s *SettlementChain) Stop() {
	s.log.Debug("Stopping settlement chain")

	if s.streamer != nil {
		s.streamer.Stop()
	}

	if s.client != nil {
		s.client.Close()
	}

	s.cancel()
	s.wg.Wait()

	s.log.Debug("Settlement chain stopped")
}

func (s *SettlementChain) PayerRegistryEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(contracts.PayerRegistryName(s.chainID))
}

func (s *SettlementChain) PayerReportManagerEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(contracts.PayerReportManagerName(s.chainID))
}
