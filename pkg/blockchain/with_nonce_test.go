package blockchain

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/blockchain/noncemanager"
	"github.com/xmtp/xmtpd/pkg/testutils"
)

type stubNonceManager struct {
	nonce    int64
	canceled bool
	consumed bool
}

func (s *stubNonceManager) GetNonce(_ context.Context) (*noncemanager.NonceContext, error) {
	n := s.nonce
	return &noncemanager.NonceContext{
		Nonce: *new(big.Int).SetInt64(n),
		Cancel: func() {
			s.canceled = true
		},
		Consume: func() error {
			s.consumed = true
			return nil
		},
	}, nil
}

func (s *stubNonceManager) FastForwardNonce(_ context.Context, _ big.Int) error {
	return nil
}

func (s *stubNonceManager) Replenish(_ context.Context, _ big.Int) error {
	return nil
}

func TestIsAlreadyKnownError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "sentinel error",
			err:  txpool.ErrAlreadyKnown,
			want: true,
		},
		{
			name: "string match from RPC",
			err:  errors.New("transaction already known"),
			want: true,
		},
		{
			name: "unrelated error",
			err:  errors.New("insufficient funds"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isAlreadyKnownError(tt.err))
		})
	}
}

func TestWithNonce_AlreadyKnownProceedsToWait(t *testing.T) {
	logger := testutils.NewLog(t)
	nm := &stubNonceManager{nonce: 1}

	dummyTx := types.NewTx(&types.LegacyTx{
		Nonce: 1,
		Gas:   21000,
	})

	waitCalled := false

	results, err := withNonce(
		context.Background(),
		logger,
		nm,
		"test_already_known",
		func(_ context.Context, _ big.Int) (*types.Transaction, error) {
			return dummyTx, txpool.ErrAlreadyKnown
		},
		func(_ context.Context, tx *types.Transaction) ([]*int, error) {
			waitCalled = true
			require.Equal(t, dummyTx.Hash(), tx.Hash())
			val := 42
			return []*int{&val}, nil
		},
	)

	require.NoError(t, err)
	require.True(t, waitCalled, "wait function should have been called for already known tx")
	require.Len(t, results, 1)
	require.Equal(t, 42, *results[0])
	require.True(t, nm.consumed, "nonce should have been consumed")
	require.False(t, nm.canceled, "nonce should not have been canceled")
}

func TestWithNonce_NilTxAlreadyKnownCancels(t *testing.T) {
	logger := testutils.NewLog(t)
	nm := &stubNonceManager{nonce: 1}

	_, err := withNonce(
		context.Background(),
		logger,
		nm,
		"test_nil_tx",
		func(_ context.Context, _ big.Int) (*types.Transaction, error) {
			return nil, txpool.ErrAlreadyKnown
		},
		func(_ context.Context, _ *types.Transaction) ([]*int, error) {
			t.Fatal("wait should not be called when tx is nil")
			return nil, nil
		},
	)

	require.Error(t, err)
	require.True(t, nm.canceled, "nonce should be canceled when tx is nil")
}

func TestWithNonce_AlreadyKnownWaitFails(t *testing.T) {
	logger := testutils.NewLog(t)
	nm := &stubNonceManager{nonce: 1}

	dummyTx := types.NewTx(&types.LegacyTx{Nonce: 1, Gas: 21000})

	_, err := withNonce(
		context.Background(),
		logger,
		nm,
		"test_wait_fails",
		func(_ context.Context, _ big.Int) (*types.Transaction, error) {
			return dummyTx, txpool.ErrAlreadyKnown
		},
		func(_ context.Context, _ *types.Transaction) ([]*int, error) {
			return nil, errors.New("wait timeout")
		},
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "wait timeout")
	require.True(t, nm.canceled, "nonce should be canceled when wait fails")
}
