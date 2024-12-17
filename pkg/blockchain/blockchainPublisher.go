package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
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
	mutexNonce             sync.Mutex
	nonce                  *uint64
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

	nonce, err := m.fetchNonce(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := m.messagesContract.AddMessage(&bind.TransactOpts{
		Context: ctx,
		Nonce:   new(big.Int).SetUint64(nonce),
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

	nonce, err := m.fetchNonce(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := m.identityUpdateContract.AddIdentityUpdate(&bind.TransactOpts{
		Context: ctx,
		Nonce:   new(big.Int).SetUint64(nonce),
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

func (m *BlockchainPublisher) fetchNonce(ctx context.Context) (uint64, error) {
	// NOTE:since pendingNonce starts at 0, and we have to return that value exactly,
	// we can't easily use Once with unsigned integers
	if m.nonce == nil {
		m.mutexNonce.Lock()
		defer m.mutexNonce.Unlock()
		if m.nonce == nil {
			// PendingNonceAt gives the next nonce that should be used
			// if we are the first thread to initialize the nonce, we want to return PendingNonce+0
			nonce, err := m.client.PendingNonceAt(ctx, m.signer.FromAddress())
			if err != nil {
				return 0, err
			}
			m.nonce = &nonce
			m.logger.Info(fmt.Sprintf("Starting server with blockchain nonce: %d", *m.nonce))
			return *m.nonce, nil
		}
	}
	// Once the nonce has been initialized we can depend on Atomic to return the next value
	next := atomic.AddUint64(m.nonce, 1)

	pending, err := m.client.PendingNonceAt(ctx, m.signer.FromAddress())
	if err != nil {
		return 0, err
	}

	// in some cases the chain nonce jumps ahead, and we need to handle this case
	// this won't catch all possible timing scenarios, but it should self-heal if the chain jumps
	if next < pending {
		m.mutexNonce.Lock()
		defer m.mutexNonce.Unlock()
		next = atomic.AddUint64(m.nonce, 1)
		if next < pending {
			m.logger.Info(
				fmt.Sprintf(
					"Skew detected. Bumping nonce tracker! Pending/Next:%d/%d",
					pending,
					next,
				),
			)
			atomic.StoreUint64(m.nonce, pending)
			return pending, nil
		}
	}
	return next, nil
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
