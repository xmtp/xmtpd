package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
Can publish to the blockchain, signing messages using the provided signer
*/
type IdentityUpdatePublisher struct {
	signer   TransactionSigner
	contract *abis.IdentityUpdates
	logger   *zap.Logger
}

func NewIdentityUpdatePublisher(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
) (*IdentityUpdatePublisher, error) {
	contract, err := abis.NewIdentityUpdates(
		common.HexToAddress(contractOptions.IdentityUpdatesContractAddress),
		client,
	)

	if err != nil {
		return nil, err
	}

	return &IdentityUpdatePublisher{
		signer: signer,
		logger: logger.Named("IdentityUpdatePublisher").
			With(zap.String("contractAddress", contractOptions.IdentityUpdatesContractAddress)),
		contract: contract,
	}, nil
}

func (g *IdentityUpdatePublisher) Publish(
	ctx context.Context,
	inboxId [32]byte,
	identityUpdate []byte,
) error {
	_, err := g.contract.AddIdentityUpdate(&bind.TransactOpts{
		Context: ctx,
		From:    g.signer.FromAddress(),
		Signer:  g.signer.SignerFunc(),
	}, inboxId, identityUpdate)

	return err
}
