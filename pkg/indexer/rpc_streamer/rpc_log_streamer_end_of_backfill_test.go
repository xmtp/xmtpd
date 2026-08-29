package rpcstreamer_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	rpcstreamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/testutils"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/blockchain"
)

func TestGetNextPageAtChainTipSkipsNextHeaderProbe(t *testing.T) {
	address := testutils.RandomAddress()
	topic := testutils.RandomLogTopic()
	head := uint64(11)
	headHeader := &types.Header{Number: newBigInt(head)}
	logMessage := types.Log{
		Address:     address,
		Topics:      []common.Hash{topic},
		Data:        []byte("tip"),
		BlockNumber: head,
	}

	mockClient := blockchainMocks.NewMockChainClient(t)
	mockClient.On("BlockNumber", mock.Anything).Return(head, nil)
	mockClient.On("HeaderByNumber", mock.Anything, newBigInt(head)).Return(headHeader, nil)
	mockClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: newBigInt(head),
		ToBlock:   newBigInt(head),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{logMessage}, nil)

	cfg := &rpcstreamer.ContractConfig{
		ID:      "testContract",
		Address: address,
		Topics:  []common.Hash{topic},
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		t.Context(),
		mockClient,
		mockClient,
		testutils.NewLog(t),
		rpcstreamer.WithContractConfig(cfg),
	)
	require.NoError(t, err)

	response, err := streamer.GetNextPage(t.Context(), cfg, head, headHeader.Hash().Bytes())

	require.ErrorIs(t, err, rpcstreamer.ErrEndOfBackfill)
	require.Equal(t, head+1, *response.NextBlockNumber)
	require.Equal(t, []types.Log{logMessage}, response.Logs)
	mockClient.AssertNotCalled(t, "HeaderByNumber", mock.Anything, newBigInt(head+1))
}

func TestGetNextPageAtChainTipDetectsCursorReorg(t *testing.T) {
	address := testutils.RandomAddress()
	topic := testutils.RandomLogTopic()
	head := uint64(11)
	headHeader := &types.Header{Number: newBigInt(head)}
	previousHeader := &types.Header{Number: newBigInt(head - 1)}

	mockClient := blockchainMocks.NewMockChainClient(t)
	mockClient.On("BlockNumber", mock.Anything).Return(head, nil)
	mockClient.On("HeaderByNumber", mock.Anything, newBigInt(head)).Return(headHeader, nil)
	mockClient.On("HeaderByNumber", mock.Anything, newBigInt(head-1)).Return(previousHeader, nil)

	cfg := &rpcstreamer.ContractConfig{
		ID:      "testContract",
		Address: address,
		Topics:  []common.Hash{topic},
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		t.Context(),
		mockClient,
		mockClient,
		testutils.NewLog(t),
		rpcstreamer.WithContractConfig(cfg),
	)
	require.NoError(t, err)

	staleHash := common.HexToHash("0xdeadbeef").Bytes()
	response, err := streamer.GetNextPage(t.Context(), cfg, head, staleHash)

	require.ErrorIs(t, err, rpcstreamer.ErrReorg)
	require.Equal(t, head-1, *response.NextBlockNumber)
	require.Equal(t, previousHeader.Hash().Bytes(), response.NextBlockHash)
	require.Empty(t, response.Logs)
	mockClient.AssertNotCalled(t, "HeaderByNumber", mock.Anything, newBigInt(head+1))
	mockClient.AssertNotCalled(t, "FilterLogs", mock.Anything, mock.Anything)
}

func newBigInt(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}
