package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xmtp/xmtpd/pkg/proto/identity/associations"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"google.golang.org/protobuf/proto"
)

func Marshal(t *testing.T, msg proto.Message) []byte {
	bytes, err := proto.Marshal(msg)
	require.NoError(t, err)
	return bytes
}

func UnmarshalUnsignedOriginatorEnvelope(
	t *testing.T,
	bytes []byte,
) *message_api.UnsignedOriginatorEnvelope {
	unsignedOriginatorEnvelope := &message_api.UnsignedOriginatorEnvelope{}
	err := proto.Unmarshal(
		bytes,
		unsignedOriginatorEnvelope,
	)
	require.NoError(t, err)
	return unsignedOriginatorEnvelope
}

func CreateClientEnvelope() *message_api.ClientEnvelope {
	return &message_api.ClientEnvelope{
		Payload: nil,
		Aad: &message_api.AuthenticatedData{
			TargetOriginator: 1,
			TargetTopic:      []byte{0x5},
			LastSeen:         &message_api.VectorClock{},
		},
	}
}

func CreatePayerEnvelope(
	t *testing.T,
	clientEnv ...*message_api.ClientEnvelope,
) *message_api.PayerEnvelope {
	if len(clientEnv) == 0 {
		clientEnv = append(clientEnv, CreateClientEnvelope())
	}
	clientEnvBytes, err := proto.Marshal(clientEnv[0])
	require.NoError(t, err)

	return &message_api.PayerEnvelope{
		UnsignedClientEnvelope: clientEnvBytes,
		PayerSignature:         &associations.RecoverableEcdsaSignature{},
	}
}

func CreateOriginatorEnvelope(
	t *testing.T,
	originatorNodeID uint32,
	originatorSequenceID uint64,
	payerEnv ...*message_api.PayerEnvelope,
) *message_api.OriginatorEnvelope {
	if len(payerEnv) == 0 {
		payerEnv = append(payerEnv, CreatePayerEnvelope(t))
	}

	unsignedEnv := &message_api.UnsignedOriginatorEnvelope{
		OriginatorNodeId:     originatorNodeID,
		OriginatorSequenceId: originatorSequenceID,
		OriginatorNs:         0,
		PayerEnvelope:        payerEnv[0],
	}

	unsignedBytes, err := proto.Marshal(unsignedEnv)
	require.NoError(t, err)

	return &message_api.OriginatorEnvelope{
		UnsignedOriginatorEnvelope: unsignedBytes,
		Proof:                      nil,
	}
}
