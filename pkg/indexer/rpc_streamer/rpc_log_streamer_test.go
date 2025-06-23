package rpc_streamer_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	mocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: Add more test coverage.

func TestRpcLogStreamer(t *testing.T) {
	address := testutils.RandomAddress()
	topic := testutils.RandomLogTopic()
	backfillFromBlock := uint64(1)
	backfillToBlock := uint64(11)
	logMessage := types.Log{
		Address: address,
		Topics:  []common.Hash{topic},
		Data:    []byte("foo"),
	}

	mockBlock2 := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(2),
		ParentHash: common.Hash{},
	})

	mockBlock11 := types.NewBlockWithHeader(&types.Header{
		Number: big.NewInt(11),
	})

	mockClient := mocks.NewMockChainClient(t)

	// Mock HeaderByNumber call for fromBlockNumber+1 (block 2) - for reorg detection.
	mockClient.On("HeaderByNumber", mock.Anything, big.NewInt(int64(backfillFromBlock+1))).
		Return(mockBlock2.Header(), nil)

	// Mock BlockNumber call to get the highest block.
	mockClient.On("BlockNumber", mock.Anything).Return(mockBlock11.NumberU64(), nil)

	// Mock FilterLogs call
	mockClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(backfillFromBlock)),
		ToBlock:   big.NewInt(int64(backfillToBlock - 1)),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{logMessage}, nil)

	mockClient.On("HeaderByNumber", mock.Anything, big.NewInt(int64(backfillToBlock))).
		Return(mockBlock11.Header(), nil)

	cfg := &rpc_streamer.ContractConfig{
		ID:                "testContract",
		FromBlockNumber:   backfillFromBlock,
		FromBlockHash:     []byte{},
		Address:           address,
		Topics:            []common.Hash{topic},
		MaxDisconnectTime: 5 * time.Minute,
	}

	streamer, err := rpc_streamer.NewRpcLogStreamer(
		context.Background(),
		mockClient,
		testutils.NewLog(t),
		rpc_streamer.WithContractConfig(cfg),
		rpc_streamer.WithBackfillBlockPageSize(10),
	)
	require.NoError(t, err)

	response, err := streamer.GetNextPage(context.Background(), cfg, backfillFromBlock, nil)
	require.NoError(t, err)
	require.EqualValues(t, mockBlock11.NumberU64(), *response.NextBlockNumber)
	require.EqualValues(t, 1, len(response.Logs))
	require.EqualValues(t, response.Logs[0].Address, address)
}
