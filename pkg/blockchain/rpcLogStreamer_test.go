package blockchain_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/xmtp/xmtpd/pkg/blockchain"

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
	fromBlock := uint64(1)
	lastBlock := uint64(10)
	logMessage := types.Log{
		Address: address,
		Topics:  []common.Hash{topic},
		Data:    []byte("foo"),
	}

	// Create mock blocks with proper headers
	header2 := &types.Header{
		Number:     big.NewInt(2),
		ParentHash: common.Hash{}, // Empty parent hash for simplicity
	}
	mockBlock2 := types.NewBlockWithHeader(header2)

	header11 := &types.Header{
		Number: big.NewInt(11),
	}
	mockBlock11 := types.NewBlockWithHeader(header11)

	mockClient := mocks.NewMockChainClient(t)

	// Mock BlockByNumber call for fromBlockNumber+1 (block 2) - for reorg detection
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(int64(fromBlock+1))).
		Return(mockBlock2, nil)

	// Mock BlockNumber call to get the highest block
	mockClient.On("BlockNumber", mock.Anything).Return(lastBlock, nil)

	// Mock FilterLogs call
	mockClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(lastBlock)),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{logMessage}, nil)

	// Mock BlockByNumber call for toBlock+1 (block 11) - for getting next block hash
	mockClient.On("BlockByNumber", mock.Anything, big.NewInt(int64(lastBlock+1))).
		Return(mockBlock11, nil)

	cfg := blockchain.ContractConfig{
		ID:                "testContract",
		FromBlockNumber:   fromBlock,
		FromBlockHash:     []byte{},
		Address:           address,
		Topics:            []common.Hash{topic},
		MaxDisconnectTime: 5 * time.Minute,
	}

	streamer, err := blockchain.NewRpcLogStreamer(
		context.Background(),
		mockClient,
		testutils.NewLog(t),
		blockchain.WithContractConfig(cfg),
	)
	require.NoError(t, err)

	logs, nextPage, _, _, err := streamer.GetNextPage(context.Background(), cfg, fromBlock, nil)
	require.NoError(t, err)
	expectedNextPage := uint64(11)
	require.Equal(t, &expectedNextPage, nextPage)
	require.Equal(t, 1, len(logs))
	require.Equal(t, logs[0].Address, address)
}
