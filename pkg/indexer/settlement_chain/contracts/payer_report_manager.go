package contracts

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	p "github.com/xmtp/xmtpd/pkg/abi/payerreportmanager"
	"github.com/xmtp/xmtpd/pkg/db"
	c "github.com/xmtp/xmtpd/pkg/indexer/common"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

const payerReportManagerName = "payer-report-manager"

var payerReportManagerTopic = []string{
	"PayerReportSubmitted",
	"PayerReportSubsetSettled",
}

type PayerReportManager struct {
	address common.Address
	topics  []common.Hash
	logger  *zap.Logger
	c.IBlockTracker
	c.IReorgHandler
	c.ILogStorer
}

var _ c.IContract = &PayerReportManager{}

func NewPayerReportManager(
	ctx context.Context,
	client *ethclient.Client,
	db *db.Handler,
	logger *zap.Logger,
	address common.Address,
	chainID int64,
	startBlock uint64,
) (*PayerReportManager, error) {
	contract, err := payerReportManagerContract(address, client)
	if err != nil {
		return nil, err
	}

	payerReportManagerTracker, err := c.NewBlockTracker(
		ctx,
		client,
		address,
		db,
		startBlock,
	)
	if err != nil {
		return nil, err
	}

	topics, err := payerReportManagerTopics()
	if err != nil {
		return nil, err
	}

	logger = logger.Named(utils.PayerReportManagerContractLoggerName).
		With(utils.ContractAddressField(address.Hex()))

	payerReportManagerStorer, err := NewPayerReportManagerStorer(db, logger, contract)
	if err != nil {
		return nil, err
	}

	reorgHandler := NewPayerReportManagerReorgHandler(logger)

	return &PayerReportManager{
		address:       address,
		topics:        topics,
		logger:        logger,
		IBlockTracker: payerReportManagerTracker,
		IReorgHandler: reorgHandler,
		ILogStorer:    payerReportManagerStorer,
	}, nil
}

func (p *PayerReportManager) Address() common.Address {
	return p.address
}

func (p *PayerReportManager) Topics() []common.Hash {
	return p.topics
}

func (p *PayerReportManager) Logger() *zap.Logger {
	return p.logger
}

func payerReportManagerContract(
	address common.Address,
	client *ethclient.Client,
) (*p.PayerReportManager, error) {
	return p.NewPayerReportManager(
		address,
		client,
	)
}

func PayerReportManagerName(chainID int64) string {
	return fmt.Sprintf("%s-%v", payerReportManagerName, chainID)
}

func payerReportManagerTopics() ([]common.Hash, error) {
	abi, err := p.PayerReportManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	topics := make([]common.Hash, len(payerReportManagerTopic))
	for i, topic := range payerReportManagerTopic {
		topics[i], err = utils.GetEventTopic(abi, topic)
		if err != nil {
			return nil, err
		}
	}

	return topics, nil
}
