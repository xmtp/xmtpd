// Package contracts implements the payer registry and payer registry manager contracts.
// The Solidity implementations are in https://github.com/xmtp/smart-contracts/tree/main/src/settlement-chain.
package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/db"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/ledger"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const payerRegistryName = "payer-registry"

var payerRegistryTopic = []string{
	"Deposit",
	"WithdrawalRequested",
	"WithdrawalCancelled",
	"UsageSettled",
}

type PayerRegistry struct {
	address common.Address
	topics  []common.Hash
	logger  *zap.Logger
	c.IBlockTracker
	c.IReorgHandler
	c.ILogStorer
}

var _ c.IContract = &PayerRegistry{}

func NewPayerRegistry(
	ctx context.Context,
	client *ethclient.Client,
	db *db.Handler,
	logger *zap.Logger,
	address common.Address,
	chainID int64,
	startBlock uint64,
) (*PayerRegistry, error) {
	contract, err := payerRegistryContract(address, client)
	if err != nil {
		return nil, err
	}

	payerRegistryTracker, err := c.NewBlockTracker(
		ctx,
		client,
		address,
		db,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := payerRegistryTopics()
	if err != nil {
		return nil, err
	}

	logger = logger.Named(utils.PayerRegistryContractLoggerName).
		With(utils.ContractAddressField(address.Hex()))

	payerRegistryStorer, err := NewPayerRegistryStorer(
		logger,
		contract,
		ledger.NewLedger(logger, db),
	)
	if err != nil {
		return nil, err
	}

	reorgHandler := NewPayerRegistryReorgHandler(logger)

	return &PayerRegistry{
		address:       address,
		topics:        topics,
		logger:        logger,
		IBlockTracker: payerRegistryTracker,
		IReorgHandler: reorgHandler,
		ILogStorer:    payerRegistryStorer,
	}, nil
}

func (pr *PayerRegistry) Address() common.Address {
	return pr.address
}

func (pr *PayerRegistry) Topics() []common.Hash {
	return pr.topics
}

func (pr *PayerRegistry) Logger() *zap.Logger {
	return pr.logger
}

func payerRegistryContract(
	address common.Address,
	client *ethclient.Client,
) (*pr.PayerRegistry, error) {
	return pr.NewPayerRegistry(
		address,
		client,
	)
}

func PayerRegistryName(chainID int64) string {
	return fmt.Sprintf("%s-%v", payerRegistryName, chainID)
}

func payerRegistryTopics() ([]common.Hash, error) {
	abi, err := pr.PayerRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	topics := make([]common.Hash, len(payerRegistryTopic))
	for i, topic := range payerRegistryTopic {
		topics[i], err = utils.GetEventTopic(abi, topic)
		if err != nil {
			return nil, err
		}
	}

	return topics, nil
}
