// Package blockchain provides infrastructure for interacting with the blockchain.
package blockchain

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
	"github.com/xmtp/xmtpd/pkg/merkle"
	"github.com/xmtp/xmtpd/pkg/payerreport"
)

// LogStreamBuilder constructs a raw blockchain listener that can be used to listen for events across many contract event types.
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
	ethereum.BlockNumberReader
	ethereum.LogFilterer
	ethereum.ChainIDReader
	ethereum.ChainReader
}

type TransactionSigner interface {
	FromAddress() common.Address
	SignerFunc() bind.SignerFn
}

type IBlockchainPublisher interface {
	Close()
	PublishIdentityUpdate(
		ctx context.Context,
		inboxID [32]byte,
		identityUpdate []byte,
	) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error)
	BootstrapIdentityUpdates(
		ctx context.Context,
		inboxIDs [][32]byte,
		identityUpdates [][]byte,
		sequenceIDs []uint64,
	) ([]*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error)
	PublishGroupMessage(
		ctx context.Context,
		groupID [16]byte,
		message []byte,
	) (*gm.GroupMessageBroadcasterMessageSent, error)
	BootstrapGroupMessages(
		ctx context.Context,
		groupIDs [][16]byte,
		messages [][]byte,
		sequenceIDs []uint64,
	) ([]*gm.GroupMessageBroadcasterMessageSent, error)
}

type PayerReportsManager interface {
	SubmitPayerReport(
		ctx context.Context,
		report *payerreport.PayerReportWithStatus,
	) (int32, ProtocolError)
	GetReport(
		ctx context.Context,
		originatorNodeID uint32,
		index uint64,
	) (*payerreport.PayerReport, error)
	GetDomainSeparator(ctx context.Context) (common.Hash, error)
	GetReportID(
		ctx context.Context,
		payerReport *payerreport.PayerReportWithStatus,
	) (payerreport.ReportID, error)
	SettleReport(
		ctx context.Context,
		originatorNodeID uint32,
		index uint64,
		proof *merkle.MultiProof,
	) error
	SettlementSummary(
		ctx context.Context,
		originatorNodeID uint32,
		index uint64,
	) (*SettlementSummary, error)
}
