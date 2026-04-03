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

func TestLookupSelector(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "known selector only",
			input: "0xa88ee577",
			want:  ErrNoChange,
		},
		{
			name:  "known selector with ABI params",
			input: "0x84e234330000000000000000000000000000000000000000000000000000000000000001",
			want:  ErrInvalidStartSequenceID,
		},
		{
			name:  "uppercase selector",
			input: "0xA88EE577",
			want:  ErrNoChange,
		},
		{
			name:  "unknown selector",
			input: "0xdeadbeef",
			want:  nil,
		},
		{
			name:  "too short",
			input: "0xdead",
			want:  nil,
		},
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "no 0x prefix",
			input: "a88ee577",
			want:  ErrNoChange,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := lookupSelector(tc.input)
			assert.Equal(t, tc.want, result)
		})
	}
}

// mockDataError implements rpc.DataError for testing.
type mockDataError struct {
	msg  string
	data any
}

func (e *mockDataError) Error() string  { return e.msg }
func (e *mockDataError) ErrorData() any { return e.data }

func TestNewBlockchainErrorWithDataError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantIs      error
		wantMessage string
	}{
		{
			name: "extracts known selector from DataError",
			err: &mockDataError{
				msg:  "execution reverted",
				data: "0xa88ee577",
			},
			wantIs: ErrNoChange,
		},
		{
			name: "extracts selector with ABI params from DataError",
			err: &mockDataError{
				msg:  "execution reverted",
				data: "0x84e234330000000000000000000000000000000000000000000000000000000000000001",
			},
			wantIs: ErrInvalidStartSequenceID,
		},
		{
			name: "falls back to regex when DataError has unknown selector",
			err: &mockDataError{
				msg:  "execution reverted: 0xa88ee577",
				data: "0xdeadbeef",
			},
			wantIs: ErrNoChange,
		},
		{
			name: "falls back to regex when DataError data is not string",
			err: &mockDataError{
				msg:  "execution reverted: 0xa88ee577",
				data: 12345,
			},
			wantIs: ErrNoChange,
		},
		{
			name:        "falls back to regex for plain errors",
			err:         errors.New("execution reverted: 0xa88ee577"),
			wantIs:      ErrNoChange,
			wantMessage: "NoChange()",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			blockErr := NewBlockchainError(tc.err)
			require.NotNil(t, blockErr)
			assert.ErrorIs(t, blockErr, tc.wantIs)
		})
	}
}

func TestNewBlockchainErrorDataErrorTakesPrecedence(t *testing.T) {
	// DataError carries ErrInvalidStartSequenceID in data field,
	// but error message string contains ErrNoChange selector.
	// DataError path should win.
	err := &mockDataError{
		msg:  "execution reverted: 0xa88ee577", // ErrNoChange in message
		data: "0x84e23433",                     // ErrInvalidStartSequenceID in data
	}
	blockErr := NewBlockchainError(err)
	require.NotNil(t, blockErr)
	assert.ErrorIs(t, blockErr, ErrInvalidStartSequenceID)
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
