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

	mockClient := mocks.NewMockChainClient(t)
	mockClient.On("BlockNumber", mock.Anything).Return(uint64(lastBlock), nil)
	mockClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(lastBlock)),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{logMessage}, nil)

	cfg := blockchain.ContractConfig{
		ID:                "testContract",
		FromBlock:         fromBlock,
		ContractAddress:   address,
		Topics:            []common.Hash{topic},
		MaxDisconnectTime: 5 * time.Minute,
	}

	streamer, err := blockchain.NewRpcLogStreamer(
		context.Background(),
		mockClient,
		testutils.NewLog(t),
		1,
		blockchain.WithContractConfig(cfg),
		blockchain.WithBackfillBlockSize(500),
	)
	require.NoError(t, err)

	logs, nextPage, err := streamer.GetNextPage(cfg, fromBlock)
	require.NoError(t, err)
	expectedNextPage := uint64(11)
	require.Equal(t, &expectedNextPage, nextPage)
	require.Equal(t, 1, len(logs))
	require.Equal(t, logs[0].Address, address)
}
