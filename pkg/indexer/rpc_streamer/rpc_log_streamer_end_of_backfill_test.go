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
	logMessage := types.Log{
		Address:     address,
		Topics:      []common.Hash{topic},
		Data:        []byte("tip"),
		BlockNumber: head,
	}

	mockClient := blockchainMocks.NewMockChainClient(t)
	mockClient.On("BlockNumber", mock.Anything).Return(head, nil)
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

	// A 32-byte hash would trigger the reorg path below the head. At the head there is no H+1
	// header to inspect, so GetNextPage must fetch H and finish the backfill directly.
	fromBlockHash := make([]byte, common.HashLength)
	response, err := streamer.GetNextPage(t.Context(), cfg, head, fromBlockHash)

	require.ErrorIs(t, err, rpcstreamer.ErrEndOfBackfill)
	require.Equal(t, head+1, *response.NextBlockNumber)
	require.Equal(t, []types.Log{logMessage}, response.Logs)
	mockClient.AssertNotCalled(t, "HeaderByNumber", mock.Anything, mock.Anything)
}

func newBigInt(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}
