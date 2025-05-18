package blockchain

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

// Implements AppChainReader
// Mocks a blockchain using the database from the pre-decentralization XMTP node.
// Used to migrate data from the old network to the new network, and can be replaced
// by EthAppChainReader when the decentralization rollout is complete.
type DatabaseAppChainReader struct{}

func NewDatabaseAppChainReader() *DatabaseAppChainReader {
	return &DatabaseAppChainReader{}
}

func (d *DatabaseAppChainReader) FilterLogs(
	ctx context.Context,
	eventType EventType,
	fromBlock uint64,
	toBlock uint64,
) ([]types.Log, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseAppChainReader) ContractAddress(eventType EventType) (string, error) {
	return "", errors.New("method not implemented")
}

func (d *DatabaseAppChainReader) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, errors.New("method not implemented")
}

func (d *DatabaseAppChainReader) BlockByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Block, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseAppChainReader) ParseMessageSent(
	log types.Log,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseAppChainReader) ParseIdentityUpdateCreated(
	log types.Log,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	return nil, errors.New("method not implemented")
}
