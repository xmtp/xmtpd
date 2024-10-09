package envelopes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/testutils"
	"google.golang.org/protobuf/proto"
)

func TestValidOriginatorEnvelope(t *testing.T) {
	originatorNodeId := uint32(1)
	originatorSequenceId := uint64(1)

	clientEnv := testutils.CreateClientEnvelope()
	payerEnvelope := testutils.CreatePayerEnvelope(t, clientEnv)
	originatorEnvelope := testutils.CreateOriginatorEnvelope(
		t,
		originatorNodeId,
		originatorSequenceId,
		payerEnvelope,
	)

	processed, err := NewOriginatorEnvelope(originatorEnvelope)
	require.NoError(t, err)
	require.Equal(t, originatorNodeId, processed.UnsignedOriginatorEnvelope.OriginatorNodeID())
	require.Equal(
		t,
		originatorSequenceId,
		processed.UnsignedOriginatorEnvelope.OriginatorSequenceID(),
	)

	serializedClientEnv, err := proto.Marshal(clientEnv)
	require.NoError(t, err)
	serializedClientEnvAfterParse, err := processed.UnsignedOriginatorEnvelope.PayerEnvelope.ClientEnvelope.Bytes()
	require.NoError(t, err)
	require.Equal(t, serializedClientEnv, serializedClientEnvAfterParse)
}

func TestSerialize(t *testing.T) {
	originatorNodeId := uint32(1)
	originatorSequenceId := uint64(1)

	clientEnv := testutils.CreateClientEnvelope()
	payerEnvelope := testutils.CreatePayerEnvelope(t, clientEnv)
	originatorEnvelope := testutils.CreateOriginatorEnvelope(
		t,
		originatorNodeId,
		originatorSequenceId,
		payerEnvelope,
	)

	serializedFromProto, err := proto.Marshal(originatorEnvelope)
	require.NoError(t, err)

	originatorStruct, err := NewOriginatorEnvelope(originatorEnvelope)
	require.NoError(t, err)
	serializedFromStruct, err := originatorStruct.Bytes()
	require.NoError(t, err)

	require.Equal(t, serializedFromProto, serializedFromStruct)
}

func TestInvalidOriginatorEnvelope(t *testing.T) {
	_, err := NewOriginatorEnvelope(nil)
	require.Error(t, err)

	empty := &message_api.OriginatorEnvelope{}
	_, err = NewOriginatorEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidUnsignedOriginatorEnvelope(t *testing.T) {
	_, err := NewUnsignedOriginatorEnvelope(nil)
	require.Error(t, err)

	empty := &message_api.UnsignedOriginatorEnvelope{}
	_, err = NewUnsignedOriginatorEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidPayerEnvelope(t *testing.T) {
	_, err := NewPayerEnvelope(nil)
	require.Error(t, err)

	empty := &message_api.PayerEnvelope{}
	_, err = NewPayerEnvelope(empty)
	require.Error(t, err)
}

func TestInvalidClientEnvelope(t *testing.T) {
	_, err := NewClientEnvelope(nil)
	require.Error(t, err)

	empty := &message_api.ClientEnvelope{}
	_, err = NewClientEnvelope(empty)
	require.Error(t, err)
}
