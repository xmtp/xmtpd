package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

type IFundsAdmin interface{}

type fundsAdmin struct {
	client *ethclient.Client
	signer TransactionSigner
	logger *zap.Logger
}

var _ IFundsAdmin = &fundsAdmin{}

func NewFundsAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	_ config.ContractsOptions,
) (IFundsAdmin, error) {
	return &fundsAdmin{
		client: client,
		signer: signer,
		logger: logger.Named("FundsAdmin"),
	}, nil
}
