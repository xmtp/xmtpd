package rpcstreamer_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"

	rpcstreamer "github.com/xmtp/xmtpd/pkg/indexer/rpc_streamer"
	"github.com/xmtp/xmtpd/pkg/testutils"
	blockchainMocks "github.com/xmtp/xmtpd/pkg/testutils/mocks/blockchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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

	mockClient := blockchainMocks.NewMockChainClient(t)

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

	cfg := &rpcstreamer.ContractConfig{
		ID:                "testContract",
		FromBlockNumber:   backfillFromBlock,
		FromBlockHash:     []byte{},
		Address:           address,
		Topics:            []common.Hash{topic},
		MaxDisconnectTime: 5 * time.Minute,
	}

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		context.Background(),
		mockClient,
		mockClient,
		testutils.NewLog(t),
		rpcstreamer.WithContractConfig(cfg),
		rpcstreamer.WithBackfillBlockPageSize(10),
	)
	require.NoError(t, err)

	response, err := streamer.GetNextPage(context.Background(), cfg, backfillFromBlock, nil)
	require.NoError(t, err)
	require.Equal(t, mockBlock11.NumberU64(), *response.NextBlockNumber)
	require.Len(t, response.Logs, 1)
	require.Equal(t, response.Logs[0].Address, address)
}

// testSubscription is a minimal ethereum.Subscription for testing.
type testSubscription struct {
	errCh chan error
}

func (s *testSubscription) Err() <-chan error { return s.errCh }
func (s *testSubscription) Unsubscribe()      { close(s.errCh) }

// TestBridgingBackfill verifies that the streamer fetches logs in the gap
// between the HTTP head (where initial backfill ends) and the WS head
// (where the subscription starts delivering).
func TestBridgingBackfill(t *testing.T) {
	address := testutils.RandomAddress()
	topic := testutils.RandomLogTopic()

	// HTTP head is at block 10, WS head is at block 15.
	// Backfill starts at block 1 with page size 20, so initial backfill ends at HTTP head (10).
	// Bridging backfill should fetch blocks 11-15 to cover the gap.
	httpHead := uint64(10)
	wsHead := uint64(15)

	initialLog := types.Log{
		Address:     address,
		Topics:      []common.Hash{topic},
		Data:        []byte("initial"),
		BlockNumber: 5,
	}
	gapLog := types.Log{
		Address:     address,
		Topics:      []common.Hash{topic},
		Data:        []byte("gap"),
		BlockNumber: 12,
	}

	mockBlock2 := types.NewBlockWithHeader(&types.Header{
		Number:     big.NewInt(2),
		ParentHash: common.Hash{},
	})

	// HTTP client: first call returns httpHead, subsequent calls return wsHead.
	httpClient := blockchainMocks.NewMockChainClient(t)
	httpClient.On("BlockNumber", mock.Anything).Return(httpHead, nil).Once()

	httpClient.On("HeaderByNumber", mock.Anything, big.NewInt(int64(2))).
		Return(mockBlock2.Header(), nil)

	// Initial backfill: blocks 1-10, returns initialLog.
	httpClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(int64(httpHead)),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{initialLog}, nil).Once()

	// After initial backfill ends (ErrEndOfBackfill with nil NextBlockNumber),
	// backfillFromBlockNumber stays at 1. The bridging loop re-fetches from block 1
	// but now the HTTP client reports wsHead as the highest block.
	httpClient.On("BlockNumber", mock.Anything).Return(wsHead, nil)

	// Bridging backfill: re-fetches from block 1 to wsHead (15), returns both logs.
	// The gapLog at block 12 would have been missed without bridging.
	httpClient.On("FilterLogs", mock.Anything, ethereum.FilterQuery{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(int64(wsHead)),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{topic}},
	}).Return([]types.Log{initialLog, gapLog}, nil)

	// WS client returns wsHead. Called multiple times: validate + watch.
	wsClient := blockchainMocks.NewMockChainClient(t)
	wsClient.On("BlockNumber", mock.Anything).Return(wsHead, nil)

	// SubscribeFilterLogs is called twice: once in validateWatcher, once in watchContract.
	// Each needs its own subscription since validateWatcher unsubscribes.
	validationSub := &testSubscription{errCh: make(chan error, 1)}
	watchSub := &testSubscription{errCh: make(chan error, 1)}
	wsClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(validationSub, nil).Once()
	wsClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
		Return(watchSub, nil).Once()

	cfg := &rpcstreamer.ContractConfig{
		ID:                "testContract",
		FromBlockNumber:   1,
		FromBlockHash:     []byte{},
		Address:           address,
		Topics:            []common.Hash{topic},
		MaxDisconnectTime: 5 * time.Minute,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	streamer, err := rpcstreamer.NewRPCLogStreamer(
		ctx,
		httpClient,
		wsClient,
		testutils.NewLog(t),
		rpcstreamer.WithContractConfig(cfg),
		rpcstreamer.WithBackfillBlockPageSize(20),
	)
	require.NoError(t, err)

	err = streamer.Start()
	require.NoError(t, err)

	ch := streamer.GetEventChannel("testContract")
	require.NotNil(t, ch)

	// Collect logs from the channel. We expect to see the gapLog, which proves
	// that the bridging backfill fetched logs beyond the initial HTTP head.
	var receivedLogs []types.Log
	timeout := time.After(5 * time.Second)
	for {
		select {
		case log := <-ch:
			if len(log.Data) > 0 {
				receivedLogs = append(receivedLogs, log)
				// Once we see the gap log, we know bridging worked.
				if string(log.Data) == "gap" {
					goto done
				}
			}
		case <-timeout:
			t.Fatalf("timed out waiting for gap log, got %d logs", len(receivedLogs))
		}
	}
done:

	// The gap log must be present, proving the bridging backfill covered the gap.
	var foundGap bool
	for _, l := range receivedLogs {
		if string(l.Data) == "gap" {
			foundGap = true
			break
		}
	}
	require.True(t, foundGap, "expected gap log from bridging backfill")

	cancel()
	streamer.Stop()
}
