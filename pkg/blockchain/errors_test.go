package blockchain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockchainErrorHashes(t *testing.T) {
	for errorCode, errorSignature := range protocolErrorsDictionary {
		// Compute keccak256 hash of the error signature.
		hash := crypto.Keccak256Hash([]byte(errorSignature.Error()))

		// Take first 4 bytes (error selector) and format as hex.
		selector := fmt.Sprintf("0x%x", hash[:4])

		// Verify the computed selector matches the error code.
		require.Equal(t, errorCode, selector,
			"error code mismatch for signature, expected: %s, got: %s",
			errorSignature, errorCode, selector)
	}
}

func TestTryExtractProtocolError_ValidSelector(t *testing.T) {
	// Standard format from debug_traceTransaction
	msg, err := tryExtractProtocolError(errors.New("execution reverted: 0xa88ee577"))
	require.NoError(t, err)
	assert.Equal(t, ErrNoChange, msg)
}

func TestTryExtractProtocolError_SelectorAtEndOfString(t *testing.T) {
	msg, err := tryExtractProtocolError(errors.New("0xa88ee577"))
	require.NoError(t, err)
	assert.Equal(t, ErrNoChange, msg)
}

func TestTryExtractProtocolError_UppercaseSelector(t *testing.T) {
	msg, err := tryExtractProtocolError(errors.New("execution reverted: 0xA88EE577"))
	require.NoError(t, err)
	assert.Equal(t, ErrNoChange, msg)
}

func TestTryExtractProtocolError_DoesNotMatchAddress(t *testing.T) {
	// An Ethereum address starts with 0x followed by 40 hex chars.
	// The first 8 hex chars happen to match a protocol error selector,
	// but should NOT be extracted because they're part of a longer hex string.
	_, err := tryExtractProtocolError(
		errors.New("invalid sender 0xa88ee577deadbeefcafebabe1234567890abcdef"),
	)
	assert.Error(t, err)
}

func TestTryExtractProtocolError_DoesNotMatchHash(t *testing.T) {
	// A transaction hash is 0x + 64 hex chars.
	_, err := tryExtractProtocolError(
		errors.New("tx 0xa88ee577deadbeefcafebabe1234567890abcdef1234567890abcdef12345678 failed"),
	)
	assert.Error(t, err)
}

func TestTryExtractProtocolError_NoHexCode(t *testing.T) {
	_, err := tryExtractProtocolError(errors.New("connection timeout"))
	assert.ErrorIs(t, err, ErrCodeNotFound)
}

func TestTryExtractProtocolError_UnknownSelector(t *testing.T) {
	_, err := tryExtractProtocolError(errors.New("execution reverted: 0xdeadbeef"))
	assert.ErrorIs(t, err, ErrCodeNotInDic)
}

func TestBlockchainError_IsNoChange(t *testing.T) {
	be := NewBlockchainError(errors.New("execution reverted: 0xa88ee577"))
	assert.True(t, be.IsNoChange())
}

func TestBlockchainError_AddressDoesNotTriggerIsNoChange(t *testing.T) {
	// The address starts with the NoChange selector bytes, but should not match
	be := NewBlockchainError(
		errors.New("call from 0xa88ee577deadbeefcafebabe1234567890abcdef failed"),
	)
	assert.False(t, be.IsNoChange())
}
