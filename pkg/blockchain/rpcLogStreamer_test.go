package blockchain_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"

	mocks "github.com/xmtp/xmtpd/pkg/mocks/blockchain"
	"github.com/xmtp/xmtpd/pkg/testutils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func buildStreamer(
	t *testing.T,
	client blockchain.ChainClient,
	fromBlock uint64,
	address common.Address,
	topic common.Hash,
) (*blockchain.RpcLogStreamer, chan types.Log) {
	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	channel := make(chan types.Log)
	cfg := blockchain.ContractConfig{
		FromBlock:       fromBlock,
		ContractAddress: address,
		Topics:          []common.Hash{topic},
		EventChannel:    channel,
	}
	return blockchain.NewRpcLogStreamer(
		context.Background(),
		client,
		log,
		[]blockchain.ContractConfig{cfg},
	), channel
}

func TestBuilder(t *testing.T) {
	rpcUrl := anvil.StartAnvil(t, false)
	testclient, err := blockchain.NewChainClient(
		context.Background(),
		rpcUrl,
	)
	require.NoError(t, err)
	builder := blockchain.NewRpcLogStreamBuilder(
		context.Background(),
		testclient,
		testutils.NewLog(t),
	)

	listenerChannel, _ := builder.ListenForContractEvent(
		1,
		testutils.RandomAddress(),
		[]common.Hash{testutils.RandomLogTopic()}, 5*time.Minute,
	)
	require.NotNil(t, listenerChannel)

	streamer, err := builder.Build()
	require.NoError(t, err)
	require.NotNil(t, streamer)
}

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

	streamer, _ := buildStreamer(t, mockClient, fromBlock, address, topic)

	cfg := blockchain.ContractConfig{
		FromBlock:       fromBlock,
		ContractAddress: address,
		Topics:          []common.Hash{topic},
		EventChannel:    make(chan types.Log),
	}

	logs, nextPage, err := streamer.GetNextPage(cfg, fromBlock)
	require.NoError(t, err)
	expectedNextPage := uint64(11)
	require.Equal(t, &expectedNextPage, nextPage)
	require.Equal(t, 1, len(logs))
	require.Equal(t, logs[0].Address, address)
}
