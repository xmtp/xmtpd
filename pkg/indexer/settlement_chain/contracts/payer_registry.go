package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	pr "github.com/xmtp/xmtpd/pkg/abi/payerregistry"
	"github.com/xmtp/xmtpd/pkg/db/queries"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
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
	querier *queries.Queries,
	logger *zap.Logger,
	address common.Address,
	chainID int,
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
		querier,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := payerRegistryTopics()
	if err != nil {
		return nil, err
	}

	logger = logger.Named("payer-registry").
		With(zap.String("contractAddress", address.Hex()))

	payerRegistryStorer, err := NewPayerRegistryStorer(querier, logger, contract)
	if err != nil {
		return nil, err
	}

	reorgHandler := c.NewChainReorgHandler(ctx, client, querier)

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

func PayerRegistryName(chainID int) string {
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
