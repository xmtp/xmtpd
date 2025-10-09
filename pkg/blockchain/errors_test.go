package blockchain

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
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
