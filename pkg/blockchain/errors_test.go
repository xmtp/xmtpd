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

func TestTryExtractProtocolError(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   error
		wantMatch error
	}{
		{
			name:      "selector after text",
			input:     "execution reverted: 0xa88ee577",
			wantMatch: ErrNoChange,
		},
		{
			name:      "selector with ABI params",
			input:     "0x84e23433" + "0000000000000000000000000000000000000000000000000000000000000001",
			wantMatch: ErrInvalidStartSequenceID,
		},
		{
			name:      "uppercase selector",
			input:     "execution reverted: 0xA88EE577",
			wantMatch: ErrNoChange,
		},
		{
			name:    "rejects address (40 hex chars)",
			input:   "invalid sender 0xa88ee577deadbeefcafebabe1234567890abcdef",
			wantErr: ErrCodeNotInDic,
		},
		{
			name:    "rejects tx hash (64 hex chars)",
			input:   "tx 0xa88ee577deadbeefcafebabe1234567890abcdef1234567890abcdef12345678 failed",
			wantErr: ErrCodeNotInDic,
		},
		{
			name:    "no hex code",
			input:   "connection timeout",
			wantErr: ErrCodeNotFound,
		},
		{
			name:    "unknown selector",
			input:   "execution reverted: 0xdeadbeef",
			wantErr: ErrCodeNotInDic,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := tryExtractProtocolError(errors.New(tc.input))
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantMatch, msg)
			}
		})
	}
}
