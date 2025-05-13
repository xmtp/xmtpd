package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	gm "github.com/xmtp/xmtpd/pkg/abi/groupmessagebroadcaster"
	iu "github.com/xmtp/xmtpd/pkg/abi/identityupdatebroadcaster"
)

type DatabaseChainClient struct{}

func NewDatabaseChainClient() *DatabaseChainClient {
	return &DatabaseChainClient{}
}

func (d *DatabaseChainClient) BlockNumber(ctx context.Context) (uint64, error) {
	return 0, fmt.Errorf("BlockNumber not implemented")
}

func (d *DatabaseChainClient) FilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
) ([]types.Log, error) {
	return nil, fmt.Errorf("FilterLogs not implemented")
}

func (d *DatabaseChainClient) BlockByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Block, error) {
	return nil, fmt.Errorf("BlockByNumber not implemented")
}

func (d *DatabaseChainClient) ParseMessageSent(
	log types.Log,
) (*gm.GroupMessageBroadcasterMessageSent, error) {
	return nil, fmt.Errorf("ParseMessageSent not implemented")
}

func (d *DatabaseChainClient) ParseIdentityUpdateCreated(
	log types.Log,
) (*iu.IdentityUpdateBroadcasterIdentityUpdateCreated, error) {
	return nil, fmt.Errorf("ParseIdentityUpdateCreated not implemented")
}
