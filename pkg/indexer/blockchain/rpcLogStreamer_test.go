package blockchain

import (
	"context"
	big "math/big"
	"testing"

	"github.com/xmtp/xmtpd/pkg/mocks"
	"github.com/xmtp/xmtpd/pkg/testutils"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func buildStreamer(
	t *testing.T,
	client ChainClient,
	fromBlock int,
	address common.Address,
	topic common.Hash,
) (*RpcLogStreamer, chan types.Log) {
	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	channel := make(chan types.Log)
	cfg := contractConfig{
		fromBlock:       fromBlock,
		contractAddress: address,
		topics:          []common.Hash{topic},
		channel:         channel,
	}
	return NewRpcLogStreamer(client, log, []contractConfig{cfg}), channel
}

func TestBuilder(t *testing.T) {
	testclient, err := NewClient(context.Background(), testutils.GetContractsOptions(t).RpcUrl)
	require.NoError(t, err)
	builder := NewRpcLogStreamBuilder(testclient, testutils.NewLog(t))

	listenerChannel := builder.ListenForContractEvent(
		1,
		testutils.RandomAddress(),
		[]common.Hash{testutils.RandomLogTopic()},
	)
	require.NotNil(t, listenerChannel)

	streamer, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, streamer)
}

func TestRpcLogStreamer(t *testing.T) {
	address := testutils.RandomAddress()
	topic := testutils.RandomLogTopic()
	fromBlock := 1
	lastBlock := 10
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

	streamer, _ := buildStreamer(t, mockClient, fromBlock, address, topic)

	cfg := contractConfig{
		fromBlock:       fromBlock,
		contractAddress: address,
		topics:          []common.Hash{topic},
		channel:         make(chan types.Log),
	}

	logs, nextPage, err := streamer.getNextPage(cfg, fromBlock)
	require.NoError(t, err)
	expectedNextPage := 11
	require.Equal(t, &expectedNextPage, nextPage)
	require.Equal(t, 1, len(logs))
	require.Equal(t, logs[0].Address, address)
}
