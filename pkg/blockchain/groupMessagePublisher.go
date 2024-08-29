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
type GroupMessagePublisher struct {
	signer   TransactionSigner
	contract *abis.GroupMessages
	logger   *zap.Logger
}

func NewGroupMessagePublisher(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
) (*GroupMessagePublisher, error) {
	contract, err := abis.NewGroupMessages(
		common.HexToAddress(contractOptions.MessagesContractAddress),
		client,
	)

	if err != nil {
		return nil, err
	}

	return &GroupMessagePublisher{
		signer: signer,
		logger: logger.Named("GroupMessagePublisher").
			With(zap.String("contractAddress", contractOptions.MessagesContractAddress)),
		contract: contract,
	}, nil
}

func (g *GroupMessagePublisher) Publish(
	ctx context.Context,
	groupID [32]byte,
	message []byte,
) error {
	_, err := g.contract.AddMessage(&bind.TransactOpts{
		Context: ctx,
		From:    g.signer.FromAddress(),
		Signer:  g.signer.SignerFunc(),
	}, groupID, message)

	return err
}
