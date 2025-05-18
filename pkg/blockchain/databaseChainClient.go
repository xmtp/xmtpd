package blockchain

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

type DatabaseChainClient struct{}

func NewDatabaseChainClient() *DatabaseChainClient {
	return &DatabaseChainClient{}
}

func (d *DatabaseChainClient) FilterLogs(
	ctx context.Context,
	eventType EventType,
	fromBlock uint64,
	toBlock uint64,
) ([]types.Log, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseChainClient) ContractAddress(eventType EventType) (string, error) {
	return "", errors.New("method not implemented")
}

func (d *DatabaseChainClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, errors.New("method not implemented")
}

func (d *DatabaseChainClient) BlockByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Block, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseChainClient) ParseMessageSent(
	log types.Log,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	return nil, errors.New("method not implemented")
}

func (d *DatabaseChainClient) ParseIdentityUpdateCreated(
	log types.Log,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	return nil, errors.New("method not implemented")
}
