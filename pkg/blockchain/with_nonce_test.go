package blockchain

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
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

func TestWithNonce_AlreadyKnown(t *testing.T) {
	logger := testutils.NewLog(t)
	nm := &stubNonceManager{nonce: 1}

	dummyTx := types.NewTx(&types.LegacyTx{
		Nonce: 1,
		Gas:   21000,
	})

	createCalled := false
	waitCalled := false

	results, err := withNonce(
		context.Background(),
		logger,
		nm,
		"test_already_known",
		func(_ context.Context, _ big.Int) (*types.Transaction, error) {
			createCalled = true
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
	require.True(t, createCalled, "create function should have been called")
	require.True(t, waitCalled, "wait function should have been called for already known tx")
	require.Len(t, results, 1)
	require.Equal(t, 42, *results[0])
	require.True(t, nm.consumed, "nonce should have been consumed, not canceled")
	require.False(t, nm.canceled, "nonce should not have been canceled")
}

func TestWithNonce_AlreadyKnownStringMatch(t *testing.T) {
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
		"test_already_known_string",
		func(_ context.Context, _ big.Int) (*types.Transaction, error) {
			return dummyTx, &rpcError{msg: "transaction already known"}
		},
		func(_ context.Context, _ *types.Transaction) ([]*int, error) {
			waitCalled = true
			val := 1
			return []*int{&val}, nil
		},
	)

	require.NoError(t, err)
	require.True(t, waitCalled)
	require.Len(t, results, 1)
	require.True(t, nm.consumed)
	require.False(t, nm.canceled)
}

// rpcError simulates an RPC error that contains "already known" in its message.
type rpcError struct {
	msg string
}

func (e *rpcError) Error() string {
	return e.msg
}

func TestWithNonce_UnknownErrorStillCancels(t *testing.T) {
	logger := testutils.NewLog(t)

	t.Run("nil tx with already known still cancels", func(t *testing.T) {
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
	})

	t.Run("unrelated error cancels nonce", func(t *testing.T) {
		nm := &stubNonceManager{nonce: 1}

		dummyTx := types.NewTx(&types.LegacyTx{Nonce: 1, Gas: 21000})
		_, err := withNonce(
			context.Background(),
			logger,
			nm,
			"test_unrelated_error",
			func(_ context.Context, _ big.Int) (*types.Transaction, error) {
				return dummyTx, &rpcError{msg: "some other error"}
			},
			func(_ context.Context, _ *types.Transaction) ([]*int, error) {
				t.Fatal("wait should not be called for unrelated errors")
				return nil, nil
			},
		)

		require.Error(t, err)
		require.True(t, nm.canceled, "nonce should be canceled for unrelated errors")
	})
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
			return nil, &rpcError{msg: "wait timeout"}
		},
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "wait timeout")
	require.True(t, nm.canceled, "nonce should be canceled when wait fails")
}
