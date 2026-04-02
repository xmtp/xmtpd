package blockchain

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

type mockSigner struct{}

func (m *mockSigner) FromAddress() common.Address {
	return common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
}

func (m *mockSigner) SignerFunc() bind.SignerFn {
	return func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		return tx, nil
	}
}

func TestExecuteWithRetry(t *testing.T) {
	tests := []struct {
		name          string
		setupCtx      func() (context.Context, context.CancelFunc)
		txFunc        func(callCount int32) error
		wantCalls     int32
		wantMaxCalls  int32 // if set, assert LessOrEqual instead of Equal
		wantErrSubstr string
	}{
		{
			name: "retries on underpriced then returns other error",
			txFunc: func(count int32) error {
				if count <= 2 {
					return errors.New("replacement transaction underpriced")
				}
				return errors.New("some other error")
			},
			wantCalls:     3,
			wantErrSubstr: "some other error",
		},
		{
			name: "no retry on non-retryable errors",
			txFunc: func(count int32) error {
				return errors.New("execution reverted")
			},
			wantCalls:     1,
			wantErrSubstr: "execution reverted",
		},
		{
			name: "exhausts retries on persistent underpriced",
			txFunc: func(count int32) error {
				return errors.New("transaction underpriced")
			},
			wantCalls:     int32(executeTxMaxRetries + 1),
			wantErrSubstr: "underpriced",
		},
		{
			name: "respects context cancellation",
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			txFunc: func(count int32) error {
				return errors.New("transaction underpriced")
			},
			wantMaxCalls:  2,
			wantErrSubstr: "context canceled",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var (
				ctx       = t.Context()
				cancel    context.CancelFunc
				logger    = testutils.NewLog(t)
				signer    *mockSigner
				callCount atomic.Int32
			)

			// Replace ctx and cancel if setupCtx is provided.
			if tc.setupCtx != nil {
				ctx, cancel = tc.setupCtx()
			}

			txFunc := func(opts *bind.TransactOpts) (*types.Transaction, error) {
				count := callCount.Add(1)

				if cancel != nil && count == 1 {
					cancel()
				}

				return nil, tc.txFunc(count)
			}

			opts := &bind.TransactOpts{
				Context: ctx,
				From:    signer.FromAddress(),
				Signer:  signer.SignerFunc(),
			}

			_, err := executeTransaction(ctx, logger, opts, txFunc)
			require.NotNil(t, err)
			assert.Contains(t, err.Error(), tc.wantErrSubstr)

			if tc.wantMaxCalls > 0 {
				assert.LessOrEqual(t, callCount.Load(), tc.wantMaxCalls)
			} else {
				assert.Equal(t, tc.wantCalls, callCount.Load())
			}
		})
	}
}
