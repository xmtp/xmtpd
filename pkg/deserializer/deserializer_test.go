package deserializer_test

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/deserializer"
	testutils "github.com/xmtp/xmtpd/pkg/testutils/envelopes"
)

func TestDeserializeMessagesFromCSV(t *testing.T) {
	f, err := os.Open("testdata/payloadSample.csv")
	require.NoError(t, err)
	defer func() {
		_ = f.Close()
	}()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	require.NoError(t, err)
	if len(records) < 1 {
		t.Fatal("CSV is empty or missing header")
	}

	for i, row := range records[1:] {
		rowNum := i + 1

		hexdata := row[0]
		isCommit := false
		if len(row) > 1 {
			isCommit, _ = strconv.ParseBool(row[1])
		}

		raw, err := hex.DecodeString(hexdata)
		require.NoError(t, err)

		var msg deserializer.MlsMessageIn
		if err := msg.TLSDeserialize(bytes.NewReader(raw)); err != nil {
			t.Errorf("row %d: failed to deserialize: %v", rowNum, err)
			continue
		}

		pub, ok := msg.Body.(*deserializer.PrivateMessageIn)
		require.True(t, ok)

		if isCommit {
			require.EqualValues(
				t,
				deserializer.ContentTypeCommit,
				pub.ContentType,
				"row %d: expected commit message, got non-commit",
				rowNum,
			)
		} else {
			require.NotEqualValues(t, deserializer.ContentTypeCommit, pub.ContentType,
				"row %d: expected non-commit message, got commit", rowNum)
		}
	}
}

func TestMinimalHex(t *testing.T) {
	tests := []struct {
		name         string
		hexPayload   string
		expectedType deserializer.ContentType
	}{
		{
			name:         "Commit",
			hexPayload:   testutils.MINIMAL_COMMIT_PAYLOAD,
			expectedType: deserializer.ContentTypeCommit,
		},
		{
			name:         "Application",
			hexPayload:   testutils.MINIMAL_APPLICATION_PAYLOAD,
			expectedType: deserializer.ContentTypeApplication,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := hex.DecodeString(tt.hexPayload)
			require.NoError(t, err, "Failed to decode hex payload")

			var msg deserializer.MlsMessageIn
			err = msg.TLSDeserialize(bytes.NewReader(raw))
			require.NoError(t, err, "Failed to deserialize MlsMessageIn")

			priv, ok := msg.Body.(*deserializer.PrivateMessageIn)
			require.True(t, ok, "Expected PrivateMessageIn type")

			require.EqualValues(t, tt.expectedType, priv.ContentType, "Unexpected content_type")
		})
	}
}
