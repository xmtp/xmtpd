package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abi/rateregistry"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
*
A RatesAdmin is a struct responsible for calling admin functions on the RatesRegistry contract
*
*/
type RatesAdmin struct {
	client   *ethclient.Client
	signer   TransactionSigner
	contract *rateregistry.RateRegistry
	logger   *zap.Logger
}

func NewRatesAdmin(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractsOptions config.ContractsOptions,
) (*RatesAdmin, error) {
	contract, err := rateregistry.NewRateRegistry(
		common.HexToAddress(contractsOptions.SettlementChain.RateRegistryAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &RatesAdmin{
		signer:   signer,
		client:   client,
		logger:   logger.Named("RatesAdmin"),
		contract: contract,
	}, nil
}

func (r *RatesAdmin) Contract() *rateregistry.RateRegistry {
	return r.contract
}
