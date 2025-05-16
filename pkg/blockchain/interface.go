package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

// Construct a raw blockchain listener that can be used to listen for events across many contract event types
type LogStreamBuilder interface {
	ListenForContractEvent(
		fromBlock uint64,
		contractAddress common.Address,
		topic common.Hash,
	) <-chan types.Log
	Build() (LogStreamer, error)
}

type LogStreamer interface {
	Start(ctx context.Context) error
}

type ChainClient interface {
	// From ethereum.BlockNumberReader
	BlockNumber(ctx context.Context) (uint64, error)

	// From ethereum.LogFilterer
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(
		ctx context.Context,
		q ethereum.FilterQuery,
		ch chan<- types.Log,
	) (ethereum.Subscription, error)

	// From ethereum.ChainIDReader
	ChainID(ctx context.Context) (*big.Int, error)

	// Only the BlockByNumber method from ethereum.ChainReader
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type TransactionSigner interface {
	FromAddress() common.Address
	SignerFunc() bind.SignerFn
}

type IBlockchainPublisher interface {
	PublishIdentityUpdate(
		ctx context.Context,
		inboxId [32]byte,
		identityUpdate []byte,
	) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error)
	PublishGroupMessage(
		ctx context.Context,
		groupdId [32]byte,
		message []byte,
	) (*gm.GroupMessageBroadcasterMessageSent, error)
}
