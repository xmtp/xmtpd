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
	"github.com/xmtp/xmtpd/pkg/indexer/settlement_chain/contracts"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"go.uber.org/zap"
)

const (
	// TODO: Modify after introducing changes to rpc log streamer.
	lagFromHighestBlock = 0
)

type SettlementChain struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	client        *ethclient.Client
	log           *zap.Logger
	streamer      *blockchain.RpcLogStreamer
	payerRegistry *contracts.PayerRegistry
	chainID       int
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
	)
	if err != nil {
		cancel()
		client.Close()
		return nil, err
	}

	payerRegistryLatestBlockNumber, _ := payerRegistry.GetLatestBlock()

	streamer, err := blockchain.NewRpcLogStreamer(
		ctxwc,
		client,
		log,
		cfg.ChainID,
		blockchain.WithLagFromHighestBlock(lagFromHighestBlock),
		blockchain.WithContractConfig(
			blockchain.ContractConfig{
				ID:                contracts.PayerRegistryName(cfg.ChainID),
				FromBlock:         payerRegistryLatestBlockNumber,
				ContractAddress:   payerRegistry.Address(),
				Topics:            payerRegistry.Topics(),
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

	return &SettlementChain{
		ctx:           ctxwc,
		cancel:        cancel,
		client:        client,
		log:           chainLogger,
		streamer:      streamer,
		chainID:       cfg.ChainID,
		payerRegistry: payerRegistry,
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
				s.streamer.Client(),
				s.PayerRegistryEventChannel(),
				s.PayerRegistryReorgChannel(),
				s.payerRegistry,
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

func (s *SettlementChain) PayerRegistryReorgChannel() chan uint64 {
	return s.streamer.GetReorgChannel(contracts.PayerRegistryName(s.chainID))
}
