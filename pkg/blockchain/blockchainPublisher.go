package blockchain

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/xmtp/xmtpd/pkg/abis"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
Can publish to the blockchain, signing messages using the provided signer
*/
type BlockchainPublisher struct {
	signer                 TransactionSigner
	client                 *ethclient.Client
	messagesContract       *abis.GroupMessages
	identityUpdateContract *abis.IdentityUpdates
	logger                 *zap.Logger
}

func NewBlockchainPublisher(
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
) (*BlockchainPublisher, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
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

	return &BlockchainPublisher{
		signer: signer,
		logger: logger.Named("GroupBlockchainPublisher").
			With(zap.String("contractAddress", contractOptions.MessagesContractAddress)),
		messagesContract:       messagesContract,
		identityUpdateContract: identityUpdateContract,
		client:                 client,
	}, nil
}

func (m *BlockchainPublisher) PublishGroupMessage(
	ctx context.Context,
	groupID [32]byte,
	message []byte,
) (*abis.GroupMessagesMessageSent, error) {
	if len(message) == 0 {
		return nil, errors.New("message is empty")
	}
	tx, err := m.messagesContract.AddMessage(&bind.TransactOpts{
		Context: ctx,
		From:    m.signer.FromAddress(),
		Signer:  m.signer.SignerFunc(),
	}, groupID, message)
	if err != nil {
		return nil, err
	}

	receipt, err := WaitForTransaction(
		ctx,
		m.logger,
		m.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if err != nil {
		return nil, err
	}

	if receipt == nil {
		return nil, errors.New("transaction receipt is nil")
	}

	return findLog(receipt, m.messagesContract.ParseMessageSent, "no message sent log found")
}

func (m *BlockchainPublisher) PublishIdentityUpdate(
	ctx context.Context,
	inboxId [32]byte,
	identityUpdate []byte,
) (*abis.IdentityUpdatesIdentityUpdateCreated, error) {
	if len(identityUpdate) == 0 {
		return nil, errors.New("identity update is empty")
	}
	tx, err := m.identityUpdateContract.AddIdentityUpdate(&bind.TransactOpts{
		Context: ctx,
		From:    m.signer.FromAddress(),
		Signer:  m.signer.SignerFunc(),
	}, inboxId, identityUpdate)
	if err != nil {
		return nil, err
	}

	receipt, err := WaitForTransaction(
		ctx,
		m.logger,
		m.client,
		2*time.Second,
		250*time.Millisecond,
		tx.Hash(),
	)
	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, errors.New("transaction receipt is nil")
	}

	return findLog(
		receipt,
		m.identityUpdateContract.ParseIdentityUpdateCreated,
		"no identity update log found",
	)
}

func findLog[T any](
	receipt *types.Receipt,
	parse func(types.Log) (*T, error),
	errorMsg string,
) (*T, error) {
	for _, logEntry := range receipt.Logs {
		if logEntry == nil {
			continue
		}
		event, err := parse(*logEntry)
		if err != nil {
			continue
		}
		return event, nil
	}

	return nil, errors.New(errorMsg)
}
