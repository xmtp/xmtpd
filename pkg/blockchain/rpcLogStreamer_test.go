package blockchain_test

import (
	"context"
	"testing"
	"time"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/config"
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
	reader blockchain.AppChainReader,
	eventType blockchain.EventType,
	fromBlock uint64,
) (*blockchain.RpcLogStreamer, chan types.Log) {
	log, err := zap.NewDevelopment()
	require.NoError(t, err)
	channel := make(chan types.Log)
	cfg := blockchain.ContractConfig{
		EventType:    eventType,
		FromBlock:    fromBlock,
		EventChannel: channel,
	}
	return blockchain.NewRpcLogStreamer(
		context.Background(),
		reader,
		log,
		[]blockchain.ContractConfig{cfg},
	), channel
}

func TestBuilder(t *testing.T) {
	cfg := config.AppChainOptions{
		RpcURL:                           anvil.StartAnvil(t, false),
		GroupMessageBroadcasterAddress:   testutils.RandomAddress().Hex(),
		IdentityUpdateBroadcasterAddress: testutils.RandomAddress().Hex(),
	}
	testclient, err := blockchain.NewAppChainReader(
		context.Background(),
		cfg,
	)
	require.NoError(t, err)
	builder := blockchain.NewRpcLogStreamBuilder(
		context.Background(),
		testclient,
		testutils.NewLog(t),
	)

	listenerChannel, _ := builder.ListenForContractEvent(
		blockchain.EventTypeMessageSent,
		1,
		5*time.Minute,
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

	mockClient := mocks.NewMockAppChainReader(t)
	mockClient.On("BlockNumber", mock.Anything).Return(uint64(lastBlock), nil)
	mockClient.On("FilterLogs", mock.Anything, blockchain.EventTypeMessageSent, fromBlock, lastBlock).
		Return([]types.Log{logMessage}, nil)
	mockClient.On("ContractAddress", blockchain.EventTypeMessageSent).Return(address.Hex(), nil)

	streamer, _ := buildStreamer(t, mockClient, blockchain.EventTypeMessageSent, fromBlock)

	cfg := blockchain.ContractConfig{
		FromBlock:    fromBlock,
		EventType:    blockchain.EventTypeMessageSent,
		EventChannel: make(chan types.Log),
	}

	logs, nextPage, err := streamer.GetNextPage(cfg, fromBlock)
	require.NoError(t, err)
	expectedNextPage := uint64(11)
	require.Equal(t, &expectedNextPage, nextPage)
	require.Equal(t, 1, len(logs))
	require.Equal(t, logs[0].Address, address)
}
