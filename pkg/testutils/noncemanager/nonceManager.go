package noncemanager

import (
	"container/heap"
	"context"
	"math/big"
	"sync"

	"github.com/xmtp/xmtpd/pkg/blockchain"
	"go.uber.org/zap"
)

type Int64Heap []int64

func (h *Int64Heap) Len() int           { return len(*h) }
func (h *Int64Heap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *Int64Heap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *Int64Heap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

func (h *Int64Heap) Pop() interface{} {
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

func (tm *TestNonceManager) GetNonce(ctx context.Context) (*blockchain.NonceContext, error) {
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

	tm.logger.Debug("Generated Nonce", zap.Int64("nonce", nonce))

	return &blockchain.NonceContext{
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
