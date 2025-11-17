package blockchain_test

import (
	"container/heap"
	"context"
	"math/big"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"github.com/xmtp/xmtpd/pkg/testutils/anvil"
	"go.uber.org/zap"
)

type Int64Heap []int64

func (h *Int64Heap) Len() int           { return len(*h) }
func (h *Int64Heap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *Int64Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Int64Heap) Push(x any) {
	*h = append(*h, x.(int64))
}

func (h *Int64Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[0] // Get the smallest element
	*h = old[1:n]
	return x
}

func (h *Int64Heap) Peek() int64 {
	if len(*h) == 0 {
		return -1 // Return an invalid value if empty
	}
	return (*h)[0]
}

type TestNonceManager struct {
	mu        sync.Mutex
	nonce     int64
	logger    *zap.Logger
	abandoned Int64Heap
}

func NewTestNonceManager(logger *zap.Logger) *TestNonceManager {
	return &TestNonceManager{logger: logger}
}

func (tm *TestNonceManager) GetNonce(ctx context.Context) (*noncemanager.NonceContext, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	var nonce int64
	if tm.abandoned.Len() > 0 {
		nonce = heap.Pop(&tm.abandoned).(int64)
	} else {
		nonce = tm.nonce
		tm.nonce++
	}

	return &noncemanager.NonceContext{
		Nonce: *new(big.Int).SetInt64(nonce),
		Cancel: func() {
			tm.mu.Lock()
			defer tm.mu.Unlock()
			tm.abandoned.Push(nonce)
		}, // No-op
		Consume: func() error {
			return nil // No-op
		},
	}, nil
}

func (tm *TestNonceManager) FastForwardNonce(ctx context.Context, nonce big.Int) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.nonce = nonce.Int64()

	return nil
}

func (tm *TestNonceManager) Replenish(ctx context.Context, nonce big.Int) error {
	return nil
}

func buildPublisher(t *testing.T) *blockchain.BlockchainPublisher {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	logger := testutils.NewLog(t)

	wsURL, rpcURL := anvil.StartAnvil(t, false)
	contractsOptions := testutils.NewContractsOptions(t, rpcURL, wsURL)

	signer, err := blockchain.NewPrivateKeySigner(
		testutils.GetPayerOptions(t).PrivateKey,
		contractsOptions.AppChain.ChainID,
	)
	require.NoError(t, err)

	client, err := blockchain.NewRPCClient(
		ctx,
		contractsOptions.SettlementChain.RPCURL,
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		client.Close()
	})

	nonceManager := NewTestNonceManager(logger)

	publisher, err := blockchain.NewBlockchainPublisher(
		ctx,
		logger,
		client,
		signer,
		contractsOptions,
		nonceManager,
	)
	require.NoError(t, err)

	return publisher
}

func TestPublishIdentityUpdate(t *testing.T) {
	publisher := buildPublisher(t)

	tests := []struct {
		name           string
		inboxID        [32]byte
		identityUpdate []byte
		ctx            context.Context
		wantErr        bool
	}{
		{
			name:           "cancelled context",
			inboxID:        testutils.RandomInboxIDBytes(),
			identityUpdate: testutils.RandomBytes(100),
			ctx:            testutils.CancelledContext(),
			wantErr:        true,
		},
		{
			name:           "empty update",
			inboxID:        testutils.RandomInboxIDBytes(),
			identityUpdate: []byte{},
			ctx:            context.Background(),
			wantErr:        true,
		},
		{
			name:           "happy path",
			inboxID:        testutils.RandomInboxIDBytes(),
			identityUpdate: testutils.RandomBytes(104),
			ctx:            context.Background(),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessage, err := publisher.PublishIdentityUpdate(
				tt.ctx,
				tt.inboxID,
				tt.identityUpdate,
			)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessage)
			require.Equal(t, tt.inboxID, logMessage.InboxId)
			require.Equal(t, tt.identityUpdate, logMessage.Update)
			require.Greater(t, logMessage.SequenceId, uint64(0))
			require.NotNil(t, logMessage.Raw.TxHash)
		})
	}
}

func TestBootstrapIdentityUpdate(t *testing.T) {
	publisher := buildPublisher(t)

	tests := []struct {
		name        string
		inboxIDs    [][32]byte
		updates     [][]byte
		sequenceIDs []uint64
		ctx         context.Context
		wantErr     bool
	}{
		// TODO(borja): Add happy path after including app chain parameter registry.
		{
			name:        "fail when contract is not paused",
			inboxIDs:    [][32]byte{testutils.RandomInboxIDBytes()},
			updates:     [][]byte{testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on array length mismatch",
			inboxIDs:    [][32]byte{testutils.RandomInboxIDBytes()},
			updates:     [][]byte{testutils.RandomBytes(100), testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1, 2, 3},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on empty message",
			inboxIDs:    [][32]byte{testutils.RandomInboxIDBytes()},
			updates:     [][]byte{},
			sequenceIDs: []uint64{1},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on cancelled context",
			inboxIDs:    [][32]byte{testutils.RandomInboxIDBytes()},
			updates:     [][]byte{testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1},
			ctx:         testutils.CancelledContext(),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessages, err := publisher.BootstrapIdentityUpdates(
				tt.ctx,
				tt.inboxIDs,
				tt.updates,
				tt.sequenceIDs,
			)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessages)
			require.Equal(t, len(tt.inboxIDs), len(logMessages))
			for i, logMessage := range logMessages {
				require.Equal(t, tt.inboxIDs[i], logMessage.InboxId)
				require.Equal(t, tt.updates[i], logMessage.Update)
				require.Equal(t, tt.sequenceIDs[i], logMessage.SequenceId)
				require.NotNil(t, logMessage.Raw.TxHash)
			}
		})
	}
}

func TestPublishGroupMessage(t *testing.T) {
	publisher := buildPublisher(t)

	tests := []struct {
		name    string
		groupID [16]byte
		message []byte
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "happy path",
			groupID: testutils.RandomGroupID(),
			message: testutils.RandomBytes(100),
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name:    "empty message",
			groupID: testutils.RandomGroupID(),
			message: []byte{},
			ctx:     context.Background(),
			wantErr: true,
		},
		{
			name:    "cancelled context",
			groupID: testutils.RandomGroupID(),
			message: testutils.RandomBytes(100),
			ctx:     testutils.CancelledContext(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessage, err := publisher.PublishGroupMessage(tt.ctx, tt.groupID, tt.message)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessage)
			require.Equal(t, tt.groupID, logMessage.GroupId)
			require.Equal(t, tt.message, logMessage.Message)
			require.Greater(t, logMessage.SequenceId, uint64(0))
			require.NotNil(t, logMessage.Raw.TxHash)
		})
	}
}

func TestBootstrapGroupMessages(t *testing.T) {
	publisher := buildPublisher(t)

	tests := []struct {
		name        string
		groupIDs    [][16]byte
		messages    [][]byte
		sequenceIDs []uint64
		ctx         context.Context
		wantErr     bool
	}{
		// TODO(borja): Add happy path after including app chain parameter registry.
		{
			name:        "fail when contract is not paused",
			groupIDs:    [][16]byte{testutils.RandomGroupID()},
			messages:    [][]byte{testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on array length mismatch",
			groupIDs:    [][16]byte{testutils.RandomGroupID()},
			messages:    [][]byte{testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1, 2, 3},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on empty message",
			groupIDs:    [][16]byte{testutils.RandomGroupID()},
			messages:    [][]byte{},
			sequenceIDs: []uint64{1},
			ctx:         context.Background(),
			wantErr:     true,
		},
		{
			name:        "fail on cancelled context",
			groupIDs:    [][16]byte{testutils.RandomGroupID()},
			messages:    [][]byte{testutils.RandomBytes(100)},
			sequenceIDs: []uint64{1},
			ctx:         testutils.CancelledContext(),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logMessages, err := publisher.BootstrapGroupMessages(
				tt.ctx,
				tt.groupIDs,
				tt.messages,
				tt.sequenceIDs,
			)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, logMessages)
			require.Equal(t, len(tt.groupIDs), len(logMessages))
			for i, logMessage := range logMessages {
				require.Equal(t, tt.groupIDs[i], logMessage.GroupId)
				require.Equal(t, tt.messages[i], logMessage.Message)
				require.Equal(t, tt.sequenceIDs[i], logMessage.SequenceId)
				require.NotNil(t, logMessage.Raw.TxHash)
			}
		})
	}
}

func TestPublishGroupMessageConcurrent(t *testing.T) {
	publisher := buildPublisher(t)

	const parallelRuns = 100
	var wg sync.WaitGroup
	errSet := sync.Map{}

	for i := 0; i < parallelRuns; i++ {
		wg.Go(func() {
			_, err := publisher.PublishGroupMessage(
				context.Background(),
				testutils.RandomGroupID(),
				testutils.RandomBytes(100),
			)
			if err != nil {
				errSet.Store(err.Error(), struct{}{})
			}
		})
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Collect and print unique errors
	var uniqueErrors []string
	errSet.Range(func(key, value any) bool {
		uniqueErrors = append(uniqueErrors, key.(string))
		return true
	})

	if len(uniqueErrors) > 0 {
		t.Errorf("Errors encountered:\n%s", strings.Join(uniqueErrors, "\n"))
	}
}
