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
	"github.com/xmtp/xmtpd/contracts/pkg/groupmessages"
	"github.com/xmtp/xmtpd/contracts/pkg/identityupdates"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
Can publish to the blockchain, signing messages using the provided signer
*/
type BlockchainPublisher struct {
	signer                 TransactionSigner
	client                 *ethclient.Client
	messagesContract       *groupmessages.GroupMessages
	identityUpdateContract *identityupdates.IdentityUpdates
	logger                 *zap.Logger
	mutexNonce             sync.Mutex
	nonce                  uint64
}

func NewBlockchainPublisher(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
) (*BlockchainPublisher, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	messagesContract, err := groupmessages.NewGroupMessages(
		common.HexToAddress(contractOptions.MessagesContractAddress),
		client,
	)

	if err != nil {
		return nil, err
	}
	identityUpdateContract, err := identityupdates.NewIdentityUpdates(
		common.HexToAddress(contractOptions.IdentityUpdatesContractAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(ctx, signer.FromAddress())
	if err != nil {
		return nil, err
	}

	// The nonce is the next ID to be used, not the current highest
	// The nonce member variable represents the last recenty used, so it is pending-1
	if nonce > 0 {
		nonce--
	}

	logger.Info(fmt.Sprintf("Starting server with blockchain nonce: %d", nonce))

	return &BlockchainPublisher{
		signer: signer,
		logger: logger.Named("GroupBlockchainPublisher").
			With(zap.String("contractAddress", contractOptions.MessagesContractAddress)),
		messagesContract:       messagesContract,
		identityUpdateContract: identityUpdateContract,
		client:                 client,
		nonce:                  nonce,
	}, nil
}

func (m *BlockchainPublisher) PublishGroupMessage(
	ctx context.Context,
	groupID [32]byte,
	message []byte,
) (*groupmessages.GroupMessagesMessageSent, error) {
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
) (*identityupdates.IdentityUpdatesIdentityUpdateCreated, error) {
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

// once fetchNonce returns a nonce, it must be used
// otherwise the chain might see a gap and deadlock
func (m *BlockchainPublisher) fetchNonce(ctx context.Context) (uint64, error) {
	pending, err := m.client.PendingNonceAt(ctx, m.signer.FromAddress())
	if err != nil {
		return 0, err
	}

	next := atomic.AddUint64(&m.nonce, 1)

	m.logger.Debug(
		"Generated nonce",
		zap.Uint64("pending_nonce", pending),
		zap.Uint64("atomic_nonce", next),
	)

	if next >= pending {
		// normal case scenario
		return next, nil

	}

	// in some cases the chain nonce jumps ahead, and we need to handle this case
	// this won't catch all possible timing scenarios, but it should self-heal if the chain jumps
	m.mutexNonce.Lock()
	defer m.mutexNonce.Unlock()
	currentNonce := atomic.LoadUint64(&m.nonce)
	if currentNonce < pending {
		m.logger.Info(
			"Nonce skew detected",
			zap.Uint64("pending_nonce", pending),
			zap.Uint64("current_nonce", currentNonce),
		)
		atomic.StoreUint64(&m.nonce, pending)
		return pending, nil
	}

	return atomic.AddUint64(&m.nonce, 1), nil
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
