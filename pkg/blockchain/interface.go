package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

// Each event type maps to a specific contract address and topic
type EventType int

const (
	EventTypeMessageSent EventType = iota
	EventTypeIdentityUpdateCreated
)

// Construct a raw blockchain listener that can be used to listen for events across many contract event types
type LogStreamBuilder interface {
	ListenForContractEvent(
		eventType EventType,
		fromBlock uint64,
	) <-chan types.Log
	Build() (LogStreamer, error)
}

type LogStreamer interface {
	Start(ctx context.Context) error
}

type AppChainReader interface {
	FilterLogs(
		ctx context.Context,
		eventType EventType,
		fromBlock uint64,
		toBlock uint64,
	) ([]types.Log, error)

	ContractAddress(eventType EventType) (string, error)

	// Matches ethereum.BlockNumberReader
	BlockNumber(ctx context.Context) (uint64, error)

	// Matches ethereum.ChainReader
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)

	// Matches contract ABI's
	ParseMessageSent(log types.Log) (*gm.GroupMessageBroadcasterMessageSent, error)
	ParseIdentityUpdateCreated(
		log types.Log,
	) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error)
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
