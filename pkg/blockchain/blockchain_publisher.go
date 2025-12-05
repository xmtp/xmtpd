package blockchain

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/blockchain/oracle"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/metrics"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"
	"go.uber.org/zap"
)

var ErrNoLogsFound = errors.New("no logs found")

// 200KB max payload + ABI encoding + safety margin.
const defaultGasLimit = uint64(6_000_000)

// BlockchainPublisher can publish to the blockchain, signing messages using the provided signer.
type BlockchainPublisher struct {
	signer                TransactionSigner
	oracle                oracle.BlockchainOracle
	nonceManager          noncemanager.NonceManager
	logger                *zap.Logger
	client                *ethclient.Client
	replenishCancel       context.CancelFunc
	wg                    sync.WaitGroup
	groupMessageABI       abi.ABI
	identityUpdateABI     abi.ABI
	groupMessageAddress   common.Address
	identityUpdateAddress common.Address
}

func NewBlockchainPublisher(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
	nonceManager noncemanager.NonceManager,
	oracle oracle.BlockchainOracle,
) (*BlockchainPublisher, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	groupMessageABI, err := abi.JSON(strings.NewReader(gm.GroupMessageBroadcasterMetaData.ABI))
	if err != nil {
		return nil, errors.New("failed to parse GroupMessageBroadcaster ABI: " + err.Error())
	}

	identityUpdateABI, err := abi.JSON(strings.NewReader(iu.IdentityUpdateBroadcasterMetaData.ABI))
	if err != nil {
		return nil, errors.New("failed to parse IdentityUpdateBroadcaster ABI: " + err.Error())
	}

	nonce, err := client.PendingNonceAt(ctx, signer.FromAddress())
	if err != nil {
		return nil, err
	}

	logger.Info("starting blockchain publisher with blockchain nonce", utils.NonceField(nonce))

	err = nonceManager.FastForwardNonce(ctx, *new(big.Int).SetUint64(nonce))
	if err != nil {
		return nil, err
	}

	replenishCtx, cancel := context.WithCancel(ctx)

	publisherLogger := logger.Named(utils.BlockchainPublisherLoggerName).
		With(utils.ContractAddressField(contractOptions.AppChain.GroupMessageBroadcasterAddress))

	publisher := BlockchainPublisher{
		signer:          signer,
		logger:          publisherLogger,
		client:          client,
		nonceManager:    nonceManager,
		replenishCancel: cancel,
		groupMessageAddress: common.HexToAddress(
			contractOptions.AppChain.GroupMessageBroadcasterAddress,
		),
		identityUpdateAddress: common.HexToAddress(
			contractOptions.AppChain.IdentityUpdateBroadcasterAddress,
		),
		groupMessageABI:   groupMessageABI,
		identityUpdateABI: identityUpdateABI,
		oracle:            oracle,
	}

	tracing.GoPanicWrap(
		replenishCtx,
		&publisher.wg,
		"replenish-nonces", func(innerCtx context.Context) {
			ticker := time.NewTicker(10 * time.Second)
			for {
				select {
				case <-innerCtx.Done():
					return
				case <-ticker.C:
					nonce, err := client.PendingNonceAt(innerCtx, signer.FromAddress())
					if err != nil {
						logger.Error("error getting pending nonce", zap.Error(err))
						continue
					}
					err = nonceManager.Replenish(innerCtx, *new(big.Int).SetUint64(nonce))
					if err != nil {
						logger.Error("error replenishing nonce", zap.Error(err))
					}
				}
			}
		},
	)

	publisher.oracle.Start()

	return &publisher, nil
}

// sendRawTransaction packs the calldata, creates a transaction, signs it, and sends it.
func (m *BlockchainPublisher) sendRawTransaction(
	ctx context.Context,
	to common.Address,
	data []byte,
	nonce *big.Int,
) (*types.Transaction, error) {
	gasPrice := m.oracle.GetGasPrice()

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce.Uint64(),
		To:       &to,
		Value:    big.NewInt(0),
		Gas:      defaultGasLimit,
		GasPrice: big.NewInt(gasPrice),
		Data:     data,
	})

	signedTx, err := m.signer.SignerFunc()(m.signer.FromAddress(), tx)
	if err != nil {
		return nil, errors.New("failed to sign transaction: " + err.Error())
	}

	err = m.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

func (m *BlockchainPublisher) PublishGroupMessage(
	ctx context.Context,
	groupID [16]byte,
	message []byte,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	if len(message) == 0 {
		return nil, errors.New("message is empty")
	}

	logs, err := withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		"publish_group_message",
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			data, err := m.groupMessageABI.Pack("addMessage", groupID, message)
			if err != nil {
				return nil, errors.New("failed to pack addMessage: " + err.Error())
			}

			return m.sendRawTransaction(ctx, m.groupMessageAddress, data, &nonce)
		},
		func(ctx context.Context, transaction *types.Transaction) ([]*gm.GroupMessageBroadcasterMessageSent, error) {
			receipt, err := WaitForTransaction(
				ctx,
				m.logger,
				m.client,
				2*time.Second,
				250*time.Millisecond,
				transaction.Hash(),
			)
			if err != nil {
				return nil, err
			}

			if receipt == nil {
				return nil, errors.New("transaction receipt is nil")
			}

			return findGroupMessageLogs(receipt, &m.groupMessageABI, 1)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(logs) != 1 {
		return nil, ErrNoLogsFound
	}

	return logs[0], nil
}

func (m *BlockchainPublisher) BootstrapGroupMessages(
	ctx context.Context,
	groupIDs [][16]byte,
	messages [][]byte,
	sequenceIDs []uint64,
) ([]*gm.GroupMessageBroadcasterMessageSent, error) {
	if len(messages) != len(groupIDs) || len(messages) != len(sequenceIDs) {
		return nil, errors.New("messages, groupIDs, and sequenceIDs must have the same length")
	}

	return withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		"bootstrap_group_messages",
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			data, err := m.groupMessageABI.Pack(
				"bootstrapMessages",
				groupIDs,
				messages,
				sequenceIDs,
			)
			if err != nil {
				return nil, errors.New("failed to pack bootstrapMessages: " + err.Error())
			}

			return m.sendRawTransaction(ctx, m.groupMessageAddress, data, &nonce)
		},
		func(ctx context.Context, transaction *types.Transaction) ([]*gm.GroupMessageBroadcasterMessageSent, error) {
			receipt, err := WaitForTransaction(
				ctx,
				m.logger,
				m.client,
				2*time.Second,
				250*time.Millisecond,
				transaction.Hash(),
			)
			if err != nil {
				return nil, err
			}

			if receipt == nil {
				return nil, errors.New("transaction receipt is nil")
			}

			return findGroupMessageLogs(receipt, &m.groupMessageABI, len(groupIDs))
		},
	)
}

func (m *BlockchainPublisher) PublishIdentityUpdate(
	ctx context.Context,
	inboxID [32]byte,
	identityUpdate []byte,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	if len(identityUpdate) == 0 {
		return nil, errors.New("identity update is empty")
	}

	logs, err := withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		"publish_identity_update",
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			data, err := m.identityUpdateABI.Pack("addIdentityUpdate", inboxID, identityUpdate)
			if err != nil {
				return nil, errors.New("failed to pack addIdentityUpdate: " + err.Error())
			}

			return m.sendRawTransaction(ctx, m.identityUpdateAddress, data, &nonce)
		},
		func(ctx context.Context, transaction *types.Transaction) ([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
			receipt, err := WaitForTransaction(
				ctx,
				m.logger,
				m.client,
				2*time.Second,
				250*time.Millisecond,
				transaction.Hash(),
			)
			if err != nil {
				return nil, err
			}

			if receipt == nil {
				return nil, errors.New("transaction receipt is nil")
			}

			return findIdentityUpdateLogs(receipt, &m.identityUpdateABI, 1)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(logs) != 1 {
		return nil, ErrNoLogsFound
	}

	return logs[0], nil
}

func (m *BlockchainPublisher) BootstrapIdentityUpdates(
	ctx context.Context,
	inboxIDs [][32]byte,
	identityUpdates [][]byte,
	sequenceIDs []uint64,
) ([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	if len(identityUpdates) != len(inboxIDs) || len(identityUpdates) != len(sequenceIDs) {
		return nil, errors.New(
			"identityUpdates, inboxIDs, and sequenceIDs must have the same length",
		)
	}

	return withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		"bootstrap_identity_updates",
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			data, err := m.identityUpdateABI.Pack(
				"bootstrapIdentityUpdates",
				inboxIDs,
				identityUpdates,
				sequenceIDs,
			)
			if err != nil {
				return nil, errors.New("failed to pack bootstrapIdentityUpdates: " + err.Error())
			}

			return m.sendRawTransaction(ctx, m.identityUpdateAddress, data, &nonce)
		},
		func(ctx context.Context, transaction *types.Transaction) ([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
			receipt, err := WaitForTransaction(
				ctx,
				m.logger,
				m.client,
				2*time.Second,
				250*time.Millisecond,
				transaction.Hash(),
			)
			if err != nil {
				return nil, err
			}

			if receipt == nil {
				return nil, errors.New("transaction receipt is nil")
			}

			return findIdentityUpdateLogs(receipt, &m.identityUpdateABI, len(inboxIDs))
		},
	)
}

// findGroupMessageLogs parses MessageSent events from the receipt logs.
func findGroupMessageLogs(
	receipt *types.Receipt,
	contractABI *abi.ABI,
	expectedEventCount int,
) ([]*gm.GroupMessageBroadcasterMessageSent, error) {
	events := make([]*gm.GroupMessageBroadcasterMessageSent, 0, expectedEventCount)

	messageSentEvent := contractABI.Events["MessageSent"]

	for _, logEntry := range receipt.Logs {
		if logEntry == nil {
			continue
		}

		// Check if this log matches the MessageSent event signature
		if len(logEntry.Topics) == 0 || logEntry.Topics[0] != messageSentEvent.ID {
			continue
		}

		event := &gm.GroupMessageBroadcasterMessageSent{
			Raw: *logEntry,
		}

		// Parse indexed parameters from topics
		// Topic[0] is the event signature
		// Topic[1] is groupId (bytes16, indexed) - left-aligned in 32-byte topic
		// Topic[2] is sequenceId (uint64, indexed)
		if len(logEntry.Topics) >= 3 {
			copy(event.GroupId[:], logEntry.Topics[1][0:16])
			event.SequenceId = new(big.Int).SetBytes(logEntry.Topics[2][:]).Uint64()
		}

		// Parse non-indexed parameters from data
		// message (bytes, non-indexed)
		if len(logEntry.Data) > 0 {
			unpacked, err := contractABI.Unpack("MessageSent", logEntry.Data)
			if err == nil && len(unpacked) > 0 {
				if msg, ok := unpacked[0].([]byte); ok {
					event.Message = msg
				}
			}
		}

		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, ErrNoLogsFound
	}

	return events, nil
}

// findIdentityUpdateLogs parses IdentityUpdateCreated events from the receipt logs.
func findIdentityUpdateLogs(
	receipt *types.Receipt,
	contractABI *abi.ABI,
	expectedEventCount int,
) ([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	events := make([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, 0, expectedEventCount)

	identityUpdateEvent := contractABI.Events["IdentityUpdateCreated"]

	for _, logEntry := range receipt.Logs {
		if logEntry == nil {
			continue
		}

		// Check if this log matches the IdentityUpdateCreated event signature
		if len(logEntry.Topics) == 0 || logEntry.Topics[0] != identityUpdateEvent.ID {
			continue
		}

		event := &iu.IdentityUpdateBroadcasterIdentityUpdateCreated{
			Raw: *logEntry,
		}

		// Parse indexed parameters from topics
		// Topic[0] is the event signature
		// Topic[1] is inboxId (bytes32, indexed)
		// Topic[2] is sequenceId (uint64, indexed)
		if len(logEntry.Topics) >= 3 {
			copy(event.InboxId[:], logEntry.Topics[1][:])
			event.SequenceId = new(big.Int).SetBytes(logEntry.Topics[2][:]).Uint64()
		}

		// Parse non-indexed parameters from data
		// update (bytes, non-indexed)
		if len(logEntry.Data) > 0 {
			unpacked, err := contractABI.Unpack("IdentityUpdateCreated", logEntry.Data)
			if err == nil && len(unpacked) > 0 {
				if update, ok := unpacked[0].([]byte); ok {
					event.Update = update
				}
			}
		}

		events = append(events, event)
	}

	if len(events) == 0 {
		return nil, ErrNoLogsFound
	}

	return events, nil
}

func withNonce[T any](ctx context.Context,
	logger *zap.Logger,
	nonceManager noncemanager.NonceManager,
	payloadType string,
	create func(context.Context, big.Int) (*types.Transaction, error),
	wait func(context.Context, *types.Transaction) ([]*T, error),
) ([]*T, error) {
	var tx *types.Transaction
	var nonceContext *noncemanager.NonceContext
	var err error

	for {
		nonceContext, err = nonceManager.GetNonce(ctx)
		if err != nil {
			return nil, err
		}
		nonce := nonceContext.Nonce

		tx, err = metrics.MeasureBroadcastTransaction(
			payloadType,
			func() (*types.Transaction, error) {
				return create(ctx, nonce)
			},
		)
		if err != nil {
			if errors.Is(err, core.ErrNonceTooLow) ||
				strings.Contains(
					err.Error(),
					"nonce too low",
				) ||
				strings.Contains(err.Error(), "replacement transaction underpriced") {
				logger.Debug(
					"nonce already used, consume and move on",
					utils.NonceField(nonce.Uint64()),
					zap.Error(err),
				)

				err = nonceContext.Consume()
				if err != nil {
					nonceContext.Cancel()
					return nil, err
				}
				continue
			}

			if strings.Contains(
				err.Error(),
				"nonce too high",
			) {
				// we have been hammering the blockchain too hard
				// back off for a little bit
				logger.Debug(
					"nonce too high, back off for a little bit",
					utils.NonceField(nonce.Uint64()),
					zap.Error(err),
				)
				utils.RandomSleep(ctx, 500*time.Millisecond)
				nonceContext.Cancel()
				continue
			}

			nonceContext.Cancel()
			return nil, err
		}
		break
	}

	val, err := metrics.MeasureWaitForTransaction(func() ([]*T, error) {
		return wait(ctx, tx)
	})
	if err != nil {
		nonceContext.Cancel()
		return nil, err
	}

	err = nonceContext.Consume()
	if err != nil {
		nonceContext.Cancel()
		return nil, err
	}

	return val, nil
}

func (m *BlockchainPublisher) Close() {
	m.logger.Info("closing")
	m.replenishCancel()
	m.wg.Wait()
	m.oracle.Stop()
	m.logger.Info("closed")
}
