package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/xmtp/xmtpd/pkg/tracing"
	"github.com/xmtp/xmtpd/pkg/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/config"
	"go.uber.org/zap"
)

/*
Can publish to the blockchain, signing messages using the provided signer
*/
type BlockchainPublisher struct {
	signer                 TransactionSigner
	client                 *ethclient.Client
	messagesContract       *gm.GroupMessageBroadcaster
	identityUpdateContract *iu.IdentityUpdateBroadcaster
	logger                 *zap.Logger
	nonceManager           noncemanager.NonceManager
	replenishCancel        context.CancelFunc
	wg                     sync.WaitGroup
}

func NewBlockchainPublisher(
	ctx context.Context,
	logger *zap.Logger,
	client *ethclient.Client,
	signer TransactionSigner,
	contractOptions config.ContractsOptions,
	nonceManager noncemanager.NonceManager,
) (*BlockchainPublisher, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	messagesContract, err := gm.NewGroupMessageBroadcaster(
		common.HexToAddress(contractOptions.AppChain.GroupMessageBroadcasterAddress),
		client,
	)
	if err != nil {
		return nil, err
	}
	identityUpdateContract, err := iu.NewIdentityUpdateBroadcaster(
		common.HexToAddress(contractOptions.AppChain.IdentityUpdateBroadcasterAddress),
		client,
	)
	if err != nil {
		return nil, err
	}

	nonce, err := client.PendingNonceAt(ctx, signer.FromAddress())
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Starting server with blockchain nonce: %d", nonce))

	err = nonceManager.FastForwardNonce(ctx, *new(big.Int).SetUint64(nonce))
	if err != nil {
		return nil, err
	}

	replenishCtx, cancel := context.WithCancel(ctx)

	streamerLogger := logger.Named("GroupBlockchainPublisher").
		With(zap.String("contractAddress", contractOptions.AppChain.GroupMessageBroadcasterAddress))

	publisher := BlockchainPublisher{
		signer:                 signer,
		logger:                 streamerLogger,
		messagesContract:       messagesContract,
		identityUpdateContract: identityUpdateContract,
		client:                 client,
		nonceManager:           nonceManager,
		replenishCancel:        cancel,
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

	return &publisher, nil
}

func (m *BlockchainPublisher) PublishGroupMessage(
	ctx context.Context,
	groupID [16]byte,
	message []byte,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	if len(message) == 0 {
		return nil, errors.New("message is empty")
	}

	return withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			return m.messagesContract.AddMessage(&bind.TransactOpts{
				Context: ctx,
				Nonce:   &nonce,
				From:    m.signer.FromAddress(),
				Signer:  m.signer.SignerFunc(),
			}, groupID, message)
		},
		func(ctx context.Context, transaction *types.Transaction) (*gm.GroupMessageBroadcasterMessageSent, error) {
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

			return findLog(
				receipt,
				m.messagesContract.ParseMessageSent,
				"no message sent log found",
			)
		},
	)
}

func (m *BlockchainPublisher) PublishIdentityUpdate(
	ctx context.Context,
	inboxId [32]byte,
	identityUpdate []byte,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	if len(identityUpdate) == 0 {
		return nil, errors.New("identity update is empty")
	}

	return withNonce(
		ctx,
		m.logger,
		m.nonceManager,
		func(ctx context.Context, nonce big.Int) (*types.Transaction, error) {
			return m.identityUpdateContract.AddIdentityUpdate(&bind.TransactOpts{
				Context: ctx,
				Nonce:   &nonce,
				From:    m.signer.FromAddress(),
				Signer:  m.signer.SignerFunc(),
			}, inboxId, identityUpdate)
		},
		func(ctx context.Context, transaction *types.Transaction) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
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

			return findLog(
				receipt,
				m.identityUpdateContract.ParseIdentityUpdateCreated,
				"no message sent log found",
			)
		},
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

func withNonce[T any](ctx context.Context,
	logger *zap.Logger,
	nonceManager noncemanager.NonceManager,
	create func(context.Context, big.Int) (*types.Transaction, error),
	wait func(context.Context, *types.Transaction) (*T, error),
) (*T, error) {
	var tx *types.Transaction
	var nonceContext *noncemanager.NonceContext
	var err error

	for {
		nonceContext, err = nonceManager.GetNonce(ctx)
		if err != nil {
			return nil, err
		}
		nonce := nonceContext.Nonce
		tx, err = create(ctx, nonce)
		if err != nil {
			if errors.Is(err, core.ErrNonceTooLow) ||
				strings.Contains(
					err.Error(),
					"nonce too low",
				) ||
				strings.Contains(err.Error(), "replacement transaction underpriced") {
				logger.Debug(
					"Nonce already used, consuming and moving on...",
					zap.Uint64("nonce", nonce.Uint64()),
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
					"Nonce too high, backing off...",
					zap.Uint64("nonce", nonce.Uint64()),
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

	defer func() {
		if err != nil {
			nonceContext.Cancel()
		}
	}()

	val, err := wait(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = nonceContext.Consume()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (m *BlockchainPublisher) Close() {
	m.logger.Info("closing blockchain publisher")
	m.replenishCancel()
	m.wg.Wait()
}
