package blockchain

import (
	"context"
	"time"

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
type MessagePublisher struct {
	signer                 TransactionSigner
	client                 *ethclient.Client
	messagesContract       *abis.GroupMessages
	identityUpdateContract *abis.IdentityUpdates
	logger                 *zap.Logger
}

func NewMessagePublisher(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
) (*MessagePublisher, error) {
	messagesContract, err := abis.NewGroupMessages(
		common.HexToAddress(contractOptions.MessagesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}
	identityUpdateContract, err := abis.NewIdentityUpdates(
		common.HexToAddress(contractOptions.IdentityUpdatesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	return &MessagePublisher{
		signer: signer,
		logger: logger.Named("GroupMessagePublisher").
			With(zap.String("contractAddress", contractOptions.MessagesContractAddress)),
		messagesContract:       messagesContract,
		identityUpdateContract: identityUpdateContract,
	}, nil
}

func (m *MessagePublisher) PublishGroupMessage(
	ctx context.Context,
	groupID [32]byte,
	message []byte,
) error {
	tx, err := m.messagesContract.AddMessage(&bind.TransactOpts{
		Context: ctx,
		From:    m.signer.FromAddress(),
		Signer:  m.signer.SignerFunc(),
	}, groupID, message)
	if err != nil {
		return err
	}

	return WaitForTransaction(
		ctx,
		m.logger,
		m.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
}

func (m *MessagePublisher) PublishIdentityUpdate(
	ctx context.Context,
	inboxId [32]byte,
	identityUpdate []byte,
) error {
	tx, err := m.identityUpdateContract.AddIdentityUpdate(&bind.TransactOpts{
		Context: ctx,
		From:    m.signer.FromAddress(),
		Signer:  m.signer.SignerFunc(),
	}, inboxId, identityUpdate)
	if err != nil {
		return err
	}

	return WaitForTransaction(
		ctx,
		m.logger,
		m.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
}