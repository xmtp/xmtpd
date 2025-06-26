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

	client, err := blockchain.NewClient(ctxwc, cfg.WssURL)
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

	streamer, err := streamer.NewRPCLogStreamer(
		ctxwc,
		client,
		chainLogger,
		streamer.WithLagFromHighestBlock(lagFromHighestBlock),
		streamer.WithContractConfig(
			&streamer.ContractConfig{
				ID:                payerRegistry.ID(cfg.ChainID),
				Contract:          payerRegistry,
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		streamer.WithContractConfig(
			&streamer.ContractConfig{
				ID:                payerReportManager.ID(cfg.ChainID),
				Contract:          payerReportManager,
				MaxDisconnectTime: cfg.MaxChainDisconnectTime,
			},
		),
		streamer.WithBackfillBlockPageSize(cfg.BackfillBlockPageSize),
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
	return s.streamer.GetEventChannel(s.payerRegistry.ID(s.chainID))
}

func (s *SettlementChain) PayerReportManagerEventChannel() <-chan types.Log {
	return s.streamer.GetEventChannel(s.payerReportManager.ID(s.chainID))
}
